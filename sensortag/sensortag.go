package sensortag

import (
	"fmt"

	"github.com/ghouscht/go-sensortag/uuid"

	"github.com/muka/go-bluetooth/api"
	"github.com/pkg/errors"
)

// SensorTag is the data structure to represent a TI Sensortag.
type SensorTag struct {
	device      *api.Device
	Temperature Sensor
	Humidity    Sensor
	Barometer   Sensor
	Optical     Sensor
	IO          InputOutput
}

// New creates and initializes a new SensorTag instance
func New(dev *api.Device) (*SensorTag, error) {
	tag := new(SensorTag)

	if !dev.IsConnected() {
		return nil, fmt.Errorf("device is not connected")
	}

	tag.device = dev

	// Sensor initialization
	tempSensor, err := tag.NewSensorConfig(uuid.Temperature)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize temperature sensor")
	}
	tag.Temperature = NewTemperature(tempSensor)

	humSensor, err := tag.NewSensorConfig(uuid.Humidity)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize humidity sensor")
	}
	tag.Humidity = NewHumidity(humSensor)

	optSensor, err := tag.NewSensorConfig(uuid.Optical)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize optical sensor")
	}
	tag.Optical = NewOptical(optSensor)

	// i/o is a bit special...
	io, err := tag.NewIO(uuid.IO)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize i/o")
	}
	tag.IO = io

	return tag, nil
}
