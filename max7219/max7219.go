package max7219

import (
	"fmt"

	"github.com/pkg/errors"
	"periph.io/x/periph/conn"
	"periph.io/x/periph/conn/physic"
	"periph.io/x/periph/conn/spi"
)

const (
	REG_DECODE       = 0x09 // "decode mode" register
	REG_INTENSITY    = 0x0a // "intensity" register
	REG_SCAN_LIMIT   = 0x0b // "scan limit" register
	REG_SHUTDOWN     = 0x0c // "shutdown" register
	REG_DISPLAY_TEST = 0x0f // "display test" register

	INTENSITY_MIN = 0x00 // minimum display intensity
	INTENSITY_MAX = 0x0f // maximum display intensity
)

type Dev struct {
	c conn.Conn
}

// NewSPI returns instance of MAX7219 which is connected in spi port, p
func NewSPI(p spi.Port) (*Dev, error) {
	c, err := p.Connect(5*physic.MegaHertz, spi.Mode0, 8)
	if err != nil {
		return nil, err
	}

	d := &Dev{
		c: c,
	}
	// d.Shutdown(false)
	// if err := d.init(); err != nil {
	// 	return nil, err
	// }
	return d, nil
}

func (d *Dev) init() error {
	{
		if err := d.write(REG_SCAN_LIMIT, 7); err != nil { // set up to scan all eight digits
			return errors.Wrap(err, "fail to init")
		}
		if err := d.write(REG_DECODE, 0x00); err != nil { // set to "no decode" for all digits
			return errors.Wrap(err, "fail to init")
		}
		if err := d.Shutdown(false); err != nil { // select normal operation (i.e. not shutdown)
			return errors.Wrap(err, "fail to init")
		}
		if err := d.DisplayTest(false); err != nil { // select normal operation (i.e. not test mode)
			return errors.Wrap(err, "fail to init")
		}
		if err := d.Clear(); err != nil { // clear all digits
			return errors.Wrap(err, "fail to init")
		}
		if err := d.SetBrightness(INTENSITY_MAX); err != nil { // set to maximum intensity
			return errors.Wrap(err, "fail to init")
		}
	}
	return nil
}

func (d *Dev) Shutdown(shutdown bool) error {
	// TODO: TBD
	return nil
}

func (d *Dev) DisplayTest(onoff bool) error {
	// TODO: TBD
	return nil
}

func (d *Dev) Clear() error {
	// TODO: TBD
	return nil
}

func (d *Dev) SetBrightness(intensity byte) error {
	intensity = intensity & 0x0f
	// TODO: TBD
	return nil
}

func (d *Dev) String() string {
	return fmt.Sprintf("max7219.Dev{%s}", d.c)
}

// ---

func (d *Dev) write(reg, data byte) error {
	// TODO: TBD
	return nil
}

func (d *Dev) sendByte(data byte) error {
	// TODO: TBD
	return nil
}

func (d *Dev) lookupCode(ch byte) byte {
	// TODO: TBD
	return 0x00
}
