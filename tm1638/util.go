// Copyright 2015, Homin Lee <homin.lee@suapapa.net>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tm1638

import "strconv"

func min(v1, v2 byte) byte {
	min := v1
	if v2 < min {
		return v2
	}
	return v1
}

func b(s string) byte {
	u, err := strconv.ParseUint(s, 2, 8)
	if err != nil {
		panic(err)
	}
	return byte(u)
}
