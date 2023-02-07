package main

import "machine"

const adcMax = 1<<16 - 1 - 3000 // don't know why 3000, supposed to be 2^12 scaled up to uint16
const adcRef = 3.3

const digitalToAnalog = float64(adcRef) / adcMax

func init() {
	machine.InitADC()
}

// -----------------------------------------------------------------------------

const (
	batCellMin = 3.0  // mV
	batCellMax = 4.25 // mV
	batCellCnt = 2    // cell count

	// 47kOhm and 75kOhm resistors;
	// Vout = Vin * R2 / (R1+R2); R1 = 47, R2 = 75, Vin = 8.5, Vout = 3.27
	// https://www.digikey.se/en/resources/conversion-calculators/conversion-calculator-voltage-divider
	batSensorAnalogDiv = float64(47) / (75 + 47)
	batSensorAnalogMin = batCellMin * batCellCnt * batSensorAnalogDiv
	batSensorAnalogMax = batCellMax * batCellCnt * batSensorAnalogDiv

	batSensorMin = batSensorAnalogMin / digitalToAnalog
	batSensorMax = batSensorAnalogMax / digitalToAnalog
)

var batterySensor = sensor{
	pin: batterySensorPin,
	min: batSensorMin,
	max: batSensorMax,
	domain: func(sr sensorReading) float64 {
		return sr.analog / (batSensorAnalogDiv * batCellCnt)
	},
	percent: func(sr sensorReading) float64 {
		return 100 * sr.fraction
	},
}

// -----------------------------------------------------------------------------

var moistureSensor = sensor{
	pin:    moistureSensorPin,
	min:    moistureSensorMin,
	max:    moistureSensorMax,
	domain: nil,
	percent: func(sr sensorReading) float64 {
		return 100 * (1 - sr.fraction) // 0 - wet, 100 - dry, so have to inverse
	},
}

// -----------------------------------------------------------------------------

type sensor struct {
	pin machine.Pin
	min float64
	max float64

	domain  func(sr sensorReading) float64
	percent func(sr sensorReading) float64

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
	return
}
