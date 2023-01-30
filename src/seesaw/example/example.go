package main

import (
	"machine"
	"plantbot/src/seesaw/soil"
	"time"
)

func main() {
	machine.I2C0.Configure(machine.I2CConfig{})
	moist := soil.New(machine.I2C0)

	time.Sleep(3 * time.Second)

	for {
		m, err := moist.Read()
		if err != nil {
			println(err.Error())
		}
		println(m)
		time.Sleep(time.Second)
	}
}
