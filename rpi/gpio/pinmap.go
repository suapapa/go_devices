// Copyright 2016, Homin Lee <homin.lee@suapapa.net>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gpio

import "strconv"

// PinMapFunc is a function for convert pin name to number
type PinMapFunc func(string) (int, error)

var defaultPinMap PinMapFunc = func(n string) (int, error) {
	return strconv.Atoi(n)
}
