package components

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

func DrawNavItem(gtx layout.Context, th *material.Theme, btn *widget.Clickable, title string, isSelected bool) layout.Dimensions {
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
