package rabbitmq

import (
	"encoding/json"
	"log"

	"processing-service/internal/model"
	"processing-service/internal/mongodb"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	channel *amqp.Channel
}

func NewConsumer(url string) *Consumer {
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Fatal(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}

	return &Consumer{channel: ch}
}

func (c *Consumer) Consume(queue string, mongoClient *mongodb.Client) {
	msgs, err := c.channel.Consume(
		queue,
		"",
		true, // auto-ack
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	forever := make(chan bool)

	go func() {
		for msg := range msgs {
			var telemetry model.Telemetry

			err := json.Unmarshal(msg.Body, &telemetry)
			if err != nil {
				log.Println("JSON error:", err)
				continue
			}

			_, err = mongoClient.Collection.InsertOne(nil, telemetry)
			if err != nil {
				log.Println("Mongo insert error:", err)
				continue
			}

			log.Println("Saved to MongoDB:", telemetry.DeviceID)
		}
	}()

	log.Println("Waiting for messages...")
	<-forever
}
