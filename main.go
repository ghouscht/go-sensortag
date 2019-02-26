package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ghouscht/go-sensortag/sensortag"
	"github.com/ghouscht/go-sensortag/sensortag/conversion"
	"github.com/muka/go-bluetooth/api"
	"github.com/pkg/errors"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("missing tag address as argument")
	}
	tagAddr := os.Args[1]

	dev, err := api.GetDeviceByAddress(tagAddr)
	if err != nil {
		log.Fatal(err)
	}

	if dev == nil {
		panic("Device not found")
	}

	if err := dev.Connect(); err != nil {
		log.Fatal(errors.Wrap(err, "failed to connect"))
	}

	sensorTag, err := sensortag.New(dev)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to create sensortag instance"))
	}

	for {
		temp, err := sensorTag.Temperature.Read()
		if err != nil {
			log.Fatal(err)
		}

		hum, err := sensorTag.Humidity.Read()
		if err != nil {
			log.Fatal(err)
		}

		baro, err := sensorTag.Barometer.Read()
		if err != nil {
			log.Fatal(err)
		}

		optical, err := sensorTag.Optical.Read()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Temp: %f°C\n", conversion.AmbientDegreeCelsius(temp))
		fmt.Printf("Object Temp: %f°C\n", conversion.ObjectDegreeCelsius(temp))
		fmt.Printf("Humidity: %f%%\n", conversion.HumidityRelative(hum))
		fmt.Printf("Humidity Temp: %f°C\n", conversion.HumidityDegreeCelsius(hum))
		fmt.Printf("Pressure: %fhPa\n", conversion.BarometerPressure(baro))
		fmt.Printf("Optical: %#vLux\n", conversion.OpticalLux(optical))

		time.Sleep(3 * time.Second)
	}

	if err := dev.Disconnect(); err != nil {
		log.Fatal(errors.Wrap(err, "failed to disconnect"))
	}
}
