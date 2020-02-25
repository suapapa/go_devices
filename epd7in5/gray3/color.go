package gray3

import "image/color"

// Gray3 implements a 3 kinds of gray color.
type Gray3 uint8

// Possible colors
const (
	White Gray3 = iota
	Gray
	Black
)

// RGBA returns either all white, gray or black.
func (g Gray3) RGBA() (uint32, uint32, uint32, uint32) {
	switch g {
	case White:
		return 0xffff, 0xffff, 0xffff, 0xffff
	case Gray:
		return 0xff00, 0xff00, 0xff00, 0xffff
	case Black:
		return 0, 0, 0, 0xffff
	}
	return 0, 0, 0, 0xffff
}

// Gray3Model is color model for white, gray or black color.
var Gray3Model = color.ModelFunc(gray3Model)

func gray3Model(c color.Color) color.Color {
	if _, ok := c.(Gray3); ok {
		return c
	}
	r, g, b, _ := c.RGBA()
	y := (19595*r + 38470*g + 7471*b + 1<<15) >> 24

	if y < 64 {
		return Black
	} else if y < 192 {
		return Gray
	}
	return White
}
