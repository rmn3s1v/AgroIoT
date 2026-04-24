package mqtt

import (
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Client struct {
	client mqtt.Client
}

func NewClient(broker string, cliendID string) *Client {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(cliendID)
	opts.SetAutoReconnect(true)

	opts.OnConnect = func(c mqtt.Client) {
		log.Println("Connected to MQTT broker")
	}

	opts.OnConnectionLost = func(c mqtt.Client, err error) {
		log.Println("Connection lost: ", err)
	}

	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
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
