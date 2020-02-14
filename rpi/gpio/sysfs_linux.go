// Copyright 2016, Homin Lee <homin.lee@suapapa.net>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build linux

package gpio

import (
	"fmt"

	"github.com/davecheney/gpio"
	gpio_driver "github.com/goiot/exp/gpio/driver"
)

// Sysfs implements github.com/goiot/exp/gpio/driver.Opener
type Sysfs struct {
	// PinMap will is used to conver pin name to number
	PinMap map[string]int
}

// Open returns github.com/goiot/exp/gpio/driver.Conn
func (m *Sysfs) Open() (gpio_driver.Conn, error) {
	conn := &sysfsConn{
		pinMap: make(map[string]gpio.Pin),
	}

	for n, v := range m.PinMap {
		p, err := gpio.OpenPin(v, gpio.ModeInput)
		if err != nil {
			conn.Close()
			return nil, err
		}
		conn.pinMap[n] = p
	}

	return conn, nil
}

// implements github.com/goiot/exp/gpio/driver.Conn
type sysfsConn struct {
	pinMap map[string]gpio.Pin
}

// Value returns the value of the pin. 0 for low values, 1 for high.
func (c *sysfsConn) Value(pin string) (int, error) {
	p, ok := c.pinMap[pin]
	if !ok {
		return 0, fmt.Errorf("gpio: unknown gpio, %s", pin)
	}

	if p.Get() {
		return 1, nil
	}

	return 0, nil
}

// SetValue sets the vaule of the pin. 0 for low values, 1 for high.
func (c *sysfsConn) SetValue(pin string, v int) error {
	p, ok := c.pinMap[pin]
	if !ok {
		return fmt.Errorf("gpio: unknown gpio, %s", pin)
	}

	if v == 1 {
		p.Set()
	} else {
		p.Clear()
	}
	return nil
}

// SetDirection sets the direction of the pin.
func (c *sysfsConn) SetDirection(pin string, dir gpio_driver.Direction) error {
	p, ok := c.pinMap[pin]
	if !ok {
		return fmt.Errorf("gpio: unknown gpio, %s", pin)
	}

	switch dir {
	case gpio_driver.In:
		p.SetMode(gpio.ModeInput)
	case gpio_driver.Out:
		p.SetMode(gpio.ModeOutput)
	default:
		return fmt.Errorf("rpi: uknnown gpiodir, %v", dir)
	}
	return nil
}

// Map maps virtual gpio pin name to a physical pin number
func (c *sysfsConn) Map(virtual string, physical int) {
	panic("not implemented")
}

// Close closes the connection and free the underlying resources.
func (c *sysfsConn) Close() error {
	for _, p := range c.pinMap {
		p.Close()
	}
	return nil
}
