package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/ghouscht/go-sensortag/sensortag"
	"github.com/muka/go-bluetooth/api"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	connCheckInterval time.Duration
)

var connectCmd = &cobra.Command{
	Use:   "connect ADDRESS",
	Short: "connects to a sensortag by address",
	Long:  ``,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tagAddr := args[0]

		// channel to signal graceful stop
		stopC := make(chan os.Signal, 1)
		signal.Notify(stopC, os.Interrupt, syscall.SIGKILL, syscall.SIGTERM)

		// we never notice if we lose connection to the sensortag, setup a ticker
		// to regularly check if the connection is still ok
		ticker := time.NewTicker(connCheckInterval)

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
			case event, ok := <-events: // print events to stdout
				if output, err := json.Marshal(event); err != nil {
					log.Error(err)
				} else {
					fmt.Printf("%s\n", output)
				}

				if !ok { // events chan was closed, exit
					events = nil
					stopC <- syscall.SIGTERM
				}
			case <-ticker.C: // check connection
				// if the sensortag isn't connected anymore exit with error
				if !dev.IsConnected() {
					log.Fatalf("lost connection to sensortag %s", tagAddr)
				}
			case sig := <-stopC: // handle stop signal
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
	connectCmd.Flags().DurationVar(&connCheckInterval, "conn-check-interval", 10*time.Second, "interval to check if the sensortag is still connected")
	rootCmd.AddCommand(connectCmd)
}
