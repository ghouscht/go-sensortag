package sensortag

import (
	"encoding/binary"

	"github.com/pkg/errors"
)

func DegreeCelsius(b []byte) float64 {
	return float64(binary.LittleEndian.Uint16(b[2:])) / 128.0
}

func (tag *SensorTag) newTemperatureSensor() (*Sensor, error) {
	cfg, err := tag.device.GetCharByUUID(uuid["TemperatureConfig"])
	if err != nil {
		return nil, errors.Wrap(err, "failed to get TemperatureConfig characteristic")
	}

	data, err := tag.device.GetCharByUUID(uuid["TemperatureData"])
	if err != nil {
		return nil, errors.Wrap(err, "failed to get TemperatureData characteristic")
	}

	period, err := tag.device.GetCharByUUID(uuid["TemperaturePeriod"])
	if err != nil {
		return nil, errors.Wrap(err, "failed to get TemperaturePeriod characteristic")
	}

	return &Sensor{
		cfg:    cfg,
		data:   data,
		period: period,
	}, nil
}
