package events

import (
	"context"
	"encoding/json"
	"time"
)

// EventType represents all event types in the system
type EventType string

const (
	// Trip Events
	EventTypeTripCreated    EventType = "trip.created"
	EventTypeTripAccepted   EventType = "trip.accepted"
	EventTypeTripCancelled  EventType = "trip.cancelled"
	EventTypeTripCompleted  EventType = "trip.completed"
	EventTypeTripInProgress EventType = "trip.in_progress"

	// Booking Events
	EventTypeBookingCreated   EventType = "booking.created"
	EventTypeBookingConfirmed EventType = "booking.confirmed"
	EventTypeBookingCancelled EventType = "booking.cancelled"

	// Driver Events
	EventTypeDriverOnline   EventType = "driver.online"
	EventTypeDriverOffline  EventType = "driver.offline"
	EventTypeDriverLocation EventType = "driver.location"
	EventTypeDriverStatus   EventType = "driver.status"

	// Rider Events
	EventTypeRiderCreated EventType = "rider.created"
	EventTypeRiderUpdated EventType = "rider.updated"

	// Payment Events
	EventTypePaymentInitiated EventType = "payment.initiated"
	EventTypePaymentCompleted EventType = "payment.completed"
	EventTypePaymentFailed    EventType = "payment.failed"
	EventTypePaymentRefunded  EventType = "payment.refunded"

	// Bidding Events
	EventTypeBiddingStarted EventType = "bidding.started"
	EventTypeBiddingEnded   EventType = "bidding.ended"
	EventTypeBidSubmitted   EventType = "bid.submitted"

	// Identity Events
	EventTypeUserCreated EventType = "user.created"
	EventTypeUserUpdated EventType = "user.updated"
	EventTypeUserDeleted EventType = "user.deleted"

	// Notification Events
	EventTypeDriverNotification EventType = "driver.notification"
	EventTypeRiderNotification  EventType = "rider.notification"
)

// Topic represents Dapr pub/sub topic names
type Topic string

const (
	TopicTripEvents          Topic = "trip.events"
	TopicBookingEvents       Topic = "booking.events"
	TopicDriverEvents        Topic = "driver.events"
	TopicRiderEvents         Topic = "rider.events"
	TopicPaymentEvents       Topic = "payment.events"
	TopicBiddingEvents       Topic = "bidding.events"
	TopicIdentityEvents      Topic = "identity.events"
	TopicDriverNotifications Topic = "driver.notifications"
	TopicRiderNotifications  Topic = "rider.notifications"
)

// BaseEvent is the standard event structure for all events
type BaseEvent struct {
	Type      EventType         `json:"type"`
	Service   string            `json:"service"`
	Timestamp time.Time         `json:"timestamp"`
	Payload   json.RawMessage   `json:"payload"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

// EventPublisher provides standardized event publishing
type EventPublisher interface {
	Publish(ctx context.Context, topic Topic, event BaseEvent) error
}

// NewBaseEvent creates a new base event
func NewBaseEvent(eventType EventType, service string, payload interface{}) (*BaseEvent, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return &BaseEvent{
		Type:      eventType,
		Service:   service,
		Timestamp: time.Now(),
		Payload:   payloadBytes,
		Metadata:  make(map[string]string),
	}, nil
}

// GetTopicForEventType returns the appropriate topic for an event type
func GetTopicForEventType(eventType EventType) Topic {
	switch {
	case eventType == EventTypeTripCreated || eventType == EventTypeTripAccepted ||
		eventType == EventTypeTripCancelled || eventType == EventTypeTripCompleted ||
		eventType == EventTypeTripInProgress:
		return TopicTripEvents
	case eventType == EventTypeBookingCreated || eventType == EventTypeBookingConfirmed ||
		eventType == EventTypeBookingCancelled:
		return TopicBookingEvents
	case eventType == EventTypeDriverOnline || eventType == EventTypeDriverOffline ||
		eventType == EventTypeDriverLocation || eventType == EventTypeDriverStatus:
		return TopicDriverEvents
	case eventType == EventTypeRiderCreated || eventType == EventTypeRiderUpdated:
		return TopicRiderEvents
	case eventType == EventTypePaymentInitiated || eventType == EventTypePaymentCompleted ||
		eventType == EventTypePaymentFailed || eventType == EventTypePaymentRefunded:
		return TopicPaymentEvents
	case eventType == EventTypeBiddingStarted || eventType == EventTypeBiddingEnded ||
		eventType == EventTypeBidSubmitted:
		return TopicBiddingEvents
	case eventType == EventTypeUserCreated || eventType == EventTypeUserUpdated ||
		eventType == EventTypeUserDeleted:
		return TopicIdentityEvents
	case eventType == EventTypeDriverNotification:
		return TopicDriverNotifications
	case eventType == EventTypeRiderNotification:
		return TopicRiderNotifications
	default:
		return TopicTripEvents // Default fallback
	}
}
