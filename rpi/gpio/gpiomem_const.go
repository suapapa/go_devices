// Copyright 2016, Homin Lee <homin.lee@suapapa.net>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gpio

const (
	bcm2835gpFsel0   = 0x0000 // GPIO Function Select 0
	bcm2835gpFsel1   = 0x0004 // GPIO Function Select 1
	bcm2835gpFsel2   = 0x0008 // GPIO Function Select 2
	bcm2835gpFsel3   = 0x000c // GPIO Function Select 3
	bcm2835gpFsel4   = 0x0010 // GPIO Function Select 4
	bcm2835gpFsel5   = 0x0014 // GPIO Function Select 5
	bcm2835gpSet0    = 0x001c // GPIO Pin Output Set 0
	bcm2835gpSet1    = 0x0020 // GPIO Pin Output Set 1
	bcm2835gpClr0    = 0x0028 // GPIO Pin Output Clear 0
	bcm2835gpClr1    = 0x002c // GPIO Pin Output Clear 1
	bcm2835gpLev0    = 0x0034 // GPIO Pin Level 0
	bcm2835gpLev1    = 0x0038 // GPIO Pin Level 1
	bcm2835gpEds0    = 0x0040 // GPIO Pin Event Detect Status 0
	bcm2835gpEds1    = 0x0044 // GPIO Pin Event Detect Status 1
	bcm2835gpREN0    = 0x004c // GPIO Pin Rising Edge Detect Enable 0
	bcm2835gpREN1    = 0x0050 // GPIO Pin Rising Edge Detect Enable 1
	bcm2835gpFEN0    = 0x0048 // GPIO Pin Falling Edge Detect Enable 0
	bcm2835gpFEN1    = 0x005c // GPIO Pin Falling Edge Detect Enable 1
	bcm2835gpHEN0    = 0x0064 // GPIO Pin High Detect Enable 0
	bcm2835gpHEN1    = 0x0068 // GPIO Pin High Detect Enable 1
	bcm2835gpLEN0    = 0x0070 // GPIO Pin Low Detect Enable 0
	bcm2835gpLEN1    = 0x0074 // GPIO Pin Low Detect Enable 1
	bcm2835gpAREN0   = 0x007c // GPIO Pin Async. Rising Edge Detect 0
	bcm2835gpAREN1   = 0x0080 // GPIO Pin Async. Rising Edge Detect 1
	bcm2835gpAFEN0   = 0x0088 // GPIO Pin Async. Falling Edge Detect 0
	bcm2835gpAFEN1   = 0x008c // GPIO Pin Async. Falling Edge Detect 1
	bcm2835gpPUD     = 0x0094 // GPIO Pin Pull-up/down Enable
	bcm2835gpPUDCLK0 = 0x0098 // GPIO Pin Pull-up/down Enable Clock 0
	bcm2835gpPUDCLK1 = 0x009c // GPIO Pin Pull-up/down Enable Clock 1

	bcm2835gpioFselINPT        = 0x0 // Input
	bcm2835gpioFselOUTP        = 0x1 // Output
	bcm2835gpioFselALT0        = 0x4 // Alternate function 0
	bcm2835gpioFselALT1        = 0x5 // Alternate function 1
	bcm2835gpioFselALT2        = 0x6 // Alternate function 2
	bcm2835gpioFselALT3        = 0x7 // Alternate function 3
	bcm2835gpioFselALT4        = 0x3 // Alternate function 4
	bcm2835gpioFselALT5        = 0x2 // Alternate function 5
	bcm2835gpioFselMASK uint32 = 0x7
)
