package sensortag

import (
	"github.com/godbus/dbus"
	"github.com/muka/go-bluetooth/bluez/profile"
	"github.com/pkg/errors"
)

type InputOutput struct {
	cfg  *profile.GattCharacteristic1
	data *profile.GattCharacteristic1
}

func (tag *SensorTag) NewIO(uuid UUID) (*InputOutput, error) {
	cfg, err := tag.device.GetCharByUUID(uuid.Config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get config characteristic")
	}

	data, err := tag.device.GetCharByUUID(uuid.Data)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get data characteristic")
	}

	return &InputOutput{
		cfg:  cfg,
		data: data,
	}, nil
}

func (io *InputOutput) enable() error {
	options := make(map[string]dbus.Variant)
	if err := io.cfg.WriteValue([]byte{1}, options); err != nil {
		return errors.Wrap(err, "failed to enable")
	}
	return nil
}

func (io *InputOutput) disable() error {
	options := make(map[string]dbus.Variant)
	if err := io.cfg.WriteValue([]byte{0}, options); err != nil {
		return errors.Wrap(err, "failed to enable")
	}
	return nil
}

func (io *InputOutput) Write(data []byte) error {
	options := make(map[string]dbus.Variant)
	if err := io.data.WriteValue(data, options); err != nil {
		return errors.Wrap(err, "failed to write data")
	}

	if err := io.enable(); err != nil {
		return err
	}

	return nil
}

func (io *InputOutput) Read() ([]byte, error) {
	options := make(map[string]dbus.Variant)
	b, err := io.data.ReadValue(options)
	if err != nil {
		return nil, err
	}
	return b, nil
}
