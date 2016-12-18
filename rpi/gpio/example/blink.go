package main

import (
	"log"
	"time"

	gpio_driver "github.com/goiot/exp/gpio/driver"
	rpi_gpio "github.com/suapapa/go_devices/rpi/gpio"
)

func main() {
	o := rpi_gpio.Mem{PinMap: rpi_gpio.DefaultPinMap}
	dev, err := o.Open()
	if err != nil {
		panic(err)
	}
	defer dev.Close()

	err = dev.SetDirection("BCM16", gpio_driver.Out)
	if err != nil {
		panic(err)
	}

	log.Println("start blinking...")
	tC := time.After(10 * time.Second)
loop:
	for {
		select {
		case <-tC:
			break loop
		default:
			dev.SetValue("BCM16", 1)
			time.Sleep(time.Second)
			dev.SetValue("BCM16", 0)
			time.Sleep(time.Second)

		}
	}
	log.Println("bye bye~")
}
