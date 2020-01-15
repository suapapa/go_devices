// Copyright 2019, Homin Lee <homin.lee@suapapa.net>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package epd2in13bc // import "github.com/suapapa/go_devices/epd2in13bc"

import (
	"fmt"
	"log"
	"time"

	"github.com/goiot/exp/gpio"
	gpio_driver "github.com/goiot/exp/gpio/driver"
	"golang.org/x/exp/io/spi"
	spi_driver "golang.org/x/exp/io/spi/driver"
)

// Display represents a epd2in13bc e-ink display
type Display struct {
	spiDev  *spi.Device // recommend 4000000Hz, spimode 0
	gpioDev *gpio.Device

	w, h int
}

// Open opens a epd2in13bc display in SPI mode
// gpio device should have PinRST pin for Reset() and
// PinDC pin for selecting data/cmd
func Open(bus spi_driver.Opener, ctr gpio_driver.Opener) (*Display, error) {
	spiDev, err := spi.Open(bus)
	if err != nil {
		return nil, err
	}
	// spiDev.SetCSChange(false) // TODO

	gpioDev, err := gpio.Open(ctr)
	if err != nil {
		return nil, err
	}

	// log.Println("setup pins")
	if err = gpioDev.SetDirection(PinRST, gpio.Out); err != nil {
		return nil, err
	}
	if err = gpioDev.SetDirection(PinDC, gpio.Out); err != nil {
		return nil, err
	}
	if err = gpioDev.SetDirection(PinBusy, gpio.In); err != nil {
		return nil, err
	}

	disp := &Display{
		spiDev:  spiDev,
		gpioDev: gpioDev,
		w:       epd2in13bcWidth,
		h:       epd2in13bcHeight,
	}

	disp.Init()

	return disp, nil
}

// Close closes all devices in Display
func (d *Display) Close() {
	d.Sleep()

	d.gpioDev.SetValue(PinRST, 0)
	d.gpioDev.SetValue(PinDC, 0)

	if d.spiDev != nil {
		d.spiDev.Close()
	}

	if d.gpioDev != nil {
		d.gpioDev.Close()
	}
}

// Reset does H/W reset if pinRst is not nil
func (d *Display) Reset() error {
	if d.gpioDev == nil {
		return fmt.Errorf("epd2in13bc: no gpio device. skip Reset")
	}

	d.gpioDev.SetValue(PinRST, 1)
	time.Sleep(200 * time.Millisecond)
	d.gpioDev.SetValue(PinRST, 0)
	time.Sleep(10 * time.Millisecond)
	d.gpioDev.SetValue(PinRST, 1)
	time.Sleep(200 * time.Millisecond)

	return nil
}

// Clear clears display with given color(8 bits)
func (d *Display) Clear(b, ry byte) {
	lineWidth := (d.w + 7) / 8
	buf := make([]byte, lineWidth*d.h)
	for i := range buf {
		buf[i] = b
	}
	d.sendCmd(0x10)
	// d.sendDatas(buf)
	for _, v := range buf {
		d.sendData(v)
	}
	d.sendCmd(0x92)

	for i := range buf {
		buf[i] = ry
	}
	d.sendCmd(0x13)
	// d.sendDatas(buf)
	for _, v := range buf {
		d.sendData(v)
	}
	d.sendCmd(0x92)

	d.sendCmd(0x12)
	d.waitTillNotBusy()
}

// Sleep makes display sleep
func (d *Display) Sleep() {
	d.sendCmd(0x02) // POWER_OFF
	d.waitTillNotBusy()
	d.sendCmd(0x07)  // DEEP_SLEEP
	d.sendData(0xA5) // check code
}

func (d *Display) sendCmd(c byte) (err error) {
	if err = d.gpioDev.SetValue(PinDC, 0); err != nil {
		log.Fatal("fail to set gpio value", err)
		return
	}
	time.Sleep(100 * time.Microsecond) // TODO: need it?
	if err = d.spiDev.Tx([]byte{c}, nil); err != nil {
		log.Fatal("fail to send cmd", err)
		return
	}
	return
}

func (d *Display) sendData(b byte) (err error) {
	if err = d.gpioDev.SetValue(PinDC, 1); err != nil {
		log.Fatal("fail to set gpio value", err)
		return
	}
	time.Sleep(100 * time.Microsecond) // TODO: need it?
	if err = d.spiDev.Tx([]byte{b}, nil); err != nil {
		log.Fatal("fail to send data", err)
		return
	}
	return
}

func (d *Display) sendDatas(bs []byte) (err error) {
	if err = d.gpioDev.SetValue(PinDC, 1); err != nil {
		log.Fatal("fail to set gpio value", err)
		return
	}
	time.Sleep(100 * time.Microsecond) // TODO: need it?
	if err = d.spiDev.Tx(bs, nil); err != nil {
		log.Fatal("fail to send datas", err)
		return
	}
	return
}

func (d *Display) waitTillNotBusy() {
	var v int // 1: idle, 0: busy
	var err error
	time.Sleep(1 * time.Second)
	for {
		if v, err = d.gpioDev.Value(PinBusy); err == nil && v == 1 {
			break
		}
		if err != nil {
			panic(err)
		}
		log.Println("busy")
		time.Sleep(100 * time.Millisecond)
	}
	log.Println("idle")
}
