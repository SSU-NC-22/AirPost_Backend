package mqtt

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	paho "github.com/eclipse/paho.mqtt.golang"
)

// defaultBroker is used when MQTT_BROKER is unset (dev/docker-compose).
const defaultBroker = "tcp://localhost:1883"

// StatusHandler is invoked for every delivery status message received.
type StatusHandler func(DeliveryStatus)

// TelemetryHandler is invoked with a drone's live position parsed from its data/<id> telemetry.
type TelemetryHandler func(droneID int, lat, lon, alt float64)

// deviceTelemetry is the device telemetry envelope every drone/station publishes on data/<node_id>.
// For a drone, values = [lat, lon, alt, velocity, battery, status].
type deviceTelemetry struct {
	NodeID string    `json:"node_id"`
	Values []float64 `json:"values"`
}

// droneNumericID turns a "DRO51"-style node id into its numeric id (51); ok=false for non-drones.
func droneNumericID(nodeID string) (int, bool) {
	if !strings.HasPrefix(nodeID, "DRO") {
		return 0, false
	}
	id, err := strconv.Atoi(strings.TrimPrefix(nodeID, "DRO"))
	if err != nil {
		return 0, false
	}
	return id, true
}

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

// SubscribeFleetTelemetry registers handler for every drone's live position, parsed from the
// data/<node_id> telemetry stream (the same MQTT the Sink forwards to Kafka). Station/non-drone
// nodes and malformed/short payloads are skipped. The backend fleet monitor feeds these positions.
func (c *Client) SubscribeFleetTelemetry(handler TelemetryHandler) error {
	token := c.mqtt.Subscribe("data/#", 1, func(_ paho.Client, m paho.Message) {
		var t deviceTelemetry
		if err := json.Unmarshal(m.Payload(), &t); err != nil {
			return
		}
		id, ok := droneNumericID(t.NodeID)
		if !ok || len(t.Values) < 3 {
			return
		}
		lat, lon, alt := t.Values[0], t.Values[1], t.Values[2]
		handler(id, lat, lon, alt)
	})
	token.Wait()
	return token.Error()
}

// PublishHold sends a deconfliction advisory to one drone: hold in place (true) or proceed (false).
// The drone's onboard node honours it on command/downlink/Hold/DRO<id>.
func (c *Client) PublishHold(droneID int, hold bool) error {
	topic := fmt.Sprintf("command/downlink/Hold/DRO%d", droneID)
	payload, err := json.Marshal(map[string]bool{"hold": hold})
	if err != nil {
		return err
	}
	token := c.mqtt.Publish(topic, 1, false, payload)
	token.Wait()
	return token.Error()
}

// Disconnect cleanly closes the broker connection.
func (c *Client) Disconnect() {
	c.mqtt.Disconnect(250)
}
