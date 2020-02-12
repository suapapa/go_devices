package epd7in5_test

import (
	"github.com/suapapa/go_devices/epd7in5"
	rpi_gpio "github.com/suapapa/go_devices/rpi/gpio"
	"golang.org/x/exp/io/spi"
)

func ExampleOpen() {
	d, err := epd7in5.Open(
		&spi.Devfs{Dev: "/dev/spidev0.0"},
		&rpi_gpio.Sysfs{
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
}
