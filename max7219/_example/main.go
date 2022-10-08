// Copyright 2020, Homin Lee <homin.lee@suapapa.net>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"time"

	"github.com/suapapa/go_devices/max7219"
	"periph.io/x/periph/conn/spi/spireg"
	"periph.io/x/periph/host"
)

func main() {
	_, err := host.Init()
	chk(err)

	bus, err := spireg.Open("")
	chk(err)

	dev, err := max7219.New(bus)
	chk(err)

	dev.DisplayTest(true)
	time.Sleep(5 * time.Second)
	dev.DisplayTest(false)

	dev.Shutdown(true)
}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}
