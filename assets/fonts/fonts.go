package fonts

import (
	_ "embed"

	"gioui.org/font"
	"gioui.org/font/opentype"
	"gioui.org/text"
)

//go:embed GoogleSans-Regular.ttf
var googleSansRegular []byte

//go:embed GoogleSans-Medium.ttf
var googleSansMedium []byte

//go:embed GoogleSans-Bold.ttf
var googleSansBold []byte

func NewShaper() *text.Shaper {
	regularFace, err := opentype.Parse(googleSansRegular)
	if err != nil {
		panic(err)
	}

	mediumFace, err := opentype.Parse(googleSansMedium)
	if err != nil {
		panic(err)
	}

	boldFace, err := opentype.Parse(googleSansBold)
	if err != nil {
		panic(err)
	}

	// Đưa Face vào font.FontFace và thiết lập Weight tương ứng
	return text.NewShaper(
		text.NoSystemFonts(), // Tuỳ chọn thêm: Tắt font hệ thống nếu bạn chỉ muốn dùng Google Sans
		text.WithCollection([]font.FontFace{
			{
				Font: regularFace.Font(),
				Face: regularFace,
			},
			{
				Font: mediumFace.Font(),
				Face: mediumFace,
			},
			{
				Font: boldFace.Font(),
				Face: boldFace,
			},
		}),
	)
}
