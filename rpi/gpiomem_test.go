package rpi

import gpio_driver "github.com/goiot/exp/gpio/driver"

// assert that rpigpio.pin implements gpio_driver.Opener
var _ gpio_driver.Opener = &GpioMem{}

// assert that rpigpio.pin implements gpio_driver.Conn
var _ gpio_driver.Conn = &gpioMemConn{}
