// Copyright 2019, Homin Lee <homin.lee@suapapa.net>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package epd2in13bc

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

// DrawBuffer draws buffer to display. b for black channel, ry for color channel.
func (d *Display) DrawBuffer(b []byte, ry []byte) error {
	// check if buffer is proper w*h
	lw := (d.w + 7) / 8
	if len(b) != lw*d.h {
		return fmt.Errorf("unexpect buffer size for bw channel, %d", len(b))
	}
	if len(ry) != lw*d.h {
		return fmt.Errorf("unexpect buffer size for color channel, %d", len(b))
	}

	d.sendCmd(0x10)
	d.sendDatas(b)
	d.sendCmd(0x92)

	d.sendCmd(0x13)
	d.sendDatas(ry)
	d.sendCmd(0x92)

	d.sendCmd(0x12)
	d.waitTillNotBusy()

	return nil
}

// # logging.debug("bufsiz = ",int(self.width/8) * self.height)
// buf = [0xFF] * (int(self.width/8) * self.height)
// image_monocolor = image.convert('1')
// imwidth, imheight = image_monocolor.size
// pixels = image_monocolor.load()
// # logging.debug("imwidth = %d, imheight = %d",imwidth,imheight)
// if(imwidth == self.width and imheight == self.height):
//     logging.debug("Vertical")
//     for y in range(imheight):
//         for x in range(imwidth):
//             # Set the bits for the column of pixels at the current position.
//             if pixels[x, y] == 0:
//                 buf[int((x + y * self.width) / 8)] &= ~(0x80 >> (x % 8))
// elif(imwidth == self.height and imheight == self.width):
//     logging.debug("Horizontal")
//     for y in range(imheight):
//         for x in range(imwidth):
//             newx = y
//             newy = self.height - x - 1
//             if pixels[x, y] == 0:
//                 buf[int((newx + newy*self.width) / 8)] &= ~(0x80 >> (y % 8))
// return buf

// Image2Buffer returns monochrome image buffer from image.Image
func (d *Display) Image2Buffer(img image.Image) ([]byte, []byte, error) {
	lw := (d.w + 7) / 8
	b := make([]byte, lw*d.h) // black buffer
	for i := range b {
		b[i] = 0xFF // fill white
	}
	ry := make([]byte, lw*d.h) // red or yellow buffer
	for i := range ry {
		ry[i] = 0xFF // fill white
	}

	imgW, imgH := img.Bounds().Dx(), img.Bounds().Dy()
	if imgW == d.w && imgH == d.h { // vertical
		for y := 0; y < imgH; y++ {
			for x := 0; x < imgW; x++ {
				xx := imgW - x
				switch checkColor(img.At(x, y)) {
				case black:
					b[xx/8+y*lw] &= ^(0x80 >> (xx % 8))
				case accent:
					ry[xx/8+y*lw] &= ^(0x80 >> (xx % 8))
				}
			}
		}
		return b, ry, nil
	} else if imgW == d.h && imgH == d.w { // Horizontal
		for y := 0; y < imgH; y++ {
			for x := 0; x < imgW; x++ {
				newX := y
				newY := d.h - x - 1
				switch checkColor(img.At(x, y)) {
				case black:
					b[newX/8+newY*lw] &= ^(0x80 >> (y % 8))
				case accent:
					ry[newX/8+newY*lw] &= ^(0x80 >> (y % 8))
				}
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
