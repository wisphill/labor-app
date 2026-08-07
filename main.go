package main

import (
	"labor-app/cmd/server_check"
	"labor-app/ui/components"
	"labor-app/ui/layouts"
	"labor-app/ui/state"
	"log"
	"os"
	"sync"
	"time"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

func main() {
	go func() {
		w := new(app.Window)
		w.Option(
			app.Title("Labor management"),
			app.Size(unit.Dp(1024), unit.Dp(768)),
			app.MinSize(unit.Dp(820), unit.Dp(600)),
			app.Maximized.Option(),
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
	var appState state.AppState

	// handle global app configuration
	appState.NameInput.SingleLine = true

	hostItems := []*state.HostState{
		{Name: "Google DNS", Address: "8.8.8.8"},
		{Name: "Local Router", Address: "192.168.1.1"},
		{Name: "Desktop Office", Address: "Yuu.local"},
		{Name: "Server Offline Test", Address: "192.168.1.250"},
	}

	// background worker to check the hosts
	go pingToServer(hostItems, w)

	// handle frame events and other events
	for {
		e := w.Event()
		switch e := e.(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			if appState.BtnServer.Clicked(gtx) {
				appState.SelectedTab = components.TabServer
			} else if appState.BtnWSL.Clicked(gtx) {
				appState.SelectedTab = components.TabWSL
			}

			// set the layout
			layout.UniformInset(unit.Dp(0)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				// Vertical flex layout
				return layout.Flex{
					Axis: layout.Horizontal,
				}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Min.X = gtx.Dp(unit.Dp(220))
						gtx.Constraints.Max.X = gtx.Dp(unit.Dp(220))
						return layouts.DrawSidebar(gtx, th, &appState)
					}),

					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{
							Axis: layout.Vertical,
						}.Layout(gtx,
							layout.Rigid(layout.Spacer{Height: unit.Dp(16)}.Layout),
							layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
								gtx.Constraints.Min.X = gtx.Dp(unit.Dp(600))
								gtx.Constraints.Min.Y = gtx.Dp(unit.Dp(600))
								return drawMain(gtx, &appState, hostItems, th)
							}),
						)
					}),
				)
			})

			// draw frame to the gpu
			e.Frame(gtx.Ops)
		}
	}
}

func drawMain(gtx layout.Context, appState *state.AppState, hostItems []*state.HostState, th *material.Theme) layout.Dimensions {
	if appState.SelectedTab == components.TabServer {
		return layouts.DrawServerContent(gtx, appState, hostItems, th)
	} else if appState.SelectedTab == components.TabWSL {
		return layouts.DrawWSLNodeContent(gtx, appState, th)
	}

	return layout.Dimensions{}
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
				online, rtt := server_check.PingOS(address, 1500*time.Millisecond)

				h.Mu.Lock()
				h.IsOnline = online
				h.PingRTT = rtt
				h.Mu.Unlock()
			}(host, addr)
		}
		wg.Wait()

		w.Invalidate()
		time.Sleep(10 * time.Second)
	}
}
