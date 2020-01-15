// Copyright 2019, Homin Lee <homin.lee@suapapa.net>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package epd2in13bc

type inkColor uint8

// epd2in13bc module should be connected in following pins
const (
	PinRST  = "RST"  // 17 for Rpi
	PinDC   = "DC"   // 25 for Rpi
	PinBusy = "BUSY" // 24 for Rpi

	// 104x212
	epd2in13bcWidth  = 104
	epd2in13bcHeight = 212

	// epd2in13bc use 3 colors black, white and accent color (red or yellow)
	black  inkColor = 0
	white  inkColor = 1
	accent inkColor = 2
)
