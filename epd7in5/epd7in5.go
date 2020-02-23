package epd7in5

import (
	"fmt"
	"image"
	"image/color"

	"github.com/suapapa/go_devices/epd7in5/gray3"
	"periph.io/x/periph/conn"
	"periph.io/x/periph/conn/gpio"
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

func (d *Dev) String() string {
	return fmt.Sprintf("epd7in5.Dev{%s, %s, %s}", d.c, d.dc, d.rect.Max)
}

// Halt turns off the display
func (d *Dev) Halt() error {
	return nil
}

// ColorModel implements display.Drawer
func (d *Dev) ColorModel() color.Model {
	return gray3.Gray3Model
}

// Bounds implements display.Drawer
func (d *Dev) Bounds() image.Rectangle {
	return image.Rectangle{
		Min: image.Point{0, 0},
		Max: image.Point{w, h},
	}
}

// Draw implements display.Drawer
func Draw(dstRect image.Rectangle, src image.Image, srcPtr image.Point) error {
	return nil
}
