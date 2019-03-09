package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/muka/go-bluetooth/api"
	"github.com/muka/go-bluetooth/emitter"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var discoverCmd = &cobra.Command{
	Use:   "discover ADAPTER",
	Short: "discover bluetooth devices",
	Long:  ``,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		adapterID := args[0]

		// stop signal chan
		stopC := make(chan os.Signal, 1)
		signal.Notify(stopC, os.Interrupt, syscall.SIGKILL, syscall.SIGTERM)

		//clean up connection on exit
		defer api.Exit()

		devices, err := api.GetDevices()
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		fmt.Println("Cached devices:")
		for _, dev := range devices {
			showDeviceInfo(&dev)
		}

		fmt.Println("Discovered devices:")
		err = discoverDevices(adapterID)
		if err != nil {
			log.Fatal(err)
		}

		// wait for stop signal
		<-stopC
	},
}

func discoverDevices(adapterID string) error {
	if err := api.StartDiscovery(); err != nil {
		return errors.Wrap(err, "failed to start discovery")
	}

	err := api.On("discovery", emitter.NewCallback(func(ev emitter.Event) {
		discoveryEvent := ev.GetData().(api.DiscoveredDeviceEvent)
		dev := discoveryEvent.Device
		showDeviceInfo(dev)
	}))

	return err
}

func showDeviceInfo(dev *api.Device) {
	if dev == nil {
		return
	}
	props, err := dev.GetProperties()
	if err != nil {
		log.Errorf("failed to get properties from device %s: %s", dev.Path, err)
		return
	}
	fmt.Printf("%s rssi=%d %s\n", props.Address, props.RSSI, props.Name)
}

func init() {
	rootCmd.AddCommand(discoverCmd)
}
