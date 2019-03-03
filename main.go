package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ghouscht/go-sensortag/sensortag"
	"github.com/muka/go-bluetooth/api"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var (
	log *zap.SugaredLogger
)

func main() {
	// logger setup
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync() // flushes buffer, if any
	log = logger.Sugar()

	if len(os.Args) < 2 {
		log.Fatal("missing tag address as argument")
	}
	tagAddr := os.Args[1]

	stopC := make(chan os.Signal, 1)
	signal.Notify(stopC, os.Interrupt, syscall.SIGKILL, syscall.SIGTERM)

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
	period := []byte{0xFF} // == 2.5s

	tempC, err := sensorTag.Temperature.StartNotify(period)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to enable notifications for temperature sensor"))
	}

	humC, err := sensorTag.Humidity.StartNotify(period)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to enable notifications for humidity sensor"))
	}

	optC, err := sensorTag.Optical.StartNotify(period)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to enable notifications for humidity sensor"))
	}

	baroC, err := sensorTag.Barometer.StartNotify(period)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to enable notifications for barometer sensor"))
	}

main:
	for {
		select {
		case t := <-tempC:
			printData(t)
		case h := <-humC:
			printData(h)
		case o := <-optC:
			printData(o)
		case b := <-baroC:
			printData(b)
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

func printData(data interface{}) {
	if output, err := json.Marshal(data); err != nil {
		log.Error(err)
	} else {
		fmt.Println(string(output))
	}
}
