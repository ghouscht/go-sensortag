package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/ghouscht/go-sensortag/sensortag"
	"github.com/muka/go-bluetooth/api"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func main() {
	// logger setup
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync() // flushes buffer, if any
	log := logger.Sugar()

	if len(os.Args) < 2 {
		log.Fatal("missing tag address as argument")
	}
	tagAddr := os.Args[1]

	stopC := make(chan os.Signal, 1)
	signal.Notify(stopC, os.Interrupt)

	log.Infow(
		"connecting...",
		"tag", tagAddr,
	)
	dev, err := api.GetDeviceByAddress(tagAddr)
	if err != nil {
		log.Fatal(err)
	}

	if dev == nil {
		// TODO retry getDeviceByAddress
		log.Fatal("device not found")
	}

	if err := dev.Connect(); err != nil {
		// TODO retry connect
		log.Fatal(errors.Wrap(err, "failed to connect"))
	}

	log.Infow(
		"connected!",
		"tag", tagAddr,
	)

	sensorTag, err := sensortag.New(dev)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to create sensortag instance"))
	}

	// enable the green LED, to signal connection
	if err := sensorTag.IO.Write([]byte{0x02}); err != nil {
		log.Errorw("Error failed to enable IOs, %s\n", err)
	}

	// Sensors
	if err := sensorTag.Temperature.SetPeriod([]byte{0xFF}); err != nil {
		log.Fatal(errors.Wrap(err, "failed to set period for temperature reading"))
	}

	tempC, err := sensorTag.Temperature.StartNotify()
	if err != nil {
		panic(err)
	}

	if err := sensorTag.Humidity.SetPeriod([]byte{0xFF}); err != nil {
		log.Fatal(errors.Wrap(err, "failed to set period for humidity reading"))
	}

	humC, err := sensorTag.Humidity.StartNotify()
	if err != nil {
		panic(err)
	}

	if err := sensorTag.Optical.SetPeriod([]byte{0xFF}); err != nil {
		log.Fatal(errors.Wrap(err, "failed to set period for humidity reading"))
	}

	optC, err := sensorTag.Optical.StartNotify()
	if err != nil {
		panic(err)
	}

main:
	for {
		select {
		case t := <-tempC:
			fmt.Printf("The current ambient temperature is: %f °C\n", t)
		case h := <-humC:
			fmt.Printf("The current humidity is: %f %%\n", h)
		case o := <-optC:
			fmt.Printf("The current ambient light is: %f Lux\n", o)
		case <-stopC:
			log.Infow(
				"disconnecting...",
				"tag", tagAddr,
			)

			if err := dev.Disconnect(); err != nil {
				log.Fatal(errors.Wrap(err, "failed to disconnect"))
			}
			break main
		}
	}
}
