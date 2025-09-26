package redis

import (
	"context"
	"log"
	"os"
	"time"
)

// InitializeAllRedisServices initializes all Redis services for the application
func InitializeAllRedisServices() {
	log.Println("🚀 Initializing Redis services...")

	// Initialize core services
	InitializeRedisServices()

	// Test services
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Test default service
	if defaultService := GetService("default"); defaultService != nil {
		if err := defaultService.Ping(ctx); err != nil {
			log.Printf("❌ Default Redis service ping failed: %v", err)
		} else {
			log.Println("✅ Default Redis service is healthy")
		}
	}

	log.Println("🎉 Redis services initialization completed!")
}

// SetupRedisEnvironment sets up Redis environment variables with defaults
func SetupRedisEnvironment() {
	// Set default Redis configuration if not already set
	if os.Getenv("REDIS_HOST") == "" {
		os.Setenv("REDIS_HOST", "localhost")
	}
	if os.Getenv("REDIS_PORT") == "" {
		os.Setenv("REDIS_PORT", "6379")
	}
	if os.Getenv("REDIS_PASSWORD") == "" {
		os.Setenv("REDIS_PASSWORD", "")
	}
	if os.Getenv("REDIS_DB") == "" {
		os.Setenv("REDIS_DB", "0")
	}
	if os.Getenv("REDIS_POOL_SIZE") == "" {
		os.Setenv("REDIS_POOL_SIZE", "10")
	}
	if os.Getenv("REDIS_MIN_IDLE") == "" {
		os.Setenv("REDIS_MIN_IDLE", "5")
	}

	// Set service-specific database configurations
	if os.Getenv("REDIS_CACHE_DB") == "" {
		os.Setenv("REDIS_CACHE_DB", "1")
	}
	if os.Getenv("REDIS_SESSION_DB") == "" {
		os.Setenv("REDIS_SESSION_DB", "2")
	}
	if os.Getenv("REDIS_GEO_DB") == "" {
		os.Setenv("REDIS_GEO_DB", "3")
	}
	if os.Getenv("REDIS_PUBSUB_DB") == "" {
		os.Setenv("REDIS_PUBSUB_DB", "4")
	}

	// Set TTL configurations
	if os.Getenv("REDIS_CACHE_TTL_SECONDS") == "" {
		os.Setenv("REDIS_CACHE_TTL_SECONDS", "3600") // 1 hour
	}
	if os.Getenv("REDIS_SESSION_TTL_SECONDS") == "" {
		os.Setenv("REDIS_SESSION_TTL_SECONDS", "86400") // 24 hours
	}
	if os.Getenv("REDIS_GEO_TTL_SECONDS") == "" {
		os.Setenv("REDIS_GEO_TTL_SECONDS", "7200") // 2 hours
	}
	if os.Getenv("REDIS_DEFAULT_TTL_SECONDS") == "" {
		os.Setenv("REDIS_DEFAULT_TTL_SECONDS", "1800") // 30 minutes
	}

	log.Println("🔧 Redis environment variables configured")
}

// GetServiceHealthStatus returns health status for Redis services
func GetServiceHealthStatus() map[string]interface{} {
	status := map[string]interface{}{
		"timestamp": time.Now().Unix(),
		"healthy":   true,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check default service
	if defaultService := GetService("default"); defaultService != nil {
		healthChecker := NewHealthChecker(defaultService)
		if err := healthChecker.CheckHealth(ctx); err != nil {
			status["healthy"] = false
			status["error"] = err.Error()
		}
	}

	return status
}

// CleanupRedisServices gracefully shuts down all Redis services
func CleanupRedisServices() {
	log.Println("🔄 Cleaning up Redis services...")
	CloseAllServices()
	log.Println("✅ Redis services cleanup completed")
}
