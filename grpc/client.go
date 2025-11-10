package grpc

import (
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
	namespace  string
}

// NewGRPCClient creates a new gRPC client with Dapr integration
func NewGRPCClient(namespace string) (*GRPCClient, error) {
	daprClient, err := client.NewClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create Dapr client: %w", err)
	}

	return &GRPCClient{
		daprClient: daprClient,
		conns:      make(map[string]*grpc.ClientConn),
		namespace:  namespace,
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

	var target string
	if c.namespace != "" {
		// Kubernetes DNS name for headless service with namespace
		target = fmt.Sprintf("%s.%s.svc.cluster.local:%s", config.Name, c.namespace, config.Port)
	} else {
		// Localhost for local development
		target = fmt.Sprintf("localhost:%s", config.Port)
	}

	conn, err := grpc.Dial(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s: %w", serviceName, err)
	}

	// Cache the connection
	c.conns[serviceName] = conn

	log.Printf("✅ Connected to %s service on %s", serviceName, target)
	return conn, nil
}

// Close closes all connections
func (c *GRPCClient) Close() error {
	var lastErr error

	for serviceName, conn := range c.conns {
		if conn != nil {
			if err := conn.Close(); err != nil {
				log.Printf("Error closing connection to %s: %v", serviceName, err)
				lastErr = err
			}
		}
	}

	if c.daprClient != nil {
		c.daprClient.Close()
	}

	return lastErr
}
