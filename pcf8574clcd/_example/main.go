// Copyright 2020, Homin Lee <homin.lee@suapapa.net>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/suapapa/go_devices/pcf8574clcd"
	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/host"
)

func main() {
	_, err := host.Init()
	chk(err)

	bus, err := i2creg.Open("")
	chk(err)

	dev, err := pcf8574clcd.New(
		bus,
		0x27, /* pcf8574.DefaultAddr */
		16, 2,
	)
	chk(err)

	err = dev.BackLight(true)
	chk(err)

	err = dev.SetCursor(0, 0)
	chk(err)
	err = dev.Write("Hello~")
	chk(err)

	err = dev.SetCursor(0, 1)
	chk(err)
	err = dev.Write("RaspberryPi!!")
	chk(err)

}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}
