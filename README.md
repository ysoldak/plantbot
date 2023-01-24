# PlantBot

_Battery powered, automatic plant watering robot with soil moisture sensor and pump._

Robot can be controlled and monitored with [Blynk](https://blynk.io) service.

Robot's code is written on Go, compiled for Nano RP2040 Connect board with [TinyGo](https://tinygo.org/).

Inspired by:
- [Automatic Smart Plant Pot](https://github.com/DIY-Machines/SmartPlantV1)
- [Flaura Smart Plant Pot](https://github.com/FlauraPlantPot/Flaura)

## Soil Moisture Sensor

- DFRobot has  good quality [capacitive sensors](https://www.dfrobot.com/product-1385.html)
- Electrokit has interesting [version with I2C](https://www.electrokit.com/en/product/jordfuktighetssensor-kapacitiv-i2c/)

Before buying a capacitive soil moisture sensor, watch this [video on common failures in cheap copies](https://www.youtube.com/watch?v=IGP38bz-K48).

Each soil moisture sensor is unique and must be [calibrated](doc/calibration.md) before use.

## Battery

PlantBot expected to run from 2 Li-Ion batteries arranged in series (~7.4v).  
Use 75kOhm and 47kOhm resistors to make a voltage divider that brings battery voltage under 3.3v.

![](doc/voltage-divider-diagram.jpeg?raw=true)

## Pump

TODO

## Deep Sleep

I deep sleep, PlantBot consumes `20mW` power when running on [Arduino Nano RP2040 Connect](https://docs.arduino.cc/hardware/nano-rp2040-connect) board with power led (green) removed.
