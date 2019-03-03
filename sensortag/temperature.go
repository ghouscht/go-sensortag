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
	s.name = "AmbientTemperature"
	s.unit = "°C"

	return &Temperature{s}
}

func (t *Temperature) StartNotify(period []byte) (chan SensorEvent, error) {
	if err := t.setPeriod(period); err != nil {
		return nil, err
	}

	// enable the sensor
	if err := t.enable(); err != nil {
		return nil, err
	}

	return t.notify(ambientDegreeCelsius)
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
