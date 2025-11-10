package grpc

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/dapr/go-sdk/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

const (
	// Default connection timeout
	defaultDialTimeout = 30 * time.Second
	// Default retry attempts
	defaultMaxRetries = 3
	// Default retry backoff
	defaultRetryBackoff = 1 * time.Second
)

// ConnectionInfo stores connection metadata
type ConnectionInfo struct {
	ServiceName string
	Target      string
	State       connectivity.State
	CreatedAt   time.Time
	LastUsed    time.Time
	Conn        *grpc.ClientConn
}

// GRPCClient manages gRPC connections for service-to-service communication
type GRPCClient struct {
	daprClient   client.Client
	conns        map[string]*grpc.ClientConn
	connInfo     map[string]*ConnectionInfo
	connsMutex   sync.RWMutex
	namespace    string
	dialTimeout  time.Duration
	maxRetries   int
	retryBackoff time.Duration
	stopMonitor  chan struct{}
	monitorWg    sync.WaitGroup
}

// NewGRPCClient creates a new gRPC client with Dapr integration
func NewGRPCClient(namespace string) (*GRPCClient, error) {
	daprClient, err := client.NewClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create Dapr client: %w", err)
	}

	grpcClient := &GRPCClient{
		daprClient:   daprClient,
		conns:        make(map[string]*grpc.ClientConn),
		connInfo:     make(map[string]*ConnectionInfo),
		namespace:    namespace,
		dialTimeout:  defaultDialTimeout,
		maxRetries:   defaultMaxRetries,
		retryBackoff: defaultRetryBackoff,
		stopMonitor:  make(chan struct{}),
	}

	// Start background connection health monitor
	grpcClient.startConnectionMonitor()

	return grpcClient, nil
}

// GetServiceConnection returns a gRPC connection to the specified service
// with automatic retry and connection health checking
func (c *GRPCClient) GetServiceConnection(serviceName string) (*grpc.ClientConn, error) {
	// Check if we already have a connection and verify it's still healthy
	c.connsMutex.RLock()
	if conn, exists := c.conns[serviceName]; exists {
		state := conn.GetState()
		c.connsMutex.RUnlock()

		// If connection is ready or idle, return it
		if state == connectivity.Ready || state == connectivity.Idle {
			// Update last used time
			if info, exists := c.connInfo[serviceName]; exists {
				info.LastUsed = time.Now()
				info.State = state
			}
			return conn, nil
		}

		// Connection is in a bad state, remove it and create a new one
		log.Printf("⚠️ Connection to %s is in state %v, reconnecting...", serviceName, state)
		c.connsMutex.Lock()
		delete(c.conns, serviceName)
		delete(c.connInfo, serviceName)
		if conn != nil {
			conn.Close()
		}
		c.connsMutex.Unlock()
	} else {
		c.connsMutex.RUnlock()
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

	// Retry connection with exponential backoff
	var conn *grpc.ClientConn
	var lastErr error

	for attempt := 0; attempt < c.maxRetries; attempt++ {
		if attempt > 0 {
			backoff := c.retryBackoff * time.Duration(1<<uint(attempt-1)) // Exponential backoff
			log.Printf("🔄 Retrying connection to %s (attempt %d/%d) after %v...",
				serviceName, attempt+1, c.maxRetries, backoff)
			time.Sleep(backoff)
		}

		// Create context with timeout for each attempt
		ctx, cancel := context.WithTimeout(context.Background(), c.dialTimeout)

		// Configure dial options with keepalive and retry
		dialOptions := []grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			// Don't use WithBlock() - allow non-blocking connection attempts
			grpc.WithKeepaliveParams(keepalive.ClientParameters{
				Time:                30 * time.Second, // Send keepalive pings every 30 seconds
				Timeout:             10 * time.Second, // Wait 10 seconds for ping ack
				PermitWithoutStream: true,             // Send pings even when there are no active streams
			}),
			grpc.WithDefaultCallOptions(
				grpc.MaxCallRecvMsgSize(4*1024*1024), // 4MB max message size
				grpc.MaxCallSendMsgSize(4*1024*1024), // 4MB max message size
			),
		}

		var err error
		conn, err = grpc.DialContext(ctx, target, dialOptions...)
		cancel()

		if err != nil {
			lastErr = fmt.Errorf("failed to connect to %s at %s (attempt %d/%d): %w",
				serviceName, target, attempt+1, c.maxRetries, err)
			log.Printf("❌ %v", lastErr)
			continue
		}

		// Wait for connection to be ready (with timeout)
		// Poll connection state until ready or timeout
		readyCtx, readyCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer readyCancel()
		stateCheckTicker := time.NewTicker(100 * time.Millisecond)
		defer stateCheckTicker.Stop()

		connectionReady := false
		for !connectionReady {
			select {
			case <-readyCtx.Done():
				// Timeout reached, check final state
				state := conn.GetState()
				if state == connectivity.Ready || state == connectivity.Idle {
					connectionReady = true
				} else {
					// Connection not ready, close and retry
					conn.Close()
					lastErr = fmt.Errorf("connection to %s at %s timed out waiting for ready state (current: %v)",
						serviceName, target, state)
					log.Printf("❌ %v", lastErr)
					conn = nil
				}
				connectionReady = true // Exit loop
			case <-stateCheckTicker.C:
				state := conn.GetState()
				if state == connectivity.Ready || state == connectivity.Idle {
					connectionReady = true
				} else if state == connectivity.Shutdown {
					// Connection was shut down, close and retry
					conn.Close()
					lastErr = fmt.Errorf("connection to %s at %s was shut down",
						serviceName, target)
					log.Printf("❌ %v", lastErr)
					conn = nil
					connectionReady = true // Exit loop
				}
				// Continue waiting for Connecting or TransientFailure states
			}
		}

		if connectionReady && conn != nil {
			break
		}

		if conn == nil {
			continue
		}
	}

	if conn == nil {
		return nil, fmt.Errorf("failed to establish connection to %s after %d attempts: %w",
			serviceName, c.maxRetries, lastErr)
	}

	// Verify final connection state
	state := conn.GetState()
	if state != connectivity.Ready && state != connectivity.Idle {
		conn.Close()
		return nil, fmt.Errorf("connection to %s at %s is in state %v, expected Ready or Idle",
			serviceName, target, state)
	}

	// Cache the connection and metadata
	c.connsMutex.Lock()
	c.conns[serviceName] = conn
	c.connInfo[serviceName] = &ConnectionInfo{
		ServiceName: serviceName,
		Target:      target,
		State:       state,
		CreatedAt:   time.Now(),
		LastUsed:    time.Now(),
		Conn:        conn,
	}
	c.connsMutex.Unlock()

	log.Printf("✅ Connected to %s service on %s (state: %v)", serviceName, target, state)
	return conn, nil
}

// Close closes all connections
func (c *GRPCClient) Close() error {
	var lastErr error

	c.connsMutex.Lock()
	defer c.connsMutex.Unlock()

	for serviceName, conn := range c.conns {
		if conn != nil {
			if err := conn.Close(); err != nil {
				log.Printf("Error closing connection to %s: %v", serviceName, err)
				lastErr = err
			}
		}
	}

	// Clear the connections map
	c.conns = make(map[string]*grpc.ClientConn)
	c.connInfo = make(map[string]*ConnectionInfo)

	// Stop connection monitor
	close(c.stopMonitor)
	c.monitorWg.Wait()

	if c.daprClient != nil {
		c.daprClient.Close()
	}

	return lastErr
}

// InitializeAllConnections pre-connects to all configured services
// This ensures all connections are established and persisted upfront
func (c *GRPCClient) InitializeAllConnections() error {
	log.Printf("🔄 Initializing all gRPC connections...")

	var wg sync.WaitGroup
	var mu sync.Mutex
	var errors []error

	// Get all configured services
	for serviceName := range Services {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			conn, err := c.GetServiceConnection(name)
			if err != nil {
				mu.Lock()
				errors = append(errors, fmt.Errorf("failed to connect to %s: %w", name, err))
				mu.Unlock()
				log.Printf("❌ Failed to initialize connection to %s: %v", name, err)
			} else {
				log.Printf("✅ Initialized connection to %s", name)
				_ = conn // Connection is cached
			}
		}(serviceName)
	}

	wg.Wait()

	if len(errors) > 0 {
		return fmt.Errorf("failed to initialize some connections: %d errors", len(errors))
	}

	log.Printf("✅ All gRPC connections initialized successfully")
	return nil
}

// GetAllConnections returns information about all active connections
func (c *GRPCClient) GetAllConnections() map[string]*ConnectionInfo {
	c.connsMutex.RLock()
	defer c.connsMutex.RUnlock()

	// Create a copy of connection info
	result := make(map[string]*ConnectionInfo)
	for serviceName, info := range c.connInfo {
		if info != nil {
			// Update state from actual connection
			if conn, exists := c.conns[serviceName]; exists {
				info.State = conn.GetState()
			}
			result[serviceName] = &ConnectionInfo{
				ServiceName: info.ServiceName,
				Target:      info.Target,
				State:       info.State,
				CreatedAt:   info.CreatedAt,
				LastUsed:    info.LastUsed,
				Conn:        info.Conn,
			}
		}
	}

	return result
}

// GetConnectionInfo returns connection information for a specific service
func (c *GRPCClient) GetConnectionInfo(serviceName string) (*ConnectionInfo, error) {
	c.connsMutex.RLock()
	defer c.connsMutex.RUnlock()

	info, exists := c.connInfo[serviceName]
	if !exists {
		return nil, fmt.Errorf("connection info not found for service: %s", serviceName)
	}

	// Update state from actual connection
	if conn, exists := c.conns[serviceName]; exists {
		info.State = conn.GetState()
	}

	return &ConnectionInfo{
		ServiceName: info.ServiceName,
		Target:      info.Target,
		State:       info.State,
		CreatedAt:   info.CreatedAt,
		LastUsed:    info.LastUsed,
		Conn:        info.Conn,
	}, nil
}

// startConnectionMonitor starts a background goroutine to monitor connection health
// and automatically reconnect if connections are lost
func (c *GRPCClient) startConnectionMonitor() {
	c.monitorWg.Add(1)
	go func() {
		defer c.monitorWg.Done()
		ticker := time.NewTicker(30 * time.Second) // Check every 30 seconds
		defer ticker.Stop()

		for {
			select {
			case <-c.stopMonitor:
				return
			case <-ticker.C:
				c.checkAndReconnectConnections()
			}
		}
	}()
}

// checkAndReconnectConnections checks all connections and reconnects if needed
func (c *GRPCClient) checkAndReconnectConnections() {
	c.connsMutex.RLock()
	servicesToCheck := make([]string, 0, len(c.conns))
	for serviceName := range c.conns {
		servicesToCheck = append(servicesToCheck, serviceName)
	}
	c.connsMutex.RUnlock()

	for _, serviceName := range servicesToCheck {
		c.connsMutex.RLock()
		conn, exists := c.conns[serviceName]
		c.connsMutex.RUnlock()

		if !exists {
			continue
		}

		state := conn.GetState()
		// If connection is in a bad state, try to reconnect
		switch state {
		case connectivity.Shutdown, connectivity.TransientFailure:
			log.Printf("⚠️ Connection to %s is in state %v, attempting to reconnect...", serviceName, state)

			// Remove bad connection
			c.connsMutex.Lock()
			delete(c.conns, serviceName)
			delete(c.connInfo, serviceName)
			if conn != nil {
				conn.Close()
			}
			c.connsMutex.Unlock()

			// Attempt to reconnect (this will be done on next GetServiceConnection call)
			// Or we can proactively reconnect here
			go func(name string) {
				_, err := c.GetServiceConnection(name)
				if err != nil {
					log.Printf("❌ Failed to reconnect to %s: %v", name, err)
				} else {
					log.Printf("✅ Successfully reconnected to %s", name)
				}
			}(serviceName)
		case connectivity.Ready, connectivity.Idle:
			// Update last used time for healthy connections
			c.connsMutex.Lock()
			if info, exists := c.connInfo[serviceName]; exists {
				info.State = state
			}
			c.connsMutex.Unlock()
		}
	}
}
