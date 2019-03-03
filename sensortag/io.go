package sensortag

import (
	"github.com/ghouscht/go-sensortag/uuid"
	"github.com/godbus/dbus"
	"github.com/pkg/errors"
)

type InputOutput struct {
	*sensorConfig
}

// NewIO returns an InputOutput object to control a SensorTags LEDs. buzzer etc.
func (tag *SensorTag) NewIO(uuid uuid.UUID) (InputOutput, error) {
	io := InputOutput{}

	cfg, err := tag.device.GetCharByUUID(uuid.Config)
	if err != nil {
		return io, errors.Wrap(err, "failed to get config characteristic")
	}

	data, err := tag.device.GetCharByUUID(uuid.Data)
	if err != nil {
		return io, errors.Wrap(err, "failed to get data characteristic")
	}

	configuration := &sensorConfig{
		cfg:  cfg,
		data: data,
	}

	return InputOutput{configuration}, nil
}

// Write enables/disables I/Os
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
