// Copyright 2019, Homin Lee <homin.lee@suapapa.net>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package epd2in9

var (
	epd2in9LutFullUpdate = []byte{
		0x50, 0xAA, 0x55, 0xAA, 0x11, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0xFF, 0xFF, 0x1F, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	epd2in9LutPartUpdate = []byte{
		0x10, 0x18, 0x18, 0x08, 0x18, 0x18,
		0x08, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x13, 0x14, 0x44, 0x12,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}
)

// InitFull initialize in full mode
func (d *Display) InitFull() {
	d.init(epd2in9LutFullUpdate)
}

// InitPart initialize in full mode
func (d *Display) InitPart() {
	d.init(epd2in9LutPartUpdate)
}

func (d *Display) init(lut []byte) {
	d.Reset()

	d.sendCmd(0x01) // DRIVER_OUTPUT_CONTROL
	d.sendData(uint8(d.h - 1))
	d.sendData(uint8((d.h - 1) >> 8))
	d.sendData(0x00) // GD = 0 SM = 0 TB = 0

	d.sendCmd(0x0C) // BOOSTER_SOFT_START_CONTROL
	d.sendData(0xD7)
	d.sendData(0xD6)
	d.sendData(0x9D)

	d.sendCmd(0x2C)  // WRITE_VCOM_REGISTER
	d.sendData(0xA8) // VCOM 7C

	d.sendCmd(0x3A)  // SET_DUMMY_LINE_PERIOD
	d.sendData(0x1A) // 4 dummy lines per gate

	d.sendCmd(0x3B)  // SET_GATE_TIME
	d.sendData(0x08) // 2us per line

	d.sendCmd(0x11)  // DATA_ENTRY_MODE_SETTING
	d.sendData(0x03) // X increment Y increment

	d.sendCmd(0x32) // WRITE_LUT_REGISTER

	d.sendDatas(lut)

	d.waitTillNotBusy()
}
