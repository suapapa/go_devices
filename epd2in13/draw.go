// Copyright 2019, Homin Lee <homin.lee@suapapa.net>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package epd2in13

import (
	"fmt"
	"image"
	"image/color"
)

// DrawImage draws a image
func (d *Display) DrawImage(i image.Image) error {
	b, err := d.Image2Buffer(i)
	if err != nil {
		return err
	}
	return d.DrawBuffer(b)
}

// DrawBuffer draws buffer to display
func (d *Display) DrawBuffer(b []byte) error {
	// check if buffer is proper w*h
	lw := (d.w + 7) / 8
	if len(b) != lw*d.h {
		return fmt.Errorf("unexpect buffer size, %d", len(b))
	}

	d.sendCmd(0x24)
	d.sendDatas(b)
	d.TurnOnFull()

	return nil
}

// Image2Buffer returns monochrome image buffer from image.Image
func (d *Display) Image2Buffer(img image.Image) ([]byte, error) {
	lw := (d.w + 7) / 8
	b := make([]byte, lw*d.h)
	// fill white
	for i := range b {
		b[i] = 0xFF
	}

	imgW, imgH := img.Bounds().Dx(), img.Bounds().Dy()
	if imgW == d.w && imgH == d.h { // vertical
		for y := 0; y < imgH; y++ {
			for x := 0; x < imgW; x++ {
				if isBlackColor(img.At(x, y)) {
					xx := imgW - x
					b[xx/8+y*lw] &= ^(0x80 >> (xx % 8))
				}
			}
		}
		return b, nil
	} else if imgW == d.h && imgH == d.w { // Horizontal
		for y := 0; y < imgH; y++ {
			for x := 0; x < imgW; x++ {
				newX := y
				newY := d.h - x - 1
				if isBlackColor(img.At(x, y)) {
					newY = imgW - newY - 1
					b[newX/8+newY*lw] &= ^(0x80 >> (y % 8))
				}
			}
		}
		return b, nil
	}

	return nil, fmt.Errorf("image size should be %dx%d of %dx%d", d.w, d.h, d.h, d.w)
}

func isBlackColor(c color.Color) bool {
	r, g, b, _ := c.RGBA()
	if r != 0 || g != 0 || b != 0 {
		return false
	}
	return true
}
