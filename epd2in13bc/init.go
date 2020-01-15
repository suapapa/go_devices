// Copyright 2019, Homin Lee <homin.lee@suapapa.net>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package epd2in13bc

import "fmt"

// InitFull initialize in full mode
func (d *Display) InitFull() {
	d.Reset()

	d.sendCmd(0x06) // BOOSTER_SOFT_START
	d.sendData(0x17)
	d.sendData(0x17)
	d.sendData(0x17)

	d.sendCmd(0x04) // POWER_ON
	d.waitTillNotBusy()

	d.sendCmd(0x00) // PANEL_SETTING
	d.sendData(0x8F)

	d.sendCmd(0x50) // VCOM_AND_DATA_INTERVAL_SETTING
	d.sendData(0xF0)

	d.sendCmd(0x61) // RESOLUTION_SETTING
	d.sendData(uint8(d.w) & 0xff)
	d.sendData(uint8(d.h >> 8))
	d.sendData(uint8(d.h) & 0xff)
}

// InitPart initialize in part mode
func (d *Display) InitPart() {
	panic(fmt.Errorf("epd2in13bc not have partial update mode"))
}
