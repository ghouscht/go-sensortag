package sensortag

import (
	"encoding/binary"
)

type Temperature struct {
	*sensorConfig
}

// empty declarationto check interface
var _ Sensor = &Temperature{}

func NewTemperature(s *sensorConfig) *Temperature {
	return &Temperature{s}
}

func (t *Temperature) SetPeriod(period []byte) error {
	return t.setPeriod(period)
}

func (t *Temperature) StartNotify() (chan float64, error) {
	return t.notify(ambientDegreeCelsius, "service0022/char0023") //TODO
}

// converts the raw ambient temperature value into °C
func ambientDegreeCelsius(b []byte) float64 {
	if len(b) != 4 {
		return -1
	}
	return float64(binary.LittleEndian.Uint16(b[2:])) / 128.0
}

/*
func objectDegreeCelsius(b []byte) float64 {
	if len(b) != 4 {
		return -1
	}

	return float64(binary.LittleEndian.Uint16(b[:2])) / 128.0
}
*/
