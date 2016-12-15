// Copyright 2016, Homin Lee <homin.lee@suapapa.net>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ili9340_test

import (
	"github.com/suapapa/go_devices/ili9340"
	rpi_gpio "github.com/suapapa/go_devices/rpi/gpio"
	"golang.org/x/exp/io/spi"
)

func ExampleOpen() {
	dev, err := ili9340.Open(
		&spi.Devfs{
			Dev:      "/dev/spidev0.1",
			Mode:     spi.Mode3,
			MaxSpeed: 500000,
		},
		&rpi_gpio.Mem{
			map[string]int{
				ili9340.PinDC:  11,
				ili9340.PinRST: 12,
			},
		},
	)
	if err != nil {
		panic(err)
	}
	defer dev.Close()
}
