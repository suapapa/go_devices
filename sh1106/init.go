// Copyright 2016, Homin Lee <homin.lee@suapapa.net>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sh1106

func (l *LCD) init128x64() {
	l.sendCmd(sh1106_DISPLAYOFF)         // 0xAE
	l.sendCmd(sh1106_SETDISPLAYCLOCKDIV) // 0xD5
	l.sendCmd(0x80)                      // the suggested ratio 0x80
	l.sendCmd(sh1106_SETMULTIPLEX)       // 0xA8
	l.sendCmd(0x3F)
	l.sendCmd(sh1106_SETDISPLAYOFFSET) // 0xD3
	l.sendCmd(0x00)                    // no offset

	l.sendCmd(sh1106_SETSTARTLINE | 0x0) // line #0 0x40
	l.sendCmd(sh1106_CHARGEPUMP)         // 0x8D

	// if (vccstate == sh1106_EXTERNALVCC) {
	// l.sendCmd(0x10)
	// } else {
	l.sendCmd(0x14)
	// }

	l.sendCmd(sh1106_MEMORYMODE) // 0x20
	l.sendCmd(0x00)              // 0x0 act like ks0108
	l.sendCmd(sh1106_SEGREMAP | 0x1)
	l.sendCmd(sh1106_COMSCANDEC)
	l.sendCmd(sh1106_SETCOMPINS) // 0xDA
	l.sendCmd(0x12)
	l.sendCmd(sh1106_SETCONTRAST) // 0x81

	// if (vccstate == sh1106_EXTERNALVCC) {
	// l.sendCmd(0x9F)
	// } else {
	l.sendCmd(0xCF)
	// }

	l.sendCmd(sh1106_SETPRECHARGE) // 0xd9

	//  if (vccstate == sh1106_EXTERNALVCC) {
	// l.sendCmd(0x22)
	// } else {
	l.sendCmd(0xF1)
	// }

	l.sendCmd(sh1106_SETVCOMDETECT) // 0xDB
	l.sendCmd(0x40)
	l.sendCmd(sh1106_DISPLAYALLON_RESUME) // 0xA4
	l.sendCmd(sh1106_NORMALDISPLAY)       // 0xA6
}

func (l *LCD) init128x32() {
	panic("not implemented")
}

func (l *LCD) init96x16() {
	panic("not implemented")
}
