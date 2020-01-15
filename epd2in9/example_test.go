package epd2in9_test

import (
	"github.com/suapapa/go_devices/epd2in9"
	rpi_gpio "github.com/suapapa/go_devices/rpi/gpio"
	"golang.org/x/exp/io/spi"
)

func ExampleOpen() {
	d, err := epd2in9.Open(
		&spi.Devfs{Dev: "/dev/spidev0.0"},
		&rpi_gpio.Sysfs{
			PinMap: map[string]int{
				epd2in9.PinRST:  17,
				epd2in9.PinDC:   25,
				epd2in9.PinBusy: 24,
			},
		},
	)
	if err != nil {
		panic(err)
	}
	defer d.Close()
}
