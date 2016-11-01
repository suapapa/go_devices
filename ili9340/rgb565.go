package ili9340

import "image/color"

func rgb565(c color.Color) uint16 {
	r, g, b, _ := c.RGBA()
	// return ((r>>3)&0x1F)<<11 | ((g>>2)&0x3f)<<5 | (b>>3)&0x1F
	return ((r & 0xF8) << 8) | ((g & 0xFC) << 3) | (b >> 3)
}
