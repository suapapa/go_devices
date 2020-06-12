// Copyright 2020, Homin Lee <homin.lee@suapapa.net>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pcf8574clcd

const (
	// commands
	lcdCLEARDISPLAY   = 0x01
	lcdRETURNHOME     = 0x02
	lcdENTRYMODESET   = 0x04
	lcdDISPLAYCONTROL = 0x08
	lcdCURSORSHIFT    = 0x10
	lcdFUNCTIONSET    = 0x20
	lcdSETCGRAMADDR   = 0x40
	lcdSETDDRAMADDR   = 0x80

	// flags for display entry mode
	lcdENTRYRIGHT          = 0x00
	lcdENTRYLEFT           = 0x02
	lcdENTRYSHIFTINCREMENT = 0x01
	lcdENTRYSHIFTDECREMENT = 0x00

	// flags for display on/off control
	lcdDISPLAYON  = 0x04
	lcdDISPLAYOFF = 0x00
	lcdCURSORON   = 0x02
	lcdCURSOROFF  = 0x00
	lcdBLINKON    = 0x01
	lcdBLINKOFF   = 0x00

	// flags for display/cursor shift
	lcdDISPLAYMOVE = 0x08
	lcdCURSORMOVE  = 0x00
	lcdMOVERIGHT   = 0x04
	lcdMOVELEFT    = 0x00

	// flags for function set
	lcd8BITMODE = 0x10
	lcd4BITMODE = 0x00
	lcd2LINE    = 0x08
	lcd1LINE    = 0x00
	lcd5x10DOTS = 0x04
	lcd5x8DOTS  = 0x00

	// flags for backlight control
	lcdBACKLIGHT   = 0x08
	lcdNOBACKLIGHT = 0x00

	bitEn   = 0B00000100 // Enable bit
	bitRw   = 0B00000010 // Read/Write bit
	bitRs   = 0B00000001 // Register select bit
	bitData = bitEn | bitRs
)
