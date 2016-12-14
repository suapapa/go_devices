package sh1106_test

import (
	rpi_gpio "github.com/suapapa/go_devices/rpi/gpio"
	"github.com/suapapa/go_devices/sh1106"
	"golang.org/x/exp/io/i2c"
)

func ExampleOpenI2C() {
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
}
