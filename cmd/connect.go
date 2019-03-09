package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ghouscht/go-sensortag/sensortag"
	"github.com/muka/go-bluetooth/api"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	tagAddr string
)

var connectCmd = &cobra.Command{
	Use:   "connect ADDRESS",
	Short: "connects to a sensortag by address",
	Long:  ``,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		defer api.Exit()

		tagAddr = args[0]

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
			log.Fatal("device not found")
		}

		if err := dev.Connect(); err != nil {
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

		irtempC, err := sensorTag.IRTemperature.StartNotify(period)
		if err != nil {
			log.Fatal(errors.Wrap(err, "failed to enable notifications for ir temperature sensor"))
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

		moveC, err := sensorTag.Movement.StartNotify(period)
		if err != nil {
			log.Fatal(errors.Wrap(err, "failed to enable notifications for movement sensor"))
		}

		// blocks until signal from stopC
		dataPrinter(stopC, dev, humC, optC, baroC, irtempC, moveC)

	},
}

func dataPrinter(stop chan os.Signal, dev *api.Device, events ...chan sensortag.SensorEvent) {
	for {
		for _, event := range events {
			// setup a go routine to read each chan and print it's data to stdout
			go func(e chan sensortag.SensorEvent) {
				for {
					data := <-e
					if output, err := json.Marshal(data); err != nil {
						log.Error(err)
					} else {
						fmt.Printf("%s\n", output)
					}

				}
			}(event)
		}

		<-stop // block until we receive a stop signal
		log.Infow(
			"disconnecting...",
			"tag", tagAddr,
		)
		//TODO: clean stop

		if err := dev.Disconnect(); err != nil {
			log.Fatal(errors.Wrap(err, "failed to disconnect"))
		}
		break
	}
}

func init() {
	rootCmd.AddCommand(connectCmd)
}
