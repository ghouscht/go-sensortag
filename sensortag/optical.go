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
	s.name = "AmbientLight"
	s.unit = "Lux"

	return &Optical{s}
}

func (o *Optical) StartNotify(period []byte) (chan SensorEvent, error) {
	if err := o.setPeriod(period); err != nil {
		return nil, err
	}

	// enable the sensor
	if err := o.enable(); err != nil {
		return nil, err
	}

	return o.notify(opticalLux)
}

// converts the raw optical value into a lux value
func opticalLux(b []byte) float64 {
	if len(b) != 2 {
		return -1
	}
	raw := binary.LittleEndian.Uint16(b)

	exponent := float64((int(raw) & 0xF000) >> 12)
	mantissa := float64((int(raw) & 0x0FFF))

	return mantissa * math.Pow(2, exponent) / 100.0
}
