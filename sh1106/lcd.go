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

	w, h int
}

func OpenI2C(o i2c_driver.Opener, addr int) (*LCD, error) {
	dev, err := i2c.Open(o, addr)
	if err != nil {
		return nil, err
	}

	display := &LCD{i2cDev: dev}
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

func (l *LCD) init() {
	panic("not implemented")
}

func (l *LCD) sendCmd(c byte) {
	panic("not implemented")
}

func (l *LCD) sendData(d byte) {
	panic("not implemented")
}
