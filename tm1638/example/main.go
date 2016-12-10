// Copyright 2015, Homin Lee <homin.lee@suapapa.net>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"os/signal"
	"time"

	"github.com/suapapa/go_tm1638"
)

func main() {
	d, err := tm1638.NewTM1638(18, 23, 24)
	if err != nil {
		panic(err)
	}

	d.SetLEDs(0x0000)
	time.Sleep(1 * time.Second)
	d.DisplayError()
	time.Sleep(1 * time.Second)
	d.DisplayHexNumber(0xff4500, 0x0F, true)
	time.Sleep(1 * time.Second)
	d.DisplayDecNumber(uint64(time.Now().UnixNano()%100000000), 0, true)
	time.Sleep(1 * time.Second)
	d.DisplaySignedDecNumber(-2345678, 0, false)
	time.Sleep(1 * time.Second)
	d.DisplayBinNumber(0x45, 0xF0)
	time.Sleep(1 * time.Second)
	d.SetLEDs(0xF00F)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			d.Close()
			os.Exit(0)
		}
	}()

	for {
		time.Sleep(10 * time.Millisecond)
		keys := d.GetButton()
		d.DisplayBinNumber(keys, 0x00)
	}

}
