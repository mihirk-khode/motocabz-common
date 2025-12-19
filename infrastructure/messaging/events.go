package messaging

import (
	"context"
	"encoding/json"
	"time"
)

// Event represents a domain event
type Event struct {
	Type      string                 `json:"type"`
	Service   string                 `json:"service"`
	Timestamp time.Time              `json:"timestamp"`
	Payload   json.RawMessage        `json:"payload"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// EventPublisher publishes events
type EventPublisher interface {
	Publish(ctx context.Context, topic string, event *Event) error
}

// EventSubscriber subscribes to events
type EventSubscriber interface {
	Subscribe(ctx context.Context, topic string, handler EventHandler) error
}

// EventHandler handles incoming events
type EventHandler func(ctx context.Context, event *Event) error

// NewEvent creates a new event
func NewEvent(eventType, service string, payload interface{}) (*Event, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return &Event{
		Type:      eventType,
		Service:   service,
		Timestamp: time.Now(),
		Payload:   payloadBytes,
		Metadata:  make(map[string]interface{}),
	}, nil
}
