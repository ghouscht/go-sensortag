package sensortag

import (
	"github.com/godbus/dbus"
	"github.com/muka/go-bluetooth/bluez/profile"
	"github.com/pkg/errors"
)

// Sensor represents a sensor gatt configuration
type Sensor struct {
	cfg    *profile.GattCharacteristic1
	data   *profile.GattCharacteristic1
	period *profile.GattCharacteristic1
}

// NewSensor returns an initialized sensor.
func (tag *SensorTag) NewSensor(uuid UUID) (*Sensor, error) {
	cfg, err := tag.device.GetCharByUUID(uuid.Config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get config characteristic")
	}

	data, err := tag.device.GetCharByUUID(uuid.Data)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get data characteristic")
	}

	period, err := tag.device.GetCharByUUID(uuid.Period)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get period characteristic")
	}

	return &Sensor{
		cfg:    cfg,
		data:   data,
		period: period,
	}, nil
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
