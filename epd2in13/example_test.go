package epd2in13_test

import (
	"github.com/suapapa/go_devices/epd2in13"
	rpi_gpio "github.com/suapapa/go_devices/rpi/gpio"
	"golang.org/x/exp/io/spi"
)

func ExampleOpen() {
	d, err := epd2in13.Open(
		&spi.Devfs{Dev: "/dev/spidev0.0"},
		&rpi_gpio.Sysfs{
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
}
