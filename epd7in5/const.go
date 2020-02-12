// Copyright 2019, Homin Lee <homin.lee@suapapa.net>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package epd7in5

type inkColor byte

// epd7in5 module should be connected in following pins
const (
	PinRST  = "RST"  // 17 for Rpi
	PinDC   = "DC"   // 25 for Rpi
	PinBusy = "BUSY" // 24 for Rpi

	epd7in5Width  = 640
	epd7in5Height = 384

	black inkColor = iota
	gray
	white
)
