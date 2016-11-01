package ili9340

import (
	"image"
	"image/color"
	"log"
)

func rgb565(c color.Color) uint16 {
	r, g, b, _ := c.RGBA()
	r >>= 8
	g >>= 8
	b >>= 8
	rgb := (uint16(r&0xF8) << 8) | uint16((g&0xFC)<<3) | uint16(b>>3)
	log.Printf("rgb: 0x%06x\n", rgb)
	return rgb
}

func color2rgb565(c color.Color) (h, l uint8) {
	r, g, b, _ := c.RGBA()
	r >>= 8
	g >>= 8
	b >>= 8
	h = uint8(r&0xF8) | uint8((g&0xFC)>>5)
	l = uint8((g&0xFC)<<3) | uint8(b>>3)
	return h, l
}

func img2rbg565(img image.Image) []uint8 {
	b := img.Bounds()
	w, h := b.Dx(), b.Dy()
	buff := make([]uint8, w*h*2)
	i := 0
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			h, l := color2rgb565(img.At(x, y))
			buff[i] = h
			buff[i+1] = l
			i += 2
		}
	}
	return buff
}
