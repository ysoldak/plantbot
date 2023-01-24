# Calibrate Soil Moisture Sensor

1. Set `moistureSensorMin` or `moistureSensorMax` to `0`, flash and connect with serial console.
2. Leave the sensor in open air, record the value, this will be `moistureSensorMax`
3. Put the sensor in a glass with water, record the value, the will be `moistureSensorMin`
4. Set recorded values in `config.go`.
