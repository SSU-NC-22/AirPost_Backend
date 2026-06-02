package mqtt

import (
	"encoding/json"
	"log"
	"os"

	paho "github.com/eclipse/paho.mqtt.golang"
)

// defaultBroker is used when MQTT_BROKER is unset (dev/docker-compose).
const defaultBroker = "tcp://localhost:1883"

// StatusHandler is invoked for every delivery status message received.
type StatusHandler func(DeliveryStatus)

// Client publishes delivery requests and dispatches incoming status messages.
// It wraps the paho MQTT client so the rest of the app depends only on our
// contract types, not on the broker library.
type Client struct {
	mqtt paho.Client
}

// brokerAddr resolves the broker URL from the environment, defaulting to the
// local dev broker. Kept separate so it is reused by tests/config.
func brokerAddr() string {
	if addr := os.Getenv("MQTT_BROKER"); addr != "" {
		return addr
	}
	return defaultBroker
}

// NewClient connects to the broker and returns a ready Client. The clientID is
// used by the broker for session identity (e.g. "airpost-application").
func NewClient(clientID string) (*Client, error) {
	opts := paho.NewClientOptions().
		AddBroker(brokerAddr()).
		SetClientID(clientID).
		SetAutoReconnect(true)

	c := paho.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}
	return &Client{mqtt: c}, nil
}

// PublishRequest serializes and publishes a delivery request at QoS 1 so the
// flight service receives it at least once.
func (c *Client) PublishRequest(req DeliveryRequest) error {
	payload, err := json.Marshal(req)
	if err != nil {
		return err
	}
	token := c.mqtt.Publish(RequestTopic, 1, false, payload)
	token.Wait()
	return token.Error()
}

// SubscribeStatus registers handler for every status message on StatusTopic.
// Malformed payloads are logged and skipped rather than crashing the consumer.
func (c *Client) SubscribeStatus(handler StatusHandler) error {
	token := c.mqtt.Subscribe(StatusTopic, 1, func(_ paho.Client, m paho.Message) {
		var status DeliveryStatus
		if err := json.Unmarshal(m.Payload(), &status); err != nil {
			log.Printf("mqtt: dropping malformed status payload: %v", err)
			return
		}
		handler(status)
	})
	token.Wait()
	return token.Error()
}

// Disconnect cleanly closes the broker connection.
func (c *Client) Disconnect() {
	c.mqtt.Disconnect(250)
}
