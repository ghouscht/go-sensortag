package sensortag

import (
	"encoding/binary"
)

type IRTemperature struct {
	*sensorConfig
}

// empty declarationto check interface
var _ Sensor = &IRTemperature{}

func NewIRTemperature(s *sensorConfig) *IRTemperature {
	return &IRTemperature{s}
}

func (t *IRTemperature) StartNotify(period []byte) (<-chan SensorEvent, error) {
	if err := t.setPeriod(period); err != nil {
		return nil, err
	}

	// enable the sensor
	if err := t.enable([]byte{0x1}); err != nil {
		return nil, err
	}

	return t.notify(t.convert)
}

// converts the raw object temperature value into °C
func (t *IRTemperature) convert(data []byte) *[]SensorEvent {
	if len(data) != 4 {
		return nil
	}

	temp := float64(binary.LittleEndian.Uint16(data[2:])) / 128.0

	return &[]SensorEvent{
		SensorEvent{
			Name:  "IRTemperature",
			Unit:  "°C",
			Value: temp,
		},
	}
}
