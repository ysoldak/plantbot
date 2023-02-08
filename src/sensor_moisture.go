package main

import "time"

type moistureSensor struct {
	sensor
}

func (s moistureSensor) percent() float64 {
	return 100 * (1 - s.reading.fraction) // 0 - wet, 100 - dry, so have to inverse
}

func (s moistureSensor) calibrate() {
	var min, max uint16
	for {
		moist := s.read()
		if min == 0 || min > moist.digital {
			min = moist.digital
		}
		if max == 0 || max < moist.digital {
			max = moist.digital
		}
		println("Moisture sensor (cur, min, max): ", moist.digital, min, max)
		time.Sleep(time.Second)
	}
}

func newMoistureSensor() moistureSensor {
	return moistureSensor{
		sensor{
			pin: moistureSensorPin,
			min: moistureSensorMin,
			max: moistureSensorMax,
		},
	}
}
