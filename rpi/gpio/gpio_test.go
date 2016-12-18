// Copyright 2016, Homin Lee <homin.lee@suapapa.net>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gpio

import gpio_driver "github.com/goiot/exp/gpio/driver"

// assert that &Mem implements gpio_driver.Opener
var _ gpio_driver.Opener = &Mem{}

// assert that &memConn implements gpio_driver.Conn
var _ gpio_driver.Conn = &memConn{}

// assert that &Sysfs implements gpio_driver.Opener
var _ gpio_driver.Opener = &Sysfs{}

// assert that &sysfsConn implements gpio_driver.Conn
var _ gpio_driver.Conn = &sysfsConn{}
