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
	if l.buff[x+(y/8)*l.w]&byte(1<<(uint(y)&7)) == 0x00 {
		return color.Gray{Y: 0x00}
	}
	return color.Gray{Y: 0xFF}
}

func (l *LCD) DrawImage(i image.Image) {
	imgW, imgH := i.Bounds().Dx(), i.Bounds().Dy()

	// TODO: fix to support images of arbitary size
	if imgW != l.w || imgH != l.h {
		panic("image should be size of 128x64")
	}

	for y := 0; y < imgH; y++ {
		for x := 0; x < imgW; x++ {
			r, g, b, _ := i.At(x, y).RGBA()
			if r != 0 || g != 0 || b != 0 {
				l.DrawPixel(x, y, true)
			} else {
				l.DrawPixel(x, y, false)
			}
		}
	}
}
