package epd7in5

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"time"

	"github.com/suapapa/go_devices/epdcolor"
	"periph.io/x/conn/v3"
	"periph.io/x/conn/v3/display"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/conn/v3/spi"
	"periph.io/x/host/v3/rpi"
)

const (
	w = 640
	h = 384
)

// Dev is open handle to display controller.
type Dev struct {
	c           conn.Conn
	dc, cs, rst gpio.PinOut
	busy        gpio.PinIO

	rect image.Rectangle
}

// NewSPI returns a Dev object that communicates over SPI to epd7in5 E-Paper display controller.
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
	return fmt.Sprintf("epd7in5.Dev{%s, %s, %s}", d.c, d.dc, d.rect.Max)
}

// ColorModel implements display.Drawer
func (d *Dev) ColorModel() color.Model {
	return epdcolor.Gray3Model
}

// Bounds implements display.Drawer
func (d *Dev) Bounds() image.Rectangle {
	return d.rect
}

// Draw implements display.Drawer
func (d *Dev) Draw(r image.Rectangle, src image.Image, sp image.Point) error {
	var buff []byte
	if img, ok := src.(*epdcolor.Gray3Image); ok && r == d.rect && img.Rect == d.rect && sp.X == 0 && sp.Y == 0 {
		buff = img.Pix
	} else {
		grayImg := epdcolor.NewGray3Image(d.rect)
		buff = grayImg.Pix
		draw.Src.Draw(grayImg, r, src, sp)
	}
	return d.drawInternal(buff)
}

func (d *Dev) drawInternal(b []byte) error {
	db := make([]byte, 0)
	for i := 0; i < w/4*h; i++ {
		tmp1 := b[i]
		var tmp2 byte
		j := 0
		for j < 4 {
			if tmp1&0xC0 == 0xC0 {
				tmp2 = 0x03
			} else if tmp1&0xC0 == 0x00 {
				tmp2 = 0x00
			} else {
				tmp2 = 0x04
			}
			tmp2 = (tmp2 << 4) & 0xFF
			tmp1 = (tmp1 << 2) & 0xFF
			j++
			if tmp1&0xC0 == 0xC0 {
				tmp2 |= 0x03
			} else if tmp1&0xC0 == 0x00 {
				tmp2 |= 0x00
			} else {
				tmp2 |= 0x04
			}
			tmp1 = (tmp1 << 2) & 0xFF
			// d.sendData(tmp2)
			db = append(db, tmp2)
			j++
		}
	}
	// log.Println("making display buffer done db len =", len(db))

	d.sendCmd([]byte{0x10})
	d.sendData(db)
	d.sendCmd([]byte{0x12})
	d.waitUntilIdle()

	return nil
}

// Halt turns off the display (clean up)
func (d *Dev) Halt() error {
	img := epdcolor.NewGray3Image(d.rect)
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
	d.sendCmd([]byte{0x01}) // POWER_SETTING
	d.sendData([]byte{0x37, 0x00})

	d.sendCmd([]byte{0x00}) // PANEL_SETTING
	d.sendData([]byte{0xCF, 0x08})

	d.sendCmd([]byte{0x06}) // BOOSTER_SOFT_START
	d.sendData([]byte{0xc7, 0xcc, 0x28})

	d.sendCmd([]byte{0x04}) // POWER_ON
	d.waitUntilIdle()

	d.sendCmd([]byte{0x30}) // PLL_CONTROL
	d.sendData([]byte{0x3c})

	d.sendCmd([]byte{0x41}) // TEMPERATURE_CALIBRATION
	d.sendData([]byte{0x00})

	d.sendCmd([]byte{0x50}) // VCOM_AND_DATA_INTERVAL_SETTING
	d.sendData([]byte{0x77})

	d.sendCmd([]byte{0x60}) // TCON_SETTING
	d.sendData([]byte{0x22})

	d.sendCmd([]byte{0x61})                          // TCON_RESOLUTION
	d.sendData([]byte{byte(w >> 8), byte(w & 0xff)}) //source 640
	d.sendData([]byte{byte(h >> 8), byte(h & 0xff)}) //gate 384

	d.sendCmd([]byte{0x82})  // VCM_DC_SETTING
	d.sendData([]byte{0x1E}) // decide by LUT file

	d.sendCmd([]byte{0xe5}) // FLASH MODE
	d.sendData([]byte{0x03})

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
