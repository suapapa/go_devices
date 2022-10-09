// Copyright 2020, Homin Lee <homin.lee@suapapa.net>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pcf8574clcd

import (
	"fmt"
	"time"

	"periph.io/x/conn/v3"
	"periph.io/x/conn/v3/i2c"
)

// DefaultAddr is defulat address of pcf8574
var DefaultAddr = uint16(0x20)

// New returns a Dev object that communicates over I²C to a pcf8574 display
// controller.
func New(i i2c.Bus, addr uint16, row, col byte) (*Dev, error) {
	// Maximum clock speed is 1/2.5µs = 400KHz.
	return newDev(&i2c.Dev{Bus: i, Addr: addr}, row, col)
}

// Dev is an open handle to the display controller.
type Dev struct {
	columms  byte
	numLines byte

	displayFunction byte
	displayControl  byte
	displayMode     byte
	backlightVal    byte
	// Communication
	c conn.Conn
}

func (d *Dev) String() string {
	return fmt.Sprintf("pcf8574clcd.Dev{%s}", d.c)
}

// New is the common initialization code that is independent of the
// communication protocol (I²C) being used.
func newDev(c conn.Conn, col, row byte) (*Dev, error) {
	d := &Dev{
		c:            c,
		columms:      col,
		numLines:     row,
		backlightVal: lcdNOBACKLIGHT,
	}

	d.begin(col, row)

	return d, nil
}

func (d *Dev) begin(cols, lines byte) {
	d.displayFunction = lcd4BITMODE | lcd1LINE | lcd5x8DOTS
	if lines > 1 {
		d.displayFunction |= lcd2LINE
	}

	// // for some 1 line displays you can select a 10 pixel high font
	// if (dotsize != 0) && (lines == 1) {
	// 	d.displayFunction |= lcd5x10DOTS
	// }

	// SEE PAGE 45/46 FOR INITIALIZATION SPECIFICATION!
	// according to datasheet, we need at least 40ms after power rises above 2.7V
	// before sending commands. Arduino can turn on way befer 4.5V so we'll wait 50
	time.Sleep(50 * time.Millisecond)

	// Now we pull both RS and R/W low to begin commands
	d.expanderWrite(d.backlightVal) // reset expanderand turn backlight off (Bit 8 =1)
	time.Sleep(1000 * time.Millisecond)

	//put the LCD into 4 bit mode
	// this is according to the hitachi HD44780 datasheet
	// figure 24, pg 46

	// we start in 8bit mode, try to set 4 bit mode
	d.write4bits(0x03 << 4)
	time.Sleep(4500 * time.Microsecond)

	// second try
	d.write4bits(0x03 << 4)
	time.Sleep(150 * time.Microsecond)

	// finally, set to 4-bit interface
	d.write4bits(0x02 << 4)

	// set # lines, font size, etc.
	d.cmd(lcdFUNCTIONSET | d.displayFunction)

	// turn the display on with no cursor or blinking default
	d.displayControl = lcdDISPLAYON | lcdCURSOROFF | lcdBLINKOFF
	d.Display(true)

	// clear it off
	d.Clear()

	// Initialize to default text direction (for roman languages)
	d.displayMode = lcdENTRYLEFT | lcdENTRYSHIFTDECREMENT

	// set the entry mode
	d.cmd(lcdENTRYMODESET | d.displayMode)

	d.Home()
}

func (d *Dev) Write(str string) error {
	for _, v := range []byte(str) {
		err := d.send(v, bitData)
		if err != nil {
			return err
		}
	}
	return nil
}

// Clear clear display and sets cursor position to zero
func (d *Dev) Clear() error {
	err := d.cmd(lcdCLEARDISPLAY)
	time.Sleep(2000 * time.Microsecond)
	return err
}

// Home sets cursor position to zero
func (d *Dev) Home() error {
	err := d.cmd(lcdRETURNHOME)
	time.Sleep(2000 * time.Microsecond)
	return err
}

// SetCursor sets cursor to given col, row position
func (d *Dev) SetCursor(col, row byte) error {
	rowOffset := []byte{0x00, 0x40, 0x14, 0x54}
	if row > d.numLines {
		row = d.numLines - 1
	}
	return d.cmd(lcdSETDDRAMADDR | (col + rowOffset[row]))
}

// Display turns on/off display (quickly)
func (d *Dev) Display(on bool) error {
	if on {
		d.displayControl |= lcdDISPLAYON
	} else {
		d.displayControl &= ^byte(lcdDISPLAYON)
	}
	return d.cmd(lcdDISPLAYCONTROL | d.displayControl)
}

// Cursor turns on/off underline cursor
func (d *Dev) Cursor(on bool) error {
	if on {
		d.displayControl |= lcdCURSORON
	} else {
		d.displayControl &= ^byte(lcdCURSORON)
	}
	return d.cmd(lcdDISPLAYCONTROL | d.displayControl)
}

// Blink turns on/off ther blinking cursor
func (d *Dev) Blink(on bool) error {
	if on {
		d.displayControl |= lcdBLINKON
	} else {
		d.displayControl &= ^byte(lcdBLINKON)
	}
	return d.cmd(lcdDISPLAYCONTROL | d.displayControl)
}

// ScrollLeft scroll display left
func (d *Dev) ScrollLeft() error {
	return d.cmd(lcdCURSORSHIFT | lcdDISPLAYMOVE | lcdMOVELEFT)
}

// ScrollRight scroll display right
func (d *Dev) ScrollRight() error {
	return d.cmd(lcdCURSORSHIFT | lcdDISPLAYMOVE | lcdMOVERIGHT)
}

// Left2Right flows text left to right
func (d *Dev) Left2Right() error {
	d.displayMode |= lcdENTRYLEFT
	return d.cmd(lcdENTRYMODESET | d.displayMode)
}

// Right2Left flows text right to left
func (d *Dev) Right2Left() error {
	d.displayMode &= ^byte(lcdENTRYLEFT)
	return d.cmd(lcdENTRYMODESET | d.displayMode)
}

// AutoScroll will 'right  justrify' of 'left justify' text from cursor
func (d *Dev) AutoScroll(on bool) error {
	if on {
		d.displayMode |= lcdENTRYSHIFTINCREMENT
	} else {
		d.displayMode &= ^byte(lcdENTRYSHIFTINCREMENT)
	}
	return d.cmd(lcdENTRYMODESET | d.displayMode)
}

// CreateChar allows to fill the first 8 CGRAM location with custom characters
func (d *Dev) CreateChar(location byte, charMap []byte) error {
	location &= 0x07
	err := d.cmd(lcdSETCGRAMADDR | (location << 3))
	if err != nil {
		return err
	}
	for _, v := range charMap {
		err := d.send(v, bitRs)
		if err != nil {
			return err
		}
	}
	return nil
}

// BackLight sets backlight brightness
func (d *Dev) BackLight(on bool) error {
	if on {
		d.backlightVal = lcdBACKLIGHT
	} else {
		d.backlightVal = lcdNOBACKLIGHT
	}
	return d.expanderWrite(0)
}

/************ low level data pushing commands **********/
func (d *Dev) cmd(value byte) error {
	return d.send(value, 0)
}

func (d *Dev) send(value, mode byte) error {
	highNib := value & 0xf0
	lowNib := (value << 4) & 0xf0
	err := d.write4bits(highNib | mode)
	if err != nil {
		return err
	}
	return d.write4bits(lowNib | mode)
}

func (d *Dev) write4bits(value byte) error {
	err := d.expanderWrite(value)
	if err != nil {
		return err
	}
	return d.pulseEnable(value)
}

func (d *Dev) expanderWrite(data byte) error {
	return d.c.Tx([]byte{(data | d.backlightVal)}, nil)
}

func (d *Dev) pulseEnable(data byte) error {
	var err error
	err = d.expanderWrite(data | bitEn)
	if err != nil {
		return err
	}
	time.Sleep(1 * time.Microsecond)
	err = d.expanderWrite(data & ^byte(bitEn))
	if err != nil {
		return err
	}
	time.Sleep(50 * time.Microsecond)
	return nil
}

// func (d *Dev) sendCmd4(c byte) error {
// 	buf := make([]byte, 4)
// 	buf[0] = (c >> 4) | 0x10
// 	buf[1] = buf[0] & 0xEF
// 	buf[2] = (c & 0x0F)
// 	buf[3] = buf[2] & 0xEF
// 	return d.c.Tx(buf, nil)
// }
