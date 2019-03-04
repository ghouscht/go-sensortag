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

func (h *Humidity) StartNotify(period []byte) (chan SensorEvent, error) {
	if err := h.setPeriod(period); err != nil {
		return nil, err
	}

	// enable the sensor
	if err := h.enable([]byte{0x1}); err != nil {
		return nil, err
	}

	return h.notify(h.convert)
}

// converts the raw humidity value into a percent value
func (h *Humidity) convert(data []byte) *[]SensorEvent {
	if len(data) != 4 {
		return nil
	}

	temp := -40 + ((165 * float64(binary.LittleEndian.Uint16(data[:2]))) / 65536.0)
	hum := float64(binary.LittleEndian.Uint16(data[2:])) * 100 / 65536.0

	return &[]SensorEvent{
		SensorEvent{
			Name:  "AmbientTemperature",
			Unit:  "°C",
			Value: temp,
		},
		SensorEvent{
			Name:  "Humidity",
			Unit:  "%",
			Value: hum,
		},
	}
}
