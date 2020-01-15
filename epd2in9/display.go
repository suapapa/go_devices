// Copyright 2019, Homin Lee <homin.lee@suapapa.net>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package epd2in9 // import "github.com/suapapa/go_devices/epd2in9"

import (
	"fmt"
	"time"

	"github.com/goiot/exp/gpio"
	gpio_driver "github.com/goiot/exp/gpio/driver"
	"golang.org/x/exp/io/spi"
	spi_driver "golang.org/x/exp/io/spi/driver"
)

// Display represents a epd2in9 e-ink display
type Display struct {
	spiDev  *spi.Device // recommend 4000000Hz, spimode 0
	gpioDev *gpio.Device

	w, h int
}

// Open opens a epd2in9 display in SPI mode
// gpio device should have PinRST pin for Reset() and
// PinDC pin for selecting data/cmd
func Open(bus spi_driver.Opener, ctr gpio_driver.Opener) (*Display, error) {
	spiDev, err := spi.Open(bus)
	if err != nil {
		return nil, err
	}
	spiDev.SetCSChange(false)
	spiDev.SetMode(spi.Mode0)

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
	// if err = gpioDev.SetDirection(PinCS, gpio.Out); err != nil {
	// 	return nil, err
	// }
	if err = gpioDev.SetDirection(PinBusy, gpio.In); err != nil {
		return nil, err
	}

	disp := &Display{
		spiDev:  spiDev,
		gpioDev: gpioDev,
		w:       epd2in9Width,
		h:       epd2in9Height,
	}

	disp.InitFull()

	return disp, nil
}

// Close closes all devices in Display
func (d *Display) Close() {
	d.Sleep()

	// d.gpioDev.SetValue(PinRST, 1)
	// d.gpioDev.SetValue(PinDC, 0)

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
		return fmt.Errorf("epd2in9: no gpio device. skip Reset")
	}

	d.gpioDev.SetValue(PinRST, 1)
	time.Sleep(200 * time.Millisecond)
	d.gpioDev.SetValue(PinRST, 0)
	time.Sleep(10 * time.Millisecond)
	d.gpioDev.SetValue(PinRST, 1)
	time.Sleep(200 * time.Millisecond)

	return nil
}

// Sleep makes display sleep
func (d *Display) Sleep() {
	d.sendCmd(0x10) // DEEP_SLEEP_MODE
	d.sendData(0x01)
}

// turnOn turns on full screen
func (d *Display) turnOn() {
	d.sendCmd(0x22)
	d.sendData(0xC4)
	d.sendCmd(0x20)
	d.sendCmd(0xFF)

	d.waitTillNotBusy()
}

func (d *Display) sendCmd(c byte) (err error) {
	if err = d.gpioDev.SetValue(PinDC, 0); err != nil {
		return
	}
	// time.Sleep(10 * time.Millisecond)
	if err = d.spiDev.Tx([]byte{c}, nil); err != nil {
		return
	}
	return
}

func (d *Display) sendData(b byte) (err error) {
	if err = d.gpioDev.SetValue(PinDC, 1); err != nil {
		return
	}
	// time.Sleep(10 * time.Millisecond)
	if err = d.spiDev.Tx([]byte{b}, nil); err != nil {
		return
	}
	return
}

func (d *Display) sendDatas(bs []byte) (err error) {
	if err = d.gpioDev.SetValue(PinDC, 1); err != nil {
		return
	}
	// time.Sleep(10 * time.Millisecond)
	if err = d.spiDev.Tx(bs, nil); err != nil {
		return
	}
	return
}

func (d *Display) waitTillNotBusy() {
	var v int // 0: idle, 1: busy
	var err error
	for {
		if v, err = d.gpioDev.Value(PinBusy); err == nil && v == 0 {
			// log.Println("idle")
			return
		}
		if err != nil {
			panic(err)
		}
		// log.Println("busy")
		time.Sleep(200 * time.Millisecond)
	}
}
