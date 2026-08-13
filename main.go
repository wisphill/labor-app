package main

/*
#cgo darwin LDFLAGS: -framework Cocoa
#include "center_mac.h"
*/
import "C"
import (
	"fmt"
	server "labor-app/cmd/host"
	"labor-app/config"
	"labor-app/ui/layouts"
	"labor-app/ui/state"
	uitray "labor-app/ui/tray"
	"log"
	"os"
	"sync"
	"time"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget/material"

	"github.com/gogpu/systray"
)

func main() {
	C.installWindowCentering()

	tray := systray.New()
	uitray.SetupTray(tray)

	go func() {
		w := new(app.Window)
		w.Option(
			app.Title("Laboratory management"),
			app.Size(unit.Dp(820), unit.Dp(404)),
			app.MinSize(unit.Dp(820), unit.Dp(404)),
			app.MaxSize(unit.Dp(820), unit.Dp(404)),
		)
		if err := run(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()

	app.Main()
}

func run(w *app.Window) error {
	th := material.NewTheme()
	var ops op.Ops

	if err := config.EnsureConfig(); err != nil {
		log.Fatal(err)
	}

	if err := config.Load(); err != nil {
		log.Fatal(err)
	}

	hostStates := []*state.HostState{
		{
			Name:    "Main Server (Yuu, Kubernetes, WSL, Window Server)",
			Address: "Yuu.local",
		},
	}
	singlePageApp := layouts.NewSinglePageApp(hostStates)

	// background worker to check the hosts
	go pingToServer(hostStates, w)
	go fetchWSLNodes(w, hostStates)
	go startLogListener(w, singlePageApp)

	// handle frame events and other events
	for {
		e := w.Event()
		switch e := e.(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			// set the layout
			layout.UniformInset(unit.Dp(0)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return singlePageApp.Layout(gtx, th)
			})

			// draw frame to the gpu
			e.Frame(gtx.Ops)
		}
	}
}

func pingToServer(hostItems []*state.HostState, w *app.Window) {
	for {
		var wg sync.WaitGroup
		for _, host := range hostItems {
			host.Mu.Lock()
			addr := host.Address
			host.Mu.Unlock()

			wg.Add(1)
			go func(h *state.HostState, address string) {
				defer wg.Done()
				online, rtt := server.PingOS(address, 1500*time.Millisecond)

				h.Mu.Lock()
				h.IsOnline = online
				h.PingRTT = rtt
				h.Mu.Unlock()
			}(host, addr)
		}
		wg.Wait()

		w.Invalidate()
		time.Sleep(3 * time.Second)
	}
}

func fetchWSLNodes(window *app.Window, hostItems []*state.HostState) {
	for {
		var wg sync.WaitGroup
		for _, host := range hostItems {
			host.Mu.Lock()
			addr := host.Address
			host.Mu.Unlock()

			wg.Add(1)
			go func(h *state.HostState, address string) {
				defer wg.Done()
				wslNodes, err := server.GetRunningWSLNodes()
				if err != nil {
					fmt.Println("Error while getting the WSL nodes")
					return
				}

				h.Mu.Lock()
				h.Wsls = make([]*state.WSLState, 0)
				for _, wslNode := range wslNodes {
					h.Wsls = append(h.Wsls, &state.WSLState{
						Name: wslNode,
					})
				}

				h.Mu.Unlock()
			}(host, addr)
		}
		wg.Wait()
		window.Invalidate()
		time.Sleep(3 * time.Second)
	}
}

func startLogListener(window *app.Window, pageApp *layouts.SinglePageApp) {
	go func() {
		// Vòng lặp range tự động thoát khi logChan bị gọi close()
		for msg := range pageApp.LogChan {
			if msg == "" {
				pageApp.ShowLogBar = false // Gửi "" -> Ẩn log bar
			} else {
				pageApp.DisplayedLogMsg = msg
				pageApp.ShowLogBar = true // Gửi text -> Hiện log bar
			}

			window.Invalidate() // Đánh thức Gio vẽ lại UI
		}

		// Khi channel bị close(logChan) hoàn toàn -> Tự ẩn log bar
		pageApp.ShowLogBar = false
		if window != nil {
			window.Invalidate()
		}
	}()
}
