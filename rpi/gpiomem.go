package rpi

import gpio_driver "github.com/goiot/exp/gpio/driver"

// GpioMem implements github.com/goiot/exp/gpio/driver.Opener
type GpioMem struct{}

// Open returns github.com/goiot/exp/gpio/driver.Conn
func (m GpioMem) Open() (gpio_driver.Conn, error) {
	// TODO: do mmap
	// TODO: make register buffers
	// TODO: make mem map
	conn := &gpioMemConn{}

	return conn, nil
}

type pin struct {
	idx    int
	offset byte
}

// implements github.com/goiot/exp/gpio/driver.Conn
type gpioMemConn struct {
	mem map[string]pin
}

// Value returns the value of the pin. 0 for low values, 1 for high.
func (m *gpioMemConn) Value(pin string) (int, error) {
	// TBD
	return 0, nil
}

// SetValue sets the vaule of the pin. 0 for low values, 1 for high.
func (m *gpioMemConn) SetValue(pin string, v int) error {
	// TBD
	return nil
}

// SetDirection sets the direction of the pin.
func (m *gpioMemConn) SetDirection(pin string, dir gpio_driver.Direction) error {
	// TBD
	return nil
}

// Map is not implemented
func (m *gpioMemConn) Map(vitual string, physical int) {
	// TBD
	panic("rpigpio: not implemented")
}

// Close closes the connection and free the underlying resources.
func (m *gpioMemConn) Close() error {
	// TBD
	return nil
}
