package rabbitmq

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"processing-service/internal/model"
	"processing-service/internal/mongodb"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	channel *amqp.Channel
}

func NewConsumer(url string) *Consumer {

	var conn *amqp.Connection
	var err error

	for {
		conn, err = amqp.Dial(url)

		if err != nil {
			log.Println("RabbitMQ not ready, retrying...")
			time.Sleep(5 * time.Second)
			continue
		}

		break
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}

	return &Consumer{
		channel: ch,
	}
}

func (c *Consumer) Consume(queue string, mongoClient *mongodb.Client) {

	_, err := c.channel.QueueDeclare(
		queue,
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Fatal("Queue declare error:", err)
	}

	msgs, err := c.channel.Consume(
		queue,
		"",
		true,
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

			_, err = mongoClient.Collection.InsertOne(
				context.Background(),
				telemetry,
			)

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
