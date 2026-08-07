package components

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

// Vẽ chấm tròn trạng thái (Green = Online, Red = Offline, Gray = Tắt)
func DrawStatusBadge(gtx layout.Context, isOnline bool) layout.Dimensions {
	size := gtx.Dp(14)

	badgeColor := color.NRGBA{R: 180, G: 180, B: 180, A: 255} // Xám (Tắt)
	if isOnline {
		badgeColor = color.NRGBA{R: 46, G: 204, B: 113, A: 255} // Xanh (Online)
	} else {
		badgeColor = color.NRGBA{R: 231, G: 76, B: 60, A: 255} // Đỏ (Offline)
	}

	defer clip.Ellipse{
		Min: image.Pt(0, 0),
		Max: image.Pt(size, size),
	}.Op(gtx.Ops).Push(gtx.Ops).Pop()

	paint.ColorOp{Color: badgeColor}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	return layout.Dimensions{Size: image.Pt(size, size)}
}
