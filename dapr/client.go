package dapr

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/dapr/go-sdk/client"
)

// DaprClient wraps the Dapr client with additional helper methods
type DaprClient struct {
	client client.Client
}

// NewDaprClient creates a new Dapr client
func NewDaprClient() (*DaprClient, error) {
	c, err := client.NewClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create Dapr client: %w", err)
	}

	return &DaprClient{
		client: c,
	}, nil
}

// InvokeService invokes a method on another service
func (d *DaprClient) InvokeService(ctx context.Context, serviceName, method string, payload interface{}) ([]byte, error) {
	// Convert payload to JSON
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	content := &client.DataContent{
		ContentType: "application/json",
		Data:        data,
	}

	resp, err := d.client.InvokeMethodWithContent(ctx, serviceName, method, "post", content)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke %s.%s: %w", serviceName, method, err)
	}

	log.Printf("📞 Invoked %s.%s", serviceName, method)
	return resp, nil
}

// InvokeServiceWithResponse invokes a service method and unmarshals the response
func (d *DaprClient) InvokeServiceWithResponse(ctx context.Context, serviceName, method string, payload interface{}, response interface{}) error {
	data, err := d.InvokeService(ctx, serviceName, method, payload)
	if err != nil {
		return err
	}

	if response != nil {
		if err := json.Unmarshal(data, response); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return nil
}

// PublishEvent publishes an event to a topic
func (d *DaprClient) PublishEvent(ctx context.Context, pubsubName, topic string, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal event payload: %w", err)
	}

	err = d.client.PublishEvent(ctx, pubsubName, topic, data)
	if err != nil {
		return fmt.Errorf("failed to publish event to %s/%s: %w", pubsubName, topic, err)
	}

	log.Printf("📢 Published event to %s/%s", pubsubName, topic)
	return nil
}

// GetState retrieves state from the state store
func (d *DaprClient) GetState(ctx context.Context, storeName, key string) ([]byte, error) {
	item, err := d.client.GetState(ctx, storeName, key, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get state %s from store %s: %w", key, storeName, err)
	}

	return item.Value, nil
}

// SaveState saves state to the state store
func (d *DaprClient) SaveState(ctx context.Context, storeName, key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal state value: %w", err)
	}

	err = d.client.SaveState(ctx, storeName, key, data, nil)
	if err != nil {
		return fmt.Errorf("failed to save state %s to store %s: %w", key, storeName, err)
	}

	log.Printf("💾 Saved state %s to store %s", key, storeName)
	return nil
}

// DeleteState deletes state from the state store
func (d *DaprClient) DeleteState(ctx context.Context, storeName, key string) error {
	err := d.client.DeleteState(ctx, storeName, key, nil)
	if err != nil {
		return fmt.Errorf("failed to delete state %s from store %s: %w", key, storeName, err)
	}

	log.Printf("🗑️ Deleted state %s from store %s", key, storeName)
	return nil
}

// HealthCheck checks if Dapr is healthy
func (d *DaprClient) HealthCheck(ctx context.Context) error {
	// Try to invoke a simple method to check if Dapr is responsive
	_, err := d.client.InvokeMethod(ctx, "dapr", "health", "get")
	if err != nil {
		return fmt.Errorf("Dapr health check failed: %w", err)
	}

	return nil
}

// Close closes the Dapr client
func (d *DaprClient) Close() error {
	if d.client != nil {
		d.client.Close()
	}
	return nil
}

// Service-specific helper methods

// InvokePaymentService invokes the payment service
func (d *DaprClient) InvokePaymentService(ctx context.Context, method string, payload interface{}, response interface{}) error {
	return d.InvokeServiceWithResponse(ctx, "payment-service", method, payload, response)
}

// InvokeTripService invokes the trip service
func (d *DaprClient) InvokeTripService(ctx context.Context, method string, payload interface{}, response interface{}) error {
	return d.InvokeServiceWithResponse(ctx, "trip-service", method, payload, response)
}

// InvokeIdentityService invokes the identity service
func (d *DaprClient) InvokeIdentityService(ctx context.Context, method string, payload interface{}, response interface{}) error {
	return d.InvokeServiceWithResponse(ctx, "identity-service", method, payload, response)
}

// InvokeDriverService invokes the driver service
func (d *DaprClient) InvokeDriverService(ctx context.Context, method string, payload interface{}, response interface{}) error {
	return d.InvokeServiceWithResponse(ctx, "driver-service", method, payload, response)
}

// InvokeRiderService invokes the rider service
func (d *DaprClient) InvokeRiderService(ctx context.Context, method string, payload interface{}, response interface{}) error {
	return d.InvokeServiceWithResponse(ctx, "rider-service", method, payload, response)
}

// Event publishing helpers

// PublishTripEvent publishes a trip-related event
func (d *DaprClient) PublishTripEvent(ctx context.Context, eventType string, payload interface{}) error {
	return d.PublishEvent(ctx, "kafka-pubsub", "trip.events", map[string]interface{}{
		"type":      eventType,
		"timestamp": time.Now().Unix(),
		"data":      payload,
	})
}

// PublishPaymentEvent publishes a payment-related event
func (d *DaprClient) PublishPaymentEvent(ctx context.Context, eventType string, payload interface{}) error {
	return d.PublishEvent(ctx, "kafka-pubsub", "payment.events", map[string]interface{}{
		"type":      eventType,
		"timestamp": time.Now().Unix(),
		"data":      payload,
	})
}

// PublishIdentityEvent publishes an identity-related event
func (d *DaprClient) PublishIdentityEvent(ctx context.Context, eventType string, payload interface{}) error {
	return d.PublishEvent(ctx, "kafka-pubsub", "identity.events", map[string]interface{}{
		"type":      eventType,
		"timestamp": time.Now().Unix(),
		"data":      payload,
	})
}
