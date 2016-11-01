package ili9340

import (
	"image"
	"image/color"
	"log"
	"time"

	"golang.org/x/exp/io/spi"
	"golang.org/x/exp/io/spi/driver"

	"github.com/davecheney/gpio"
)

// LCD reperesents TFT-LCD panel which using ili9340 controller
type LCD struct {
	dev  *spi.Device
	dc   gpio.Pin
	w, h int
}

// Open connets to the passed driver and sets things up
func Open(o driver.Opener, dc gpio.Pin) (*LCD, error) {
	/*
		dev, err := spi.Open(&spi.Devfs{
			Dev:      "/dev/spidev0.1",
			Mode:     spi.Mode3,
			MaxSpeed: 500000,
		})
	*/
	dc.SetMode(gpio.ModeInput)
	dc.SetMode(gpio.ModeOutput)

	device, err := spi.Open(o)
	if err != nil {
		return nil, err
	}
	device.SetCSChange(true)

	lcd := &LCD{
		dev: device,
		dc:  dc,
	}

	lcd.reset()
	lcd.init()
	lcd.Rotate(0)

	return lcd, nil
}

// DrawPixel sets a pixel at a position x, y.
func (l *LCD) DrawPixel(x, y int, c color.Color) {
	rgb := rgb565(c)
	l.PushColor(x, y, x+1, y+1, rgb)
}

// DrawRect draw a rect at x,y
func (l *LCD) DrawRect(x, y, w, h int, c color.Color) {
	rgb := rgb565(c)
	l.PushColor(x, y, x+w-1, y+h-1, rgb)
}

// SetImage draws an image on the display starting from x, y.
func (l *LCD) SetImage(x, y int, img image.Image) {
	imgW, imgH := img.Bounds().Dx(), img.Bounds().Dy()
	endX, endY := x+imgW, y+imgH

	if endX >= l.w {
		endX = l.w
	}
	if endY >= l.h {
		endY = l.h
	}

	var imgX, imgY int
	for i := x; i < endX; i++ {
		imgY = 0
		for j := y; j < endY; j++ {
			l.DrawPixel(i, j, img.At(imgX, imgY))
			imgY++
		}
		imgX++
	}
}

// Width returns the display width.
func (l *LCD) Width() int { return l.w }

// Height returns the display height.
func (l *LCD) Height() int { return l.h }

// Close takes care of cleaning things up.
func (l *LCD) Close() {
	l.dev.Close()
	l.dc.Close()
}

// PushColor sets a color to window
func (l *LCD) PushColor(xs, ys, xe, ye int, c uint16) {
	/* Column address */
	l.writeReg(0x2A,
		uint8(xs>>8), uint8(xs&0xFF),
		uint8(xe>>8), uint8(xe&0xFF),
	)

	/* Row address */
	l.writeReg(0x2B,
		uint8(ys>>8), uint8(ys&0xFF),
		uint8(ye>>8), uint8(ye&0xFF),
	)

	/* Memory write */
	l.writeReg(0x2C,
		uint8(c>>8), uint8(c),
	)
}

// Rotate rotates display
func (l *LCD) Rotate(r int) {
	var val byte
	switch r {
	case 270:
		val = ili9340_MADCTL_MV
		l.w = ili9340_TFTHEIGHT
		l.h = ili9340_TFTWIDTH
	case 180:
		val = ili9340_MADCTL_MY
		l.w = ili9340_TFTWIDTH
		l.h = ili9340_TFTHEIGHT
	case 90:
		val = ili9340_MADCTL_MV | ili9340_MADCTL_MY | ili9340_MADCTL_MX
		l.w = ili9340_TFTHEIGHT
		l.h = ili9340_TFTWIDTH
	default:
		val = ili9340_MADCTL_MX
		l.w = ili9340_TFTWIDTH
		l.h = ili9340_TFTHEIGHT
	}

	l.writeReg(0x36, val|ili9340_MADCTL_BGR)
}

// Invert invert colors on the display
func (l *LCD) Invert(on bool) {
	if on {
		l.writeReg(ili9340_INVON)
	} else {
		l.writeReg(ili9340_INVOFF)
	}
}

func (l *LCD) init() {
	l.writeReg(0xEF, 0x03, 0x80, 0x02)
	l.writeReg(0xCF, 0x00, 0xC1, 0x30)
	l.writeReg(0xED, 0x64, 0x03, 0x12, 0x81)
	l.writeReg(0xE8, 0x85, 0x00, 0x78)
	l.writeReg(0xCB, 0x39, 0x2C, 0x00, 0x34, 0x02)
	l.writeReg(0xF7, 0x20)
	l.writeReg(0xEA, 0x00, 0x00)

	/* Power Control 1 */
	l.writeReg(0xC0, 0x23)

	/* Power Control 2 */
	l.writeReg(0xC1, 0x10)

	/* VCOM Control 1 */
	l.writeReg(0xC5, 0x3e, 0x28)

	/* VCOM Control 2 */
	l.writeReg(0xC7, 0x86)

	/* COLMOD: Pixel Format Set */
	/* 16 bits/pixel */
	l.writeReg(0x3A, 0x55)

	/* Frame Rate Control */
	/* Division ratio = fosc, Frame Rate = 79Hz */
	l.writeReg(0xB1, 0x00, 0x18)

	/* Display Function Control */
	l.writeReg(0xB6, 0x08, 0x82, 0x27)

	/* Gamma Function Disable */
	l.writeReg(0xF2, 0x00)

	/* Gamma curve selected  */
	l.writeReg(0x26, 0x01)

	/* Positive Gamma Correction */
	l.writeReg(0xE0,
		0x0F, 0x31, 0x2B, 0x0C, 0x0E, 0x08, 0x4E, 0xF1,
		0x37, 0x07, 0x10, 0x03, 0x0E, 0x09, 0x00)

	/* Negative Gamma Correction */
	l.writeReg(0xE1,
		0x00, 0x0E, 0x14, 0x03, 0x11, 0x07, 0x31, 0xC1,
		0x48, 0x08, 0x0F, 0x0C, 0x31, 0x36, 0x0F)

	/* Sleep OUT */
	l.writeReg(0x11)

	time.Sleep(120 * time.Millisecond)

	/* Display ON */
	l.writeReg(0x29)
}

func (l *LCD) reset() {
	// l.rst.Set()
	// time.Sleep(5 * time.Millisecond)
	// l.rst.Clear()
	// time.Sleep(20 * time.Millisecond)
	// l.rst.Set()
	// time.Sleep(120 * time.Millisecond)
}

func (l *LCD) writeReg(vs ...uint8) {
	if vs == nil || len(vs) == 0 {
		return
	}

	l.dc.Clear()
	log.Println("write command:", vs[:1])
	l.dev.Tx(vs[:1], nil)

	if len(vs) > 1 {
		l.dc.Set()
		log.Println("write data:", vs[1:])
		l.dev.Tx(vs[1:], nil)
	}
}
