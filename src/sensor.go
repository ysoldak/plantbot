package main

import "machine"

const adcMax = 1<<16 - 1 - 3000 // don't know why 3000, supposed to be 2^12 scaled up to uint16
const adcRef = 3.3

const digitalToAnalog = float64(adcRef) / adcMax

func init() {
	machine.InitADC()
}

// -----------------------------------------------------------------------------

type sensor struct {
	pin machine.Pin
	min float64
	max float64

	reading sensorReading

	adc machine.ADC
}

type sensorReading struct {
	digital  uint16
	analog   float64
	fraction float64
}

func (s *sensor) configure() {
	s.adc = machine.ADC{Pin: s.pin}
	s.adc.Configure(machine.ADCConfig{})
}

func (s *sensor) read() (sr sensorReading) {
	sr = sensorReading{}
	sr.digital = s.adc.Get()
	sr.analog = float64(sr.digital) * digitalToAnalog
	if float64(sr.digital) > s.min {
		sr.fraction = (float64(sr.digital) - s.min) / (s.max - s.min)
	}
	s.reading = sr
	return
}
