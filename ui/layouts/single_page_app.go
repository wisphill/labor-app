package layouts

import (
	"fmt"
	"image"
	"image/color"
	server "labor-app/cmd/host"
	"labor-app/ui/components"
	"labor-app/ui/state"
	"log"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type SinglePageApp struct {
	hosts        []*state.HostState
	list         widget.List
	shutdownIcon *components.SVGRenderer
	serverIcon   *components.SVGRenderer

	wslList layout.List
}

func NewSinglePageApp(hostStates []*state.HostState) *SinglePageApp {
	// Khởi tạo ở ngoài vòng lặp sự kiện (Event Loop)
	shutdownIcon, err := components.LoadSVG("assets/shutdown.svg", 24, 24, color.NRGBA{R: 255, G: 255, B: 255, A: 255})
	if err != nil {
		log.Fatalf("Error whil loading SVG: %v", err)
	}
	serverIcon, err := components.LoadSVG("assets/server.svg", 30, 30, color.NRGBA{R: 0, G: 0, B: 0, A: 255})
	if err != nil {
		log.Fatalf("Error whil loading SVG: %v", err)
	}

	return &SinglePageApp{
		list: widget.List{
			List: layout.List{
				Axis: layout.Vertical,
			},
		},
		hosts:        hostStates,
		shutdownIcon: shutdownIcon,
		serverIcon:   serverIcon,
	}
}

func (app *SinglePageApp) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	// setting global configs for the layout
	app.wslList.Axis = layout.Horizontal

	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return material.List(th, &app.list).Layout(
				gtx,
				len(app.hosts),
				func(gtx layout.Context, i int) layout.Dimensions {
					return app.layoutHostRow(gtx, th, app.hosts[i])
				},
			)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{
				Left: unit.Dp(16),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				lbl := material.H6(th, "Running WSL Nodes")
				lbl.Font.Weight = 600
				return lbl.Layout(gtx)
			})
		}),
		// WSL nodes - horizontal
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return app.layoutWSLNodes(gtx, th)
		}),
	)
}

// Layout vẽ từng dòng Server
func (app *SinglePageApp) layoutHostRow(gtx layout.Context, th *material.Theme, host *state.HostState) layout.Dimensions {
	host.Mu.Lock()
	isOnline := host.IsOnline
	rtt := host.PingRTT
	name := host.Name
	address := host.Address
	host.Mu.Unlock()

	isPowerButtonClicked := host.BtnPower.Clicked(gtx) // consume the event
	if isPowerButtonClicked {
		if isOnline {
			go server.TurnOffServer()
		} else if !isOnline {
			go server.TurnOnServer()
		}
	}

	return layout.Inset{
		Top: unit.Dp(10), Bottom: unit.Dp(10),
		Left: unit.Dp(16), Right: unit.Dp(16),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{
			Axis:      layout.Horizontal,
			Alignment: layout.Middle,
		}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return components.DrawStatusBadge(gtx, isOnline)
			}),

			layout.Rigid(layout.Spacer{
				Width: unit.Dp(10),
			}.Layout),

			// Name + IP
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{
					Axis: layout.Vertical,
				}.Layout(gtx,
					// ┌──────────────────────┐
					// │ phil             🖥 │
					// └──────────────────────┘
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{
							Axis:      layout.Horizontal,
							Alignment: layout.Middle,
						}.Layout(gtx,

							// Server name
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								lbl := material.Body2(th, name)
								lbl.Font.Weight = font.Medium
								lbl.TextSize = unit.Sp(12)
								lbl.LineHeight = unit.Sp(15)
								return lbl.Layout(gtx)
							}),

							// Khoảng cách name -> icon
							layout.Rigid(layout.Spacer{
								Width: unit.Dp(5),
							}.Layout),

							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return layout.Inset{
									Right: unit.Dp(4),
								}.Layout(gtx, app.serverIcon.Layout)
							}),
						)
					}),

					// IP address
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{
							Top: unit.Dp(-5),
						}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							lbl := material.Caption(th, address)
							lbl.Color = color.NRGBA{
								R: 130,
								G: 130,
								B: 130,
								A: 255,
							}
							lbl.LineHeight = unit.Sp(14)

							return lbl.Layout(gtx)
						})
					}),
				)
			}),

			// C. Ping RTT
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				pingStr := "Timeout"
				pingColor := color.NRGBA{R: 220, G: 50, B: 50, A: 255} // Red

				if isOnline {
					pingStr = fmt.Sprintf("%d ms", rtt.Milliseconds())
					pingColor = color.NRGBA{R: 40, G: 160, B: 60, A: 255} // Green
				}

				lbl := material.Body2(th, pingStr)
				lbl.Color = pingColor
				lbl.TextSize = unit.Sp(12)
				return layout.Inset{Right: unit.Dp(16)}.Layout(gtx, lbl.Layout)
			}),

			// D. power button
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				btn := material.ButtonLayout(th, &host.BtnPower)

				if isOnline {
					btn.Background = color.NRGBA{R: 231, G: 76, B: 60, A: 255}
				} else {
					btn.Background = color.NRGBA{R: 46, G: 204, B: 113, A: 255}
				}

				circuitText := "Shutdown"
				if !isOnline {
					circuitText = "Turn on"
				}

				return btn.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.UniformInset(unit.Dp(0)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{
							Alignment: layout.Middle,
						}.Layout(gtx,
							layout.Rigid(layout.Spacer{
								Width: unit.Dp(6),
							}.Layout),
							// Icon
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return layout.Inset{
									Right: unit.Dp(4),
								}.Layout(gtx, app.shutdownIcon.Layout)
							}),

							// Text
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return layout.Inset{
									Top: unit.Dp(6),
								}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
									lbl := material.Body2(th, circuitText)
									lbl.Color = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
									lbl.Font.Weight = font.Medium

									return lbl.Layout(gtx)
								})
							}),
							layout.Rigid(layout.Spacer{
								Width: unit.Dp(6),
							}.Layout),
						)
					})
				})
			}),
		)
	})
}

func (app *SinglePageApp) layoutWSLNodes(
	gtx layout.Context,
	th *material.Theme,
) layout.Dimensions {
	allWslNodes := make([]*state.WSLState, 0)
	for _, host := range app.hosts {
		allWslNodes = append(allWslNodes, host.Wsls...)
	}

	fmt.Println("kakakakaka")
	fmt.Println(allWslNodes)
	return app.wslList.Layout(
		gtx,
		len(allWslNodes),
		func(gtx layout.Context, i int) layout.Dimensions {
			return app.layoutWSLNode(gtx, th, allWslNodes[i])
		},
	)
}

func (app *SinglePageApp) layoutWSLNode(
	gtx layout.Context,
	th *material.Theme,
	wslNodeState *state.WSLState,
) layout.Dimensions {
	return layout.Inset{
		Right: unit.Dp(8),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return widget.Border{
			Color: color.NRGBA{
				R: 220,
				G: 220,
				B: 220,
				A: 255,
			},
			Width:        unit.Dp(1),
			CornerRadius: unit.Dp(6),
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {

			return layout.Inset{
				Top:    unit.Dp(8),
				Bottom: unit.Dp(8),
				Left:   unit.Dp(12),
				Right:  unit.Dp(12),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {

				return layout.Flex{
					Axis: layout.Vertical,
				}.Layout(
					gtx,

					// Name
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{
							Axis: layout.Horizontal,
						}.Layout(
							gtx,

							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return layout.Inset{
									Right: unit.Dp(6),
								}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
									return layout.Dimensions{
										Size: image.Pt(8, 8),
									}
								})
							}),

							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								lbl := material.Body2(th, wslNodeState.Name)
								lbl.TextSize = unit.Sp(14)
								lbl.Font.Weight = font.Medium

								return lbl.Layout(gtx)
							}),
						)
					}),

					// Status
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						lbl := material.Caption(th, "Running • WSL 2")
						lbl.TextSize = unit.Sp(12)

						return layout.Inset{
							Top: unit.Dp(4),
						}.Layout(gtx, lbl.Layout)
					}),

					// D. power button
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						btn := material.ButtonLayout(th, &wslNodeState.BtnPower)

						if true {
							btn.Background = color.NRGBA{R: 231, G: 76, B: 60, A: 255}
						} else {
							btn.Background = color.NRGBA{R: 46, G: 204, B: 113, A: 255}
						}

						return btn.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return layout.UniformInset(unit.Dp(0)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								return layout.Flex{
									Alignment: layout.Middle,
								}.Layout(gtx,
									layout.Rigid(layout.Spacer{
										Width: unit.Dp(6),
									}.Layout),
									// Icon
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										return layout.Inset{
											Right: unit.Dp(4),
										}.Layout(gtx, app.shutdownIcon.Layout)
									}),
								)
							})
						})
					}),
				)
			})
		})
	})
}
