// Copyright 2015-2016, Homin Lee <homin.lee@suapapa.net>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tm1638

// TM1638 module should be conncted in following gpio pins
const (
	PinSTB  = "STB"
	PinCLK  = "CLK"
	PinDATA = "DATA"
)

// Color is type for LED colors
type Color byte

// Available colors for leds
const (
	Off Color = iota
	Green
	Red
)
