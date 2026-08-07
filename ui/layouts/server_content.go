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

type ServerTab struct {
	hosts []*state.HostState // manage state of items
	list  widget.List
}

func DrawServerContent(gtx layout.Context,
	state *state.AppState,
	hostItemsState []*state.HostState,
	th *material.Theme) layout.Dimensions {
	servers := NewServerTab(hostItemsState)
	serverLayout := servers.Layout(gtx, th)
	return layout.UniformInset(unit.Dp(24)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return serverLayout
	})
}

func NewServerTab(hostItemsState []*state.HostState) *ServerTab {
	st := &ServerTab{
		list: widget.List{
			List: layout.List{
				Axis: layout.Vertical,
			},
		},
		hosts: hostItemsState,
	}

	return st
}

// Layout hiển thị Server Tab
func (st *ServerTab) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return material.List(th, &st.list).Layout(gtx, len(st.hosts), func(gtx layout.Context, i int) layout.Dimensions {
		return st.layoutHostRow(gtx, th, st.hosts[i])
	})
}

// Layout cho từng dòng Host
func (st *ServerTab) layoutHostRow(gtx layout.Context, th *material.Theme, host *state.HostState) layout.Dimensions {
	host.Mu.Lock()
	isOnline := host.IsOnline
	rtt := host.PingRTT
	host.Mu.Unlock()

	// Khung Container từng hàng
	return layout.Inset{
		Top: unit.Dp(8), Bottom: unit.Dp(8),
		Left: unit.Dp(16), Right: unit.Dp(16),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{
			Axis:      layout.Horizontal,
			Alignment: layout.Middle,
		}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return components.DrawStatusBadge(gtx, isOnline)
			}),

			// Space
			layout.Rigid(layout.Spacer{Width: unit.Dp(12)}.Layout),

			// 2. Tên Server & IP
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(material.Body1(th, host.Name).Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						lbl := material.Caption(th, host.Address)
						lbl.Color = color.NRGBA{R: 120, G: 120, B: 120, A: 255}
						return lbl.Layout(gtx)
					}),
				)
			}),

			// 3. Số Ping (Latency ms)
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				pingStr := "OFF"
				pingColor := color.NRGBA{R: 150, G: 150, B: 150, A: 255}

				if isOnline {
					pingStr = fmt.Sprintf("%d ms", rtt.Milliseconds())
					pingColor = color.NRGBA{R: 0, G: 150, B: 0, A: 255} // Xanh lá
				} else {
					pingStr = "Timeout"
					pingColor = color.NRGBA{R: 200, G: 0, B: 0, A: 255} // Đỏ
				}

				lbl := material.Body2(th, pingStr)
				lbl.Color = pingColor
				return layout.Inset{Right: unit.Dp(16)}.Layout(gtx, lbl.Layout)
			}),
		)
	})
}
