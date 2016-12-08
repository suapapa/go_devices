// Copyright 2016, Homin Lee <homin.lee@suapapa.net>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sh1106

import (
	"image"
	"image/color"
)

// Implemet image.Image interface:
// func (l *LCD) ColorModel() color.Model {
// 	return color.GrayModel
// }
//
// func (l *LCD) Bounds() image.Rectangle {
// 	return image.Rect(0, 0, int(l.w), int(l.h))
// }
//
// func (l *LCD) At(x, y int) color.Color {
// 	if l.buff[x+(y/8)*l.w]&byte(1<<(uint(y)&7)) == 0x00 {
// 		return color.Gray{Y: 0x00}
// 	}
// 	return color.Gray{Y: 0xFF}
// }

// DrawImage draws a image to the internal buffer
func (l *LCD) DrawImage(i image.Image) {
	l.DrawImageWithColorCvt(
		i,
		func(c color.Color) bool {
			r, g, b, _ := c.RGBA()
			if r != 0 || g != 0 || b != 0 {
				return true
			}
			return false
		},
	)
}

// DrawImageWithColorCvt draws a image to internal bufer
// using cc as color.Color to bool converter
func (l *LCD) DrawImageWithColorCvt(i image.Image, cc func(color.Color) bool) {
	imgW, imgH := i.Bounds().Dx(), i.Bounds().Dy()

	// TODO: fix to support images of arbitary size
	if imgW != l.w || imgH != l.h {
		panic("image should be size of 128x64")
	}

	for y := 0; y < imgH; y++ {
		for x := 0; x < imgW; x++ {
			l.DrawPixel(x, y, cc(i.At(x, y)))
		}
	}
}
