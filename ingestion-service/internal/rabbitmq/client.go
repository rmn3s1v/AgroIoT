package rabbitmq

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Client struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewClient(url string) *Client {
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("Failed to open a channel:", err)
	}

	return &Client{
		conn:    conn,
		channel: ch,
	}
}

func (c *Client) Publish(queueName string, body []byte) {
	_, err := c.channel.QueueDeclare(
		queueName,
		true,  // durable
		false, // auto delete
		false, // exclusive
		false,
		nil,
	)
	if err != nil {
		log.Println("Queue declare error:", err)
		return
	}

	err = c.channel.Publish(
		"",        // exchange
		queueName, // routing key
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	if err != nil {
		log.Println("Publish error:", err)
		return
	}

	log.Println("Message sent to RabbitMQ")
}
