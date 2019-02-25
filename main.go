package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ghouscht/go-sensortag/sensortag"
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

	temp, err := sensorTag.Temperature.Read()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Temp: %f°C\n", sensortag.DegreeCelsius(temp))

	/*
		if err := dev.Disconnect(); err != nil {
			log.Fatal(errors.Wrap(err, "failed to disconnect"))
		}
	*/
}
