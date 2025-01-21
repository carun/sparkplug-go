package main

import (
	"log"

	"google.golang.org/protobuf/proto" // Correct import for proto.Marshal

	myProto "sparkbuf/proto"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Configuring MQTT
func connectToMqttBroker() mqtt.Client {
	opts := mqtt.NewClientOptions().AddBroker("mqtt://127.0.0.1:1883")
	opts.SetClientID("myid")
	opts.SetUsername("")
	opts.SetPassword("")

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	return client
}

// Sending a Message via MQTT
func sendMsg(client mqtt.Client, data []byte) {

	// Publish the message to the topic
	if token := client.Publish("your/topic", 0, false, data); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	} else {
		log.Print("Message published successfully!")
	}
}

func main() {
	mqttClient := connectToMqttBroker()
	defer mqttClient.Disconnect(250)

	// Create a Protobuf message
	nbirth := &myProto.Payload{}

	// Serialize the Protobuf message
	data, err := proto.Marshal(nbirth)
	if err != nil {
		log.Fatalf("Failed to marshal message: %v", err)
	}

	sendMsg(mqttClient, data)
}
