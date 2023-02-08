package main

import (
	"errors"
	"machine"
	"time"
)

var pumpPwmGroup = machine.PWM7

func init() {
	// See https: //docs.arduino.cc/static/97fbbd7b8b0d0f12efd7fc3e3cd5de35/ABX00053-full-pinout.pdf
	switch pumpPwmPin {
	case machine.D2:
		pumpPwmGroup = machine.PWM4
	case machine.D3:
		pumpPwmGroup = machine.PWM7
	case machine.D4, machine.D5:
		pumpPwmGroup = machine.PWM0
	case machine.D6, machine.D7:
		pumpPwmGroup = machine.PWM1
	case machine.D8, machine.D9:
		pumpPwmGroup = machine.PWM2
	default:
		pumpPwmGroup = nil
	}
}

type pump struct {
	powerPin machine.Pin
	pwmPin   machine.Pin

	delay    time.Duration
	duration time.Duration

	period  uint32
	dutyOff uint32
	dutyRun uint32

	channel  uint8
	valueOff uint32
	valueRun uint32
}

func newPump() pump {
	return pump{
		powerPin: pumpPowerPin,
		pwmPin:   pumpPwmPin,
		delay:    pumpPwmDelay,
		duration: pumpRunDuration,
		period:   pumpPwmPeriod,
		dutyOff:  pumpPwmDutyOff,
		dutyRun:  pumpPwmDutyRun,
	}
}

func (p *pump) configure() error {

	if p.powerPin != machine.NoPin {
		p.powerPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
		p.powerPin.Low()
	}

	if pumpPwmGroup == nil {
		return nil
	}

	err := pumpPwmGroup.Configure(machine.PWMConfig{
		Period: uint64(p.period) * 1000, // period in nanoseconds, we have it in microseconds
	})
	if err != nil {
		return errors.New("failed to configure PWM")
	}

	mult := pumpPwmGroup.Top() / p.period // this value represents 1us in PWM
	p.valueOff = p.dutyOff * mult
	p.valueRun = p.dutyRun * mult

	// Configure the two channels we'll use as outputs.
	p.channel, err = pumpPwmGroup.Channel(p.pwmPin)
	if err != nil {
		return errors.New("failed to configure PWM channel")
	}
	pumpPwmGroup.Set(p.channel, p.valueOff)
	return nil
}

func (p *pump) work() {
	if p.powerPin != machine.NoPin {
		p.powerPin.High() // power on
	}

	if pumpPwmGroup != nil {
		time.Sleep(p.delay)                     // init
		pumpPwmGroup.Set(p.channel, p.valueRun) // start
	}

	time.Sleep(p.duration) // run

	if pumpPwmGroup != nil {
		pumpPwmGroup.Set(p.channel, p.valueOff) // stop
	}

	if p.powerPin != machine.NoPin {
		p.powerPin.Low() // power off
	}
}
