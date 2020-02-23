package gray3

import "image/color"

type Gray3 struct {
	Y uint8
}

func (c Gray3) RGBA() (r, g, b, a uint32) {
	y := uint32(c.Y)
	y |= y << 8
	return y, y, y, 0xffff
}

var (
	Gray3Model color.Model = color.ModelFunc(gray3Model)
)

func gray3Model(c color.Color) color.Color {
	if _, ok := c.(Gray3); ok {
		return c
	}
	r, g, b, _ := c.RGBA()
	y := (19595*r + 38470*g + 7471*b + 1<<15) >> 24

	return Gray3{uint8(y)}
}
