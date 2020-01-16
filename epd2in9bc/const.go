// Copyright 2019, Homin Lee <homin.lee@suapapa.net>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package epd2in9bc

type inkColor byte

// epd2in9bc module should be connected in following pins
const (
	PinRST  = "RST"  // 17 for Rpi
	PinDC   = "DC"   // 25 for Rpi
	PinBusy = "BUSY" // 24 for Rpi

	epd2in9bcWidth  = 128
	epd2in9bcHeight = 296

	black inkColor = iota
	white
	accent
)
