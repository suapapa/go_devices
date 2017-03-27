# devices in Go

This repository has 
implementation of interface, `golang.org/x/exp/gpio/...` and,
several devices which based on `golang.org/x/exp/io/...`.

----

This package implements `gpio` and `gpio/driver`:
* rpi/gpio

These packages communicates i2c or spi:
* nunchuck
* max72xx
* ili9340
* sh1106

This package communicates with it's own protocol using gpio:
* tm1638

You can find `example` in each directories.

# install

    $ go get -u github.com/suapapa/go_devices/...

# references
* [A Proposal: Peripheral I/O for Go](http://go-talks.appspot.com/github.com/rakyll/talks/pio/pio.slide)

# license
See `LICENSE` file.