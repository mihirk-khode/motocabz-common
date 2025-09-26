# Central Redis Service

A comprehensive Redis service implementation for the Motocabz application that provides centralized Redis operations across all microservices.

## Features

- **Centralized Redis Management**: Single point of configuration and connection management
- **Service-Specific Adapters**: Tailored Redis operations for each microservice (Trip, Driver, Rider, Identity)
- **Connection Pooling**: Optimized connection management with configurable pool sizes
- **Health Monitoring**: Built-in health checks and monitoring capabilities
- **Pub/Sub Support**: Full publish/subscribe functionality for real-time communication
- **Geolocation Operations**: Specialized geospatial operations for driver location management
- **Atomic Operations**: Support for atomic counters and operations
- **JSON Support**: Built-in JSON serialization/deserialization
- **Key Management**: Consistent key naming with builders and prefixes
- **Multi-Database Support**: Separate databases for different service types

## Architecture

```
Common/redis/
├── redis.go          # Core Redis service implementation
├── config.go         # Configuration and setup utilities
├── init.go          # Initialization and health checks
├── example.go       # Usage examples
├── geolocation.go   # Geolocation-specific operations
└── README.md        # This file

Service-specific adapters:
├── Trip/internal/redis/service.go     # Trip service Redis operations
├── Driver/internal/redis/service.go   # Driver service Redis operations
├── Rider/internal/redis/service.go    # Rider service Redis operations
└── Identity/internal/redis/service.go # Identity service Redis operations
```

## Quick Start

### 1. Environment Setup

Set the following environment variables:

```bash
# Basic Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
REDIS_POOL_SIZE=10
REDIS_MIN_IDLE=5

# Service-specific Databases
REDIS_CACHE_DB=1
REDIS_SESSION_DB=2
REDIS_GEO_DB=3
REDIS_PUBSUB_DB=4

# TTL Configuration (in seconds)
REDIS_CACHE_TTL_SECONDS=3600      # 1 hour
REDIS_SESSION_TTL_SECONDS=86400   # 24 hours
REDIS_GEO_TTL_SECONDS=7200        # 2 hours
REDIS_DEFAULT_TTL_SECONDS=1800    # 30 minutes
```

### 2. Initialize Redis Services

```go
package main

import (
    "github.com/iamarpitzala/motocabz/common/redis"
)

func main() {
    // Setup environment and initialize all Redis services
    redis.SetupRedisEnvironment()
    redis.InitializeAllRedisServices()
    
    // Your application code here
    
    // Cleanup on shutdown
    defer redis.CleanupRedisServices()
}
```

### 3. Using the Central Redis Service

```go
import (
    "context"
    "github.com/iamarpitzala/motocabz/common/redis"
)

func example() {
    ctx := context.Background()
    
    // Get default Redis service
    redisService := redis.GetDefaultRedis()
    
    // Basic operations
    redisService.Set(ctx, "key", "value", time.Minute)
    value, _ := redisService.Get(ctx, "key")
    
    // JSON operations
    data := map[string]interface{}{"name": "John", "age": 30}
    redisService.JSONSet(ctx, "user:123", "", data)
    
    var result map[string]interface{}
    redisService.JSONGet(ctx, "user:123", "", &result)
    
    // Pub/Sub
    redisService.Publish(ctx, "channel", "message")
    pubsub := redisService.Subscribe(ctx, "channel")
    defer pubsub.Close()
}
```

### 4. Using Service-Specific Adapters

#### Trip Service

```go
import tripredis "github.com/iamarpitzala/motocabz/trip/internal/redis"

func tripExample() {
    ctx := context.Background()
    tripService := tripredis.NewTripRedisService()
    
    // Set bidding session
    biddingData := map[string]interface{}{
        "sessionID": "bidding_123",
        "tripID": "trip_456",
        "status": "active",
    }
    tripService.SetBiddingSession(ctx, "bidding_123", biddingData)
    
    // Add bid
    tripService.AddBid(ctx, "bidding_123", "driver_789", 25.50)
    
    // Get bids
    bids, _ := tripService.GetBids(ctx, "bidding_123")
}
```

#### Driver Service

```go
import driverredis "github.com/iamarpitzala/motocabz/driver/internal/redis"

func driverExample() {
    ctx := context.Background()
    driverService := driverredis.NewDriverRedisService()
    
    // Set driver location
    driverService.SetDriverLocation(ctx, "driver_123", 40.7128, -74.0060)
    
    // Set driver status
    driverService.SetDriverAvailable(ctx, "driver_123")
    
    // Add to available drivers
    driverService.AddToAvailableDrivers(ctx, "driver_123")
}
```

#### Rider Service

```go
import riderredis "github.com/iamarpitzala/motocabz/rider/internal/redis"

func riderExample() {
    ctx := context.Background()
    riderService := riderredis.NewRiderRedisService()
    
    // Set rider profile
    profile := map[string]interface{}{
        "name": "John Doe",
        "phone": "+1234567890",
    }
    riderService.SetRiderProfile(ctx, "rider_123", profile)
    
    // Set current trip
    riderService.SetRiderCurrentTrip(ctx, "rider_123", "trip_456")
}
```

#### Identity Service

```go
import identityredis "github.com/iamarpitzala/motocabz/identity/internal/redis"

func identityExample() {
    ctx := context.Background()
    identityService := identityredis.NewIdentityRedisService()
    
    // Set OTP
    identityService.SetOTP(ctx, "+1234567890", "123456", time.Minute*5)
    
    // Set user session
    sessionData := map[string]interface{}{
        "userID": "user_123",
        "role": "rider",
    }
    identityService.SetUserSession(ctx, "session_456", sessionData)
}
```

## Service Configuration

### Database Allocation

| Service | Database | Purpose |
|---------|----------|---------|
| Default | 0 | General operations |
| Cache | 1 | Caching layer |
| Session | 2 | User sessions |
| Geolocation | 3 | Driver locations |
| PubSub | 4 | Message publishing |
| Trip | 5 | Trip-specific data |
| Driver | 6 | Driver-specific data |
| Rider | 7 | Rider-specific data |
| Identity | 8 | Authentication data |

### Key Naming Conventions

Keys follow a consistent pattern with service prefixes:

```
{service}:{type}:{id}:{subtype}

Examples:
- trip:session:bidding_123:data
- driver:location:driver_456:coordinates
- rider:profile:rider_789:preferences
- identity:session:session_123:data
```

## Advanced Features

### Health Monitoring

```go
// Check health of all services
healthStatus := redis.GetServiceHealthStatus()
fmt.Printf("Services healthy: %v\n", healthStatus["healthy"])
```

### Key Builder

```go
keyBuilder := redis.NewKeyBuilder("trip")

// Build various key types
userKey := keyBuilder.BuildUserKey("user_123", "profile")
sessionKey := keyBuilder.BuildSessionKey("session_456", "data")
tripKey := keyBuilder.BuildTripKey("trip_789", "status")
```

### Batch Operations

```go
// Pipeline operations for better performance
pipe := redisService.(*redis.RedisService).client.Pipeline()
pipe.Set(ctx, "key1", "value1", time.Minute)
pipe.Set(ctx, "key2", "value2", time.Minute)
pipe.Set(ctx, "key3", "value3", time.Minute)
_, err := pipe.Exec(ctx)
```

### Geolocation Operations

```go
// Add driver location
geoManager := redis.NewGeoLocationManager(redisService.(*redis.RedisService).client)
geoManager.AddDriverLocation(ctx, "driver_123", 40.7128, -74.0060, map[string]interface{}{
    "vehicleType": "sedan",
    "rating": 4.8,
})

// Find nearby drivers
drivers, _ := geoManager.FindNearbyDrivers(ctx, 40.7128, -74.0060, 5.0, 10)
```

## Performance Considerations

1. **Connection Pooling**: Configure appropriate pool sizes based on your load
2. **TTL Management**: Set appropriate expiration times to prevent memory bloat
3. **Batch Operations**: Use pipelines for multiple operations
4. **Key Design**: Keep keys short and use consistent naming
5. **Memory Management**: Monitor Redis memory usage and set limits

## Monitoring and Maintenance


## Error Handling

All Redis operations return standard Go errors. Common error scenarios:

```go
// Handle connection errors
if err := redisService.Ping(ctx); err != nil {
    log.Printf("Redis connection failed: %v", err)
    // Handle connection failure
}

// Handle key not found
value, err := redisService.Get(ctx, "nonexistent")
if err == redis.Nil {
    // Key doesn't exist
} else if err != nil {
    // Other error
}

// Handle JSON unmarshaling errors
var data map[string]interface{}
if err := redisService.JSONGet(ctx, "key", "", &data); err != nil {
    log.Printf("JSON unmarshal failed: %v", err)
}
```

## Best Practices

1. **Always use context**: Pass context for timeout and cancellation support
2. **Set appropriate TTLs**: Prevent memory leaks with proper expiration
3. **Handle errors gracefully**: Check for Redis.Nil and connection errors
4. **Use service adapters**: Leverage service-specific operations for better organization
5. **Monitor performance**: Use health checks and metrics
6. **Clean up resources**: Always close pub/sub connections and clean up on shutdown

## Troubleshooting

### Common Issues

1. **Connection Refused**: Check Redis server status and connection parameters
2. **Memory Issues**: Monitor TTL settings and implement cleanup routines
3. **Performance Issues**: Check connection pool size and use batch operations
4. **Key Conflicts**: Ensure consistent key naming across services

### Debug Mode

Enable debug logging by setting the log level:

```go
// Enable Redis client debug logging
redisService := redis.GetDefaultRedis().(*redis.RedisService)
redisService.client.Options().OnConnect = func(ctx context.Context, cn *redis.Conn) error {
    return cn.ClientSetName(ctx, "motocabz-client").Err()
}
```

## Contributing

When adding new Redis operations:

1. Follow the existing naming conventions
2. Add appropriate error handling
3. Include TTL management where applicable
4. Update service adapters if needed
5. Add tests and documentation
6. Consider performance implications

## License

This Redis service is part of the Motocabz project and follows the same licensing terms.
