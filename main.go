package main

import (
	"labor-app/ui/components"
	"labor-app/ui/layouts"
	"labor-app/ui/state"
	"log"
	"os"

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
	var state state.AppState

	// handle global app configuration
	state.NameInput.SingleLine = true

	// handle frame events and other events
	for {
		e := w.Event()
		switch e := e.(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			if state.BtnServer.Clicked(gtx) {
				state.SelectedTab = components.TabServer
			} else if state.BtnWSL.Clicked(gtx) {
				state.SelectedTab = components.TabWSL
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
						return layouts.DrawSidebar(gtx, th, &state)
					}),

					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{
							Axis: layout.Vertical,
						}.Layout(gtx,
							layout.Rigid(layout.Spacer{Height: unit.Dp(16)}.Layout),
							layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
								gtx.Constraints.Min.X = gtx.Dp(unit.Dp(600))
								gtx.Constraints.Min.Y = gtx.Dp(unit.Dp(600))
								return drawMain(gtx, &state, th)
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

func drawMain(gtx layout.Context, state *state.AppState, th *material.Theme) layout.Dimensions {
	if state.SelectedTab == components.TabServer {
		return layouts.DrawServerContent(gtx, state, th)
	} else if state.SelectedTab == components.TabWSL {
		return layouts.DrawWSLNodeContent(gtx, state, th)
	}

	return layout.Dimensions{}
}
