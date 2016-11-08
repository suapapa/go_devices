package main

import (
	"image"
	"image/png"
	"os"

	"github.com/suapapa/go_devices/sh1106"
	"golang.org/x/exp/io/i2c"
)

func main() {
	l, err := sh1106.OpenI2C(&i2c.Devfs{Dev: "/dev/i2c-1"}, 0x3C, nil)
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
