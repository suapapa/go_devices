package nunchuck // import "github.com/suapapa/go_devices/nunchuck"

// REF: http://playground.arduino.cc/Main/WiiChuckClass

import (
	"sync"
	"time"

	"golang.org/x/exp/io/i2c"
	i2c_driver "golang.org/x/exp/io/i2c/driver"
)

const (
	addr = 0x52

	radius = 210

	dftZeroJoyX   = 124
	dftZeroJoyY   = 132
	dftZeroAccelX = 510
	dftZeroAccelY = 490
	dftZeroAccelZ = 460
)

// consts for internel buff index
const (
	idxJoyX = iota
	idxJoyY
	idxAccelX
	idxAccelY
	idxAccelZ
	idxMisc

	maskMiscAccel       = 0x03
	maskMiscAccelXShift = 2
	maskMiscAccelYShift = 4
	maskMiscAccelZShift = 6
	maskMiscBtnZ        = 0x01
	maskMiscBtnC        = 0x02
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

func (c *Controller) Close() {
	c.dev.Close()
}

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

func (c *Controller) Update() {
	c.Lock()
	defer c.Unlock()

	c.dev.Read(c.buff[:])
	c.dev.Write([]byte{0x00})
}

func (c *Controller) JoyX() int {
	c.RLock()
	defer c.Unlock()

	return int(c.buff[idxJoyX]) - c.zeroJoyX
}

func (c *Controller) JoyY() int {
	c.RLock()
	defer c.Unlock()

	return int(c.buff[idxJoyY]) - c.zeroJoyY
}

func (c *Controller) AccelX() int {
	c.RLock()
	defer c.Unlock()

	return c.accelX() - c.zeroAccelX
}

func (c *Controller) AccelY() int {
	c.RLock()
	defer c.Unlock()

	return c.accelY() - c.zeroAccelY
}

func (c *Controller) AccelZ() int {
	c.RLock()
	defer c.Unlock()

	return c.accelZ() - c.zeroAccelZ
}

func (c *Controller) BtnZ() bool {
	c.RLock()
	defer c.Unlock()

	return (c.buff[idxMisc] & maskMiscBtnZ) == 0
}

func (c *Controller) BtnC() bool {
	c.RLock()
	defer c.Unlock()

	return (c.buff[idxMisc] & maskMiscBtnC) == 0
}

func (c *Controller) String() string {
	// TODO:
	return ""
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
