// Copyright 2020, Homin Lee <homin.lee@suapapa.net>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"time"

	"github.com/suapapa/go_devices/max7219"
	"periph.io/x/conn/v3/spi/spireg"
	"periph.io/x/host/v3"
)

func main() {
	_, err := host.Init()
	chk(err)

	bus, err := spireg.Open("")
	chk(err)

	dev, err := max7219.New(bus)
	chk(err)

	dev.DisplayTest(true)
	time.Sleep(1 * time.Second)
	dev.DisplayTest(false)

	// dev.Shutdown(true)
	// time.Sleep(1 * time.Second)
	// dev.Shutdown(false)

	dev.Write(3, 1<<5)
}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}
