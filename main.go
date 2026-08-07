package main

import (
	"image"
	"image/color"
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
			app.Size(unit.Dp(400), unit.Dp(320)),
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
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
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
				return drawNavItem(gtx, th, &state.btnServer, "Server", state.selectedTab == TabServer)
			}),

			layout.Rigid(layout.Spacer{Height: unit.Dp(6)}.Layout),

			// Tab 2: WSL Node
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return drawNavItem(gtx, th, &state.btnWSL, "WSL Node", state.selectedTab == TabWSL)
			}),
		)
	})
}

func drawNavItem(gtx layout.Context, th *material.Theme, btn *widget.Clickable, title string, isSelected bool) layout.Dimensions {
	return btn.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		// Bảng màu giống giao diện mẫu
		blueColor := color.NRGBA{R: 0x00, G: 0x7A, B: 0xFF, A: 0xFF} // Màu xanh Blue
		activeBg := color.NRGBA{R: 0xC6, G: 0xC6, B: 0xC6, A: 0xFF}  // Màu xám đậm bo góc khi Active
		textColor := color.NRGBA{R: 0x1A, G: 0x1A, B: 0x1A, A: 0xFF}

		return layout.Stack{}.Layout(gtx,
			layout.Expanded(func(gtx layout.Context) layout.Dimensions {
				if isSelected {
					rr := gtx.Dp(unit.Dp(8)) // Độ bo góc 8dp
					defer clip.RRect{
						Rect: image.Rectangle{Max: gtx.Constraints.Min},
						SE:   rr, SW: rr, NW: rr, NE: rr,
					}.Push(gtx.Ops).Pop()
					paint.ColorOp{Color: activeBg}.Add(gtx.Ops)
					paint.PaintOp{}.Add(gtx.Ops)
				}
				return layout.Dimensions{Size: gtx.Constraints.Min}
			}),

			// Lớp trên: Nội dung gồm Icon và Chữ
			layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{
					Top:    unit.Dp(10),
					Bottom: unit.Dp(10),
					Left:   unit.Dp(12),
					Right:  unit.Dp(12),
				}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{
						Axis:      layout.Horizontal,
						Alignment: layout.Middle,
					}.Layout(gtx,
						// Icon đại diện màu xanh
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							size := gtx.Dp(unit.Dp(20))
							defer clip.RRect{
								Rect: image.Rectangle{Max: image.Pt(size, size)},
								SE:   gtx.Dp(4), SW: gtx.Dp(4), NW: gtx.Dp(4), NE: gtx.Dp(4),
							}.Push(gtx.Ops).Pop()
							paint.ColorOp{Color: blueColor}.Add(gtx.Ops)
							paint.PaintOp{}.Add(gtx.Ops)
							return layout.Dimensions{Size: image.Pt(size, size)}
						}),

						// Khoảng cách giữa Icon và Text
						layout.Rigid(layout.Spacer{Width: unit.Dp(12)}.Layout),

						// Tiêu đề Tab
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							lbl := material.Body1(th, title)
							lbl.Color = textColor
							lbl.TextSize = unit.Sp(15)
							return lbl.Layout(gtx)
						}),
					)
				})
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
