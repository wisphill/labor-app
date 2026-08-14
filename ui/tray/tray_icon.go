package uitray

import (
	server "labor-app/cmd/host"
	"labor-app/platform/darwin"
	"log"
	"os"

	"github.com/gogpu/systray"

	_ "embed"
)

// Embed the icon.ico (in the same folder with main.go)
//
//go:embed icon.ico
var iconBytes []byte

func SetupTray(tray *systray.SystemTray, onClickAdmin func()) {
	menu := systray.NewMenu()

	menu.Add("Open", onClickAdmin)
	menu.Add("Turn on server", func() {
		go server.TurnOnServer()
	})
	menu.Add("Turn off server", func() {
		go server.TurnOffServer()
	})

	// 1. Kiểm tra trạng thái hiện tại (nếu có hàm check)
	isAutoStart := darwin.IsStartAtLoginEnabled() // Hoặc mặc định là true/false

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
		SetIcon(iconBytes).
		SetTooltip("Laboratory management").
		SetMenu(menu).
		Show()
}
