package sensortag

import "encoding/binary"

// Barometer is a pressure sensor.
type Barometer struct {
	*sensorConfig
}

// empty declarationto check interface
var _ Sensor = &Barometer{}

// NewBarometer returns an initialized barometer sensor.
func NewBarometer(s *sensorConfig) *Barometer {
	return &Barometer{s}
}

// StartNotify enables and starts notification from the barometer.
func (b *Barometer) StartNotify(period []byte) (<-chan SensorEvent, error) {
	if err := b.setPeriod(period); err != nil {
		return nil, err
	}

	// enable the sensor
	if err := b.enable([]byte{0x1}); err != nil {
		return nil, err
	}

	return b.notify(b.convert)
}

// converts the raw barometer value into hPa
func (b *Barometer) convert(data []byte) *[]SensorEvent {
	if len(data) != 6 {
		return nil
	}

	raw := binary.LittleEndian.Uint32(data[2:])
	pressureMask := (int(uint32(raw)) >> 8) & 0x00ffffff
	pressure := float64(pressureMask) / 100.0

	return &[]SensorEvent{
		SensorEvent{
			Name:  "Barometer",
			Unit:  "hPa",
			Value: &pressure,
		},
	}
}
