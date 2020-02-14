// Copyright 2016, Homin Lee <homin.lee@suapapa.net>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !linux

package gpio

import (
	"fmt"

	gpio_driver "github.com/goiot/exp/gpio/driver"
)

// Mem implements github.com/goiot/exp/gpio/driver.Opener
type Mem struct {
	// PinMap will is used to conver pin name to number
	PinMap map[string]int
}

// Open returns github.com/goiot/exp/gpio/driver.Conn
func (m *Mem) Open() (gpio_driver.Conn, error) {
	return nil, fmt.Errorf("not implemented")
}
