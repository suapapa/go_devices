package epd7in5

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"time"

	"github.com/suapapa/go_devices/epd7in5/gray3"
	"periph.io/x/periph/conn"
	"periph.io/x/periph/conn/display"
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/physic"
	"periph.io/x/periph/conn/spi"
	"periph.io/x/periph/host/rpi"
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
func NewSPI(p spi.Port, dc, cs, rst gpio.PinOut, busy gpio.PinIO) (*Dev, error) {
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
		cs:   cs,
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
	cs := rpi.P1_24   // 8
	rst := rpi.P1_11  // 17
	busy := rpi.P1_18 // 24
	return NewSPI(p, dc, cs, rst, busy)
}

func (d *Dev) String() string {
	return fmt.Sprintf("epd7in5.Dev{%s, %s, %s}", d.c, d.dc, d.rect.Max)
}

// ColorModel implements display.Drawer
func (d *Dev) ColorModel() color.Model {
	return gray3.Gray3Model
}

// Bounds implements display.Drawer
func (d *Dev) Bounds() image.Rectangle {
	return d.rect
}

// // Draw implements display.Drawer
// func (d *Dev) Draw(r image.Rectangle, src image.Image, sp image.Point) error {
// 	if sp.X != 0 || sp.Y != 0 {
// 		return fmt.Errorf("sp should start from 0,0")
// 	}
// 	if r.Dx() != w || r.Dy() != h {
// 		return fmt.Errorf("rect shold be %dx%d", w, h)
// 	}

// 	buff, err := d.Image2Buffer(src)
// 	if err != nil {
// 		return err
// 	}
// 	return d.DrawBuffer(buff)
// }

// Draw implements display.Drawer
func (d *Dev) Draw(r image.Rectangle, src image.Image, sp image.Point) error {
	var buff []byte
	if img, ok := src.(*gray3.Image); ok && r == d.rect && img.Rect == d.rect && sp.X == 0 && sp.Y == 0 {
		buff = img.Pix
	} else {
		grayImg := gray3.NewImage(d.rect)
		buff = grayImg.Pix
		draw.Src.Draw(grayImg, r, src, sp)
	}
	return d.drawInternal(buff)
}

func (d *Dev) drawInternal(b []byte) error {
	// log.Println("DrawBuffer start")
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

	d.sendCmd(0x10)
	// TODO: need this?
	for i := 0; i < len(db); i += 4096 {
		d.sendDatas(db[i : i+4096])
	}
	d.sendCmd(0x12)
	d.waitUntilIdle()
	// log.Println("DrawBuffer end")
	return nil
}

// Image2Buffer returns monochrome image buffer from image.Image
func (d *Dev) Image2Buffer(img image.Image) ([]byte, error) {
	// log.Println("Image2Buffer")
	b := make([]byte, w*h/4)

	imgW, imgH := img.Bounds().Dx(), img.Bounds().Dy()
	if imgW == w && imgH == h {
		for y := 0; y < imgH; y++ {
			for x := 0; x < imgW; x++ {
				switch checkColor(img.At(x, y)) {
				case black:
					b[(x+y*w)/4] &= ^(0xC0 >> (x % 4 * 2))
				case gray:
					b[(x+y*w)/4] &= ^(0xC0 >> (x % 4 * 2))
					b[(x+y*w)/4] |= (0x40 >> (x % 4 * 2))
				case white:
					b[(x+y*w)/4] |= (0xC0 >> (x % 4 * 2))
				}
			}
		}
	} else if imgW == h && imgH == w {
		for y := 0; y < imgH; y++ {
			for x := 0; x < imgW; x++ {
				nx := y
				ny := h - x - 1
				switch checkColor(img.At(x, y)) {
				case black:
					b[(nx+ny*w)/4] &= ^(0xC0 >> (y % 4 * 2))
				case gray:
					b[(nx+ny*w)/4] &= ^(0xC0 >> (y % 4 * 2))
					b[(nx+ny*w)/4] |= (0x40 >> (y % 4 * 2))
				case white:
					b[(nx+ny*w)/4] |= (0xC0 >> (y % 4 * 2))
				}
			}
		}
	} else {
		return nil, fmt.Errorf("image size should be %dx%d of %dx%d", w, h, h, w)
	}

	return b, nil
}

// Halt turns off the display (clean up)
func (d *Dev) Halt() error {
	// db := make([]byte, w*h/2)
	// for i := 0; i < len(db); i++ {
	// 	db[i] = 0x33
	// }
	// d.sendCmd(0x10)
	// for i := 0; i < len(db); i += 4096 {
	// 	d.sendDatas(db[i : i+4096])
	// }
	// d.sendCmd(0x12)
	// d.waitUntilIdle()
	img := gray3.NewImage(d.rect)
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
	d.sendCmd(0x02) // power off
	d.waitUntilIdle()
	d.sendCmd(0x07) // deep sleep
	d.sendData(0xA5)
}

// Init initialize the display config. This method is already called when creating
// a device using NewSPI and NewSPIHat methods.
//
// It should be only used when you put the device to sleep and need to re-init the device.
func (d *Dev) Init() error {
	d.sendCmd(0x01) // POWER_SETTING
	d.sendData(0x37)
	d.sendData(0x00)

	d.sendCmd(0x00) // PANEL_SETTING
	d.sendData(0xCF)
	d.sendData(0x08)

	d.sendCmd(0x06) // BOOSTER_SOFT_START
	d.sendData(0xc7)
	d.sendData(0xcc)
	d.sendData(0x28)

	d.sendCmd(0x04) // POWER_ON
	d.waitUntilIdle()

	d.sendCmd(0x30) // PLL_CONTROL
	d.sendData(0x3c)

	d.sendCmd(0x41) // TEMPERATURE_CALIBRATION
	d.sendData(0x00)

	d.sendCmd(0x50) // VCOM_AND_DATA_INTERVAL_SETTING
	d.sendData(0x77)

	d.sendCmd(0x60) // TCON_SETTING
	d.sendData(0x22)

	d.sendCmd(0x61)          // TCON_RESOLUTION
	d.sendData(byte(w >> 8)) //source 640
	d.sendData(byte(w & 0xff))
	d.sendData(byte(h >> 8)) //gate 384
	d.sendData(byte(h & 0xff))

	d.sendCmd(0x82)  // VCM_DC_SETTING
	d.sendData(0x1E) // decide by LUT file

	d.sendCmd(0xe5) // FLASH MODE
	d.sendData(0x03)

	return nil
}

func (d *Dev) waitUntilIdle() {
	for d.busy.Read() == gpio.Low {
		time.Sleep(100 * time.Millisecond)
	}
}

func (d *Dev) sendData(c byte) error {
	if err := d.dc.Out(gpio.High); err != nil {
		return err
	}
	return d.c.Tx([]byte{c}, nil)
}

func (d *Dev) sendDatas(cs []byte) error {
	if err := d.dc.Out(gpio.High); err != nil {
		return err
	}
	return d.c.Tx(cs, nil)
}

func (d *Dev) sendCmd(c byte) error {
	if err := d.dc.Out(gpio.Low); err != nil {
		return err
	}
	return d.c.Tx([]byte{c}, nil)
}

type inkColor byte

const (
	black inkColor = iota
	gray
	white
)

func checkColor(c color.Color) inkColor {
	g := color.GrayModel.Convert(c).(color.Gray)

	if g.Y < 64 {
		return black
	} else if g.Y < 192 {
		return gray
	}
	return white
}

var _ display.Drawer = &Dev{}
