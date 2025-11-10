package grpc

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

const (
	defaultDialTimeout      = 12 * time.Second
	defaultMaxRetries       = 3
	defaultRetryBackoff     = 1 * time.Second
	defaultKeepaliveTime    = 60 * time.Second
	defaultKeepaliveTimeout = 20 * time.Second
	defaultReadyTimeout     = 15 * time.Second
	defaultMaxMsgSize       = 4 * 1024 * 1024
	healthCheckInterval     = 30 * time.Second
)

type Options struct {
	Namespace        string
	DialTimeout      time.Duration
	MaxRetries       int
	RetryBackoff     time.Duration
	KeepaliveTime    time.Duration
	KeepaliveTimeout time.Duration
	ReadyTimeout     time.Duration
	MaxMsgSize       int
}

type GRPCClient struct {
	conns            map[string]*grpc.ClientConn
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
	})
}

func NewGRPCClientWithOptions(opts *Options) (*GRPCClient, error) {
	if opts == nil {
		opts = &Options{}
	}

	ctx, cancel := context.WithCancel(context.Background())

	grpcClient := &GRPCClient{
		conns:            make(map[string]*grpc.ClientConn),
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

	grpcClient.startConnectionMonitor()

	return grpcClient, nil
}

func (c *GRPCClient) GetConn(serviceName string) *grpc.ClientConn {
	conn, err := c.GetServiceConnection(serviceName)
	if err != nil {
		log.Printf("Failed to get connection to %s: %v", serviceName, err)
		return nil
	}
	return conn
}

func (c *GRPCClient) GetServiceConnection(serviceName string) (*grpc.ClientConn, error) {
	return c.GetServiceConnectionWithContext(context.Background(), serviceName)
}

func (c *GRPCClient) GetServiceConnectionWithContext(ctx context.Context, serviceName string) (*grpc.ClientConn, error) {
	c.connsMutex.RLock()
	if conn, exists := c.conns[serviceName]; exists {
		state := conn.GetState()
		c.connsMutex.RUnlock()
		if state == connectivity.Ready || state == connectivity.Idle || state == connectivity.Connecting {
			return conn, nil
		}
		c.connsMutex.Lock()
		delete(c.conns, serviceName)
		if conn != nil {
			conn.Close()
		}
		c.connsMutex.Unlock()
	} else {
		c.connsMutex.RUnlock()
	}

	config, exists := GetServiceConfig(serviceName)
	if !exists {
		return nil, fmt.Errorf("service %s not found in configuration", serviceName)
	}

	var target string
	if c.namespace != "" {
		target = fmt.Sprintf("%s.%s.svc.cluster.local:%s", config.Name, c.namespace, config.Port)
	} else {
		target = fmt.Sprintf("localhost:%s", config.Port)
	}

	dialOptions := c.getDialOptions()

	var conn *grpc.ClientConn
	var lastErr error
	for attempt := 0; attempt < c.maxRetries; attempt++ {
		dialCtx, dialCancel := context.WithTimeout(ctx, c.dialTimeout)
		conn, lastErr = grpc.DialContext(dialCtx, target, dialOptions...)
		dialCancel()

		if lastErr == nil {
			break
		}
		if attempt < c.maxRetries-1 {
			time.Sleep(c.calculateExponentialBackoff(attempt))
		}
	}

	if lastErr != nil {
		return nil, fmt.Errorf("failed to dial %s at %s after %d retries: %w", serviceName, target, c.maxRetries, lastErr)
	}

	readyCtx, readyCancel := context.WithTimeout(ctx, c.readyTimeout)
	defer readyCancel()

	if err := c.waitForReady(readyCtx, conn, serviceName, target); err != nil {
		state := conn.GetState()
		if state == connectivity.Connecting {
			// allow connecting
		} else {
			conn.Close()
			return nil, fmt.Errorf("connection to %s at %s not ready: %w", serviceName, target, err)
		}
	}

	state := conn.GetState()
	if state != connectivity.Ready && state != connectivity.Idle && state != connectivity.Connecting {
		conn.Close()
		return nil, fmt.Errorf("connection to %s at %s is in state %v, expected Ready, Idle, or Connecting", serviceName, target, state)
	}

	c.connsMutex.Lock()
	c.conns[serviceName] = conn
	c.connsMutex.Unlock()

	return conn, nil
}

func (c *GRPCClient) getDialOptions() []grpc.DialOption {
	return []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                c.keepaliveTime,
			Timeout:             c.keepaliveTimeout,
			PermitWithoutStream: false,
		}),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(c.maxMsgSize),
			grpc.MaxCallSendMsgSize(c.maxMsgSize),
		),
		grpc.WithDefaultCallOptions(grpc.WaitForReady(false)),
	}
}

func (c *GRPCClient) calculateExponentialBackoff(attempt int) time.Duration {
	backoff := c.retryBackoff * time.Duration(1<<uint(attempt))
	if backoff > 10*time.Second {
		backoff = 10 * time.Second
	}
	return backoff
}

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
		}
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for ready state (current: %v): %w", state, ctx.Err())
		case <-ticker.C:
		}
	}
}

func (c *GRPCClient) Close() error {
	var lastErr error
	if c.cancel != nil {
		c.cancel()
	}
	close(c.stopMonitor)
	c.monitorWg.Wait()
	c.connsMutex.Lock()
	defer c.connsMutex.Unlock()
	for _, conn := range c.conns {
		if conn != nil {
			if err := conn.Close(); err != nil {
				if lastErr == nil {
					lastErr = err
				}
			}
		}
	}
	c.conns = make(map[string]*grpc.ClientConn)
	return lastErr
}

func (c *GRPCClient) InitializeAllConnections() error {
	return c.InitializeAllConnectionsWithContext(context.Background())
}

func (c *GRPCClient) InitializeAllConnectionsWithContext(ctx context.Context) error {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var errors []error
	for serviceName := range Services {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			_, err := c.GetServiceConnectionWithContext(ctx, name)
			if err != nil {
				mu.Lock()
				errors = append(errors, fmt.Errorf("failed to connect to %s: %w", name, err))
				mu.Unlock()
			}
		}(serviceName)
	}
	wg.Wait()
	if len(errors) > 0 {
		return fmt.Errorf("failed to initialize some connections: %d errors", len(errors))
	}
	return nil
}

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
		switch state {
		case connectivity.Shutdown, connectivity.TransientFailure:
			c.connsMutex.Lock()
			delete(c.conns, serviceName)
			if conn != nil {
				conn.Close()
			}
			c.connsMutex.Unlock()
			if state == connectivity.TransientFailure {
				go func(name string) {
					ctx, cancel := context.WithTimeout(context.Background(), c.dialTimeout)
					defer cancel()
					c.GetServiceConnectionWithContext(ctx, name)
				}(serviceName)
			}
		}
	}
}
