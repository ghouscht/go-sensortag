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

func (t *Optical) SetPeriod(period []byte) error {
	return t.setPeriod(period)
}

func (t *Optical) StartNotify() (chan float64, error) {
	return t.notify(opticalLux, "service0042/char0043") //TODO: discover service
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
