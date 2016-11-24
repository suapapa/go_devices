package nunchuck_test

import (
	"fmt"

	"github.com/suapapa/go_devices/nunchuck"
	"golang.org/x/exp/io/i2c"
)

func ExampleOpen() {
	c, err := nunchuck.Open(&i2c.Devfs{Dev: "/dev/i2c-1"})
	if err != nil {
		panic(err)
	}
	defer c.Close()

	c.Update()
	fmt.Println(c)
}
