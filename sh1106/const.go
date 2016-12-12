// Copyright 2016, Homin Lee <homin.lee@suapapa.net>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sh1106

// sh1106 module should be connected in following pins
const (
	PinDC  = "DC"
	PinRST = "RST"
)

const (
	DefaultI2CAddr = 0x3C // 011110+SA0+RW - 0x3C or 0x3D

	// 128x64
	sh1106_LCDWIDTH  = 128
	sh1106_LCDHEIGHT = 64

	// 128x32
	// sh1106_LCDWIDTH  = 128
	// sh1106_LCDHEIGHT = 32

	// 96x16
	// sh1106_LCDWIDTH  = 96
	// sh1106_LCDHEIGHT = 16

	sh1106_SETCONTRAST         = 0x81
	sh1106_DISPLAYALLON_RESUME = 0xA4
	sh1106_DISPLAYALLON        = 0xA5
	sh1106_NORMALDISPLAY       = 0xA6
	sh1106_INVERTDISPLAY       = 0xA7
	sh1106_DISPLAYOFF          = 0xAE
	sh1106_DISPLAYON           = 0xAF

	sh1106_SETDISPLAYOFFSET = 0xD3
	sh1106_SETCOMPINS       = 0xDA

	sh1106_SETVCOMDETECT = 0xDB

	sh1106_SETDISPLAYCLOCKDIV = 0xD5
	sh1106_SETPRECHARGE       = 0xD9

	sh1106_SETMULTIPLEX = 0xA8

	sh1106_SETLOWCOLUMN  = 0x00
	sh1106_SETHIGHCOLUMN = 0x10

	sh1106_SETSTARTLINE = 0x40

	sh1106_MEMORYMODE = 0x20
	sh1106_COLUMNADDR = 0x21
	sh1106_PAGEADDR   = 0x22

	sh1106_COMSCANINC = 0xC0
	sh1106_COMSCANDEC = 0xC8

	sh1106_SEGREMAP = 0xA0

	sh1106_CHARGEPUMP = 0x8D

	sh1106_EXTERNALVCC  = 0x1
	sh1106_SWITCHCAPVCC = 0x2

	// Scrolling s
	sh1106_ACTIVATE_SCROLL                      = 0x2F
	sh1106_DEACTIVATE_SCROLL                    = 0x2E
	sh1106_SET_VERTICAL_SCROLL_AREA             = 0xA3
	sh1106_RIGHT_HORIZONTAL_SCROLL              = 0x26
	sh1106_LEFT_HORIZONTAL_SCROLL               = 0x27
	sh1106_VERTICAL_AND_RIGHT_HORIZONTAL_SCROLL = 0x29
	sh1106_VERTICAL_AND_LEFT_HORIZONTAL_SCROLL  = 0x2A
)
