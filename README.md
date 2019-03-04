# go-sensortag
## Setup
A few simple seteps need to be done, prior to run `go-sensortag`.
### dbus profiles
Two dbus profiles need to be created:  

File: `/etc/dbus-1/system.d/dbus-go-bluetooth-dev.conf` owner: root group: root mode: 0644  
Content:
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

File: `/etc/dbus-1/system.d/dbus-go-bluetooth-service.conf` owner: root group: root mode: 0644  
Content:
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
Enable access to the hciconfig binary for all users:
```bash
sudo setcap 'cap_net_raw,cap_net_admin+eip' `which hciconfig`
```

## Run
First you need to find the MAC address of your sensortag:
```bash
sudo hcitool lescan
```

With the MAC address you should now be able to run `go-sensortag`:
```bash
./go-sensortag 54:6C:0E:FF:FF:FF
```
