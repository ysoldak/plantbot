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
	machine.RTC.SetInterrupt(deepSleepDuration.Nanoseconds()+1, false, nil)

	// Indicate wake up
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	led.High()

	// Init sensors
	batterySensor.configure()
	moistureSensor.configure()

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
	blynk.updateFloat(blynkNameBatterySensorPercent, bat.percent*100)
	blynk.updateFloat(blynkNameBatterySensorDomain, bat.domain)
	blynk.updateFloat(blynkNameMoistureSensorPercent, moist.percent*100)

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
