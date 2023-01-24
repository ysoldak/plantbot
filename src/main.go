package main

import (
	"device/arm"
	"machine"
	"time"
)

var Version string

var led = machine.LED

func main() {

	// This wakes up system from deep sleep later.
	// We want schedulting to be the first thing to happen on [re]boot to minimise drift.
	machine.RTC.SetInterrupt(uint32(deepSleepDuration.Seconds()), false, nil)

	// Indicate wake up
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	led.High()

	// Init sensors
	batterySensor.configure()
	moistureSensor.configure()

	if moistureSensorMin == 0 || moistureSensorMax == 0 {
		calibrate()
	}

	// Connect to Wifi
	err := setupWifi(wifiSsid, wifiPass)
	if err != nil {
		println(err.Error())
		sleep()
		arm.SystemReset()
	}

	// Notify Blynk
	blynk := newBlynk()
	blynk.sendEvent("CONNECT")

	// Read sensors
	bat := batterySensor.read()
	moist := moistureSensor.read()

	// Some tracing
	now := time.Now().Format(time.RFC3339)
	println(now, bat.digital, bat.analog, bat.percent, bat.domain)
	println(now, moist.digital, moist.percent)

	// Update values at Blynk
	blynk.updateFloat(blynkNameBatterySensorPercent, bat.percent*100)        // percent charge left
	blynk.updateFloat(blynkNameBatterySensorDomain, bat.domain)              // inferred cell voltage
	blynk.updateFloat(blynkNameMoistureSensorPercent, 100*(1-moist.percent)) // 0 - wet, 100 - dry, so have to inverse

	// Deep sleep
	sleep()

	// Reboot after deep sleep since we can't yet wake up correctly
	arm.SystemReset()
}

func sleep() {
	println("Deep sleep")
	led.Low()      // turn off led to indicate sleep state
	stopWifinina() // stop ESP32 wifi co-proc
	time.Sleep(time.Second)
	machine.Sleep() // stop main proc
}

func calibrate() {
	var min, max uint16
	for {
		moist := moistureSensor.read()
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
