package main

const (
	batCellMin   = 3.0  // volts
	batCellMax   = 4.25 // volts
	batCellCount = 2    // cells in battery

	// 47kOhm and 75kOhm resistors;
	// Vout = Vin * R2 / (R1+R2); R1 = 47, R2 = 75, Vin = 8.5, Vout = 3.27
	// https://www.digikey.se/en/resources/conversion-calculators/conversion-calculator-voltage-divider
	batSensorAnalogDiv = float64(47) / (75 + 47)
	batSensorAnalogMin = batCellMin * batCellCount * batSensorAnalogDiv
	batSensorAnalogMax = batCellMax * batCellCount * batSensorAnalogDiv

	batSensorMin = batSensorAnalogMin / digitalToAnalog
	batSensorMax = batSensorAnalogMax / digitalToAnalog
)

type batterySensor struct {
	sensor
}

func (s batterySensor) cellVoltage() float64 {
	return s.reading.analog / (batSensorAnalogDiv * batCellCount)
}

func (s batterySensor) percent() float64 {
	return 100 * s.reading.fraction
}

func newBatterySensor() batterySensor {
	return batterySensor{
		sensor{
			pin: batterySensorPin,
			min: batSensorMin,
			max: batSensorMax,
		},
	}
}
