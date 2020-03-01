package epdcolor

import (
	"image/color"
)

// WB implements a 3 kinds of gray color.
type WB uint8

// Possible colors
const (
	WBWhite WB = iota
	WBBlack
)

func (c WB) String() string {
	switch c {
	case WBWhite:
		return "WBWhite"
	case WBBlack:
		return "WBBlack"
	}
	return "Unknown"
}

// RGBA returns either all white, gray or black.
func (c WB) RGBA() (uint32, uint32, uint32, uint32) {
	switch c {
	case WBWhite:
		return 0xffff, 0xffff, 0xffff, 0xffff
	case WBBlack:
		return 0, 0, 0, 0xffff
	}
	return 0, 0, 0, 0xffff
}

// WBModel is color model for white or black color.
var WBModel = color.ModelFunc(bwModel)

func bwModel(c color.Color) color.Color {
	if _, ok := c.(WB); ok {
		return c
	}
	r, g, b, _ := c.RGBA()
	y := (19595*r + 38470*g + 7471*b + 1<<15) >> 24

	// if y < 64 {
	// 	return WBBlack
	// } else if y < 192 {
	// 	return Gray
	// }
	if y < 128 {
		return WBBlack
	}
	return WBWhite
}
