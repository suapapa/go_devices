package epd7in5v2

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"time"

	"github.com/suapapa/go_devices/epdcolor"
	"periph.io/x/periph/conn"
	"periph.io/x/periph/conn/display"
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/physic"
	"periph.io/x/periph/conn/spi"
	"periph.io/x/periph/host/rpi"
)

const (
	w = 800
	h = 480
)

// Dev is open handle to display controller.
type Dev struct {
	c           conn.Conn
	dc, cs, rst gpio.PinOut
	busy        gpio.PinIO

	rect image.Rectangle
}

// NewSPI returns a Dev object that communicates over SPI to epd7in5v2 E-Paper display controller.
func NewSPI(p spi.Port, dc, rst gpio.PinOut, busy gpio.PinIO) (*Dev, error) {
	if err := dc.Out(gpio.Low); err != nil {
		return nil, err
	}

	c, err := p.Connect(5*physic.MegaHertz, spi.Mode0, 8)
	if err != nil {
		return nil, err
	}

	d := &Dev{
		c:    c,
		dc:   dc,
		rst:  rst,
		busy: busy,
		rect: image.Rect(0, 0, w, h),
	}
	d.Reset()
	if err := d.Init(); err != nil {
		return nil, err
	}
	return d, nil
}

// NewSPIHat returns a Dev object that communicates over SPI
// and have the default config for the e-paper hat for raspberry pi
func NewSPIHat(p spi.Port) (*Dev, error) {
	dc := rpi.P1_22   // 25
	rst := rpi.P1_11  // 17
	busy := rpi.P1_18 // 24
	return NewSPI(p, dc, rst, busy)
}

func (d *Dev) String() string {
	return fmt.Sprintf("epd7in5v2.Dev{%s, %s, %s}", d.c, d.dc, d.rect.Max)
}

// ColorModel implements display.Drawer
func (d *Dev) ColorModel() color.Model {
	return epdcolor.WBModel
}

// Bounds implements display.Drawer
func (d *Dev) Bounds() image.Rectangle {
	return d.rect
}

// Draw implements display.Drawer
func (d *Dev) Draw(r image.Rectangle, src image.Image, sp image.Point) error {
	var buff []byte
	if img, ok := src.(*epdcolor.WBImage); ok && r == d.rect && img.Rect == d.rect && sp.X == 0 && sp.Y == 0 {
		buff = img.Pix
	} else {
		bwImg := epdcolor.NewWBImage(d.rect)
		buff = bwImg.Pix
		draw.Src.Draw(bwImg, r, src, sp)
	}
	return d.drawInternal(buff)
}

func (d *Dev) drawInternal(b []byte) error {
	d.sendCmd([]byte{0x13})
	d.sendData(b)
	d.sendCmd([]byte{0x12})
	time.Sleep(100 * time.Millisecond)
	d.waitUntilIdle()

	return nil
}

// Halt turns off the display (clean up)
func (d *Dev) Halt() error {
	img := epdcolor.NewWBImage(d.rect)
	return d.drawInternal(img.Pix)
}

// Reset can be also used to awaken the device
func (d *Dev) Reset() {
	d.rst.Out(gpio.Low)
	time.Sleep(200 * time.Millisecond)
	d.rst.Out(gpio.High)
	time.Sleep(200 * time.Millisecond)
}

// Sleep after this command is transmitted, the chip would enter the
// deep-sleep mode to save power.
//
// The deep sleep mode would return to standby by hardware reset.
// You can use Reset() to awaken and Init to re-initialize the device.
func (d *Dev) Sleep() {
	d.sendCmd([]byte{0x02}) // power off
	d.waitUntilIdle()
	d.sendCmd([]byte{0x07}) // deep sleep
	d.sendData([]byte{0xA5})
}

// Init initialize the display config. This method is already called when creating
// a device using NewSPI and NewSPIHat methods.
//
// It should be only used when you put the device to sleep and need to re-init the device.
func (d *Dev) Init() error {
	d.sendCmd([]byte{0x01}) // POWER SETTING
	d.sendData([]byte{0x07})
	d.sendData([]byte{0x07}) // VGH=20V,VGL=-20V
	d.sendData([]byte{0x3f}) // VDH=15V
	d.sendData([]byte{0x3f}) // VDL=-15V

	d.sendCmd([]byte{0x04}) // POWER ON
	time.Sleep(100 * time.Millisecond)
	d.waitUntilIdle()

	d.sendCmd([]byte{0x00})  // PANNEL SETTING
	d.sendData([]byte{0x1F}) // KW-3f   KWR-2F  WBROTP 0f       WBOTP 1f

	d.sendCmd([]byte{0x61})  // tres
	d.sendData([]byte{0x03}) // source 800
	d.sendData([]byte{0x20})
	d.sendData([]byte{0x01}) // gate 480
	d.sendData([]byte{0xE0})

	d.sendCmd([]byte{0x15})
	d.sendData([]byte{0x00})

	d.sendCmd([]byte{0x50}) // VCOM AND DATA INTERVAL SETTING
	d.sendData([]byte{0x10})
	d.sendData([]byte{0x07})

	d.sendCmd([]byte{0x60}) // TCON SETTING
	d.sendData([]byte{0x22})

	return nil
}

func (d *Dev) waitUntilIdle() {
	for d.busy.Read() == gpio.Low {
		time.Sleep(100 * time.Millisecond)
	}
}

func (d *Dev) sendData(c []byte) error {
	if err := d.dc.Out(gpio.High); err != nil {
		return err
	}
	return d.c.Tx(c, nil)
}

func (d *Dev) sendCmd(c []byte) error {
	if err := d.dc.Out(gpio.Low); err != nil {
		return err
	}
	return d.c.Tx(c, nil)
}

var _ display.Drawer = &Dev{}
