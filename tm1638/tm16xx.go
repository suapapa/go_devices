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
	v := intensity
	if v > 7 {
		v = 7
	}
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
	v := intensity
	if v > 7 {
		v = 7
	}
	if active {
		v |= 8
	}
	d.sendCmd(0x80 | intensity)

	d.strobe.Clear()
	d.clk.Clear()
	d.clk.Set()
	d.strobe.Set()
}

// DisplayDigit displays a digit
func (d tm16xx) DisplayDigit(digit byte, pos int, dot bool) {
	panic("not implemented")
}

// DisplayError display Error
func (d tm16xx) DisplayError() {
	panic("not implemented")
}

// ClearDigit clear digit in given position
func (d tm16xx) ClearDigit(pos int, dot bool) {
	panic("not implemented")
}

func (d tm16xx) setDisplay(val []byte) {
	panic("not implemented")
}

func (d tm16xx) sendCmd(cmd byte) {
	panic("not implemented")
}

func (d tm16xx) sendData(addr, data byte) {
	panic("not implemented")
}

func (d tm16xx) send(data byte) {
	panic("not implemented")
}

func (d tm16xx) receive() (temp byte) {
	panic("not implemented")
}

func (d tm16xx) sendChar(pos byte, data byte, dot bool) {
	panic("not implemented")
}

// Close closes all open pins
func (d tm16xx) Close() {
	panic("not implemented")
}
