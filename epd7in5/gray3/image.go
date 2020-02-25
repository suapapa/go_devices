package gray3

import (
	"image"
	"image/color"
	"image/draw"
)

// Image is a Gray3 image
type Image struct {
	// Pix holds images pixels
	Pix []byte
	// Rect is image.Rectangle
	Rect image.Rectangle
}

// NewImage returns gray3.Image instance.
func NewImage(r image.Rectangle) *Image {
	w := r.Dx()
	h := r.Dy()
	pix := make([]byte, w*h/4)
	for i := 0; i < len(pix); i++ {
		pix[i] = 0b1111_1111
	}
	return &Image{
		Pix:  pix,
		Rect: r,
	}
}

// ColorModel implements draw.Image
func (i *Image) ColorModel() color.Model {
	return Gray3Model
}

// Bounds implements draw.Image
func (i *Image) Bounds() image.Rectangle {
	return i.Rect
}

// At implements draw.Image
func (i *Image) At(x, y int) color.Color {
	pos := (x + y*i.Rect.Dx()) / 4
	shift := 6 - (x % 4 * 2)
	pix := (i.Pix[pos] >> shift) & 0b11
	switch pix {
	case 0b00:
		return Black
	case 0b10:
		return Gray
	case 0b11:
		return White
	}

	return Black
}

// Set implements draw.Image
func (i *Image) Set(x, y int, c color.Color) {
	if _, ok := c.(Gray3); !ok {
		c = gray3Model(c)
	}
	pos := (x + y*i.Rect.Dx()) / 4
	shift := x % 4 * 2
	switch c {
	case Black:
		i.Pix[pos] &= ^(0b1100_0000 >> shift)
	case Gray:
		i.Pix[pos] &= ^(0b1100_0000 >> shift)
		i.Pix[pos] |= (0b1000_0000 >> shift)
	case White:
		i.Pix[pos] |= (0b1100_0000 >> shift)
	}
}

var _ draw.Image = &Image{}
