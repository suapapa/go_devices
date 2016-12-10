// Copyright 2015, Homin Lee <homin.lee@suapapa.net>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tm1638

import "github.com/davecheney/gpio"

// tm16xx represent a tm16xx module
type tm16xx struct {
	data, clk, strobe gpio.Pin
	displays          int
}

// Newtm16xx returns point of tm16xx which is initialized given gpio numbers
func newTm16xx(data, clk, strobe int,
	activeDisplay bool, intensity byte) (*tm16xx, error) {

	var d tm16xx
	var err error

	d.data, err = gpio.OpenPin(data, gpio.ModeOutput)
	if err != nil {
		return nil, err
	}
	d.clk, err = gpio.OpenPin(clk, gpio.ModeOutput)
	if err != nil {
		return nil, err
	}
	d.strobe, err = gpio.OpenPin(strobe, gpio.ModeOutput)
	if err != nil {
		return nil, err
	}

	d.data.SetMode(gpio.ModeOutput)
	d.clk.SetMode(gpio.ModeOutput)
	d.strobe.SetMode(gpio.ModeOutput)

	d.strobe.Set()
	d.clk.Set()

	d.sendCmd(0x40)
	v := min(7, intensity)
	if activeDisplay {
		v |= 8
	}
	d.sendCmd(0x80 | v)

	d.strobe.Clear()
	d.sendCmd(0xC0)
	for i := 0; i < 16; i++ {
		d.sendCmd(0x00)
	}

	d.strobe.Set()
	return &d, nil
}

// SetupDisplay initialized the display
func (d tm16xx) SetupDisplay(active bool, intensity byte) {
	v := min(7, intensity)
	if active {
		v |= 8
	}
	d.sendCmd(0x80 | v)

	d.strobe.Clear()
	d.clk.Clear()
	d.clk.Set()
	d.strobe.Set()
}

// DisplayDigit displays a digit
func (d tm16xx) DisplayDigit(digit byte, pos int, dot bool) {
	d.sendChar(byte(pos), fontNumber[digit&0x0F], dot)
}

// DisplayError display Error
func (d tm16xx) DisplayError() {
	d.setDisplay(fontErrorData)
}

// ClearDigit clear digit in given position
func (d tm16xx) ClearDigit(pos int, dot bool) {
	d.sendChar(byte(pos), 0, dot)
}

func (d tm16xx) setDisplay(val []byte) {
	for i, c := range val {
		d.sendChar(byte(i), c, false)
	}
}

func (d tm16xx) sendCmd(cmd byte) {
	d.strobe.Clear()
	d.send(cmd)
	d.strobe.Set()
}

func (d tm16xx) sendData(addr, data byte) {
	d.sendCmd(0x44)
	d.strobe.Clear()
	d.send(0xC0 | addr)
	d.send(data)
	d.strobe.Set()
}

func (d tm16xx) send(data byte) {
	for i := 0; i < 8; i++ {
		d.clk.Clear()
		if data&1 == 0 {
			d.data.Clear()
		} else {
			d.data.Set()
		}
		data >>= 1
		d.clk.Set()
	}
}

func (d tm16xx) receive() (temp byte) {
	d.data.SetMode(gpio.ModeInput)
	d.data.Set() // TODO: is this makes data pin pull up?

	for i := 0; i < 8; i++ {
		temp >>= 1
		d.clk.Clear()
		if d.data.Get() {
			temp |= 0x80
		}
		d.clk.Set()
	}

	d.data.SetMode(gpio.ModeOutput)
	d.data.Clear()

	return
}

func (d tm16xx) sendChar(pos byte, data byte, dot bool) {
	if dot {
		data |= 0x80
	}
	d.sendData(pos<<1, data)
}

// Close closes all open pins
func (d tm16xx) Close() {
	d.clk.Close()
	d.data.Close()
	d.strobe.Close()
}
