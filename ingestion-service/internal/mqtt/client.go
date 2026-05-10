package mqtt

import (
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Client struct {
	client mqtt.Client
}

func NewClient(broker string, clientID string) *Client {

	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(clientID)
	opts.SetAutoReconnect(true)

	opts.OnConnect = func(c mqtt.Client) {
		log.Println("Connected to MQTT broker")
	}

	opts.OnConnectionLost = func(c mqtt.Client, err error) {
		log.Println("Connection lost:", err)
	}

	client := mqtt.NewClient(opts)

	for {
		token := client.Connect()

		if token.Wait() && token.Error() != nil {
			log.Println("MQTT broker not ready, retrying...")
			time.Sleep(5 * time.Second)
			continue
		}

		break
	}

	return &Client{client: client}
}

func (c *Client) Subscribe(topic string, handler mqtt.MessageHandler) {
	token := c.client.Subscribe(topic, 1, handler)
	token.Wait()

	if token.Error() != nil {
		log.Println("Subscribe error:", token.Error())
	} else {
		log.Println("Subscribe to:", topic)
	}
}
