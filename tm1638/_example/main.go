package main

import (
	"github.com/suapapa/go_devices/tm1638"
	"periph.io/x/periph/conn/gpio/gpioreg"
	"periph.io/x/periph/host"
)

func main() {
	if _, err := host.Init(); err != nil {
		panic(err)
	}

	tm1638, err := tm1638.Open(
		gpioreg.ByName("17"), // data
		gpioreg.ByName("27"), // clk
		gpioreg.ByName("22"), // stb
	)
	if err != nil {
		panic(err)
	}

	tm1638.SetString("HelloWrd")
}
