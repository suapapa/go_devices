package sh1106

import (
	"image"
	"image/color"
)

func (l *LCD) ColorModel() color.Model {
	return color.GrayModel
}

func (l *LCD) Bounds() image.Rectangle {
	return image.Rect(0, 0, int(l.w), int(l.h))
}

func (l *LCD) At(x, y int) color.Color {
	ux, uy := uint(x), uint(y)
	if l.buff[ux+(uy/8)*l.w]&byte(1<<(uy&7)) == 0x00 {
		return color.Gray{Y: 0x00}
	}
	return color.Gray{Y: 0xFF}
}
