package sensortag

import (
	"fmt"

	"github.com/muka/go-bluetooth/api"
	"github.com/pkg/errors"
)

// SensorTag is the data structure to represent a TI Sensortag.
type SensorTag struct {
	device      *api.Device
	Temperature *Sensor
	Humidity    *Sensor
	Barometer   *Sensor
	Optical     *Sensor
	IO          *InputOutput
}

// New creates and initializes a new SensorTag instance
func New(dev *api.Device) (*SensorTag, error) {
	tag := new(SensorTag)

	if !dev.IsConnected() {
		return nil, fmt.Errorf("device is not connected")
	}

	tag.device = dev
	tempSensor, err := tag.NewSensor(temperature)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize temperature sensor")
	}
	tag.Temperature = tempSensor

	humiditySensor, err := tag.NewSensor(humidity)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize humidity sensor")
	}
	tag.Humidity = humiditySensor

	baroSensor, err := tag.NewSensor(barometer)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize barometer sensor")
	}
	tag.Barometer = baroSensor

	opticalSensor, err := tag.NewSensor(optical)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize optical sensor")
	}
	tag.Optical = opticalSensor

	io, err := tag.NewIO(io)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize io")
	}
	tag.IO = io

	return tag, nil
}
