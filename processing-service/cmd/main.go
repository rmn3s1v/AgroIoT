package main

import (
	"log"

	"processing-service/internal/mongodb"
	"processing-service/internal/rabbitmq"
)

func main() {
	rabbit := rabbitmq.NewConsumer("amqp://guest:guest@rabbitmq:5672/")
	mongo := mongodb.NewClient("mongodb://mongodb:27017")

	log.Println("Processing service started")

	rabbit.Consume("telemetry_queue", mongo)
}
