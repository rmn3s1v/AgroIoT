package main

import (
	"log"

	"ingestion-service/internal/handler"
	"ingestion-service/internal/mqtt"
)

func main() {
	broker := "tcp://localhost:1883"

	client := mqtt.NewClient(broker, "ingestion-service")

	topic := "devices/+/+/telemetry"

	client.Subscribe(topic, handler.HandleTelemetry)

	log.Println("Ingestion service is running...")

	select {}
}
