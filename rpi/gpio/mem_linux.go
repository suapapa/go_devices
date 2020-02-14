// Copyright 2016, Homin Lee <homin.lee@suapapa.net>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build linux

package gpio

import (
	"fmt"

	gpio_driver "github.com/goiot/exp/gpio/driver"
	rpio "github.com/stianeikeland/go-rpio"
)

// Mem implements github.com/goiot/exp/gpio/driver.Opener
type Mem struct {
	// PinMap will is used to conver pin name to number
	PinMap map[string]int
}

// Open returns github.com/goiot/exp/gpio/driver.Conn
func (m *Mem) Open() (gpio_driver.Conn, error) {
	err := rpio.Open()
	if err != nil {
		return nil, err
	}

	conn := &memConn{
		pinMap: make(map[string]rpio.Pin),
	}

	for n, v := range m.PinMap {
		p := rpio.Pin(v)
		conn.pinMap[n] = p
	}

	return conn, nil
}

// implements github.com/goiot/exp/gpio/driver.Conn
type memConn struct {
	pinMap map[string]rpio.Pin
}

// Value returns the value of the pin. 0 for low values, 1 for high.
func (c *memConn) Value(pin string) (int, error) {
	p, ok := c.pinMap[pin]
	if !ok {
		return 0, fmt.Errorf("gpio: unknown gpio, %s", pin)
	}

	v := p.Read()

	if v == rpio.High {
		return 1, nil
	}
	return 0, nil
}

// SetValue sets the vaule of the pin. 0 for low values, 1 for high.
func (c *memConn) SetValue(pin string, v int) error {
	p, ok := c.pinMap[pin]
	if !ok {
		return fmt.Errorf("gpio: unknown gpio, %s", pin)
	}

	if v == 1 {
		p.High()
	} else {
		p.Low()
	}

	return nil
}

// SetDirection sets the direction of the pin.
func (c *memConn) SetDirection(pin string, dir gpio_driver.Direction) error {
	p, ok := c.pinMap[pin]
	if !ok {
		return fmt.Errorf("gpio: unknown gpio, %s", pin)
	}

	switch dir {
	case gpio_driver.In:
		p.Input()
	case gpio_driver.Out:
		p.Output()
	}

	return nil
}

// Map maps virtual gpio pin name to a physical pin number
func (c *memConn) Map(virtual string, physical int) {
	panic("not implemented")
}

// Close closes the connection and free the underlying resources.
func (c *memConn) Close() error {
	return rpio.Close()
}
