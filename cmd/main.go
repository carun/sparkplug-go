package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"sparkplug-go/internal/config"
	"sparkplug-go/internal/mqtt"
	"sparkplug-go/internal/sparkplug"
	pb "sparkplug-go/proto"
)

func createNodeMetrics() []*pb.Metric {
	metrics := []*pb.Metric{
		sparkplug.CreateMetric("Node Control/Rebirth", pb.DataType_Boolean, false),
		sparkplug.CreateMetric("Node Control/Reboot", pb.DataType_Boolean, false),
		sparkplug.CreateMetric("Node Control/NextServer", pb.DataType_Boolean, false),
		sparkplug.CreateMetric("Properties/Version", pb.DataType_String, "v1.0.0"),
		sparkplug.CreateMetric("Properties/Vendor", pb.DataType_String, "Your Company"),
	}

	// Add a DataSet metric example
	dataSet := sparkplug.CreateDataSet(
		[]string{"timestamp", "value", "quality"},
		[]string{"String", "Double", "Int32"},
	)
	metrics = append(metrics, sparkplug.CreateMetric("DataSet Example", pb.DataType_DataSetType, dataSet))

	// Add properties to a metric
	propMetric := sparkplug.CreateMetric("Sensor/Temperature", pb.DataType_Double, 23.45)
	propMetric.Properties = sparkplug.CreatePropertySet(map[string]interface{}{
		"engUnit": "Celsius",
		"engHigh": float64(100.0),
		"engLow":  float64(-40.0),
	})
	metrics = append(metrics, propMetric)

	return metrics
}

func createDeviceMetrics() []*pb.Metric {
	return []*pb.Metric{
		sparkplug.CreateMetric("Device Control/Rebirth", pb.DataType_Boolean, false),
		sparkplug.CreateMetric("Device Control/Reboot", pb.DataType_Boolean, false),
		sparkplug.CreateMetric("Sensors/Temperature", pb.DataType_Double, 25.6),
		sparkplug.CreateMetric("Sensors/Humidity", pb.DataType_Double, 60.0),
		sparkplug.CreateMetric("Sensors/Pressure", pb.DataType_Double, 1013.25),
		sparkplug.CreateMetric("Status/Battery", pb.DataType_Int32, int32(85)),
		sparkplug.CreateMetric("Status/Connected", pb.DataType_Boolean, true),
	}
}

func simulateDeviceData() []*pb.Metric {
	return []*pb.Metric{
		sparkplug.CreateMetric("Sensors/Temperature", pb.DataType_Double, 25.6+float64(time.Now().Second()%10)),
		sparkplug.CreateMetric("Sensors/Humidity", pb.DataType_Double, 60.0+float64(time.Now().Second()%20)),
		sparkplug.CreateMetric("Status/Battery", pb.DataType_Int32, int32(85-time.Now().Second()%10)),
	}
}

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create MQTT client
	client, err := mqtt.NewClient(cfg)
	if err != nil {
		log.Fatalf("Failed to create MQTT client: %v", err)
	}
	defer client.Disconnect()

	// Create Sparkplug message and topic builders
	msgBuilder := sparkplug.NewMessageBuilder()
	topicBuilder := sparkplug.NewTopicBuilder(
		cfg.Sparkplug.GroupID,
		cfg.Sparkplug.EdgeNodeID,
		cfg.Sparkplug.DeviceID,
		cfg.Sparkplug.ScadaHostID,
	)

	// Send NBIRTH message
	nbirthPayload, err := msgBuilder.CreateNBirth(createNodeMetrics())
	if err != nil {
		log.Fatalf("Failed to create NBIRTH payload: %v", err)
	}
	if err := client.Publish(topicBuilder.NBirthTopic(), nbirthPayload); err != nil {
		log.Fatalf("Failed to publish NBIRTH: %v", err)
	}
	log.Println("Published NBIRTH message")

	// Send DBIRTH message
	dbirthPayload, err := msgBuilder.CreateDBirth(createDeviceMetrics())
	if err != nil {
		log.Fatalf("Failed to create DBIRTH payload: %v", err)
	}
	if err := client.Publish(topicBuilder.DBirthTopic(), dbirthPayload); err != nil {
		log.Fatalf("Failed to publish DBIRTH: %v", err)
	}
	log.Println("Published DBIRTH message")

	// Create a ticker for sending NDATA and DDATA messages
	dataTicker := time.NewTicker(5 * time.Second)
	defer dataTicker.Stop()

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	log.Println("Started publishing data. Press Ctrl+C to exit.")

	for {
		select {
		case <-dataTicker.C:
			// Send NDATA message
			ndataPayload, err := msgBuilder.CreateNData([]*pb.Metric{
				sparkplug.CreateMetric("Properties/Version", pb.DataType_String, "v1.0.0"),
			})
			if err != nil {
				log.Printf("Failed to create NDATA payload: %v", err)
				continue
			}
			if err := client.Publish(topicBuilder.NDataTopic(), ndataPayload); err != nil {
				log.Printf("Failed to publish NDATA: %v", err)
			}

			// Send DDATA message with simulated data
			ddataPayload, err := msgBuilder.CreateDData(simulateDeviceData())
			if err != nil {
				log.Printf("Failed to create DDATA payload: %v", err)
				continue
			}
			if err := client.Publish(topicBuilder.DDataTopic(), ddataPayload); err != nil {
				log.Printf("Failed to publish DDATA: %v", err)
			}

		case sig := <-sigChan:
			log.Printf("Received signal: %v", sig)

			// Send NDEATH message
			ndeathPayload, err := msgBuilder.CreateNDeath()
			if err != nil {
				log.Printf("Failed to create NDEATH payload: %v", err)
			} else {
				if err := client.Publish(topicBuilder.NDeathTopic(), ndeathPayload); err != nil {
					log.Printf("Failed to publish NDEATH: %v", err)
				}
			}
			log.Println("Published NDEATH message")
			return
		}
	}
}
