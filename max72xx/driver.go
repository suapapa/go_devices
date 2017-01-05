// Copyright 2016, Homin Lee <homin.lee@suapapa.net>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package max72xx

import (
	"bytes"
	"fmt"

	"golang.org/x/exp/io/spi"
	spi_driver "golang.org/x/exp/io/spi/driver"
)

// Driver represents MAX72XX LED driver chain
type Driver struct {
	dev     *spi.Device
	Drivers int // can connect up to 8 driver
	buff    []byte
}

// Open opens MAX72XX LED driver chain
func Open(bus spi_driver.Opener, numDriver int) (*Driver, error) {
	spiDev, err := spi.Open(bus)
	if err != nil {
		return nil, err
	}

	if numDriver <= 0 || numDriver > 8 {
		numDriver = 1
	}

	driver := &Driver{
		dev:     spiDev,
		Drivers: numDriver,
		buff:    make([]byte, numDriver*8),
	}

	driver.init()

	return driver, nil
}

// Close closes MAX72XX driver
func (d *Driver) Close() {
	d.dev.Close()
}

// Display displays internel buffer con
func (d *Driver) Display() error {
	buff := make([]byte, d.Drivers*2)

	for r := 0; r < 8; r++ {
		rowNo := r + 1
		for i := 0; i < d.Drivers; i++ {
			buffIdx := d.Drivers*i + r
			outIdx := (d.Drivers - i - 1) * 2
			buff[outIdx] = byte(rowNo)
			buff[outIdx+1] = d.buff[buffIdx]
		}

		if err := d.dev.Tx(buff, nil); err != nil {
			return err
		}
	}
	return nil
}

// Clear clears internel buffer
// need to call Diaplay to up actual diaplays
func (d *Driver) Clear() {
	for i := 0; i < len(d.buff); i++ {
		d.buff[i] = 0
	}
}

func (d *Driver) init() {
	d.opAll(opDISPLAYTEST, 0)
	d.opAll(opSCANLIMIT, 7)
	d.opAll(opDECODEMODE, 0)
	for i := byte(1); i <= 8; i++ {
		d.opAll(i, 0)
	}
	d.opAll(opSHUTDOWN, 0)
}

func (d *Driver) op(idx int, op, data byte) error {
	buff := make([]byte, d.Drivers*2)
	buffIdx := (d.Drivers - idx - 1) * 2

	if buffIdx < 0 || buffIdx >= d.Drivers*2-1 {
		return fmt.Errorf("max72xx: invalid addr, %d", idx)
	}

	buff[buffIdx], buff[buffIdx+1] = op, data

	return d.dev.Tx(buff, nil)
}

func (d *Driver) opAll(op, data byte) error {
	buff := []byte{op, data}
	return d.dev.Tx(bytes.Repeat(buff, d.Drivers), nil)
}
