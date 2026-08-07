package layouts

import (
	"image/color"
	"labor-app/ui/components"
	"labor-app/ui/state"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

func DrawSidebar(gtx layout.Context, th *material.Theme, state *state.AppState) layout.Dimensions {
	// background for the sidebar
	sidebarBg := color.NRGBA{R: 0xE8, G: 0xE8, B: 0xE8, A: 0xFF}
	paint.FillShape(gtx.Ops, sidebarBg, clip.Rect{Max: gtx.Constraints.Max}.Op())

	return layout.UniformInset(unit.Dp(12)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			// Tab 1: Server
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return components.DrawNavItem(gtx, th, &state.BtnServer, "Server", state.SelectedTab == components.TabServer)
			}),

			layout.Rigid(layout.Spacer{Height: unit.Dp(6)}.Layout),

			// Tab 2: WSL Node
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return components.DrawNavItem(gtx, th, &state.BtnWSL, "WSL Node", state.SelectedTab == components.TabWSL)
			}),
		)
	})
}
