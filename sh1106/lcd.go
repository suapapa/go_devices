package sh1106

import (
	"github.com/davecheney/gpio"
	"golang.org/x/exp/io/i2c"
	i2c_driver "golang.org/x/exp/io/i2c/driver"
	"golang.org/x/exp/io/spi"
	spi_driver "golang.org/x/exp/io/spi/driver"
)

type LCD struct {
	i2cDev *i2c.Device

	spiDev *spi.Device
	pinDC  gpio.Pin

	w, h uint
	buff []byte
}

func OpenI2C(o i2c_driver.Opener, addr int) (*LCD, error) {
	dev, err := i2c.Open(o, addr)
	if err != nil {
		return nil, err
	}

	display := &LCD{i2cDev: dev}

	// TODO: support not only 128x64
	display.w = sh1106_LCDWIDTH
	display.h = sh1106_LCDHEIGHT
	display.init()

	return display, nil
}

func OpenSpi(o spi_driver.Opener, dc gpio.Pin) (*LCD, error) {
	dc.SetMode(gpio.ModeInput)
	dc.SetMode(gpio.ModeOutput)

	dev, err := spi.Open(o)
	if err != nil {
		return nil, err
	}
	dev.SetCSChange(false)

	display := &LCD{spiDev: dev, pinDC: dc}

	// TODO: support not only 128x64
	display.w = sh1106_LCDWIDTH
	display.h = sh1106_LCDHEIGHT
	display.init()

	return display, nil
}

func (l *LCD) Close() {
	if l.i2cDev != nil {
		l.i2cDev.Close()
	}

	if l.spiDev != nil {
		l.spiDev.Close()
		l.pinDC.Close()
	}
}

func (l *LCD) DrawPixel(x, y uint, p bool) {
	if x >= l.w || y >= l.h {
		return
	}

	if p { // BLACK
		l.buff[x+(y/8)*l.w] |= byte(1 << (y & 7))
	} else { // WHITE
		l.buff[x+(y/8)*l.w] &= byte(^(1 << (y & 7)))
	}
}

func (l *LCD) Display() {
	l.sendCmd(sh1106_SETLOWCOLUMN | 0x0)  // low col = 0
	l.sendCmd(sh1106_SETHIGHCOLUMN | 0x0) // hi col = 0
	l.sendCmd(sh1106_SETSTARTLINE | 0x0)  // line #0

	// height := byte(l.h) >> 3 // 64 >> 3
	width := byte(l.w) >> 3 // 132 >> 3

	mRow := byte(0)
	mCol := byte(2)

	var r byte
	for k := 0; k < len(l.buff); k += int(width) {
		// send a bunch of data in one xmission
		l.sendCmd(0xB0 + r + mRow)    //set page address
		l.sendCmd(mCol & 0xf)         //set lower column address
		l.sendCmd(0x10 | (mCol >> 4)) //set higher column address

		l.sendData(l.buff[k : k+int(width)])
		r++
	}
}

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

	l.buff = make([]byte, l.w*l.h/8)
	l.init128x64()
}

func (l *LCD) sendCmd(c byte) {
	if l.i2cDev != nil {
		l.i2cDev.Write([]byte{0x00, c})
	} else {
		l.pinDC.Clear()
		l.spiDev.Tx([]byte{c}, nil)
	}
}

func (l *LCD) sendData(d []byte) {
	if l.i2cDev != nil {
		l.i2cDev.Write([]byte{0x40})
		l.i2cDev.Write(d)
	} else {
		l.pinDC.Set()
		l.spiDev.Tx(d, nil)
	}
}
