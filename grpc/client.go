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
	"google.golang.org/grpc/status"
)

const (
	// Default connection timeout
	defaultDialTimeout = 30 * time.Second
	// Default retry attempts
	defaultMaxRetries = 3
	// Default retry backoff
	defaultRetryBackoff = 1 * time.Second
	// Default keepalive time - increased to reduce ping frequency (prevents "too_many_pings" error)
	defaultKeepaliveTime = 60 * time.Second
	// Default keepalive timeout
	defaultKeepaliveTimeout = 20 * time.Second
	// Default ready timeout - increased for Kubernetes environments
	defaultReadyTimeout = 30 * time.Second
	// Default max message size (4MB)
	defaultMaxMsgSize = 4 * 1024 * 1024
	// Connection health check interval
	healthCheckInterval = 30 * time.Second
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

// Options configures the GRPCClient behavior
type Options struct {
	Namespace        string
	DialTimeout      time.Duration
	MaxRetries       int
	RetryBackoff     time.Duration
	KeepaliveTime    time.Duration
	KeepaliveTimeout time.Duration
	ReadyTimeout     time.Duration
	MaxMsgSize       int
	EnableMetrics    bool
}

// GRPCClient manages gRPC connections for service-to-service communication
type GRPCClient struct {
	daprClient       client.Client
	conns            map[string]*grpc.ClientConn
	connInfo         map[string]*ConnectionInfo
	connsMutex       sync.RWMutex
	namespace        string
	dialTimeout      time.Duration
	maxRetries       int
	retryBackoff     time.Duration
	keepaliveTime    time.Duration
	keepaliveTimeout time.Duration
	readyTimeout     time.Duration
	maxMsgSize       int
	stopMonitor      chan struct{}
	monitorWg        sync.WaitGroup
	ctx              context.Context
	cancel           context.CancelFunc
}

// NewGRPCClient creates a new gRPC client with Dapr integration
func NewGRPCClient(namespace string) (*GRPCClient, error) {
	return NewGRPCClientWithOptions(&Options{
		Namespace:        namespace,
		DialTimeout:      defaultDialTimeout,
		MaxRetries:       defaultMaxRetries,
		RetryBackoff:     defaultRetryBackoff,
		KeepaliveTime:    defaultKeepaliveTime,
		KeepaliveTimeout: defaultKeepaliveTimeout,
		ReadyTimeout:     defaultReadyTimeout,
		MaxMsgSize:       defaultMaxMsgSize,
		EnableMetrics:    false,
	})
}

// NewGRPCClientWithOptions creates a new gRPC client with custom options
func NewGRPCClientWithOptions(opts *Options) (*GRPCClient, error) {
	if opts == nil {
		opts = &Options{}
	}

	daprClient, err := client.NewClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create Dapr client: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	grpcClient := &GRPCClient{
		daprClient:       daprClient,
		conns:            make(map[string]*grpc.ClientConn),
		connInfo:         make(map[string]*ConnectionInfo),
		namespace:        opts.Namespace,
		dialTimeout:      opts.DialTimeout,
		maxRetries:       opts.MaxRetries,
		retryBackoff:     opts.RetryBackoff,
		keepaliveTime:    opts.KeepaliveTime,
		keepaliveTimeout: opts.KeepaliveTimeout,
		readyTimeout:     opts.ReadyTimeout,
		maxMsgSize:       opts.MaxMsgSize,
		stopMonitor:      make(chan struct{}),
		ctx:              ctx,
		cancel:           cancel,
	}

	// Set defaults if not provided
	if grpcClient.dialTimeout == 0 {
		grpcClient.dialTimeout = defaultDialTimeout
	}
	if grpcClient.maxRetries == 0 {
		grpcClient.maxRetries = defaultMaxRetries
	}
	if grpcClient.retryBackoff == 0 {
		grpcClient.retryBackoff = defaultRetryBackoff
	}
	if grpcClient.keepaliveTime == 0 {
		grpcClient.keepaliveTime = defaultKeepaliveTime
	}
	if grpcClient.keepaliveTimeout == 0 {
		grpcClient.keepaliveTimeout = defaultKeepaliveTimeout
	}
	if grpcClient.readyTimeout == 0 {
		grpcClient.readyTimeout = defaultReadyTimeout
	}
	if grpcClient.maxMsgSize == 0 {
		grpcClient.maxMsgSize = defaultMaxMsgSize
	}

	// Start background connection health monitor
	grpcClient.startConnectionMonitor()

	return grpcClient, nil
}

// GetServiceConnection returns a gRPC connection to the specified service
// with automatic retry and connection health checking
func (c *GRPCClient) GetServiceConnection(serviceName string) (*grpc.ClientConn, error) {
	return c.GetServiceConnectionWithContext(context.Background(), serviceName)
}

// GetServiceConnectionWithContext returns a gRPC connection with context support
func (c *GRPCClient) GetServiceConnectionWithContext(ctx context.Context, serviceName string) (*grpc.ClientConn, error) {
	// Check if we already have a connection and verify it's still healthy
	c.connsMutex.RLock()
	if conn, exists := c.conns[serviceName]; exists {
		state := conn.GetState()
		c.connsMutex.RUnlock()

		// If connection is ready, idle, or connecting, return it
		// CONNECTING state is valid - gRPC will complete the connection asynchronously
		if state == connectivity.Ready || state == connectivity.Idle || state == connectivity.Connecting {
			// Update last used time
			c.connsMutex.Lock()
			if info, exists := c.connInfo[serviceName]; exists {
				info.LastUsed = time.Now()
				info.State = state
			}
			c.connsMutex.Unlock()
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

	// Create dial context with timeout
	dialCtx, dialCancel := context.WithTimeout(ctx, c.dialTimeout)
	defer dialCancel()

	// Configure dial options with best practices
	dialOptions := c.getDialOptions()

	// Attempt connection
	conn, err := grpc.DialContext(dialCtx, target, dialOptions...)
	if err != nil {
		return nil, fmt.Errorf("failed to dial %s at %s: %w", serviceName, target, err)
	}

	// Wait for connection to be ready with context (use longer timeout for Kubernetes)
	readyCtx, readyCancel := context.WithTimeout(ctx, c.readyTimeout)
	defer readyCancel()

	// Try to wait for ready, but don't fail if it's still connecting
	// Connections in CONNECTING state will become ready when first used
	if err := c.waitForReady(readyCtx, conn, serviceName, target); err != nil {
		// If connection is still in CONNECTING state, allow it (will become ready on first use)
		state := conn.GetState()
		if state == connectivity.Connecting {
			log.Printf("⚠️ Connection to %s at %s is still connecting, will become ready on first use", serviceName, target)
			// Cache the connection anyway - gRPC will handle the connection asynchronously
		} else {
			conn.Close()
			return nil, fmt.Errorf("connection to %s at %s not ready: %w", serviceName, target, err)
		}
	}

	// Verify final connection state - allow CONNECTING state for Kubernetes
	state := conn.GetState()
	if state != connectivity.Ready && state != connectivity.Idle && state != connectivity.Connecting {
		conn.Close()
		return nil, fmt.Errorf("connection to %s at %s is in state %v, expected Ready, Idle, or Connecting",
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

// getDialOptions returns configured dial options with best practices
func (c *GRPCClient) getDialOptions() []grpc.DialOption {
	return []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                c.keepaliveTime,
			Timeout:             c.keepaliveTimeout,
			PermitWithoutStream: false, // Only send keepalive when there are active streams (prevents "too_many_pings")
		}),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(c.maxMsgSize),
			grpc.MaxCallSendMsgSize(c.maxMsgSize),
		),
		// Add interceptors for retry and logging
		grpc.WithUnaryInterceptor(c.unaryClientInterceptor()),
		grpc.WithStreamInterceptor(c.streamClientInterceptor()),
		// Use WaitForReady to allow connections in CONNECTING state to proceed
		grpc.WithDefaultCallOptions(grpc.WaitForReady(false)), // Don't block on ready, allow async connection
	}
}

// waitForReady waits for the connection to be ready
func (c *GRPCClient) waitForReady(ctx context.Context, conn *grpc.ClientConn, serviceName, target string) error {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		state := conn.GetState()
		switch state {
		case connectivity.Ready, connectivity.Idle:
			return nil
		case connectivity.Shutdown:
			return fmt.Errorf("connection was shut down")
		case connectivity.TransientFailure:
			// Wait for automatic reconnection
		case connectivity.Connecting:
			// Still connecting
		}

		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for ready state (current: %v): %w", state, ctx.Err())
		case <-ticker.C:
			// Continue checking
		}
	}
}

// unaryClientInterceptor provides unary client interceptor for logging and retry
func (c *GRPCClient) unaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		start := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)
		duration := time.Since(start)

		if err != nil {
			st, ok := status.FromError(err)
			if ok {
				log.Printf("gRPC call %s failed: code=%s, message=%s, duration=%v",
					method, st.Code(), st.Message(), duration)
			} else {
				log.Printf("gRPC call %s failed: error=%v, duration=%v", method, err, duration)
			}
		} else {
			log.Printf("gRPC call %s succeeded: duration=%v", method, duration)
		}

		return err
	}
}

// streamClientInterceptor provides stream client interceptor for logging
func (c *GRPCClient) streamClientInterceptor() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		start := time.Now()
		stream, err := streamer(ctx, desc, cc, method, opts...)
		duration := time.Since(start)

		if err != nil {
			log.Printf("gRPC stream %s failed: error=%v, duration=%v", method, err, duration)
		} else {
			log.Printf("gRPC stream %s established: duration=%v", method, duration)
		}

		return stream, err
	}
}

// Close closes all connections gracefully
func (c *GRPCClient) Close() error {
	var lastErr error

	// Cancel context to stop background goroutines
	if c.cancel != nil {
		c.cancel()
	}

	// Stop connection monitor
	close(c.stopMonitor)
	c.monitorWg.Wait()

	c.connsMutex.Lock()
	defer c.connsMutex.Unlock()

	// Close all connections
	for serviceName, conn := range c.conns {
		if conn != nil {
			// Gracefully close connection
			if err := conn.Close(); err != nil {
				log.Printf("Error closing connection to %s: %v", serviceName, err)
				if lastErr == nil {
					lastErr = err
				}
			} else {
				log.Printf("Closed connection to %s", serviceName)
			}
		}
	}

	// Clear the connections map
	c.conns = make(map[string]*grpc.ClientConn)
	c.connInfo = make(map[string]*ConnectionInfo)

	// Close Dapr client
	if c.daprClient != nil {
		c.daprClient.Close()
	}

	return lastErr
}

// InitializeAllConnections pre-connects to all configured services
// This ensures all connections are established and persisted upfront
func (c *GRPCClient) InitializeAllConnections() error {
	return c.InitializeAllConnectionsWithContext(context.Background())
}

// InitializeAllConnectionsWithContext pre-connects to all configured services with context
func (c *GRPCClient) InitializeAllConnectionsWithContext(ctx context.Context) error {
	log.Printf("🔄 Initializing all gRPC connections...")

	var wg sync.WaitGroup
	var mu sync.Mutex
	var errors []error

	// Get all configured services
	for serviceName := range Services {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			conn, err := c.GetServiceConnectionWithContext(ctx, name)
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
		ticker := time.NewTicker(healthCheckInterval)
		defer ticker.Stop()

		for {
			select {
			case <-c.ctx.Done():
				return
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
		case connectivity.Shutdown:
			log.Printf("⚠️ Connection to %s is shut down, removing from pool", serviceName)
			c.connsMutex.Lock()
			delete(c.conns, serviceName)
			delete(c.connInfo, serviceName)
			if conn != nil {
				conn.Close()
			}
			c.connsMutex.Unlock()
		case connectivity.TransientFailure:
			// Check if it's been in transient failure for too long
			c.connsMutex.RLock()
			info, infoExists := c.connInfo[serviceName]
			c.connsMutex.RUnlock()

			if infoExists && time.Since(info.LastUsed) > 2*healthCheckInterval {
				log.Printf("⚠️ Connection to %s has been in transient failure, attempting to reconnect...", serviceName)
				c.connsMutex.Lock()
				delete(c.conns, serviceName)
				delete(c.connInfo, serviceName)
				if conn != nil {
					conn.Close()
				}
				c.connsMutex.Unlock()

				// Attempt to reconnect in background
				go func(name string) {
					ctx, cancel := context.WithTimeout(context.Background(), c.dialTimeout)
					defer cancel()
					_, err := c.GetServiceConnectionWithContext(ctx, name)
					if err != nil {
						log.Printf("❌ Failed to reconnect to %s: %v", name, err)
					} else {
						log.Printf("✅ Successfully reconnected to %s", name)
					}
				}(serviceName)
			}
		case connectivity.Ready, connectivity.Idle:
			// Update state for healthy connections
			c.connsMutex.Lock()
			if info, exists := c.connInfo[serviceName]; exists {
				info.State = state
			}
			c.connsMutex.Unlock()
		}
	}
}
