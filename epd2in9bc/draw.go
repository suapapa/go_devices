// Copyright 2019, Homin Lee <homin.lee@suapapa.net>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package epd2in9bc

import (
	"fmt"
	"image"
	"image/color"
)

// DrawImage draws a image
func (d *Display) DrawImage(i image.Image) error {
	b, ry, err := d.Image2Buffer(i)
	if err != nil {
		return err
	}
	return d.DrawBuffer(b, ry)
}

// Clear fill display with given patten in byte (8 pixel)
func (d *Display) Clear(cB, cRY byte) {
	lw := (d.w + 7) / 8
	b := make([]byte, lw*d.h)
	for i := range b {
		b[i] = cB
	}
	ry := make([]byte, lw*d.h)
	for i := range ry {
		ry[i] = cRY
	}
	d.DrawBuffer(b, ry)
}

// DrawBuffer draws buffer to display
func (d *Display) DrawBuffer(b, ry []byte) error {
	// check if buffer is proper w*h
	lw := (d.w + 7) / 8
	if len(b) != lw*d.h {
		return fmt.Errorf("unexpect buffer size, %d", len(b))
	}
	if len(ry) != lw*d.h {
		return fmt.Errorf("unexpect buffer size, %d", len(b))
	}

	// log.Println("DrawBuffer", b, ry)
	// d.init()

	d.sendCmd(0x10)
	// for _, v := range b {
	// 	d.sendData(v)
	// }
	for i := 0; i < len(b); i += 4096 {
		if vlen := len(b[i:]); vlen > 4096 {
			d.sendDatas(b[i : i+4096])
		} else {
			d.sendDatas(b[i : i+vlen])
		}
	}
	// d.sendDatas(b)
	d.sendCmd(0x13)
	// for _, v := range ry {
	// 	d.sendData(v)
	// }
	for i := 0; i < len(ry); i += 4096 {
		if vlen := len(ry[i:]); vlen > 4096 {
			d.sendDatas(ry[i : i+4096])
		} else {
			d.sendDatas(ry[i : i+vlen])
		}
	}
	// d.sendDatas(ry)
	d.sendCmd(0x12)
	d.waitTillNotBusy()

	return nil
}

// Image2Buffer returns monochrome image buffer from image.Image
func (d *Display) Image2Buffer(img image.Image) ([]byte, []byte, error) {
	b := make([]byte, (d.w/8)*d.h)
	for i := range b {
		b[i] = 0xFF // fill white
	}
	ry := make([]byte, (d.w/8)*d.h)
	for i := range ry {
		ry[i] = 0xFF // fill white
	}

	imgW, imgH := img.Bounds().Dx(), img.Bounds().Dy()
	if imgW == d.w && imgH == d.h { // vertical
		for y := 0; y < imgH; y++ {
			for x := 0; x < imgW; x++ {
				switch checkColor(img.At(x, y)) {
				case black:
					b[(x+d.w)/8] &= ^(0x80 >> (x % 8))
				case accent:
					ry[(x+d.w)/8] &= ^(0x80 >> (x % 8))
				}
				// if isBlackColor(img.At(x, y)) {
				// 	b[(x+d.w)/8] &= ^(0x80 >> (x % 8))
				// }
			}
		}
		return b, ry, nil
	} else if imgW == d.h && imgH == d.w { // horizontal
		for y := 0; y < imgH; y++ {
			for x := 0; x < imgW; x++ {
				newX := y
				newY := d.h - x - 1
				switch checkColor(img.At(x, y)) {
				case black:
					b[(newX+newY*d.w)/8] &= ^(0x80 >> (y % 8))
				case accent:
					ry[(newX+newY*d.w)/8] &= ^(0x80 >> (y % 8))
				}
				// if isBlackColor(img.At(x, y)) {
				// 	b[(newX+newY*d.w)/8] &= ^(0x80 >> (y % 8))
				// }
			}
		}
		return b, ry, nil
	}

	return nil, nil, fmt.Errorf("image size should be %dx%d of %dx%d", d.w, d.h, d.h, d.w)
}

func checkColor(c color.Color) inkColor {
	r, g, b, _ := c.RGBA()
	if r == 0 && g == 0 && b == 0 {
		return black
	}
	if r == 0xffffffff && g == 0xffffffff && b == 0xffffffff {
		return white
	}
	return accent
}
