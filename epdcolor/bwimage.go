package epdcolor

import (
	"image"
	"image/color"
	"image/draw"
)

// Image is a Gray3 image
type WBImage struct {
	// Pix holds images pixels
	Pix []byte
	// Rect is image.Rectangle
	Rect image.Rectangle
}

// NewImage returns gray3.Image instance.
func NewWBImage(r image.Rectangle) *WBImage {
	w := r.Dx()
	h := r.Dy()
	pix := make([]byte, w*h/8)
	// for i := 0; i < len(pix); i++ {
	// 	pix[i] = 0b1111_1111
	// }
	return &WBImage{
		Pix:  pix,
		Rect: r,
	}
}

// ColorModel implements draw.Image
func (i *WBImage) ColorModel() color.Model {
	return WBModel
}

// Bounds implements draw.Image
func (i *WBImage) Bounds() image.Rectangle {
	return i.Rect
}

// At implements draw.Image
func (i *WBImage) At(x, y int) color.Color {
	pos := (x + y*i.Rect.Dx()) / 8
	shift := 7 - (x % 8)
	pix := (i.Pix[pos] >> shift) & 1
	switch pix {
	case 0:
		return WBWhite
	case 1:
		return WBBlack
	}

	return WBBlack
}

// Set implements draw.Image
func (i *WBImage) Set(x, y int, c color.Color) {
	if _, ok := c.(WB); !ok {
		c = bwModel(c)
	}
	pos := (x + y*i.Rect.Dx()) / 8
	shift := x % 8
	switch c {
	case WBWhite:
		i.Pix[pos] &= ^(0b1000_0000 >> shift)
	case WBBlack:
		i.Pix[pos] |= (0b1000_0000 >> shift)
	}
}

var _ draw.Image = &WBImage{}
