package config

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
)

// IEnv provides environment variable access
type IEnv interface {
	Get(key string) string
	GetServerPort() string
	GetDBPort() string
	GetDBUsername() string
	GetDBPassword() string
	GetDBHost() string
	GetDBName() string
	GetDBSSLMODE() string
	GetGinMode() string
	GetMigrationDir() string
	GetDaprHTTPPort() string
	GetDaprGRPCPort() string
	GetDaprAppID() string
	GetDaprAppPort() string
	GetDaprPubsubName() string
	GetNameSpace() string
	GetJWTSecret() string
}

// IBaseConfig provides standard configuration interface
// All services should implement this minimum interface
type IBaseConfig interface {
	Env() IEnv
	DB() interface{} // Service-specific DB client (ent.Client, sql.DB, etc.)
	GRPCClient(serviceName string) *grpc.ClientConn
	GRPCClientWithContext(ctx context.Context, serviceName string) (*grpc.ClientConn, error)
	PubsubName() string
	Close()
}

// IRedisConfig extends base config with Redis support
type IRedisConfig interface {
	IBaseConfig
	Redis() *redis.Client
}

// IS3Config extends base config with S3 support
type IS3Config interface {
	IBaseConfig
	S3Client() *s3.Client
}

// IWebSocketConfig extends base config with WebSocket support
type IWebSocketConfig interface {
	IBaseConfig
	WSManager() IWSManager
}

// IFullConfig includes all infrastructure components
// Use this for services that need Redis, S3, and WebSocket
type IFullConfig interface {
	IRedisConfig
	IS3Config
	IWebSocketConfig
}

// IWSManager provides WebSocket connection management
// Services should implement this interface with their WebSocket manager
// See Common/websocket/websocket.go for a reference implementation
type IWSManager interface {
	// AddConnection adds a new WebSocket connection
	AddConnection(userID, userType string, conn interface{})
	// RemoveConnection removes a WebSocket connection
	RemoveConnection(userID, userType string)
	// SendMessage sends a message to a specific user
	SendMessage(userID, userType string, message interface{}) error
}
