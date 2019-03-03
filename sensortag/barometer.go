package sensortag

import "encoding/binary"

// Barometer is a pressure sensor.
type Barometer struct {
	*sensorConfig
}

// empty declarationto check interface
var _ Sensor = &Barometer{}

func NewBarometer(s *sensorConfig) *Barometer {
	s.name = "Pressure"
	s.unit = "hPa"

	return &Barometer{s}
}

func (b *Barometer) StartNotify(period []byte) (chan SensorEvent, error) {
	if err := b.setPeriod(period); err != nil {
		return nil, err
	}

	// enable the sensor
	if err := b.enable(); err != nil {
		return nil, err
	}

	return b.notify(barometerPressure)
}

// converts the raw barometer value into hPa
func barometerPressure(b []byte) float64 {
	if len(b) != 6 {
		return -1
	}

	raw := binary.LittleEndian.Uint32(b[2:])
	pressureMask := (int(uint32(raw)) >> 8) & 0x00ffffff

	return float64(pressureMask) / 100.0
}
