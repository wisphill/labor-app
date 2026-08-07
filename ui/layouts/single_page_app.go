package layouts

import (
	"fmt"
	"image/color"
	"labor-app/ui/components"
	"labor-app/ui/state"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type SinglePageApp struct {
	hosts []*state.HostState
	list  widget.List
}

func NewSinglePageApp(hostStates []*state.HostState) *SinglePageApp {
	return &SinglePageApp{
		list: widget.List{
			List: layout.List{
				Axis: layout.Vertical,
			},
		},
		hosts: hostStates,
	}
}

func (app *SinglePageApp) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return material.List(th, &app.list).Layout(gtx, len(app.hosts), func(gtx layout.Context, i int) layout.Dimensions {
		return app.layoutHostRow(gtx, th, app.hosts[i])
	})
}

// Layout vẽ từng dòng Server
func (app *SinglePageApp) layoutHostRow(gtx layout.Context, th *material.Theme, host *state.HostState) layout.Dimensions {
	// 1. Kiểm tra xự kiện Click Nút Bật / Tắt
	if host.BtnTurnOn.Clicked(gtx) {
		//cmdStr := strings.Join(host.TurnOnScript.Commands, " && ")
		// go OpenTerminalApp(cmdStr)
	}

	if host.BtnShutdown.Clicked(gtx) {
		// cmdStr := strings.Join(host.ShutdownScript.Commands, " && ")
		// go OpenTerminalApp(cmdStr)
	}

	// 2. Lấy dữ liệu an toàn từ Mutex
	host.Mu.Lock()
	isOnline := host.IsOnline
	rtt := host.PingRTT
	name := host.Name
	address := host.Address
	host.Mu.Unlock()

	// 3. Dựng Layout hàng
	return layout.Inset{
		Top: unit.Dp(10), Bottom: unit.Dp(10),
		Left: unit.Dp(16), Right: unit.Dp(16),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{
			Axis:      layout.Horizontal,
			Alignment: layout.Middle,
		}.Layout(gtx,
			// A. Đèn báo trạng thái On/Off (Xanh/Đỏ)
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return components.DrawStatusBadge(gtx, isOnline)
			}),

			layout.Rigid(layout.Spacer{Width: unit.Dp(12)}.Layout),

			// B. Tên Server & Địa chỉ IP "x.x.x.x"
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(material.Body1(th, name).Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						lbl := material.Caption(th, address)
						lbl.Color = color.NRGBA{R: 130, G: 130, B: 130, A: 255}
						return lbl.Layout(gtx)
					}),
				)
			}),

			// C. Số Ping RTT
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				pingStr := "Timeout"
				pingColor := color.NRGBA{R: 220, G: 50, B: 50, A: 255} // Red

				if isOnline {
					pingStr = fmt.Sprintf("%d ms", rtt.Milliseconds())
					pingColor = color.NRGBA{R: 40, G: 160, B: 60, A: 255} // Green
				}

				lbl := material.Body2(th, pingStr)
				lbl.Color = pingColor
				return layout.Inset{Right: unit.Dp(16)}.Layout(gtx, lbl.Layout)
			}),

			// D. Nút Bật (Turn On Terminal Script)
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				btn := material.Button(th, &host.BtnTurnOn, "⚡ Bật")
				btn.Background = color.NRGBA{R: 46, G: 204, B: 113, A: 255}
				return layout.Inset{Right: unit.Dp(8)}.Layout(gtx, btn.Layout)
			}),

			// E. Nút Tắt (Shutdown Terminal Script)
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				btn := material.Button(th, &host.BtnShutdown, "🛑 Tắt")
				btn.Background = color.NRGBA{R: 231, G: 76, B: 60, A: 255}
				return btn.Layout(gtx)
			}),
		)
	})
}
