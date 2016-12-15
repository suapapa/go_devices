// Copyright 2016, Homin Lee <homin.lee@suapapa.net>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ili9340

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"time"

	"github.com/goiot/exp/gpio"
	gpio_driver "github.com/goiot/exp/gpio/driver"
	"golang.org/x/exp/io/spi"
	spi_driver "golang.org/x/exp/io/spi/driver"
)

// LCD reperesents TFT-LCD panel which using ili9340 controller
type LCD struct {
	spiDev  *spi.Device
	gpioDev *gpio.Device

	W, H int
}

// Open connets LCD to passed SPI bus and GPIO controls
// GPIO device should have "DC" and "RST"
func Open(bus spi_driver.Opener, ctr gpio_driver.Opener) (*LCD, error) {
	spiDev, err := spi.Open(bus)
	if err != nil {
		return nil, err
	}
	spiDev.SetCSChange(false)

	gpioDev, err := gpio.Open(ctr)
	if err != nil {
		spiDev.Close()
		return nil, err
	}

	if err = gpioDev.SetDirection(PinDC, gpio.Out); err != nil {
		spiDev.Close()
		return nil, err
	}

	if err = gpioDev.SetDirection(PinRST, gpio.Out); err != nil {
		// RST pin could be skipped when no reset is requiered
		// spiDev.Close()
		// return nil, err
	}

	lcd := &LCD{
		spiDev:  spiDev,
		gpioDev: gpioDev,
	}

	lcd.Reset()
	lcd.init()
	lcd.Rotate(0)

	return lcd, nil
}

// Close takes care of cleaning things up.
func (l *LCD) Close() {
	l.spiDev.Close()
	l.gpioDev.Close()
}

// Reset resets LCD by using PinRST
func (l *LCD) Reset() {
	if err := l.gpioDev.SetDirection(PinRST, gpio.Out); err != nil {
		return
	}

	l.gpioDev.SetValue(PinRST, 1)
	time.Sleep(5 * time.Millisecond)
	l.gpioDev.SetValue(PinRST, 0)
	time.Sleep(20 * time.Millisecond)
	l.gpioDev.SetValue(PinRST, 1)
	time.Sleep(120 * time.Millisecond)
}

// DrawPixel sets a pixel at a position x, y.
func (l *LCD) DrawPixel(x, y int, c color.Color) {
	l.setWinAddr(x, y, x, y)
	l.writeData(rgb565(c))
}

// DrawRect draws a rect at x, y in given w, h
func (l *LCD) DrawRect(x, y, w, h int, c color.Color) {
	l.setWinAddr(x, y, x+w-1, y+h-1)
	l.writeData(bytes.Repeat(rgb565(c), w*h))
}

// DrawImage draws an image on the display starting from x, y.
func (l *LCD) DrawImage(x, y int, img image.Image) {
	w, h := img.Bounds().Dx(), img.Bounds().Dy()

	l.setWinAddr(x, y, x+w-1, y+h-1)
	// TODO: crop image if it oversizes LCD
	l.writeData(rgb565Img(img))
}

// Rotate rotates display
func (l *LCD) Rotate(r int) error {
	var val byte
	var w, h int
	switch r {
	case 270:
		val = regMADCTLvalMV
		w, h = height, width
	case 180:
		val = regMADCTLvalMY
		w, h = width, height
	case 90:
		val = regMADCTLvalMV | regMADCTLvalMY | regMADCTLvalMX
		w, h = height, width
	case 0:
		val = regMADCTLvalMX
		w, h = width, height
	default:
		return fmt.Errorf("ili9340: unsupported rotation degree, %d", r)
	}

	l.W, l.H = w, h
	l.writeReg(regMADCTL, val|regMADCTLvalBGR)

	return nil
}

// Invert invert colors on the display
func (l *LCD) Invert(on bool) {
	if on {
		l.writeReg(regINVON)
	} else {
		l.writeReg(regINVOFF)
	}
}

func (l *LCD) init() {
	l.writeReg(0xEF, 0x03, 0x80, 0x02)
	l.writeReg(0xCF, 0x00, 0xC1, 0x30)
	l.writeReg(0xED, 0x64, 0x03, 0x12, 0x81)
	l.writeReg(0xE8, 0x85, 0x00, 0x78)
	l.writeReg(0xCB, 0x39, 0x2C, 0x00, 0x34, 0x02)
	l.writeReg(0xF7, 0x20)
	l.writeReg(0xEA, 0x00, 0x00)

	/* Power Control 1 */
	l.writeReg(0xC0, 0x23)

	/* Power Control 2 */
	l.writeReg(0xC1, 0x10)

	/* VCOM Control 1 */
	l.writeReg(0xC5, 0x3e, 0x28)

	/* VCOM Control 2 */
	l.writeReg(0xC7, 0x86)

	/* COLMOD: Pixel Format Set */
	/* 16 bits/pixel */
	l.writeReg(0x3A, 0x55)

	/* Frame Rate Control */
	/* Division ratio = fosc, Frame Rate = 79Hz */
	l.writeReg(0xB1, 0x00, 0x18)

	/* Display Function Control */
	l.writeReg(0xB6, 0x08, 0x82, 0x27)

	/* Gamma Function Disable */
	l.writeReg(0xF2, 0x00)

	/* Gamma curve selected  */
	l.writeReg(0x26, 0x01)

	/* Positive Gamma Correction */
	l.writeReg(0xE0,
		0x0F, 0x31, 0x2B, 0x0C, 0x0E, 0x08, 0x4E, 0xF1,
		0x37, 0x07, 0x10, 0x03, 0x0E, 0x09, 0x00)

	/* Negative Gamma Correction */
	l.writeReg(0xE1,
		0x00, 0x0E, 0x14, 0x03, 0x11, 0x07, 0x31, 0xC1,
		0x48, 0x08, 0x0F, 0x0C, 0x31, 0x36, 0x0F)

	/* Sleep OUT */
	l.writeReg(0x11)

	time.Sleep(120 * time.Millisecond)

	/* Display ON */
	l.writeReg(0x29)
}

// setWinAddr set windows on LCD and return data length in the window
func (l *LCD) setWinAddr(xs, ys, xe, ye int) int {
	/* Column address */
	l.writeReg(0x2A,
		uint8(xs>>8), uint8(xs&0xFF),
		uint8(xe>>8), uint8(xe&0xFF),
	)

	/* Row address */
	l.writeReg(0x2B,
		uint8(ys>>8), uint8(ys&0xFF),
		uint8(ye>>8), uint8(ye&0xFF),
	)

	/* Memory write */
	l.writeReg(0x2C)

	return ((xe-xs)*(ye-ys) + 1) * 2
}

func (l *LCD) writeData(data []byte) {
	l.gpioDev.SetValue(PinDC, 1)
	l.spiDev.Tx(data, nil)
}

func (l *LCD) writeReg(vs ...uint8) {
	l.gpioDev.SetValue(PinDC, 0)
	l.spiDev.Tx(vs[:1], nil)

	if len(vs) > 1 {
		l.writeData(vs[1:])
	}
}
