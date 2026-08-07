package main

import (
	"image/color"
	"labor-app/ui"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type Tab int

const (
	TabServer Tab = iota
	TabWSL
)

type AppState struct {
	btn         widget.Clickable
	nameInput   widget.Editor
	selectedTab Tab
	btnServer   widget.Clickable
	btnWSL      widget.Clickable
}

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
	var state AppState

	// handle global app configuration
	state.nameInput.SingleLine = true

	// handle frame events and other events
	for {
		e := w.Event()
		switch e := e.(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			if state.btnServer.Clicked(gtx) {
				state.selectedTab = TabServer
			} else if state.btnWSL.Clicked(gtx) {
				state.selectedTab = TabWSL
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
						return drawSidebar(gtx, th, &state)
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

func drawSidebar(gtx layout.Context, th *material.Theme, state *AppState) layout.Dimensions {
	// background for the sidebar
	sidebarBg := color.NRGBA{R: 0xE8, G: 0xE8, B: 0xE8, A: 0xFF}
	paint.FillShape(gtx.Ops, sidebarBg, clip.Rect{Max: gtx.Constraints.Max}.Op())

	return layout.UniformInset(unit.Dp(12)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			// Tab 1: Server
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return ui.DrawNavItem(gtx, th, &state.btnServer, "Server", state.selectedTab == TabServer)
			}),

			layout.Rigid(layout.Spacer{Height: unit.Dp(6)}.Layout),

			// Tab 2: WSL Node
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return ui.DrawNavItem(gtx, th, &state.btnWSL, "WSL Node", state.selectedTab == TabWSL)
			}),
		)
	})
}

func drawMain(gtx layout.Context, state *AppState, th *material.Theme) layout.Dimensions {
	if state.selectedTab == TabServer {
		return layout.UniformInset(unit.Dp(24)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			input := material.Editor(th, &state.nameInput, "Enter your name 1111")
			return input.Layout(gtx)
		})
	} else if state.selectedTab == TabWSL {
		return layout.UniformInset(unit.Dp(24)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			input := material.Editor(th, &state.nameInput, "Enter your name 2222")
			return input.Layout(gtx)
		})
	}

	return layout.Dimensions{}
}
