// Copyright 2019, Homin Lee <homin.lee@suapapa.net>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package epd7in5

import "time"

func (d *Display) init() {
	// EPD hardware init start
	d.Reset()
	time.Sleep(time.Second)

	d.sendCmd(0x01) // POWER_SETTING
	d.sendData(0x37)
	d.sendData(0x00)

	d.sendCmd(0x00) // PANEL_SETTING
	d.sendData(0xCF)
	d.sendData(0x08)

	d.sendCmd(0x06) // BOOSTER_SOFT_START
	d.sendData(0xc7)
	d.sendData(0xcc)
	d.sendData(0x28)

	d.sendCmd(0x04) // POWER_ON
	d.waitTillNotBusy()

	d.sendCmd(0x30) // PLL_CONTROL
	d.sendData(0x3c)

	d.sendCmd(0x41) // TEMPERATURE_CALIBRATION
	d.sendData(0x00)

	d.sendCmd(0x50) // VCOM_AND_DATA_INTERVAL_SETTING
	d.sendData(0x77)

	d.sendCmd(0x60) // TCON_SETTING
	d.sendData(0x22)

	d.sendCmd(0x61)            // TCON_RESOLUTION
	d.sendData(byte(d.w >> 8)) //source 640
	d.sendData(byte(d.w & 0xff))
	d.sendData(byte(d.h >> 8)) //gate 384
	d.sendData(byte(d.h & 0xff))

	d.sendCmd(0x82)  // VCM_DC_SETTING
	d.sendData(0x1E) // decide by LUT file

	d.sendCmd(0xe5) // FLASH MODE
	d.sendData(0x03)
}
