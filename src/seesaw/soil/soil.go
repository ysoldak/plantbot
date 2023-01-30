package soil

import (
	"time"

	"tinygo.org/x/drivers"
)

const Address = 0x36

var readMoistureCommand = []uint8{0x0F, 0x10}

// Device wraps an I2C connection to a Soil device.
type Device struct {
	bus     drivers.I2C
	Address uint16
	buf     [2]uint8
}

// New creates a new SeeSaw Soil Sensor connection. The I2C bus must already be configured.
//
// This function only creates the Device object, it does not touch the device.
func New(bus drivers.I2C) *Device {
	return &Device{
		bus:     bus,
		Address: Address,
	}
}

// Read returns the moisture reading in range 0 to 1023
func (d *Device) Read() (t uint16, err error) {
	for retry := 0; retry < 5; retry++ {
		err = d.bus.Tx(d.Address, readMoistureCommand, nil)
		if err != nil {
			continue
		}
		time.Sleep(time.Duration(3000+retry*1000) * time.Microsecond)
		err = d.bus.Tx(d.Address, nil, d.buf[:])
		if err == nil {
			return (uint16(d.buf[0]) << 8) | uint16(d.buf[1]), nil
		}
	}
	return // returns last error
}