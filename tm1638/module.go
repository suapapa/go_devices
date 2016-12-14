// Copyright 2015-2016, Homin Lee <homin.lee@suapapa.net>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tm1638

import (
	"github.com/goiot/exp/gpio"
	gpio_driver "github.com/goiot/exp/gpio/driver"
)

// Module represents tm1638 based module
type Module struct {
	dev *gpio.Device
}

// Open opens a tm1638 Module
// gpio driver should have following pins:
//   * "DATA" : data pin
//   * "CLK" : clock pin
//   * "STB" : strobe pin
func Open(d gpio_driver.Opener) (*Module, error) {
	gpioDev, err := gpio.Open(d)
	if err != nil {
		return nil, err
	}

	if err = gpioDev.SetDirection(PinDATA, gpio.Out); err != nil {
		return nil, err
	}
	if err = gpioDev.SetDirection(PinCLK, gpio.Out); err != nil {
		return nil, err
	}
	if err = gpioDev.SetDirection(PinSTB, gpio.Out); err != nil {
		return nil, err
	}

	return &Module{
		dev: gpioDev,
	}, nil
}

// Close closes tm1638 module
func (m *Module) Close() {
	m.dev.Close()
}

// SetLed sets led in pos to given color
func (m *Module) SetLed(pos int, led Color) {
	m.sendData(byte(pos<<1)+1, byte(led))
}

// SetFND sets FND in pos to data.
// The bits in the data are displayed by mapping bellow
//  -- 0 --
// |       |
// 5       1
//  -- 6 --
// 4       2
// |       |
//  -- 3 --  .7
func (m *Module) SetFND(pos int, data byte) {
	m.sendData(byte(pos)<<1, data)
}

// SetChar sets FND in given position to given character
func (m *Module) SetChar(pos int, c rune, dot bool) {
	data, ok := font[c]
	if !ok {
		data = 0x00
	}
	if dot {
		data |= 0x80
	}
	m.sendData(byte(pos)<<1, data)
}

// SetString sets FND to given str
func (m *Module) SetString(str string) {
	i := 0
	for _, r := range str {
		m.SetChar(i, r, false)
		i++
	}
}

// GetButtons read buttons
func (m *Module) GetButtons() (keys byte) {
	m.dev.SetValue(PinSTB, 0)
	m.send(0x042)

	// TODO: why it repeats 4 time?
	for i := 0; i < 4; i++ {
		keys |= (m.receive() << uint(i))
	}
	m.dev.SetValue(PinSTB, 1)

	return
}

func (m *Module) sendData(addr, data byte) {
	m.sendCmd(0x44)
	m.dev.SetValue(PinSTB, 0)
	m.send(0xC0 | addr)
	m.send(data)
	m.dev.SetValue(PinSTB, 1)
}

func (m *Module) sendCmd(cmd byte) {
	m.dev.SetValue(PinSTB, 0)
	m.send(cmd)
	m.dev.SetValue(PinSTB, 1)
}

func (m *Module) send(data byte) {
	for i := 0; i < 8; i++ {
		m.dev.SetValue(PinCLK, 0)
		if data&1 == 0 {
			m.dev.SetValue(PinDATA, 0)
		} else {
			m.dev.SetValue(PinDATA, 1)
		}
		data >>= 1
		m.dev.SetValue(PinCLK, 1)
	}
}

func (m *Module) receive() (data byte) {
	m.dev.SetDirection(PinDATA, gpio.In)
	m.dev.SetValue(PinDATA, 1) // TODO: is this makes data pin pull up?

	for i := 0; i < 8; i++ {
		data >>= 1
		m.dev.SetValue(PinCLK, 0)
		if b, err := m.dev.Value(PinDATA); err == nil && b == 1 {
			data |= 0x80
		}
		m.dev.SetValue(PinCLK, 1)
	}

	m.dev.SetDirection(PinDATA, gpio.Out)
	m.dev.SetValue(PinDATA, 0)

	return
}
