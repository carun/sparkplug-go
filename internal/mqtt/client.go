package mqtt

import (
	"fmt"
	"time"

	"sparkplug-go/internal/config"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type Client struct {
	client MQTT.Client
	config *config.Config
}

func NewClient(cfg *config.Config) (*Client, error) {
	opts := MQTT.NewClientOptions()
	broker := fmt.Sprintf("tcp://%s:%d", cfg.MQTT.Broker, cfg.MQTT.Port)
	opts.AddBroker(broker)
	opts.SetClientID(cfg.MQTT.ClientID)

	if cfg.MQTT.Username != "" {
		opts.SetUsername(cfg.MQTT.Username)
		opts.SetPassword(cfg.MQTT.Password)
	}

	opts.SetKeepAlive(60 * time.Second)
	opts.SetDefaultPublishHandler(defaultMessageHandler)
	opts.SetConnectionLostHandler(connectionLostHandler)
	opts.SetOnConnectHandler(connectHandler)

	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	return &Client{
		client: client,
		config: cfg,
	}, nil
}

func (c *Client) Publish(topic string, payload []byte) error {
	token := c.client.Publish(topic, 0, false, payload)
	token.Wait()
	return token.Error()
}

func (c *Client) Disconnect() {
	c.client.Disconnect(250)
}

func defaultMessageHandler(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("Received message on topic: %s\nMessage: %s\n", msg.Topic(), msg.Payload())
}

func connectHandler(client MQTT.Client) {
	fmt.Println("Connected to MQTT Broker")
}

func connectionLostHandler(client MQTT.Client, err error) {
	fmt.Printf("Connection lost: %v\n", err)
}
