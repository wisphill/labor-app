package components

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

func DrawMonitorIcon(gtx layout.Context) layout.Dimensions {
	const (
		w = 18
		h = 14
	)

	dark := color.NRGBA{
		R: 45,
		G: 45,
		B: 45,
		A: 255,
	}

	// Monitor body
	monitor := clip.RRect{
		Rect: image.Rectangle{
			Min: image.Pt(0, 0),
			Max: image.Pt(w, h),
		},
		NE: 2,
		NW: 2,
		SE: 2,
		SW: 2,
	}

	paint.FillShape(
		gtx.Ops,
		dark,
		monitor.Op(gtx.Ops),
	)

	// Screen
	screen := clip.Rect{
		Min: image.Pt(2, 2),
		Max: image.Pt(w-2, h-3),
	}

	paint.FillShape(
		gtx.Ops,
		color.NRGBA{
			R: 240,
			G: 240,
			B: 240,
			A: 255,
		},
		screen.Op(),
	)

	// Stand
	stand := clip.Rect{
		Min: image.Pt(7, h),
		Max: image.Pt(11, h+2),
	}

	paint.FillShape(
		gtx.Ops,
		dark,
		stand.Op(),
	)

	// Base
	base := clip.RRect{
		Rect: image.Rectangle{
			Min: image.Pt(5, h+2),
			Max: image.Pt(13, h+3),
		},
		NE: 1,
		NW: 1,
		SE: 1,
		SW: 1,
	}

	paint.FillShape(
		gtx.Ops,
		dark,
		base.Op(gtx.Ops),
	)

	return layout.Dimensions{
		Size: image.Pt(w, h+3),
	}
}
