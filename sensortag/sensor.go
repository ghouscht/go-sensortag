package sensortag

import (
	"github.com/ghouscht/go-sensortag/uuid"
	"github.com/godbus/dbus"
	"github.com/muka/go-bluetooth/bluez/profile"
	"github.com/pkg/errors"
)

type sensorConfig struct {
	name string
	unit string

	cfg    *profile.GattCharacteristic1
	data   *profile.GattCharacteristic1
	period *profile.GattCharacteristic1
}

// Sensor is an interface for sensortag sensors.
type Sensor interface {
	StartNotify([]byte) (chan SensorEvent, error)
	//EnableNotify([]byte) error
}

// SensorEvent ...
type SensorEvent struct {
	Name  string  `json:"name"`
	Unit  string  `json:"unit"`
	Value float64 `json:"value"`
}

type conversionFunc func([]byte) float64

// NewSensorConfig returns a pointer to an initialized sensor config.
func (tag *SensorTag) NewSensorConfig(uuid uuid.UUID) (*sensorConfig, error) {
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

	return &sensorConfig{
		cfg:    cfg,
		data:   data,
		period: period,
	}, nil
}

func (s *sensorConfig) notify(conversion conversionFunc) (chan SensorEvent, error) {
	dataC := make(chan SensorEvent)

	if err := s.data.StartNotify(); err != nil {
		return nil, errors.Wrap(err, "failed to start notifications")
	}

	// get the data channel
	events, err := s.data.Register()
	if err != nil {
		return nil, err
	}

	go func() {
		for event := range events {
			if event == nil {
				// terminate
				close(dataC)
				return
			}

			// the events channel always returns all notifications (also those from sensors handled
			// by a different goroutine) that's why we need to match against the service name...
			if string(event.Path) != s.data.Path {
				continue
			}

			if len(event.Body) < 1 {
				continue
			}

			body, ok := event.Body[1].(map[string]dbus.Variant)
			if !ok {
				continue
			}

			if _, ok := body["Value"]; !ok {
				continue
			}

			rawData := body["Value"].Value().([]byte)

			dataC <- SensorEvent{
				Name:  s.name,
				Unit:  s.unit,
				Value: conversion(rawData),
			}
		}
	}()

	// start notifying
	return dataC, nil
}

func (s *sensorConfig) setPeriod(period []byte) error {
	options := make(map[string]dbus.Variant)
	if err := s.period.WriteValue(period, options); err != nil {
		return errors.Wrap(err, "failed to set period")
	}
	return nil
}

func (s *sensorConfig) enable() error {
	options := make(map[string]dbus.Variant)
	if err := s.cfg.WriteValue([]byte{0x1}, options); err != nil {
		return errors.Wrap(err, "failed to enable")
	}
	return nil
}

func (s *sensorConfig) disable() error {
	options := make(map[string]dbus.Variant)
	if err := s.cfg.WriteValue([]byte{0x0}, options); err != nil {
		return errors.Wrap(err, "failed to disable")
	}
	return nil
}
