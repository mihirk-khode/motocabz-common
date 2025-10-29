package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
	PoolSize int
	MinIdle  int
}

// IRedisService defines the interface for Redis operations
type IRedisService interface {
	// Connection management
	Ping(ctx context.Context) error
	Close() error

	// Basic operations
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Del(ctx context.Context, keys ...string) error
	Exists(ctx context.Context, keys ...string) (int64, error)
	Expire(ctx context.Context, key string, expiration time.Duration) error
	TTL(ctx context.Context, key string) (time.Duration, error)

	// Hash operations
	HGet(ctx context.Context, key, field string) (string, error)
	HSet(ctx context.Context, key string, values ...interface{}) error
	HGetAll(ctx context.Context, key string) (map[string]string, error)
	HDel(ctx context.Context, key string, fields ...string) error
	HExists(ctx context.Context, key, field string) (bool, error)

	// List operations
	LPush(ctx context.Context, key string, values ...interface{}) error
	RPush(ctx context.Context, key string, values ...interface{}) error
	LPop(ctx context.Context, key string) (string, error)
	RPop(ctx context.Context, key string) (string, error)
	LLen(ctx context.Context, key string) (int64, error)
	LRange(ctx context.Context, key string, start, stop int64) ([]string, error)

	// Set operations
	SAdd(ctx context.Context, key string, members ...interface{}) error
	SRem(ctx context.Context, key string, members ...interface{}) error
	SMembers(ctx context.Context, key string) ([]string, error)
	SIsMember(ctx context.Context, key string, member interface{}) (bool, error)
	SCard(ctx context.Context, key string) (int64, error)

	// Sorted set operations
	ZAdd(ctx context.Context, key string, members ...redis.Z) error
	ZRem(ctx context.Context, key string, members ...interface{}) error
	ZRange(ctx context.Context, key string, start, stop int64) ([]string, error)
	ZRangeByScore(ctx context.Context, key string, opt *redis.ZRangeBy) ([]string, error)
	ZCard(ctx context.Context, key string) (int64, error)
	ZScore(ctx context.Context, key, member string) (float64, error)

	// Pub/Sub operations
	Publish(ctx context.Context, channel string, message interface{}) error
	Subscribe(ctx context.Context, channels ...string) *redis.PubSub
	PSubscribe(ctx context.Context, channels ...string) *redis.PubSub

	// Geo operations (wrapper around existing geolocation service)
	GeoAdd(ctx context.Context, key string, geoLocation ...*redis.GeoLocation) error
	GeoRadius(ctx context.Context, key string, longitude, latitude float64, query *redis.GeoRadiusQuery) ([]redis.GeoLocation, error)
	GeoPos(ctx context.Context, key string, members ...string) ([]*redis.GeoPos, error)

	// Atomic operations
	Incr(ctx context.Context, key string) (int64, error)
	Decr(ctx context.Context, key string) (int64, error)
	IncrBy(ctx context.Context, key string, value int64) (int64, error)
	DecrBy(ctx context.Context, key string, value int64) (int64, error)

	// JSON operations
	JSONSet(ctx context.Context, key, path string, value interface{}) error
	JSONGet(ctx context.Context, key, path string, dest interface{}) error
	JSONDel(ctx context.Context, key, path string) error

	// Utility operations
	Keys(ctx context.Context, pattern string) ([]string, error)
	Scan(ctx context.Context, cursor uint64, match string, count int64) ([]string, uint64, error)
	FlushDB(ctx context.Context) error
	Info(ctx context.Context, section ...string) (string, error)
}

// RedisService implements the Redis service interface
type RedisService struct {
	client *redis.Client
	config RedisConfig
}

// NewRedisService creates a new Redis service instance
func NewRedisService(config RedisConfig) IRedisService {
	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", config.Host, config.Port),
		Password:     config.Password,
		DB:           config.DB,
		PoolSize:     config.PoolSize,
		MinIdleConns: config.MinIdle,
		MaxRetries:   3,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		log.Printf("----------------------------------------------------------------Error connecting to Redis: %v--------------------------------------------------------------------", err)
	}
	log.Println("----------------------------------------------------------------Redis is up!!--------------------------------------------------------------------")
	fmt.Println("----------------------------------------------------------------Redis client created successfully--------------------------------------------------------------------")

	return &RedisService{
		client: client,
		config: config,
	}
}

// NewRedisServiceWithClient creates a Redis service with an existing client
func NewRedisServiceWithClient(client *redis.Client) IRedisService {
	return &RedisService{
		client: client,
	}
}

// Ping tests the Redis connection
func (rs *RedisService) Ping(ctx context.Context) error {
	return rs.client.Ping(ctx).Err()
}

// Close closes the Redis connection
func (rs *RedisService) Close() error {
	return rs.client.Close()
}

// Get retrieves a value by key
func (rs *RedisService) Get(ctx context.Context, key string) (string, error) {
	return rs.client.Get(ctx, key).Result()
}

// Set stores a value with optional expiration
func (rs *RedisService) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return rs.client.Set(ctx, key, value, expiration).Err()
}

// Del deletes one or more keys
func (rs *RedisService) Del(ctx context.Context, keys ...string) error {
	return rs.client.Del(ctx, keys...).Err()
}

// Exists checks if keys exist
func (rs *RedisService) Exists(ctx context.Context, keys ...string) (int64, error) {
	return rs.client.Exists(ctx, keys...).Result()
}

// Expire sets expiration for a key
func (rs *RedisService) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return rs.client.Expire(ctx, key, expiration).Err()
}

// TTL returns the time to live for a key
func (rs *RedisService) TTL(ctx context.Context, key string) (time.Duration, error) {
	return rs.client.TTL(ctx, key).Result()
}

// HGet retrieves a field from a hash
func (rs *RedisService) HGet(ctx context.Context, key, field string) (string, error) {
	return rs.client.HGet(ctx, key, field).Result()
}

// HSet sets fields in a hash
func (rs *RedisService) HSet(ctx context.Context, key string, values ...interface{}) error {
	return rs.client.HSet(ctx, key, values...).Err()
}

// HGetAll retrieves all fields from a hash
func (rs *RedisService) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return rs.client.HGetAll(ctx, key).Result()
}

// HDel deletes fields from a hash
func (rs *RedisService) HDel(ctx context.Context, key string, fields ...string) error {
	return rs.client.HDel(ctx, key, fields...).Err()
}

// HExists checks if a field exists in a hash
func (rs *RedisService) HExists(ctx context.Context, key, field string) (bool, error) {
	return rs.client.HExists(ctx, key, field).Result()
}

// LPush pushes values to the left of a list
func (rs *RedisService) LPush(ctx context.Context, key string, values ...interface{}) error {
	return rs.client.LPush(ctx, key, values...).Err()
}

// RPush pushes values to the right of a list
func (rs *RedisService) RPush(ctx context.Context, key string, values ...interface{}) error {
	return rs.client.RPush(ctx, key, values...).Err()
}

// LPop pops a value from the left of a list
func (rs *RedisService) LPop(ctx context.Context, key string) (string, error) {
	return rs.client.LPop(ctx, key).Result()
}

// RPop pops a value from the right of a list
func (rs *RedisService) RPop(ctx context.Context, key string) (string, error) {
	return rs.client.RPop(ctx, key).Result()
}

// LLen returns the length of a list
func (rs *RedisService) LLen(ctx context.Context, key string) (int64, error) {
	return rs.client.LLen(ctx, key).Result()
}

// LRange returns a range of elements from a list
func (rs *RedisService) LRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return rs.client.LRange(ctx, key, start, stop).Result()
}

// SAdd adds members to a set
func (rs *RedisService) SAdd(ctx context.Context, key string, members ...interface{}) error {
	return rs.client.SAdd(ctx, key, members...).Err()
}

// SRem removes members from a set
func (rs *RedisService) SRem(ctx context.Context, key string, members ...interface{}) error {
	return rs.client.SRem(ctx, key, members...).Err()
}

// SMembers returns all members of a set
func (rs *RedisService) SMembers(ctx context.Context, key string) ([]string, error) {
	return rs.client.SMembers(ctx, key).Result()
}

// SIsMember checks if a member exists in a set
func (rs *RedisService) SIsMember(ctx context.Context, key string, member interface{}) (bool, error) {
	return rs.client.SIsMember(ctx, key, member).Result()
}

// SCard returns the cardinality of a set
func (rs *RedisService) SCard(ctx context.Context, key string) (int64, error) {
	return rs.client.SCard(ctx, key).Result()
}

// ZAdd adds members to a sorted set
func (rs *RedisService) ZAdd(ctx context.Context, key string, members ...redis.Z) error {
	return rs.client.ZAdd(ctx, key, members...).Err()
}

// ZRem removes members from a sorted set
func (rs *RedisService) ZRem(ctx context.Context, key string, members ...interface{}) error {
	return rs.client.ZRem(ctx, key, members...).Err()
}

// ZRange returns a range of members from a sorted set
func (rs *RedisService) ZRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return rs.client.ZRange(ctx, key, start, stop).Result()
}

// ZRangeByScore returns members from a sorted set by score
func (rs *RedisService) ZRangeByScore(ctx context.Context, key string, opt *redis.ZRangeBy) ([]string, error) {
	return rs.client.ZRangeByScore(ctx, key, opt).Result()
}

// ZCard returns the cardinality of a sorted set
func (rs *RedisService) ZCard(ctx context.Context, key string) (int64, error) {
	return rs.client.ZCard(ctx, key).Result()
}

// ZScore returns the score of a member in a sorted set
func (rs *RedisService) ZScore(ctx context.Context, key, member string) (float64, error) {
	return rs.client.ZScore(ctx, key, member).Result()
}

// Publish publishes a message to a channel
func (rs *RedisService) Publish(ctx context.Context, channel string, message interface{}) error {
	return rs.client.Publish(ctx, channel, message).Err()
}

// Subscribe subscribes to channels
func (rs *RedisService) Subscribe(ctx context.Context, channels ...string) *redis.PubSub {
	return rs.client.Subscribe(ctx, channels...)
}

// PSubscribe subscribes to channel patterns
func (rs *RedisService) PSubscribe(ctx context.Context, channels ...string) *redis.PubSub {
	return rs.client.PSubscribe(ctx, channels...)
}

// GeoAdd adds geospatial items
func (rs *RedisService) GeoAdd(ctx context.Context, key string, geoLocation ...*redis.GeoLocation) error {
	return rs.client.GeoAdd(ctx, key, geoLocation...).Err()
}

// GeoRadius searches for geospatial items within radius
func (rs *RedisService) GeoRadius(ctx context.Context, key string, longitude, latitude float64, query *redis.GeoRadiusQuery) ([]redis.GeoLocation, error) {
	return rs.client.GeoRadius(ctx, key, longitude, latitude, query).Result()
}

// GeoPos returns geospatial positions
func (rs *RedisService) GeoPos(ctx context.Context, key string, members ...string) ([]*redis.GeoPos, error) {
	return rs.client.GeoPos(ctx, key, members...).Result()
}

// Incr increments a key by 1
func (rs *RedisService) Incr(ctx context.Context, key string) (int64, error) {
	return rs.client.Incr(ctx, key).Result()
}

// Decr decrements a key by 1
func (rs *RedisService) Decr(ctx context.Context, key string) (int64, error) {
	return rs.client.Decr(ctx, key).Result()
}

// IncrBy increments a key by a value
func (rs *RedisService) IncrBy(ctx context.Context, key string, value int64) (int64, error) {
	return rs.client.IncrBy(ctx, key, value).Result()
}

// DecrBy decrements a key by a value
func (rs *RedisService) DecrBy(ctx context.Context, key string, value int64) (int64, error) {
	return rs.client.DecrBy(ctx, key, value).Result()
}

// JSONSet sets a JSON value
func (rs *RedisService) JSONSet(ctx context.Context, key, path string, value interface{}) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	return rs.Set(ctx, key, string(jsonData), 0)
}

// JSONGet retrieves and unmarshals a JSON value
func (rs *RedisService) JSONGet(ctx context.Context, key, path string, dest interface{}) error {
	jsonStr, err := rs.Get(ctx, key)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(jsonStr), dest)
}

// JSONDel deletes a JSON key
func (rs *RedisService) JSONDel(ctx context.Context, key, path string) error {
	return rs.Del(ctx, key)
}

// Keys returns all keys matching a pattern
func (rs *RedisService) Keys(ctx context.Context, pattern string) ([]string, error) {
	return rs.client.Keys(ctx, pattern).Result()
}

// Scan iterates over keys matching a pattern
func (rs *RedisService) Scan(ctx context.Context, cursor uint64, match string, count int64) ([]string, uint64, error) {
	return rs.client.Scan(ctx, cursor, match, count).Result()
}

// FlushDB removes all keys from the current database
func (rs *RedisService) FlushDB(ctx context.Context) error {
	return rs.client.FlushDB(ctx).Err()
}

// Info returns Redis server information
func (rs *RedisService) Info(ctx context.Context, section ...string) (string, error) {
	return rs.client.Info(ctx, section...).Result()
}

// Simple service registry
var (
	services = make(map[string]IRedisService)
	mu       sync.RWMutex
)

// GetService retrieves a Redis service by name
func GetService(name string) IRedisService {
	mu.RLock()
	defer mu.RUnlock()
	return services[name]
}

// SetService sets a Redis service by name
func SetService(name string, service IRedisService) {
	mu.Lock()
	defer mu.Unlock()
	services[name] = service
}

// CloseAllServices closes all Redis services
func CloseAllServices() {
	mu.Lock()
	defer mu.Unlock()
	for _, service := range services {
		service.Close()
	}
	services = make(map[string]IRedisService)
}

// GetDefaultRedis returns the default Redis service
func GetDefaultRedis() IRedisService {
	return GetService("default")
}

// SetDefaultRedis sets the default Redis service
func SetDefaultRedis(service IRedisService) {
	SetService("default", service)
}
