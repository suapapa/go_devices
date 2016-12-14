// Copyright 2016, Homin Lee <homin.lee@suapapa.net>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gpio

// DefaultPinMap maps Broadcom pin number string to int
// ref. https://pinout.xyz
var DefaultPinMap = map[string]int{
	"BCM2":  2,  // 3, SDA
	"BCM3":  3,  // 5, SCL
	"BCM4":  4,  // 7, GPCLK0
	"BCM14": 14, // 8, TXD
	"BCM15": 15, // 10, RXD
	"BCM17": 17, // 11
	"BCM18": 18, // 12, PWM0
	"BCM27": 27, // 13
	"BCM22": 22, // 15
	"BCM23": 23, // 16
	"BCM24": 24, // 18
	"BCM10": 10, // 19, MOSI
	"BCM9":  9,  // 21, MISO
	"BCM11": 11, // 23, SCLK
	"BCM8":  8,  // 24, CE0
	"BCM7":  7,  // 26, CD1
	"BCM0":  0,  // 27, ID_SD
	"BCM1":  1,  // 28, ID_SC
	"BCM5":  5,  // 29
	"BCM6":  6,  // 31
	"BCM12": 12, // 32, PWM0
	"BCM13": 13, // 33, PWM1
	"BCM19": 19, // 35, MISO
	"BCM16": 16, // 36
	"BCM26": 26, // 37
	"BCM20": 20, // 38, MOSI
	"BCM21": 21, // 40, SCLK
}
