package grpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/dapr/go-sdk/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
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
	// Check if we already have a connection and verify it's still healthy
	if conn, exists := c.conns[serviceName]; exists {
		// Check connection state - if it's not ready, try to reconnect
		state := conn.GetState()
		if state == connectivity.Ready || state == connectivity.Idle {
			return conn, nil
		}
		// Connection is in a bad state, remove it and create a new one
		log.Printf("⚠️ Connection to %s is in state %v, reconnecting...", serviceName, state)
		delete(c.conns, serviceName)
		if conn != nil {
			conn.Close()
		}
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

	// Configure dial options with timeout, keepalive, and retry
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dialOptions := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(), // Block until connection is established
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                10 * time.Second, // Send keepalive pings every 10 seconds
			Timeout:             3 * time.Second,  // Wait 3 seconds for ping ack before considering the connection dead
			PermitWithoutStream: true,             // Send pings even when there are no active streams
		}),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(4*1024*1024), // 4MB max message size
			grpc.MaxCallSendMsgSize(4*1024*1024), // 4MB max message size
		),
	}

	conn, err := grpc.DialContext(ctx, target, dialOptions...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s at %s: %w", serviceName, target, err)
	}

	// WithBlock() ensures connection is established, but verify state
	state := conn.GetState()
	if state != connectivity.Ready && state != connectivity.Idle {
		conn.Close()
		return nil, fmt.Errorf("connection to %s at %s is in state %v, expected Ready or Idle", serviceName, target, state)
	}

	// Cache the connection
	c.conns[serviceName] = conn

	log.Printf("✅ Connected to %s service on %s (state: %v)", serviceName, target, state)
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
