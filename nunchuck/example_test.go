// Copyright 2016, Homin Lee <homin.lee@suapapa.net>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nunchuck_test

import (
	"github.com/suapapa/go_devices/nunchuck"
	"golang.org/x/exp/io/i2c"
)

func ExampleOpen() {
	c, err := nunchuck.Open(&i2c.Devfs{Dev: "/dev/i2c-1"})
	if err != nil {
		panic(err)
	}
	defer c.Close()
}
