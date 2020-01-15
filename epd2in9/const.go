// Copyright 2019, Homin Lee <homin.lee@suapapa.net>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package epd2in9

// epd2in9 module should be connected in following pins
const (
	PinRST  = "RST"  // 17 for Rpi
	PinDC   = "DC"   // 25 for Rpi
	PinBusy = "BUSY" // 24 for Rpi
)

const (
	epd2in9Width  = 128
	epd2in9Height = 296
)
