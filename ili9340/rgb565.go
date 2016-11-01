package ili9340

import "image/color"

func rgb565(c color.Color) uint16 {
	r, g, b, _ := c.RGBA()
	return (uint16(r&0xF8) << 8) | uint16((g&0xFC)<<3) | uint16(b>>3)
}
