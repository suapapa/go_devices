// Copyright 2019, Homin Lee <homin.lee@suapapa.net>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package epd2in9bc

import "time"

// Init initialize display panel
func (d *Display) Init() {
	d.Reset()
	time.Sleep(time.Second)
	d.init()
}

func (d *Display) init() {

	d.sendCmd(0x06) // boost
	// d.sendData(0x17)
	// d.sendData(0x17)
	// d.sendData(0x17)
	d.sendDatas([]byte{0x17, 0x17, 0x17})
	d.sendCmd(0x04) // POWER_ON
	d.waitTillNotBusy()
	d.sendCmd(0x00) // PANEL_SETTING
	d.sendData(0x8F)
	d.sendCmd(0x50) // VCOM_AND_DATA_INTERVAL_SETTING
	d.sendData(0x77)
	d.sendCmd(0x61) // TCON_RESOLUTION
	// d.sendData(0x80)
	// d.sendData(0x01)
	// d.sendData(0x28)
	d.sendDatas([]byte{0x80, 0x01, 0x28})

	d.sendCmd(0x82)
	d.sendData(0x0A)

	d.waitTillNotBusy()
}
