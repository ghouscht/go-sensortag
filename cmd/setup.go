package cmd

import (
	"os"
	"os/user"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var goBluetoothDev = `
<!DOCTYPE busconfig PUBLIC "-//freedesktop//DTD D-BUS Bus Configuration 1.0//EN" "http://www.freedesktop.org/standards/dbus/1.0/busconfig.dtd">
<busconfig>
    <policy group="sudo">
        <allow own="go.bluetooth"/>
        <allow own_prefix="go.bluetooth"/>
        <allow send_destination="go.bluetooth" send_interface="go.bluetooth"/>
        <allow send_destination="go.bluetooth" send_interface="org.freedesktop.DBus.Introspectable"/>
        <allow send_destination="go.bluetooth" send_interface="org.freedesktop.DBus.Properties"/>
        <allow send_destination="go.bluetooth" send_interface="org.freedesktop.DBus.ObjectManager"/>
        <allow send_destination="go.bluetooth"/>

        <allow own="org.bluez"/>
        <allow own_prefix="org.bluez"/>
        <allow send_destination="org.bluez" send_interface="org.bluez"/>
        <allow send_destination="org.bluez" send_interface="org.freedesktop.DBus.Introspectable"/>
        <allow send_destination="org.bluez" send_interface="org.freedesktop.DBus.Properties"/>
        <allow send_destination="org.bluez" send_interface="org.freedesktop.DBus.ObjectManager"/>
        <allow send_destination="org.bluez"/>

        <allow send_interface="org.freedesktop.DBus.ObjectManager"/>
        <allow send_interface="org.freedesktop.DBus.Properties"/>
        <allow send_interface="org.freedesktop.DBus.Introspectable"/>
        <allow send_interface="org.bluez.GattCharacteristic1"/>
        <allow send_interface="org.bluez.GattDescriptor1"/>
        <allow send_interface="org.bluez.GattService1"/>
        <allow send_interface="org.bluez.GattManager1"/>
    </policy>
    <policy context="default">
    </policy>
</busconfig>`

var goBluetoothService = `
<!DOCTYPE busconfig PUBLIC "-//freedesktop//DTD D-BUS Bus Configuration 1.0//EN" "http://www.freedesktop.org/standards/dbus/1.0/busconfig.dtd">
<busconfig>
  <policy user="root">
    <allow own="go.bluetooth"/>
    <allow send_destination="go.bluetooth"/>
    <allow send_destination="org.bluez"/>
  </policy>
  <policy at_console="true">
    <allow own="go.bluetooth"/>
    <allow send_destination="go.bluetooth"/>
    <allow send_destination="org.bluez"/>
  </policy>
  <policy context="default">
    <deny send_destination="go.bluetooth"/>
  </policy>
</busconfig>`

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "setup dbus profiles",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if user, err := user.Current(); err != nil || user.Username != "root" {
			log.Fatal("setup cmd must be run as root")
		}

		profile, err := os.Create("/etc/dbus-1/system.d/dbus-go-bluetooth-dev.conf")
		if err != nil {
			log.Fatal(errors.Wrap(err, "failed to create dbus profile"))
		}
		defer profile.Close()

		if written, err := profile.Write([]byte(goBluetoothDev)); err != nil || written != len([]byte(goBluetoothDev)) {
			log.Fatal(errors.Wrap(err, "failed to write dbus profile"))
		}
		if err := profile.Chmod(0644); err != nil {
			log.Fatal(errors.Wrap(err, "failed to change mode of dbus profile"))
		}

		svc, err := os.Create("/etc/dbus-1/system.d/dbus-go-bluetooth-service.conf")
		if err != nil {
			log.Fatal(errors.Wrap(err, "failed to create dbus svc"))
		}
		defer svc.Close()

		if written, err := svc.Write([]byte(goBluetoothService)); err != nil || written != len([]byte(goBluetoothService)) {
			log.Fatal(errors.Wrap(err, "failed to write dbus svc"))
		}
		if err := svc.Chmod(0644); err != nil {
			log.Fatal(errors.Wrap(err, "failed to change mode of dbus svc"))
		}

		log.Info("setup completed")
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
}
