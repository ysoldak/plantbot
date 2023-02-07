package main

import (
	"machine"
	"time"

	"tinygo.org/x/drivers/net"
	"tinygo.org/x/drivers/wifinina"
)

var (
	spi     = machine.NINA_SPI
	adaptor *wifinina.Device
)

func setupWifi(ssid, pass string) (err error) {
	if ssid == "" || pass == "" {
		return nil
	}
	configureWifinina()
	return connectToAP(ssid, pass)
}

func configureWifinina() {
	spi.Configure(machine.SPIConfig{
		Frequency: 8 * 1e6,
		SDO:       machine.NINA_SDO,
		SDI:       machine.NINA_SDI,
		SCK:       machine.NINA_SCK,
	})

	adaptor = wifinina.New(spi,
		machine.NINA_CS,
		machine.NINA_ACK,
		machine.NINA_GPIO0,
		machine.NINA_RESETN)

	net.ActiveDevice = nil
	adaptor.Configure()
}

func stopWifinina() {
	adaptor.Stop()
}

func connectToAP(ssid, pass string) (err error) {

	time.Sleep(2 * time.Second)
	trace("Connecting to " + ssid)
	for i := 0; i < 10; i++ {
		err = adaptor.ConnectToAccessPoint(ssid, pass, 10*time.Second)
		if err == nil {
			break
		}
		println(err.Error())
		time.Sleep(1 * time.Second)
	}
	if err != nil {
		return
	}
	trace("Connected")

	time.Sleep(2 * time.Second)
	var ip wifinina.IPAddress
	for i := 0; i < 10; i++ {
		ip, _, _, err = adaptor.GetIP()
		if err == nil {
			break
		}
		println(err.Error())
		time.Sleep(1 * time.Second)
	}
	trace(ip.String())
	return
}
