package main

/*
#cgo darwin CFLAGS: -x objective-c
#cgo darwin LDFLAGS: -framework Cocoa
#import <Cocoa/Cocoa.h>
#import <objc/runtime.h>

// Con trỏ lưu lại hàm gốc của macOS
static IMP original_setActivationPolicy;

// HÀM GIẢ MẠO: Bất kể ai (kể cả Gio UI) gọi đổi Policy, ta đều ép nó về Accessory (1)
BOOL hook_setActivationPolicy(id self, SEL _cmd, NSApplicationActivationPolicy policy) {
    BOOL (*original)(id, SEL, NSApplicationActivationPolicy) = (void *)original_setActivationPolicy;
    // Bỏ qua biến policy được truyền vào, luôn bắt macOS chạy Accessory
    return original(self, _cmd, NSApplicationActivationPolicyAccessory);
}

// Hàm này sẽ đánh tráo hàm gốc của hệ điều hành
void forceAccessoryForever() {
    Method method = class_getInstanceMethod([NSApplication class], @selector(setActivationPolicy:));
    original_setActivationPolicy = method_getImplementation(method);
    method_setImplementation(method, (IMP)hook_setActivationPolicy);

    // Ép policy ngay lúc này luôn để chắc chắn
    dispatch_async(dispatch_get_main_queue(), ^{
        [NSApp setActivationPolicy:NSApplicationActivationPolicyAccessory];
    });
}

// Ép macOS focus vào cửa sổ (vì app ẩn Dock sẽ bị mất khả năng tự focus)
void forceActivateApp() {
    dispatch_async(dispatch_get_main_queue(), ^{
        [NSApp activateIgnoringOtherApps:YES];
    });
}

#include "center_mac.h"
*/
import "C"
import (
	"context"
	"fmt"
	"labor-app/config"
	"labor-app/ui/layouts"
	"labor-app/ui/state"
	uitray "labor-app/ui/tray"
	"log"
	"sync"
	"time"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget/material"

	"github.com/gogpu/systray"
)

var (
	winMutex  sync.Mutex
	activeWin *app.Window
)

func main() {
	C.forceAccessoryForever()
	C.installWindowCentering()

	appCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	host := &state.HostState{
		Name:         "Main Server (Yuu, Kubernetes, WSL, Window Server)",
		Address:      "Yuu.local",
		ServerSignal: make(chan bool),
	}

	go host.PingToServerLoop(appCtx)
	go host.FetchWSLNodesLoop(appCtx)
	go host.HandleServerSignal(appCtx)

	tray := systray.New()
	uitray.SetupTray(appCtx, host, tray, func() {
		openGioWindow(host)
	})

	app.Main()
}

// openGioWindow mở cửa sổ Gio UI hoặc đưa cửa sổ đã có lên trên cùng
func openGioWindow(host *state.HostState) {
	winMutex.Lock()

	// Nếu cửa sổ đã mở -> Kéo lên trên cùng
	if activeWin != nil {
		activeWin.Perform(system.ActionRaise)
		winMutex.Unlock()

		// Ép macOS focus
		C.forceActivateApp()
		return
	}

	// Nếu chưa mở -> Tạo cửa sổ mới
	w := new(app.Window)
	w.Option(
		app.Title("Laboratory management"),
		app.Size(unit.Dp(820), unit.Dp(404)),
		app.MinSize(unit.Dp(820), unit.Dp(404)),
		app.MaxSize(unit.Dp(820), unit.Dp(404)),
	)

	activeWin = w
	winMutex.Unlock()

	// Ép macOS focus vào cửa sổ mới
	C.forceActivateApp()

	// Chạy vòng lặp render window
	go func() {
		if err := run(w, host); err != nil {
			log.Println("Window closed with error:", err)
		}

		winMutex.Lock()
		activeWin = nil
		winMutex.Unlock()
	}()
}

func run(w *app.Window, host *state.HostState) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	th := material.NewTheme()
	th.Face = "Google Sans"
	var ops op.Ops

	if err := config.EnsureConfig(); err != nil {
		log.Fatal(err)
	}

	if err := config.Load(); err != nil {
		log.Fatal(err)
	}

	singlePageApp := layouts.NewSinglePageApp(host)

	// background worker to check the hosts
	go fetchServerUI(ctx, host, w)
	go fetchWSLUI(ctx, w, host)
	go startLogListener(w, singlePageApp)

	// handle frame events and other events
	for {
		e := w.Event()
		switch e := e.(type) {
		case app.DestroyEvent:
			cancel()
			return nil
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

func fetchServerUI(ctx context.Context, host *state.HostState, w *app.Window) {
	host.Mu.Lock()
	currentHostStatus := host.IsOnline
	host.Mu.Unlock()
	for {
		select {
		case <-ctx.Done():
			return

		default:
			time.Sleep(3 * time.Second)
			host.Mu.Lock()
			if currentHostStatus == host.IsOnline {
				host.Mu.Unlock()
				continue
			}

			currentHostStatus = host.IsOnline
			host.Mu.Unlock()
			fmt.Println("Invalidate UI")
			w.Invalidate()
		}
	}
}

func fetchWSLUI(ctx context.Context, w *app.Window, host *state.HostState) {
	host.Mu.Lock()
	currentWslList := len(host.Wsls)
	host.Mu.Unlock()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			time.Sleep(3 * time.Second)
			host.Mu.Lock()
			if currentWslList == len(host.Wsls) {
				host.Mu.Unlock()
				continue
			}

			currentWslList = len(host.Wsls)
			host.Mu.Unlock()
			w.Invalidate()
		}
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
