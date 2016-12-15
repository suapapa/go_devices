// Copyright 2016, Homin Lee <homin.lee@suapapa.net>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sh1106 // import "github.com/suapapa/go_devices/sh1106"

import (
	"fmt"
	"time"

	"github.com/goiot/exp/gpio"
	gpio_driver "github.com/goiot/exp/gpio/driver"
	"golang.org/x/exp/io/i2c"
	i2c_driver "golang.org/x/exp/io/i2c/driver"
	"golang.org/x/exp/io/spi"
	spi_driver "golang.org/x/exp/io/spi/driver"
)

// LCD represents a shll06 drived OLED display
type LCD struct {
	i2cDev  *i2c.Device
	spiDev  *spi.Device
	gpioDev *gpio.Device

	w, h int
	buff []byte
}

// Open opens a sh1106 LCD in I2C or SPI mode
// func Open(o interface{}, addr int, rst, dc *gpio.Pin) (*LCD, error) {
// 	if i2cO, ok := o.(i2c_driver.Opener); ok {
// 		return OpenI2C(i2cO, addr, rst)
// 	} else if spiO, ok := o.(spi_driver.Opener); ok {
// 		return OpenSpi(spiO, dc, rst)
// 	}
// 	return nil, errors.New("sh1106: unknown driver.Opener")
// }

// OpenI2C opens a sh1106 LCD in I2C mode
// gpio device should have PinRST pin for Reset()
// can pass nil for ctr if ignore Reset()
func OpenI2C(bus i2c_driver.Opener, ctr gpio_driver.Opener, addr int) (*LCD, error) {
	i2cDev, err := i2c.Open(bus, addr)
	if err != nil {
		return nil, err
	}

	var gpioDev *gpio.Device
	if ctr != nil {
		if gpioDev, err = gpio.Open(ctr); err != nil {
			i2cDev.Close()
			return nil, err
		}

		if err = gpioDev.SetDirection(PinRST, gpio.Out); err != nil {
			i2cDev.Close()
			return nil, err
		}
	}

	lcd := &LCD{
		i2cDev:  i2cDev,
		gpioDev: gpioDev,
		// TODO: support not only 128x64
		w: sh1106_LCDWIDTH,
		h: sh1106_LCDHEIGHT,
	}

	lcd.Reset()
	lcd.init()

	return lcd, nil
}

// OpenSpi opens a sh1106 LCD in SPI mode
// gpio device should have PinRST pin for Reset() and
// PinDC pin for selecting data/cmd
func OpenSpi(bus spi_driver.Opener, ctr gpio_driver.Opener) (*LCD, error) {
	spiDev, err := spi.Open(bus)
	if err != nil {
		return nil, err
	}
	spiDev.SetCSChange(false)

	gpioDev, err := gpio.Open(ctr)
	if err != nil {
		return nil, err
	}

	if err = gpioDev.SetDirection(PinDC, gpio.Out); err != nil {
		return nil, err
	}

	if err = gpioDev.SetDirection(PinRST, gpio.Out); err != nil {
		return nil, err
	}

	lcd := &LCD{
		spiDev:  spiDev,
		gpioDev: gpioDev,
		// TODO: support not only 128x64
		w: sh1106_LCDWIDTH,
		h: sh1106_LCDHEIGHT,
	}

	lcd.Reset()
	lcd.init()

	return lcd, nil
}

// Close closes all devices in LCD
func (l *LCD) Close() {
	if l.i2cDev != nil {
		l.i2cDev.Close()
	}

	if l.spiDev != nil {
		l.spiDev.Close()
	}

	if l.gpioDev != nil {
		l.gpioDev.Close()
	}
}

// Reset does H/W reset if pinRst is not nil
func (l *LCD) Reset() error {
	if l.gpioDev == nil {
		return fmt.Errorf("sh1106: no gpio device. skip Reset")
	}

	l.gpioDev.SetValue(PinRST, 1)
	time.Sleep(1 * time.Millisecond)
	l.gpioDev.SetValue(PinRST, 0)
	time.Sleep(10 * time.Millisecond)
	l.gpioDev.SetValue(PinRST, 1)
	return nil
}

// Clear clean internal buffer
func (l *LCD) Clear() {
	for i := range l.buff {
		l.buff[i] = 0x00
	}
}

// DrawPixel sets a dot to the internal buffer
func (l *LCD) DrawPixel(x, y int, p bool) {
	if x >= l.w || y >= l.h {
		return
	}

	if p { // white
		l.buff[x+(y/8)*l.w] |= byte(1 << (uint(y) & 7))
	} else { // black
		l.buff[x+(y/8)*l.w] &= byte(^(1 << (uint(y) & 7)))
	}
}

// Display sends what buffer contents to LCD and turn it on
func (l *LCD) Display() {
	l.sendCmd(sh1106_SETLOWCOLUMN | 0x0)  // low col = 0
	l.sendCmd(sh1106_SETHIGHCOLUMN | 0x0) // hi col = 0
	l.sendCmd(sh1106_SETSTARTLINE | 0x0)  // line #0

	height := byte(l.h) >> 3 // 64 >> 3 = 8
	width := byte(l.w) >> 3  // 132 >> 3 = 16

	row := byte(0)
	col := byte(2)

	k := 0
	for i := byte(0); i < height; i++ {
		l.sendCmd(0xB0 + i + row)    //set page address
		l.sendCmd(col & 0xf)         //set lower column address
		l.sendCmd(0x10 | (col >> 4)) //set higher column address

		for j := byte(0); j < 8; j++ {
			l.sendData(l.buff[k : k+int(width)])
			k += int(width)
		}
	}
}

// Invert flips black and white on the LCD
func (l *LCD) Invert(i bool) {
	if i {
		l.sendCmd(sh1106_INVERTDISPLAY)
	} else {
		l.sendCmd(sh1106_NORMALDISPLAY)
	}
}

func (l *LCD) init() {
	if l.w != 128 && l.h != 64 {
		panic("not implemented")
	}

	l.buff = make([]byte, (l.w*l.h+7)/8)
	l.init128x64()
	l.display(true)
}

func (l *LCD) display(on bool) {
	if on {
		l.sendCmd(0x8D)
		l.sendCmd(0x14)
		l.sendCmd(0xAF)
	} else {
		l.sendCmd(0x8D)
		l.sendCmd(0x10)
		l.sendCmd(0xAE)
	}
}

func (l *LCD) sendCmd(c byte) (err error) {
	if l.i2cDev != nil {
		err = l.i2cDev.Write([]byte{0x00, c})
	} else {
		if err = l.gpioDev.SetValue(PinDC, 0); err != nil {
			return
		}
		err = l.spiDev.Tx([]byte{c}, nil)
	}
	return
}

func (l *LCD) sendData(d []byte) (err error) {
	if l.i2cDev != nil {
		err = l.i2cDev.Write(append([]byte{0x40}, d...))
	} else {
		if err = l.gpioDev.SetValue(PinDC, 1); err != nil {
			return
		}
		err = l.spiDev.Tx(d, nil)
	}
	return
}
