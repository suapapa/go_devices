// Copyright 2019, Homin Lee <homin.lee@suapapa.net>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package epd2in13

// epd2in13 module should be connected in following pins
const (
	PinRST  = "RST"  // 17 for Rpi
	PinDC   = "DC"   // 25 for Rpi
	PinCS   = "CS"   // 8 for Rpi
	PinBusy = "BUSY" // 24 for Rpi
)

const (
	// 122x250
	epd2in13Width  = 122
	epd2in13Height = 250

	epd2in13FullUpdate = 0
	epd2in13PartUpdate = 1
)
