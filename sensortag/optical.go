package sensortag

import (
	"encoding/binary"
	"math"
)

type Optical struct {
	*sensorConfig
}

// empty declarationto check interface
var _ Sensor = &Optical{}

func NewOptical(s *sensorConfig) *Optical {
	return &Optical{s}
}

func (o *Optical) StartNotify(period []byte) (<-chan SensorEvent, error) {
	if err := o.setPeriod(period); err != nil {
		return nil, err
	}

	// enable the sensor
	if err := o.enable([]byte{0x1}); err != nil {
		return nil, err
	}

	return o.notify(o.convert)
}

// converts the raw optical value into a lux value
func (o *Optical) convert(data []byte) *[]SensorEvent {
	if len(data) != 2 {
		return nil
	}
	raw := binary.LittleEndian.Uint16(data)

	exponent := float64((int(raw) & 0xF000) >> 12)
	mantissa := float64((int(raw) & 0x0FFF))
	lux := mantissa * math.Pow(2, exponent) / 100.0

	return &[]SensorEvent{
		SensorEvent{
			Name:  "AmbientLight",
			Unit:  "Lux",
			Value: lux,
		},
	}
}
