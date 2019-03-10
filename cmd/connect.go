package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/ghouscht/go-sensortag/sensortag"
	"github.com/muka/go-bluetooth/api"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var connectCmd = &cobra.Command{
	Use:   "connect ADDRESS",
	Short: "connects to a sensortag by address",
	Long:  ``,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tagAddr := args[0]

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
			log.Fatal(errors.Wrap(err, "failed to enable notifications for optical sensor"))
		}

		baroC, err := sensorTag.Barometer.StartNotify(period)
		if err != nil {
			log.Fatal(errors.Wrap(err, "failed to enable notifications for barometer sensor"))
		}

		moveC, err := sensorTag.Movement.StartNotify(period)
		if err != nil {
			log.Fatal(errors.Wrap(err, "failed to enable notifications for movement sensor"))
		}

		events := merge(irtempC, humC, optC, baroC, moveC)
	main:
		for {
			select {
			case event, ok := <-events:
				if output, err := json.Marshal(event); err != nil {
					log.Error(err)
				} else {
					fmt.Printf("%s\n", output)
				}

				if !ok { // events chan was closed, exit
					events = nil
					stopC <- syscall.SIGTERM
				}
			case sig := <-stopC:
				log.Infow(
					"disconnecting...",
					"tag", tagAddr,
					"signal", sig,
				)
				if err := dev.Disconnect(); err != nil {
					log.Fatal(errors.Wrap(err, "failed to disconnect"))
				}

				log.Infow(
					"disconnected",
					"tag", tagAddr,
				)
				break main
			}
		}
	},
}

func merge(events ...<-chan sensortag.SensorEvent) <-chan sensortag.SensorEvent {
	out := make(chan sensortag.SensorEvent)

	var wg sync.WaitGroup
	wg.Add(len(events))

	for _, event := range events {
		go func(event <-chan sensortag.SensorEvent) {
			for v := range event {
				out <- v
			}
			wg.Done()
		}(event)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func init() {
	rootCmd.AddCommand(connectCmd)
}
