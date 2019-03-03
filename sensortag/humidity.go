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
	s.name = "Humidity"
	s.unit = "%"

	return &Humidity{s}
}

func (h *Humidity) StartNotify(period []byte) (chan SensorEvent, error) {
	if err := h.setPeriod(period); err != nil {
		return nil, err
	}

	// enable the sensor
	if err := h.enable(); err != nil {
		return nil, err
	}

	return h.notify(humidityRelative)
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
