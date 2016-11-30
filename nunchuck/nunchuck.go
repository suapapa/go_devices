// Copyright 2016, Homin Lee <homin.lee@suapapa.net>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nunchuck // import "github.com/suapapa/go_devices/nunchuck"

import (
	"fmt"
	"math"
	"sync"
	"time"

	"golang.org/x/exp/io/i2c"
	i2c_driver "golang.org/x/exp/io/i2c/driver"
)

const (
	// for internel buff index
	idxJoyX = iota
	idxJoyY
	idxAccelX
	idxAccelY
	idxAccelZ
	idxMisc

	// I2C address
	addr = 0x52

	// default zero values
	dftZeroJoyX   = 124
	dftZeroJoyY   = 132
	dftZeroAccelX = 510
	dftZeroAccelY = 490
	dftZeroAccelZ = 460

	// for bitmask
	maskMiscAccel       = 0x03
	maskMiscAccelXShift = 2
	maskMiscAccelYShift = 4
	maskMiscAccelZShift = 6
	maskMiscBtnZ        = 0x01
	maskMiscBtnC        = 0x02

	radius = 210
)

// Controller represents a nunchuck controller
type Controller struct {
	dev  *i2c.Device
	buff [6]byte

	zeroJoyX, zeroJoyY                 int
	zeroAccelX, zeroAccelY, zeroAccelZ int

	sync.RWMutex
}

// Open opens a nunchuck controller
func Open(o i2c_driver.Opener) (*Controller, error) {
	dev, err := i2c.Open(o, addr)
	if err != nil {
		return nil, err
	}

	con := &Controller{
		dev:        dev,
		zeroJoyX:   dftZeroJoyX,
		zeroJoyY:   dftZeroJoyY,
		zeroAccelX: dftZeroAccelX,
		zeroAccelY: dftZeroAccelY,
		zeroAccelZ: dftZeroAccelZ,
	}

	err = con.init()
	if err != nil {
		return nil, err
	}

	return con, nil
}

// Close closed the controller
func (c *Controller) Close() {
	c.dev.Close()
}

// Update reads data from the controller
func (c *Controller) Update() (err error) {
	c.Lock()
	defer c.Unlock()

	err = c.dev.Read(c.buff[:])
	if err != nil {
		return
	}
	c.dev.Write([]byte{0x00})
	return
}

// Calibrate set fix zero position for joystick and accelometer
func (c *Controller) Calibrate() {
	c.Update()

	c.Lock()
	defer c.Unlock()

	c.zeroJoyX = int(c.buff[idxJoyX])
	c.zeroJoyY = int(c.buff[idxJoyY])
	c.zeroAccelX = c.accelX()
	c.zeroAccelX = c.accelY()
	c.zeroAccelX = c.accelZ()
}

// JoyX returns joystic X axis
func (c *Controller) JoyX() int {
	c.RLock()
	defer c.RUnlock()

	return int(c.buff[idxJoyX]) - c.zeroJoyX
}

// JoyY returns joystic Y axis
func (c *Controller) JoyY() int {
	c.RLock()
	defer c.RUnlock()

	return int(c.buff[idxJoyY]) - c.zeroJoyY
}

// AccelX returns accelometer X axis
func (c *Controller) AccelX() int {
	c.RLock()
	defer c.RUnlock()

	return c.accelX() - c.zeroAccelX
}

// AccelY returns accelometer Y axis
func (c *Controller) AccelY() int {
	c.RLock()
	defer c.RUnlock()

	return c.accelY() - c.zeroAccelY
}

// AccelZ returns accelometer Z axis
func (c *Controller) AccelZ() int {
	c.RLock()
	defer c.RUnlock()

	return c.accelZ() - c.zeroAccelZ
}

// BtnZ returns button Z is pressed
func (c *Controller) BtnZ() bool {
	c.RLock()
	defer c.RUnlock()

	return (c.buff[idxMisc] & maskMiscBtnZ) == 0
}

// BtnC returns button C is pressed
func (c *Controller) BtnC() bool {
	c.RLock()
	defer c.RUnlock()

	return (c.buff[idxMisc] & maskMiscBtnC) == 0
}

// Roll returns roll in degrees
func (c *Controller) Roll() float64 {
	aX, aZ := float64(c.AccelX()), float64(c.AccelZ())
	return math.Atan2(aX, aZ) / math.Pi * 180
}

// Pitch returns pitch in degrees
func (c *Controller) Pitch() float64 {
	aY := float64(c.AccelY())
	return math.Acos(aY/radius) / math.Pi * 180
}

// String implemets stringer interface
func (c *Controller) String() string {
	return fmt.Sprintf(
		"joy(X: %v, Y: %v) accel(X: %v, Y: %v, Z: %v) Btn(Z: %v, C: %v)",
		c.JoyX(), c.JoyY(),
		c.AccelX(), c.AccelY(), c.AccelZ(),
		c.BtnZ(), c.BtnC(),
	)
}

func (c *Controller) init() (err error) {
	err = c.dev.Write([]byte{0xf0, 0x55})
	if err != nil {
		return
	}
	time.Sleep(1 * time.Millisecond)
	err = c.dev.Write([]byte{0xfb, 0x00})
	return
}

func (c *Controller) accelX() int {
	return int(c.buff[idxAccelX])<<2 +
		int((c.buff[idxMisc]>>maskMiscAccelXShift)&maskMiscAccel)
}

func (c *Controller) accelY() int {
	return int(c.buff[idxAccelY])<<2 +
		int((c.buff[idxMisc]>>maskMiscAccelYShift)&maskMiscAccel)
}

func (c *Controller) accelZ() int {
	return int(c.buff[idxAccelZ])<<2 +
		int((c.buff[idxMisc]>>maskMiscAccelZShift)&maskMiscAccel)
}
