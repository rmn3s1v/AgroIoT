package main

import (
	"log"

	"ingestion-service/internal/handler"
	"ingestion-service/internal/mqtt"
	"ingestion-service/internal/rabbitmq"
)

func main() {
	broker := "tcp://localhost:1883"

	mqttClient := mqtt.NewClient(broker, "ingestion-service")

	rabbit := rabbitmq.NewClient("amqp://guest:guest@localhost:5672/")

	topic := "devices/+/+/telemetry"

	mqttClient.Subscribe(topic, handler.HandleTelemetry(rabbit))

	log.Println("Ingestion service is running...")

	select {}
}
