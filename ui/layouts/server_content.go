package layouts

import (
	"labor-app/ui/state"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

func DrawServerContent(gtx layout.Context, state *state.AppState, th *material.Theme) layout.Dimensions {
	return layout.UniformInset(unit.Dp(24)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		input := material.Editor(th, &state.NameInput, "Enter your name 1111")
		return input.Layout(gtx)
	})
}
