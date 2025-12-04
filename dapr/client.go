package dapr

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	dapr "github.com/dapr/go-sdk/client"
)

// ServiceClient wraps Dapr client with service-to-service communication capabilities
// Use cases:
// 1. Pub/Sub: Publish events to topics for asynchronous communication
// 2. State Management: Store and retrieve state from Redis state store
// 3. Secret Management: Retrieve secrets from secret stores
type ServiceClient struct {
	client dapr.Client
}

// NewDaprClient creates a new Dapr service client
// The Dapr client connects to the Dapr sidecar using environment variables:
// - DAPR_GRPC_PORT: gRPC port (default: 50001)
// - DAPR_HTTP_PORT: HTTP port (default: 3500)
// Make sure Dapr sidecar is running before calling this function.
func NewDaprClient() (*ServiceClient, error) {
	client, err := dapr.NewClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create Dapr client: %w. Make sure Dapr sidecar is running. Use 'dapr run' or the provided run-with-dapr script", err)
	}
	return &ServiceClient{client: client}, nil
}

// Close gracefully closes the Dapr client
func (s *ServiceClient) Close() error {
	if s.client != nil {
		s.client.Close()
	}
	return nil
}

// GetRawClient returns the underlying Dapr client for advanced use cases
func (s *ServiceClient) GetRawClient() dapr.Client {
	return s.client
}

// Pub/Sub Methods

// PublishEvent publishes an event to a Dapr pub/sub topic
// Use case: Asynchronous event-driven communication between services
// Example: Publishing trip events, driver notifications, payment events
func (s *ServiceClient) PublishEvent(ctx context.Context, pubsubName, topic string, data interface{}) error {
	var payload []byte
	var err error

	// Handle different data types
	switch v := data.(type) {
	case []byte:
		payload = v
	case string:
		payload = []byte(v)
	default:
		payload, err = json.Marshal(data)
		if err != nil {
			return fmt.Errorf("failed to marshal event data: %w", err)
		}
	}

	err = s.client.PublishEvent(ctx, pubsubName, topic, payload)
	if err != nil {
		return fmt.Errorf("failed to publish event to topic %s on pubsub %s: %w", topic, pubsubName, err)
	}

	log.Printf("✅ Successfully published event to topic %s on pubsub %s", topic, pubsubName)
	return nil
}

// State Management Methods

// SaveState saves state to a Dapr state store
// Use case: Caching, session management, temporary data storage
// Example: Storing trip session data, driver location cache, payment session data
func (s *ServiceClient) SaveState(ctx context.Context, storeName, key string, value interface{}) error {
	var data []byte
	var err error

	// Handle different data types
	switch v := value.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		data, err = json.Marshal(value)
		if err != nil {
			return fmt.Errorf("failed to marshal state value: %w", err)
		}
	}

	err = s.client.SaveState(ctx, storeName, key, data, nil)
	if err != nil {
		return fmt.Errorf("failed to save state to store %s with key %s: %w", storeName, key, err)
	}

	return nil
}

// GetState retrieves state from a Dapr state store
// Use case: Retrieving cached data, session information
// Example: Getting trip session data, driver location cache, payment session data
func (s *ServiceClient) GetState(ctx context.Context, storeName, key string) ([]byte, error) {
	item, err := s.client.GetState(ctx, storeName, key, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get state from store %s with key %s: %w", storeName, key, err)
	}

	return item.Value, nil
}

// GetStateWithMetadata retrieves state from a Dapr state store with metadata
func (s *ServiceClient) GetStateWithMetadata(ctx context.Context, storeName, key string, metadata map[string]string) ([]byte, map[string]string, error) {
	item, err := s.client.GetState(ctx, storeName, key, metadata)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get state from store %s with key %s: %w", storeName, key, err)
	}

	return item.Value, item.Metadata, nil
}

// DeleteState deletes state from a Dapr state store
// Use case: Removing cached data, cleaning up sessions
// Example: Clearing trip session data, removing expired sessions, payment session cleanup
func (s *ServiceClient) DeleteState(ctx context.Context, storeName, key string) error {
	err := s.client.DeleteState(ctx, storeName, key, nil)
	if err != nil {
		return fmt.Errorf("failed to delete state from store %s with key %s: %w", storeName, key, err)
	}

	return nil
}

// Secret Management Methods

// GetSecret retrieves a secret from a Dapr secret store
// Use case: Retrieving sensitive configuration data
// Example: API keys, database passwords, third-party service credentials, payment gateway keys
func (s *ServiceClient) GetSecret(ctx context.Context, storeName, key string) (map[string]string, error) {
	secrets, err := s.client.GetSecret(ctx, storeName, key, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get secret from store %s with key %s: %w", storeName, key, err)
	}

	return secrets, nil
}

// GetSecretWithMetadata retrieves a secret from a Dapr secret store with metadata
func (s *ServiceClient) GetSecretWithMetadata(ctx context.Context, storeName, key string, metadata map[string]string) (map[string]string, map[string]string, error) {
	secrets, err := s.client.GetSecret(ctx, storeName, key, metadata)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get secret from store %s with key %s: %w", storeName, key, err)
	}

	return secrets, metadata, nil
}

