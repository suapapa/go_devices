package main

import (
	"log"
	"time"

	gpio_driver "github.com/goiot/exp/gpio/driver"
	rpi_gpio "github.com/suapapa/go_devices/rpi/gpio"
)

func main() {
	o := rpi_gpio.Sysfs{PinMap: rpi_gpio.DefaultPinMap}
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
	tC := time.After(1 * time.Second)
	cnt := 0
loop:
	for {
		select {
		case <-tC:
			break loop
		default:
			if cnt%2 == 0 {
				dev.SetValue("BCM16", 1)
			} else {
				dev.SetValue("BCM16", 0)

			}
			cnt++
		}
	}
	log.Println("bye bye~", cnt)
}
