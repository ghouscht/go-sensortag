package sensortag

import (
	"fmt"

	"github.com/godbus/dbus"
	"github.com/muka/go-bluetooth/api"
	"github.com/muka/go-bluetooth/bluez/profile"
	"github.com/pkg/errors"
)

// SensorTag is the data structure to represent a TI Sensortag.
type SensorTag struct {
	device      *api.Device
	Temperature *Sensor
}

// New creates and initializes a new SensorTag instance
func New(dev *api.Device) (*SensorTag, error) {
	tag := new(SensorTag)

	if !dev.IsConnected() {
		return nil, fmt.Errorf("device is not connected")
	}

	tag.device = dev
	tempSensor, err := tag.newTemperatureSensor()
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize temperature sensor")
	}
	tag.Temperature = tempSensor

	return tag, nil
}

// Sensor represents a sensor gatt configuration
type Sensor struct {
	cfg    *profile.GattCharacteristic1
	data   *profile.GattCharacteristic1
	period *profile.GattCharacteristic1
}

func (s *Sensor) Read() ([]byte, error) {
	if err := s.enable(); err != nil {
		return nil, err
	}

	options := make(map[string]dbus.Variant)
	b, err := s.data.ReadValue(options)
	if err != nil {
		return nil, err
	}

	// disable sensor again to save energy
	s.disable()

	return b, nil
}

func (s *Sensor) enable() error {
	options := make(map[string]dbus.Variant)
	if err := s.cfg.WriteValue([]byte{1}, options); err != nil {
		return errors.Wrap(err, "failed to enable")
	}
	return nil
}

func (s *Sensor) disable() error {
	options := make(map[string]dbus.Variant)
	if err := s.cfg.WriteValue([]byte{0}, options); err != nil {
		return errors.Wrap(err, "failed to disable")
	}
	return nil
}
