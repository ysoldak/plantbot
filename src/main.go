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
	if sleepDeep {
		machine.RTC.SetInterrupt(uint32(sleepDuration.Seconds()), false, nil)
	}

	// Indicate wake up
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	led.High()

	// Init sensors
	batterySensor := newBatterySensor()
	batterySensor.configure()

	moistureSensor := newMoistureSensor()
	moistureSensor.configure()

	// Calibrate moisture sensor, optional
	if moistureSensorMin == 0 || moistureSensorMax == 0 {
		moistureSensor.calibrate() // locks forever; shall adjust config and re-flash
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

	// Init pump
	pump := newPump()
	pump.configure()

	for {
		led.High()

		// Read sensors
		bat := batterySensor.read()
		batPercent := batterySensor.percent()
		batCellVolt := batterySensor.cellVoltage()

		moist := moistureSensor.read()
		moistPercent := moistureSensor.percent()

		// Some tracing
		now := time.Now().Format(time.RFC3339)
		println(now, bat.digital, bat.analog, batPercent, batCellVolt)
		println(now, moist.digital, moistPercent)

		// Update values at Blynk
		blynk.updateFloat(blynkNameBatterySensorPercent, batPercent)    // percent charge left
		blynk.updateFloat(blynkNameBatterySensorCellVolt, batCellVolt)  // inferred cell voltage
		blynk.updateFloat(blynkNameMoistureSensorPercent, moistPercent) // 0% - air, 100% - water

		// Pump water when needed
		if moistPercent < pumpThreshold {
			// do pumping
			pump.work()
			// notify Blynk
			blynk.sendEvent("PUMP")
		}

		// Either light or deep sleep
		sleep()
	}
}

func sleep() {
	led.Low() // turn off led to indicate sleep state

	if !sleepDeep {
		println("Light sleep for 5 seconds")
		time.Sleep(sleepDuration)
		return
	}

	println("Deep sleep")
	stopWifinina() // stop ESP32 wifi co-proc
	time.Sleep(time.Second)
	machine.Sleep()   // stop main proc
	arm.SystemReset() // reboot after deep sleep since we can't yet wake up correctly
}
