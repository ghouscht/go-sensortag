package sensortag

import (
	"encoding/binary"
)

type Humidity struct {
	*sensorConfig
}

// empty declarationto check interface
var _ Sensor = &Humidity{}

func NewHumidity(s *sensorConfig) *Humidity {
	return &Humidity{s}
}

func (t *Humidity) SetPeriod(period []byte) error {
	return t.setPeriod(period)
}

func (t *Humidity) StartNotify() (chan float64, error) {
	return t.notify(humidityRelative, "service002a/char002b") //TODO
}

// converts the raw humidity value into a percent value
func humidityRelative(b []byte) float64 {
	if len(b) != 4 {
		return -1
	}

	return float64(binary.LittleEndian.Uint16(b[2:])) * 100 / 65536.0
}

/*
func HumidityDegreeCelsius(b []byte) float64 {
	if len(b) != 4 {
		return -1
	}

	return -40 + ((165 * float64(binary.LittleEndian.Uint16(b[:2]))) / 65536.0)
}
*/
