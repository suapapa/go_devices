package max7219

import (
	"fmt"

	"github.com/pkg/errors"
	"periph.io/x/conn/v3"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/conn/v3/spi"
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

// New returns instance of MAX7219 which is connected in spi port, p
func New(p spi.Port) (*Dev, error) {
	c, err := p.Connect(10*physic.MegaHertz, spi.Mode0, 8)
	if err != nil {
		return nil, err
	}

	d := &Dev{
		c: c,
	}
	if err := d.init(); err != nil {
		return nil, err
	}
	return d, nil
}

func (d *Dev) String() string {
	return fmt.Sprintf("max7219.Dev{%s}", d.c)
}

// Shutdown takes display out of shutdown
func (d *Dev) Shutdown(shutdown bool) error {
	var v byte
	if !shutdown {
		v = 1
	}
	return d.Write(REG_SHUTDOWN, v)
}

// DisplayTest start and stop a display test
func (d *Dev) DisplayTest(onoff bool) error {
	var v byte
	if onoff {
		v = 1
	}
	return d.Write(REG_DISPLAY_TEST, v)
}

// Clear clears the display
func (d *Dev) Clear() error {
	for i := 0; i < 8; i++ {
		if err := d.Write(byte(i+1), 0x00); err != nil {
			return errors.Wrap(err, "fail to clear")
		}
	}
	return nil
}

// SetBrightness sets the LED display brightness
func (d *Dev) SetBrightness(intensity byte) error {
	intensity &= 0x00f
	return d.Write(REG_INTENSITY, intensity)
}

// Write writes a data to reg
func (d *Dev) Write(reg, data byte) error {
	buf := []byte{reg, data}
	return d.c.Tx(buf, nil)
}

// WriteChar write a ASCII character to a FND display
func (d *Dev) WriteChar(reg, char byte) error {
	buf := []byte{reg, d.lookupCode(char)}
	return d.c.Tx(buf, nil)
}

/*
*********************************************************************************************************
* LED Segments:         a
*                     ----
*                   f|    |b
*                    |  g |
*                     ----
*                   e|    |c
*                    |    |
*                     ----  o dp
*                       d
*   Register bits:
*      bit:  7  6  5  4  3  2  1  0
*           dp  a  b  c  d  e  f  g
*********************************************************************************************************
 */
func (d *Dev) lookupCode(ch byte) byte {
	font := map[byte]byte{
		// {' ': 0x00,
		'0': 0x7e,
		'1': 0x30,
		'2': 0x6d,
		'3': 0x79,
		'4': 0x33,
		'5': 0x5b,
		'6': 0x5f,
		'7': 0x70,
		'8': 0x7f,
		'9': 0x7b,
		'A': 0x77,
		'B': 0x1f,
		'C': 0x4e,
		'D': 0x3d,
		'E': 0x4f,
		'F': 0x47,
	}

	v, ok := font[ch]
	if !ok {
		return 0
	}
	return v
}

func (d *Dev) init() error {
	{
		if err := d.Write(REG_SCAN_LIMIT, 7); err != nil { // set up to scan all eight digits
			return errors.Wrap(err, "fail to init")
		}
		if err := d.Write(REG_DECODE, 0x00); err != nil { // set to "no decode" for all digits
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
