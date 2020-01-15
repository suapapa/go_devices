// Copyright 2019, Homin Lee <homin.lee@suapapa.net>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package epd2in9

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

// Clear fill display with given patten in byte (8 pixel)
func (d *Display) Clear(c byte) {
	lw := (d.w + 7) / 8
	b := make([]byte, lw*d.h)
	for i := range b {
		b[i] = c
	}
	d.DrawBuffer(b)
}

// DrawBuffer draws buffer to display
func (d *Display) DrawBuffer(b []byte) error {
	// check if buffer is proper w*h
	lw := (d.w + 7) / 8
	if len(b) != lw*d.h {
		return fmt.Errorf("unexpect buffer size, %d", len(b))
	}

	d.setWindow(0, 0, d.w-1, d.h-1)
	for j := 0; j < d.h; j++ {
		d.setCursor(0, j)
		d.sendCmd(0x24) // WRITE_RAM
		for i := 0; i < lw; i++ {
			d.sendData(b[i+j*lw])
		}
	}
	d.turnOn()

	return nil
}

func (d *Display) setWindow(xStart, yStart, xEnd, yEnd int) {
	d.sendCmd(0x44) // SET_RAM_X_ADDRESS_START_END_POSITION
	// x point must be the multiple of 8 or the last 3 bits will be ignored
	d.sendData(uint8(xStart >> 3))
	d.sendData(uint8(xEnd >> 3))
	d.sendCmd(0x45) // SET_RAM_Y_ADDRESS_START_END_POSITION
	d.sendData(uint8(yStart))
	d.sendData(uint8(yStart >> 8))
	d.sendData(uint8(yEnd))
	d.sendData(uint8(yEnd >> 8))
}

func (d *Display) setCursor(x, y int) {
	d.sendCmd(0x4E) // SET_RAM_X_ADDRESS_COUNTER
	// x point must be the multiple of 8 or the last 3 bits will be ignored
	d.sendData(uint8(x >> 3))
	d.sendCmd(0x4F) // SET_RAM_Y_ADDRESS_COUNTER
	d.sendData(uint8(y))
	d.sendData(uint8(y >> 8))
	d.waitTillNotBusy()
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
	} else if imgW == d.h && imgH == d.w { // horizontal
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
