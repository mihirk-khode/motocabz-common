package redis

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

// EnvRedisConfig loads Redis configuration from environment variables
type EnvRedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
	PoolSize int
	MinIdle  int
}

// LoadFromEnv loads Redis configuration from environment variables
func LoadFromEnv() RedisConfig {
	host := getEnvOrDefault("REDIS_HOST", "localhost")
	port := getEnvOrDefault("REDIS_PORT", "6379")
	password := getEnvOrDefault("REDIS_PASSWORD", "")
	db := getEnvIntOrDefault("REDIS_DB", 0)
	poolSize := getEnvIntOrDefault("REDIS_POOL_SIZE", 10)
	minIdle := getEnvIntOrDefault("REDIS_MIN_IDLE", 5)

	return RedisConfig{
		Host:     host,
		Port:     port,
		Password: password,
		DB:       db,
		PoolSize: poolSize,
		MinIdle:  minIdle,
	}
}

// getEnvOrDefault returns environment variable value or default
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvIntOrDefault returns environment variable as int or default
func getEnvIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// InitializeRedisService initializes and configures a Redis service
func InitializeRedisService(config RedisConfig) IRedisService {
	service := NewRedisService(config)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := service.Ping(ctx); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Printf("✅ Redis service initialized successfully (Host: %s:%s, DB: %d)",
		config.Host, config.Port, config.DB)

	return service
}

// InitializeRedisFromEnv initializes Redis service from environment variables
func InitializeRedisFromEnv() IRedisService {
	config := LoadFromEnv()
	return InitializeRedisService(config)
}

// InitializeRedisServices initializes all Redis services
func InitializeRedisServices() {
	// Add default service
	defaultService := InitializeRedisFromEnv()
	SetService("default", defaultService)

	// Add cache service (same Redis, different DB)
	cacheConfig := LoadFromEnv()
	cacheConfig.DB = getEnvIntOrDefault("REDIS_CACHE_DB", 1)
	cacheService := InitializeRedisService(cacheConfig)
	SetService("cache", cacheService)

	// Add session service (same Redis, different DB)
	sessionConfig := LoadFromEnv()
	sessionConfig.DB = getEnvIntOrDefault("REDIS_SESSION_DB", 2)
	sessionService := InitializeRedisService(sessionConfig)
	SetService("session", sessionService)

	// Set global default
	SetDefaultRedis(defaultService)

	log.Printf("✅ Redis services initialized")
}

// RedisServiceConfig represents a service-specific Redis configuration
type RedisServiceConfig struct {
	ServiceName string
	DB          int
	Prefix      string
	TTL         time.Duration
}

// GetServiceConfig returns configuration for a specific service
func GetServiceConfig(serviceName string) RedisServiceConfig {
	switch serviceName {
	case "cache":
		return RedisServiceConfig{
			ServiceName: "cache",
			DB:          getEnvIntOrDefault("REDIS_CACHE_DB", 1),
			Prefix:      "cache:",
			TTL:         time.Duration(getEnvIntOrDefault("REDIS_CACHE_TTL_SECONDS", 3600)) * time.Second,
		}
	case "session":
		return RedisServiceConfig{
			ServiceName: "session",
			DB:          getEnvIntOrDefault("REDIS_SESSION_DB", 2),
			Prefix:      "session:",
			TTL:         time.Duration(getEnvIntOrDefault("REDIS_SESSION_TTL_SECONDS", 86400)) * time.Second,
		}
	case "geolocation":
		return RedisServiceConfig{
			ServiceName: "geolocation",
			DB:          getEnvIntOrDefault("REDIS_GEO_DB", 3),
			Prefix:      "geo:",
			TTL:         time.Duration(getEnvIntOrDefault("REDIS_GEO_TTL_SECONDS", 7200)) * time.Second,
		}
	case "pubsub":
		return RedisServiceConfig{
			ServiceName: "pubsub",
			DB:          getEnvIntOrDefault("REDIS_PUBSUB_DB", 4),
			Prefix:      "pubsub:",
			TTL:         0, // Pub/sub doesn't need TTL
		}
	default:
		return RedisServiceConfig{
			ServiceName: serviceName,
			DB:          0,
			Prefix:      fmt.Sprintf("%s:", serviceName),
			TTL:         time.Duration(getEnvIntOrDefault("REDIS_DEFAULT_TTL_SECONDS", 1800)) * time.Second,
		}
	}
}

// GetServiceRedis returns a Redis service for a specific service type
func GetServiceRedis(serviceName string) IRedisService {
	// Try to get existing service first
	if service := GetService(serviceName); service != nil {
		return service
	}

	// Create new service if not exists
	config := GetServiceConfig(serviceName)
	redisConfig := LoadFromEnv()
	redisConfig.DB = config.DB
	service := InitializeRedisService(redisConfig)
	SetService(serviceName, service)

	return service
}

// RedisKeyBuilder helps build consistent Redis keys
type RedisKeyBuilder struct {
	prefix  string
	service string
}

// NewKeyBuilder creates a new Redis key builder
func NewKeyBuilder(serviceName string) *RedisKeyBuilder {
	config := GetServiceConfig(serviceName)
	return &RedisKeyBuilder{
		prefix:  config.Prefix,
		service: serviceName,
	}
}

// BuildKey builds a Redis key with the configured prefix
func (kb *RedisKeyBuilder) BuildKey(parts ...string) string {
	key := kb.prefix
	for i, part := range parts {
		if i > 0 {
			key += ":"
		}
		key += part
	}
	return key
}

// BuildUserKey builds a user-specific key
func (kb *RedisKeyBuilder) BuildUserKey(userID string, parts ...string) string {
	allParts := append([]string{"user", userID}, parts...)
	return kb.BuildKey(allParts...)
}

// BuildSessionKey builds a session-specific key
func (kb *RedisKeyBuilder) BuildSessionKey(sessionID string, parts ...string) string {
	allParts := append([]string{"session", sessionID}, parts...)
	return kb.BuildKey(allParts...)
}

// BuildDriverKey builds a driver-specific key
func (kb *RedisKeyBuilder) BuildDriverKey(driverID string, parts ...string) string {
	allParts := append([]string{"driver", driverID}, parts...)
	return kb.BuildKey(allParts...)
}

// BuildDriverSubscriptionKey builds a driver subscription key
func (kb *RedisKeyBuilder) BuildDriverSubscriptionKey(driverID string) string {
	return kb.BuildDriverKey(driverID, "subscription")
}

// BuildRiderKey builds a rider-specific key
func (kb *RedisKeyBuilder) BuildRiderKey(riderID string, parts ...string) string {
	allParts := append([]string{"rider", riderID}, parts...)
	return kb.BuildKey(allParts...)
}

// BuildTripKey builds a trip-specific key
func (kb *RedisKeyBuilder) BuildTripKey(tripID string, parts ...string) string {
	allParts := append([]string{"trip", tripID}, parts...)
	return kb.BuildKey(allParts...)
}

// BuildBiddingKey builds a bidding session-specific key
func (kb *RedisKeyBuilder) BuildBiddingKey(sessionID string, parts ...string) string {
	allParts := append([]string{"bidding", sessionID}, parts...)
	return kb.BuildKey(allParts...)
}

func (kb *RedisKeyBuilder) BuildFareKey(sessionID string, parts ...string) string {
	allParts := append([]string{"fare", sessionID}, parts...)
	return kb.BuildKey(allParts...)
}

// RedisHealthChecker provides health check functionality for Redis
type RedisHealthChecker struct {
	service IRedisService
}

// NewHealthChecker creates a new Redis health checker
func NewHealthChecker(service IRedisService) *RedisHealthChecker {
	return &RedisHealthChecker{
		service: service,
	}
}

// CheckHealth performs a comprehensive health check
func (rhc *RedisHealthChecker) CheckHealth(ctx context.Context) error {
	// Test basic connectivity
	if err := rhc.service.Ping(ctx); err != nil {
		return fmt.Errorf("ping failed: %w", err)
	}

	// Test basic operations
	testKey := "health:check:test"
	testValue := "test_value"

	// Test SET
	if err := rhc.service.Set(ctx, testKey, testValue, time.Minute); err != nil {
		return fmt.Errorf("set operation failed: %w", err)
	}

	// Test GET
	value, err := rhc.service.Get(ctx, testKey)
	if err != nil {
		return fmt.Errorf("get operation failed: %w", err)
	}

	if value != testValue {
		return fmt.Errorf("value mismatch: expected %s, got %s", testValue, value)
	}

	// Test DEL
	if err := rhc.service.Del(ctx, testKey); err != nil {
		return fmt.Errorf("del operation failed: %w", err)
	}

	return nil
}

// GetHealthStatus returns detailed health status
func (rhc *RedisHealthChecker) GetHealthStatus(ctx context.Context) map[string]interface{} {
	status := map[string]interface{}{
		"timestamp": time.Now().Unix(),
		"healthy":   false,
		"errors":    []string{},
	}

	if err := rhc.CheckHealth(ctx); err != nil {
		status["errors"] = append(status["errors"].([]string), err.Error())
		return status
	}

	// Get Redis info
	info, err := rhc.service.Info(ctx, "memory", "stats")
	if err != nil {
		status["errors"] = append(status["errors"].([]string), fmt.Sprintf("info failed: %v", err))
	} else {
		status["redis_info"] = info
	}

	status["healthy"] = true
	return status
}
