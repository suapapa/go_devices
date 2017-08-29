// Copyright 2016, Homin Lee <homin.lee@suapapa.net>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gpio

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"

	gpio_driver "github.com/goiot/exp/gpio/driver"
)

// Mem implements github.com/goiot/exp/gpio/driver.Opener
type Mem struct {
	// PinMap will is used to conver pin name to number
	PinMap map[string]int
}

// Open returns github.com/goiot/exp/gpio/driver.Conn
func (m *Mem) Open() (gpio_driver.Conn, error) {
	conn := &memConn{}
	err := conn.mmap()
	if err != nil {
		return nil, err
	}

	conn.pinMap = m.PinMap

	return conn, nil
}

// implements github.com/goiot/exp/gpio/driver.Conn
type memConn struct {
	buf                         []byte
	gpfsel, gpset, gpclr, gplev []*uint32
	pinMap                      map[string]int
}

// Value returns the value of the pin. 0 for low values, 1 for high.
func (c *memConn) Value(pin string) (int, error) {
	p, ok := c.pinMap[pin]
	if !ok {
		return 0, fmt.Errorf("gpio: unknown gpio, %s", pin)
	}

	offset, shift := p/32, byte(p%32)
	v := 0
	if *c.gplev[offset]&(1<<shift) == (1 << shift) {
		v = 1
	}
	return v, nil
}

// SetValue sets the vaule of the pin. 0 for low values, 1 for high.
func (c *memConn) SetValue(pin string, v int) error {
	p, ok := c.pinMap[pin]
	if !ok {
		return fmt.Errorf("gpio: unknown gpio, %s", pin)
	}

	offset, shift := p/32, byte(p%32)
	if v == 1 {
		*c.gpset[offset] = (1 << shift)
	} else {
		*c.gpclr[offset] = (1 << shift)
	}
	return nil
}

// SetDirection sets the direction of the pin.
func (c *memConn) SetDirection(pin string, dir gpio_driver.Direction) error {
	p, ok := c.pinMap[pin]
	if !ok {
		return fmt.Errorf("gpio: unknown gpio, %s", pin)
	}

	offset, shift := p/10, uint32(p%10)*3
	mask := bcm2835gpioFselMASK << shift
	var mode uint32
	switch dir {
	case gpio_driver.In:
		mode = bcm2835gpioFselINPT
	case gpio_driver.Out:
		mode = bcm2835gpioFselOUTP
	default:
		return fmt.Errorf("rpi: uknnown gpiodir, %v", dir)
	}

	v := *c.gpfsel[offset]
	v &= ^mask
	v |= mode << shift

	*c.gpfsel[offset] = v & mask

	return nil
}

// Map maps virtual gpio pin name to a physical pin number
func (c *memConn) Map(virtual string, physical int) {
	c.pinMap[virtual] = physical
}

// Close closes the connection and free the underlying resources.
func (c *memConn) Close() error {
	return syscall.Munmap(c.buf)
}

func (c *memConn) mmap() error {
	f, err := os.OpenFile("/dev/gpiomem", os.O_RDWR|os.O_SYNC, 0)
	if err != nil {
		return err
	}
	defer f.Close()

	buf, err := syscall.Mmap(int(f.Fd()),
		0, 4*1024,
		syscall.PROT_READ|syscall.PROT_WRITE,
		syscall.MAP_SHARED,
	)
	if err != nil {
		return err
	}

	c.gpfsel = []*uint32{
		(*uint32)(unsafe.Pointer(&buf[bcm2835gpFsel0])),
		(*uint32)(unsafe.Pointer(&buf[bcm2835gpFsel1])),
		(*uint32)(unsafe.Pointer(&buf[bcm2835gpFsel2])),
		(*uint32)(unsafe.Pointer(&buf[bcm2835gpFsel3])),
		(*uint32)(unsafe.Pointer(&buf[bcm2835gpFsel4])),
		(*uint32)(unsafe.Pointer(&buf[bcm2835gpFsel5])),
	}
	c.gpset = []*uint32{
		(*uint32)(unsafe.Pointer(&buf[bcm2835gpSet0])),
		(*uint32)(unsafe.Pointer(&buf[bcm2835gpSet1])),
	}
	c.gpclr = []*uint32{
		(*uint32)(unsafe.Pointer(&buf[bcm2835gpClr0])),
		(*uint32)(unsafe.Pointer(&buf[bcm2835gpClr1])),
	}
	c.gplev = []*uint32{
		(*uint32)(unsafe.Pointer(&buf[bcm2835gpLev0])),
		(*uint32)(unsafe.Pointer(&buf[bcm2835gpLev1])),
	}

	c.buf = buf
	return nil
}
