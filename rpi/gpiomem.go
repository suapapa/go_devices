package rpi

import (
	"fmt"
	"os"
	"strconv"
	"syscall"
	"unsafe"

	gpio_driver "github.com/goiot/exp/gpio/driver"
)

// GpioMem implements github.com/goiot/exp/gpio/driver.Opener
type GpioMem struct{}

// Open returns github.com/goiot/exp/gpio/driver.Conn
func (m *GpioMem) Open() (gpio_driver.Conn, error) {
	conn := &gpioMemConn{}
	err := conn.mmap()
	if err != nil {
		return nil, err
	}

	// Need to map all gpios through 0 to 53?
	for i := 0; i < 54; i++ {
		conn.name[strconv.Itoa(i)] = i
	}

	return conn, nil
}

// implements github.com/goiot/exp/gpio/driver.Conn
type gpioMemConn struct {
	buf                         []byte
	gpfsel, gpset, gpclr, gplev []*uint32
	name                        map[string]int
}

// Value returns the value of the pin. 0 for low values, 1 for high.
func (c *gpioMemConn) Value(pin string) (int, error) {
	if p, ok := c.name[pin]; ok {
		offset, shift := p/32, byte(p%32)
		v := 0
		if *c.gplev[offset]&(1<<shift) == (1 << shift) {
			v = 1
		}
		return v, nil
	}

	return 0, fmt.Errorf("rpi: unknown gpio name, %s", pin)
}

// SetValue sets the vaule of the pin. 0 for low values, 1 for high.
func (c *gpioMemConn) SetValue(pin string, v int) error {
	if p, ok := c.name[pin]; ok {
		offset, shift := p/32, byte(p%32)
		*c.gpset[offset] = (1 << shift)
		return nil
	}

	return fmt.Errorf("rpi: unknown gpio name, %s", pin)
}

// SetDirection sets the direction of the pin.
func (c *gpioMemConn) SetDirection(pin string, dir gpio_driver.Direction) error {
	if p, ok := c.name[pin]; ok {
		offset, shift := p/10, uint32(p%10)*3
		mask := bcm2835gpioFselMASK << shift
		var mode uint32
		switch dir {
		case gpio_driver.In:
			mode = bcm2835gpioFselINPT
		case gpio_driver.Out:
			mode = bcm2835gpioFselOUTP
		default:
			return fmt.Errorf("rpi: uknnown gpiodir, %v", dir)
		}

		v := *c.gpfsel[offset]
		v &= ^mask
		v |= mode << shift

		*c.gpfsel[offset] = v & mask

		return nil
	}
	return fmt.Errorf("rpi: unknown gpio name, %s", pin)
}

// Map maps virtual gpio pin name to a physical pin number
func (c *gpioMemConn) Map(virtual string, physical int) {
	c.name[virtual] = physical
}

// Close closes the connection and free the underlying resources.
func (c *gpioMemConn) Close() error {
	return syscall.Munmap(c.buf)
}

func (c *gpioMemConn) mmap() error {
	f, err := os.OpenFile("/dev/gpiomem", os.O_RDWR|os.O_SYNC, 0)
	if err != nil {
		return err
	}
	defer f.Close()

	buf, err := syscall.Mmap(int(f.Fd()),
		0, 4*1024,
		syscall.PROT_READ|syscall.PROT_WRITE,
		syscall.MAP_SHARED,
	)
	if err != nil {
		return err
	}

	c.gpfsel = []*uint32{
		(*uint32)(unsafe.Pointer(&buf[bcm2835gpFsel0])),
		(*uint32)(unsafe.Pointer(&buf[bcm2835gpFsel1])),
		(*uint32)(unsafe.Pointer(&buf[bcm2835gpFsel2])),
		(*uint32)(unsafe.Pointer(&buf[bcm2835gpFsel3])),
		(*uint32)(unsafe.Pointer(&buf[bcm2835gpFsel4])),
		(*uint32)(unsafe.Pointer(&buf[bcm2835gpFsel5])),
	}
	c.gpset = []*uint32{
		(*uint32)(unsafe.Pointer(&buf[bcm2835gpSet0])),
		(*uint32)(unsafe.Pointer(&buf[bcm2835gpSet1])),
	}
	c.gpclr = []*uint32{
		(*uint32)(unsafe.Pointer(&buf[bcm2835gpClr0])),
		(*uint32)(unsafe.Pointer(&buf[bcm2835gpClr1])),
	}
	c.gplev = []*uint32{
		(*uint32)(unsafe.Pointer(&buf[bcm2835gpLev0])),
		(*uint32)(unsafe.Pointer(&buf[bcm2835gpLev1])),
	}

	c.buf = buf
	return nil
}
