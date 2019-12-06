package main

import (
	"image"
	"image/png"
	"os"

	"github.com/suapapa/go_devices/epd2in13"
	rpi_gpio "github.com/suapapa/go_devices/rpi/gpio"
	"golang.org/x/exp/io/spi"
)

func main() {
	d, err := epd2in13.Open(
		&spi.Devfs{Dev: "/dev/spidev0.0"},
		&rpi_gpio.Mem{
			PinMap: map[string]int{
				epd2in13.PinRST:  17,
				epd2in13.PinDC:   25,
				epd2in13.PinCS:   8,
				epd2in13.PinBusy: 24,
			},
		},
	)
	if err != nil {
		panic(err)
	}
	defer d.Close()

	img, err := openPNG("gopher_250x122.png")
	if err != nil {
		panic(err)
	}

	d.DrawImage(img)
	d.TurnOnFull()
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
