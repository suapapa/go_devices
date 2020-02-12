package main

import (
	"image"
	"image/png"
	"log"
	"os"
	"time"

	"github.com/suapapa/go_devices/epd7in5"
	rpi_gpio "github.com/suapapa/go_devices/rpi/gpio"
	"golang.org/x/exp/io/spi"
)

func main() {
	imageFileName := os.Args[1]

	log.Println("init sequence...")
	d, err := epd7in5.Open(
		&spi.Devfs{
			Dev:      "/dev/spidev0.0",
			Mode:     spi.Mode0,
			MaxSpeed: 400000,
		},
		&rpi_gpio.Mem{
			PinMap: map[string]int{
				epd7in5.PinRST:  17,
				epd7in5.PinDC:   25,
				epd7in5.PinBusy: 24,
			},
		},
	)
	if err != nil {
		panic(err)
	}
	defer d.Close()

	log.Println("Clear...")
	d.Clear()
	time.Sleep(1 * time.Second)

	log.Println("draw image...")
	img, err := openPNG(imageFileName)
	if err != nil {
		panic(err)
	}
	err = d.DrawImage(img)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("exit...")
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
