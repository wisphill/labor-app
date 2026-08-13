package components

import (
	"bytes"
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/paint"
	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
)

// SVGRenderer lưu trữ hình ảnh đã được xử lý để vẽ tốc độ cao
type SVGRenderer struct {
	imageOp paint.ImageOp
	size    image.Point
}

// LoadSVG đọc file .svg, hỗ trợ đổi màu (tintColor) linh hoạt
// Nếu tintColor = nil, giữ nguyên màu gốc của file SVG.
// LoadSVG đọc file .svg, hỗ trợ đổi màu (tintColor) bao gồm cả kênh Alpha (độ trong suốt)
// Nếu tintColor = nil, giữ nguyên màu gốc của file SVG.
func LoadSVG(data []byte, width, height int, tintColor color.Color) (*SVGRenderer, error) {
	// 1. Parse file SVG
	icon, err := oksvg.ReadIconStream(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	// 2. Thiết lập kích thước muốn render
	icon.SetTarget(0, 0, float64(width), float64(height))

	// 3. Tạo một khung canvas trong suốt (RGBA)
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	scanner := rasterx.NewScannerGV(width, height, img, img.Bounds())
	dasher := rasterx.NewDasher(width, height, scanner)

	// 4. Vẽ SVG lên canvas
	icon.Draw(dasher, 1)

	// 5. Đổi màu và xử lý Alpha nếu có truyền vào
	if tintColor != nil {
		// Chuyển đổi màu sang dạng NRGBA không pre-multiplied để dễ tính toán
		nrgba := color.NRGBAModel.Convert(tintColor).(color.NRGBA)
		tr := uint32(nrgba.R)
		tg := uint32(nrgba.G)
		tb := uint32(nrgba.B)
		ta := uint32(nrgba.A) // Alpha của tintColor (0 - 255)

		// Duyệt qua từng pixel để phủ màu và kết hợp độ mờ
		for i := 0; i < len(img.Pix); i += 4 {
			origAlpha := uint32(img.Pix[i+3]) // Độ mờ viền chống răng cưa của SVG
			if origAlpha > 0 {
				// Tính alpha cuối cùng = (Độ mờ gốc của SVG * Độ mờ của tintColor) / 255
				finalAlpha := (origAlpha * ta) / 255

				// Ghi đè theo chuẩn Premultiplied Alpha của GioUI
				img.Pix[i] = uint8((tr * finalAlpha) / 255)
				img.Pix[i+1] = uint8((tg * finalAlpha) / 255)
				img.Pix[i+2] = uint8((tb * finalAlpha) / 255)
				img.Pix[i+3] = uint8(finalAlpha)
			}
		}
	}

	// 6. Cache thành paint.ImageOp để GioUI vẽ siêu mượt ở 60FPS
	op := paint.NewImageOp(img)

	return &SVGRenderer{
		imageOp: op,
		size:    image.Pt(width, height),
	}, nil
}

// Layout vẽ SVG ra màn hình UI
func (s *SVGRenderer) Layout(gtx layout.Context) layout.Dimensions {
	s.imageOp.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	return layout.Dimensions{
		Size: s.size,
	}
}
