// Copyright 2016, Homin Lee <homin.lee@suapapa.net>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package max72xx

import (
	"golang.org/x/exp/io/spi"
	spi_driver "golang.org/x/exp/io/spi/driver"
)

// Driver represents MAX72XX LED driver chain
type Driver struct {
	dev *spi.Device
}

// Open opens MAX72XX LED driver chain
func Open(bus spi_driver.Opener) (*Driver, error) {
	spiDev, err := spi.Open(bus)
	if err != nil {
		return nil, err
	}

	driver := &Driver{
		dev: spiDev,
	}

	return driver, nil
}

// Close closes MAX72XX driver
func (d *Driver) Close() {
	d.dev.Close()
}

