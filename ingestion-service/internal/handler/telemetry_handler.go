package handler

import (
	"encoding/json"
	"log"

	"ingestion-service/internal/model"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func HandleTelemetry(client mqtt.Client, msg mqtt.Message) {
	log.Println("Received message on topic:", msg.Topic())

	var telemetry model.Telemetry

	err := json.Unmarshal(msg.Payload(), &telemetry)
	if err != nil {
		log.Println("Invalid JSON:", err)
		return
	}

	log.Printf("Parsed telemetry: %+v\n", telemetry)

	// TODO: validation
	// TODO: send to RabbitMQ
}
