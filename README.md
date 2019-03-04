# go-sensortag
## Setup
### dbus profiles
File: `/etc/dbus-1/system.d/dbus-go-bluetooth-dev.conf`
```
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
</busconfig>
```

File: `/etc/dbus-1/system.d/dbus-go-bluetooth-service.conf`
```
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
</busconfig>
```

### hciconfig acces
```bash
sudo setcap 'cap_net_raw,cap_net_admin+eip' `which hciconfig`
```
