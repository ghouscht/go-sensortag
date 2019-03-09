# go-sensortag
## Run
Make sure bluetooth is running:
```bash
sudo service bluetooth status
```

First you need to find the MAC address of your sensortag (assuming `hci0` as adapter):
```bash
./go-sensortag discover hci0
```

With the MAC address you should now be able to connect to your sensortag:
```bash
./go-sensortag connect 54:6C:0E:FF:FF:FF
```
