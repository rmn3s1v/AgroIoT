package handler

import (
	"encoding/json"
	"log"

	"ingestion-service/internal/model"
	"ingestion-service/internal/rabbitmq"
	"ingestion-service/internal/validator"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func HandleTelemetry(rabbit *rabbitmq.Client) mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {

		log.Println("Received message on topic:", msg.Topic())

		var telemetry model.Telemetry

		err := json.Unmarshal(msg.Payload(), &telemetry)
		if err != nil {
			log.Println("Invalid JSON:", err)
			return
		}

		err = validator.ValidateTelemetry(telemetry)
		if err != nil {
			log.Println("Validation error:", err)
			return
		}

		normalized := validator.NormalizeMetrics(telemetry.Metrics)

		event := map[string]interface{}{
			"device_id":   telemetry.DeviceID,
			"device_type": telemetry.DeviceType,
			"timestamp":   telemetry.Timestamp,
			"metrics":     normalized,
		}

		body, err := json.Marshal(event)
		if err != nil {
			log.Println("JSON marshal error:", err)
			return
		}

		// 🚀 отправка в RabbitMQ
		rabbit.Publish("telemetry_queue", body)
	}
}
