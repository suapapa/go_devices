package ili9340

// Copyright 2016, Homin Lee <homin.lee@suapapa.net>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"bytes"
	"fmt"
	"image/color"
	"testing"
)

func TestRGB565(t *testing.T) {
	testCases := []struct {
		c color.Color
		b []byte
	}{
		{c: color.Black, b: []byte{0x00, 0x00}},
		{c: color.White, b: []byte{0xff, 0xff}},
		{c: color.RGBA{R: 0xff}, b: []byte{0xf8, 0x00}},
		{c: color.RGBA{G: 0xff}, b: []byte{0x07, 0xe0}},
		{c: color.RGBA{B: 0xff}, b: []byte{0x00, 0x1f}},
	}
	for _, tc := range testCases {
		t.Run(
			fmt.Sprintf("%T%v to 0x%04x", tc.c, tc.c, tc.b),
			func(t *testing.T) {
				if got := rgb565(tc.c); bytes.Compare(got, tc.b) != 0 {
					t.Errorf("got 0x%04x; want 0x%04x", got, tc.b)
				}
			},
		)
	}
}
