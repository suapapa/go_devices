package epd2in13bc_test

import (
	"github.com/suapapa/go_devices/epd2in13bc"
	rpi_gpio "github.com/suapapa/go_devices/rpi/gpio"
	"golang.org/x/exp/io/spi"
)

func ExampleOpen() {
	d, err := epd2in13bc.Open(
		&spi.Devfs{Dev: "/dev/spidev0.0"},
		&rpi_gpio.Sysfs{
			PinMap: map[string]int{
				epd2in13bc.PinRST:  17,
				epd2in13bc.PinDC:   25,
				epd2in13bc.PinBusy: 24,
			},
		},
	)
	if err != nil {
		panic(err)
	}
	defer d.Close()
}
