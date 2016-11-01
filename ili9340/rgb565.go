package ili9340

import (
	"image/color"
	"log"
)

func rgb565(c color.Color) uint16 {
	r, g, b, _ := c.RGBA()
	rgb := (uint16(r&0xF8) << 8) | uint16((g&0xFC)<<3) | uint16(b>>3)
	log.Printf("rgb: 0x%06x\n", rgb)
	return rgb
}
