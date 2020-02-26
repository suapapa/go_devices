package bw

import (
	"image/color"
)

// BW implements a 3 kinds of gray color.
type BW uint8

// Possible colors
const (
	White BW = iota
	Black
)

func (c BW) String() string {
	switch c {
	case White:
		return "White"
	case Black:
		return "Black"
	}
	return "Unknown"
}

// RGBA returns either all white, gray or black.
func (c BW) RGBA() (uint32, uint32, uint32, uint32) {
	switch c {
	case White:
		return 0xffff, 0xffff, 0xffff, 0xffff
	case Black:
		return 0, 0, 0, 0xffff
	}
	return 0, 0, 0, 0xffff
}

// BWModel is color model for white or black color.
var BWModel = color.ModelFunc(bwModel)

func bwModel(c color.Color) color.Color {
	if _, ok := c.(BW); ok {
		return c
	}
	r, g, b, _ := c.RGBA()
	y := (19595*r + 38470*g + 7471*b + 1<<15) >> 24

	// if y < 64 {
	// 	return Black
	// } else if y < 192 {
	// 	return Gray
	// }
	if y < 128 {
		return Black
	}
	return White
}
