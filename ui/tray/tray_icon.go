package uitray

import (
	"context"
	"fmt"
	"labor-app/platform/darwin"
	"labor-app/ui/state"
	"log"
	"os"
	"time"

	"github.com/gogpu/systray"

	_ "embed"
)

// Embed the icon.ico (in the same folder with main.go)
//
//go:embed icon.ico
var iconBytes []byte

func SetupTray(ctx context.Context, host *state.HostState, tray *systray.SystemTray, onClickAdmin func()) {
	menu := systray.NewMenu()

	menu.Add("Open", onClickAdmin)
	serverItem := menu.Add("Server is loading", func() {
		host.Mu.Lock()
		isOnline := host.IsOnline
		host.Mu.Unlock()

		if isOnline {
			host.ServerSignal <- false
		} else {
			host.ServerSignal <- true
			fmt.Println("Clicked to turn on the server nowwww")
		}
	})

	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for _ = range ticker.C {
			select {
			case <-ctx.Done():
				return
			default:
				host.Mu.Lock()
				isOnline := host.IsOnline
				host.Mu.Unlock()
				if isOnline {
					serverItem.SetLabel("Turn off server")
				} else {
					serverItem.SetLabel("Turn on server")
				}
			}
		}
	}()

	isAutoStart := darwin.IsStartAtLoginEnabled()
	var autoStartItem *systray.MenuItem

	// 2. Tạo menu item checkbox
	autoStartItem = menu.AddCheckbox("Start application on login", isAutoStart, func() {
		go func() {
			if isAutoStart {
				if err := darwin.DisableStartAtLogin(); err != nil {
					log.Printf("[AutoStart] Disable error: %v", err)
					return
				}

				autoStartItem.SetChecked(false)
				isAutoStart = false
				log.Println("[AutoStart] Disabled")
			} else {
				if err := darwin.EnableStartAtLogin(); err != nil {
					log.Printf("[AutoStart] Enable error: %v", err)
					return
				}
				autoStartItem.SetChecked(true)
				isAutoStart = true
				log.Println("[AutoStart] Enabled")
			}
		}()
	})
	menu.Add("Quit", func() {
		os.Exit(0)
	})

	tray.
		SetTemplateIcon(iconBytes).
		SetTooltip("Laboratory management").
		SetMenu(menu).
		Show()
}
