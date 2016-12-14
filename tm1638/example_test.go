package tm1638_test

import (
	"github.com/suapapa/go_devices/rpi/gpio"
	"github.com/suapapa/go_devices/tm1638"
)

func ExampleOpen() {
	m, err := tm1638.Open(
		&gpio.Mem{
			PinMap: map[string]int{
				tm1638.PinCLK:  18,
				tm1638.PinDATA: 23,
				tm1638.PinSTB:  24,
			},
		},
	)
	if err != nil {
		panic(err)
	}
	defer m.Close()
}
