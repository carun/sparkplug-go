package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	MQTT struct {
		Broker   string
		Port     int
		ClientID string
		Username string
		Password string
	}
	Sparkplug struct {
		GroupID     string
		EdgeNodeID  string
		DeviceID    string
		ScadaHostID string
	}
}

func LoadConfig() (*Config, error) {
	viper.SetDefault("mqtt.broker", "localhost")
	viper.SetDefault("mqtt.port", 1883)
	viper.SetDefault("mqtt.clientid", "sparkplug-go-publisher")
	viper.SetDefault("sparkplug.groupid", "Sparkplug B Devices")
	viper.SetDefault("sparkplug.edgenodeid", "Node-001")
	viper.SetDefault("sparkplug.deviceid", "Device-001")
	viper.SetDefault("sparkplug.scadahostid", "scada-host-1")

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	var config Config
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
