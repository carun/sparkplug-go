# Prerequisites

1. Mosquitto MQTT broker
1. [MQTT Explorer](https://mqtt-explorer.com/)
1. Go 1.23.5

```bash
sudo apt install g++ protobuf-compiler protoc-gen-go mosquitto
```

## Build and run

```bash
make clean && make
./build/sp-pub
```
