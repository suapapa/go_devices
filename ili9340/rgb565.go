// Copyright 2016, Homin Lee <homin.lee@suapapa.net>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ili9340

import (
	"image"
	"image/color"
)

// rgb565 converts color.Color to 2 bytes length []byte
func rgb565(c color.Color) []byte {
	r, g, b, _ := c.RGBA()
	r >>= 8
	g >>= 8
	b >>= 8
	ret := make([]byte, 2)
	ret[0] = byte(r&0xF8) | byte((g&0xFC)>>5)
	ret[1] = byte((g&0xFC)<<3) | byte(b>>3)
	return ret
}

// rgb565Img converts image.Image to 2*w*h bytes lenght of []byte
func rgb565Img(img image.Image) []byte {
	b := img.Bounds()
	w, h := b.Dx(), b.Dy()
	buff := make([]byte, 0, w*h*2)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			buff = append(buff, rgb565(img.At(x, y))...)
		}
	}
	return buff
}
