package epdcolor

import (
	"image"
	"image/color"
	"image/draw"
)

// Image is a Gray3 image
type BWImage struct {
	// Pix holds images pixels
	Pix []byte
	// Rect is image.Rectangle
	Rect image.Rectangle
}

// NewImage returns gray3.Image instance.
func NewBWImage(r image.Rectangle) *BWImage {
	w := r.Dx()
	h := r.Dy()
	pix := make([]byte, w*h/8)
	for i := 0; i < len(pix); i++ {
		pix[i] = 0b1111_1111
	}
	return &BWImage{
		Pix:  pix,
		Rect: r,
	}
}

// ColorModel implements draw.Image
func (i *BWImage) ColorModel() color.Model {
	return BWModel
}

// Bounds implements draw.Image
func (i *BWImage) Bounds() image.Rectangle {
	return i.Rect
}

// At implements draw.Image
func (i *BWImage) At(x, y int) color.Color {
	pos := (x + y*i.Rect.Dx()) / 8
	shift := 7 - (x % 8)
	pix := (i.Pix[pos] >> shift) & 1
	switch pix {
	case 0:
		return BWBlack
	case 1:
		return BWWhite
	}

	return BWBlack
}

// Set implements draw.Image
func (i *BWImage) Set(x, y int, c color.Color) {
	if _, ok := c.(BW); !ok {
		c = bwModel(c)
	}
	pos := (x + y*i.Rect.Dx()) / 8
	shift := x % 8
	switch c {
	case BWBlack:
		i.Pix[pos] &= ^(0b1000_0000 >> shift)
	case BWWhite:
		i.Pix[pos] |= (0b1000_0000 >> shift)
	}
}

var _ draw.Image = &BWImage{}
