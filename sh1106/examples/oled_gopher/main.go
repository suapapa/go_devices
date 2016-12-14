package main

import (
	"image"
	"image/png"
	"os"

	rpi_gpio "github.com/suapapa/go_devices/rpi/gpio"
	"github.com/suapapa/go_devices/sh1106"
	"golang.org/x/exp/io/i2c"
)

func main() {
	l, err := sh1106.OpenI2C(
		&i2c.Devfs{Dev: "/dev/i2c-1"},
		&rpi_gpio.Mem{
			PinMap: map[string]int{
				sh1106.PinRST: 14,
			},
		},
		sh1106.DefaultI2CAddr, // 0x3C
	)
	if err != nil {
		panic(err)
	}
	defer l.Close()

	img, err := openPNG("gopher-side_128x64.png")
	if err != nil {
		panic(err)
	}

	l.DrawImage(img)
	l.Display()
}

func openPNG(filename string) (image.Image, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, err := png.Decode(f)
	if err != nil {
		return nil, err
	}

	return img, nil
}
