package grpc

import (
	"context"
	"fmt"
	"log"

	"github.com/dapr/go-sdk/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// GRPCClient manages gRPC connections for service-to-service communication
type GRPCClient struct {
	daprClient client.Client
	conns      map[string]*grpc.ClientConn
}

// NewGRPCClient creates a new gRPC client with Dapr integration
func NewGRPCClient() (*GRPCClient, error) {
	daprClient, err := client.NewClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create Dapr client: %w", err)
	}

	return &GRPCClient{
		daprClient: daprClient,
		conns:      make(map[string]*grpc.ClientConn),
	}, nil
}

// GetServiceConnection returns a gRPC connection to the specified service
func (c *GRPCClient) GetServiceConnection(serviceName string) (*grpc.ClientConn, error) {
	// Check if we already have a connection
	if conn, exists := c.conns[serviceName]; exists {
		return conn, nil
	}

	// Get service configuration
	config, exists := GetServiceConfig(serviceName)
	if !exists {
		return nil, fmt.Errorf("service %s not found in configuration", serviceName)
	}

	// Create connection using Dapr service invocation
	// In Dapr, we use the service name directly for service-to-service calls
	target := fmt.Sprintf("localhost:%s", config.Port)

	conn, err := grpc.Dial(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s: %w", serviceName, err)
	}

	// Cache the connection
	c.conns[serviceName] = conn

	log.Printf("✅ Connected to %s service on %s", serviceName, config.Port)
	return conn, nil
}

// InvokeService invokes a service method using Dapr
func (c *GRPCClient) InvokeService(ctx context.Context, serviceName, method string, data []byte) ([]byte, error) {
	// Use Dapr's service invocation for inter-service communication
	content := &client.DataContent{
		ContentType: "application/json",
		Data:        data,
	}

	resp, err := c.daprClient.InvokeMethodWithContent(ctx, serviceName, method, "post", content)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke %s.%s: %w", serviceName, method, err)
	}

	return resp, nil
}

// PublishEvent publishes an event to a topic using Dapr
func (c *GRPCClient) PublishEvent(ctx context.Context, topic string, data []byte) error {
	pubsubName := "kafka-pubsub" // Default pubsub component

	err := c.daprClient.PublishEvent(ctx, pubsubName, topic, data)
	if err != nil {
		return fmt.Errorf("failed to publish event to %s: %w", topic, err)
	}

	log.Printf("📢 Published event to topic: %s", topic)
	return nil
}

// Close closes all connections
func (c *GRPCClient) Close() error {
	var lastErr error

	for serviceName, conn := range c.conns {
		if err := conn.Close(); err != nil {
			log.Printf("Error closing connection to %s: %v", serviceName, err)
			lastErr = err
		}
	}

	if c.daprClient != nil {
		c.daprClient.Close()
	}

	return lastErr
}

// GetPaymentServiceClient returns a payment service client
func (c *GRPCClient) GetPaymentServiceClient(ctx context.Context) (interface{}, error) {
	conn, err := c.GetServiceConnection(PaymentService)
	if err != nil {
		return nil, err
	}

	// TODO: Return the actual payment service client
	// This would require importing the payment service protobuf client
	log.Printf("Payment service client connected")
	return conn, nil
}

// GetTripServiceClient returns a trip service client
func (c *GRPCClient) GetTripServiceClient(ctx context.Context) (interface{}, error) {
	conn, err := c.GetServiceConnection(TripService)
	if err != nil {
		return nil, err
	}

	// TODO: Return the actual trip service client
	log.Printf("Trip service client connected")
	return conn, nil
}

// GetIdentityServiceClient returns an identity service client
func (c *GRPCClient) GetIdentityServiceClient(ctx context.Context) (interface{}, error) {
	conn, err := c.GetServiceConnection(IdentityService)
	if err != nil {
		return nil, err
	}

	// TODO: Return the actual identity service client
	log.Printf("Identity service client connected")
	return conn, nil
}

// GetDriverServiceClient returns a driver service client
func (c *GRPCClient) GetDriverServiceClient(ctx context.Context) (interface{}, error) {
	conn, err := c.GetServiceConnection(DriverService)
	if err != nil {
		return nil, err
	}

	// TODO: Return the actual driver service client
	log.Printf("Driver service client connected")
	return conn, nil
}

// GetRiderServiceClient returns a rider service client
func (c *GRPCClient) GetRiderServiceClient(ctx context.Context) (interface{}, error) {
	conn, err := c.GetServiceConnection(RiderService)
	if err != nil {
		return nil, err
	}

	// TODO: Return the actual rider service client
	log.Printf("Rider service client connected")
	return conn, nil
}
