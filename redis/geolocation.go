package redis

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

// GeoLocation represents a geographical location with metadata
type GeoLocation struct {
	Member    string                 `json:"member"`
	Latitude  float64                `json:"latitude"`
	Longitude float64                `json:"longitude"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// DriverLocation represents a driver's current location and status
type DriverLocation struct {
	DriverID    string    `json:"driverId"`
	Latitude    float64   `json:"latitude"`
	Longitude   float64   `json:"longitude"`
	Status      string    `json:"status"` // available, busy, offline
	LastSeen    time.Time `json:"lastSeen"`
	VehicleType string    `json:"vehicleType,omitempty"`
	Rating      float64   `json:"rating,omitempty"`
	Distance    float64   `json:"distance,omitempty"` // calculated distance from rider
}

// IGeoLocationManager defines the interface for Redis geolocation operations
type IGeoLocationManager interface {
	// Driver location management
	AddDriverLocation(ctx context.Context, driverID string, lat, lng float64, metadata map[string]interface{}) error
	UpdateDriverLocation(ctx context.Context, driverID string, lat, lng float64) error
	RemoveDriverLocation(ctx context.Context, driverID string) error
	GetDriverLocation(ctx context.Context, driverID string) (*DriverLocation, error)

	// Driver discovery
	FindNearbyDrivers(ctx context.Context, lat, lng float64, radius float64, limit int) ([]DriverLocation, error)
	FindAvailableDrivers(ctx context.Context, lat, lng float64, radius float64, limit int) ([]DriverLocation, error)
	FindDriversByVehicleType(ctx context.Context, lat, lng float64, radius float64, vehicleType string, limit int) ([]DriverLocation, error)

	// Driver status management
	SetDriverStatus(ctx context.Context, driverID string, status string) error
	GetDriverStatus(ctx context.Context, driverID string) (string, error)
	GetAvailableDriversCount(ctx context.Context) (int64, error)

	// Batch operations
	AddMultipleDriverLocations(ctx context.Context, locations []GeoLocation) error
	GetMultipleDriverLocations(ctx context.Context, driverIDs []string) ([]DriverLocation, error)

	// Health and monitoring
	Ping(ctx context.Context) error
	GetStats(ctx context.Context) (map[string]interface{}, error)
}

// GeoLocationManager implements Redis geolocation operations
type GeoLocationManager struct {
	client    *redis.Client
	keyPrefix string
}

// NewGeoLocationManager creates a new Redis geolocation manager
func NewGeoLocationManager(client *redis.Client) IGeoLocationManager {
	return &GeoLocationManager{
		client:    client,
		keyPrefix: "motocabz:geo:",
	}
}

// Constants for Redis keys
const (
	DriverLocationKey = "drivers:location"
	DriverStatusKey   = "drivers:status"
	DriverMetadataKey = "drivers:metadata"
	DriverLastSeenKey = "drivers:lastseen"
)

// AddDriverLocation adds or updates a driver's location
func (gm *GeoLocationManager) AddDriverLocation(ctx context.Context, driverID string, lat, lng float64, metadata map[string]interface{}) error {
	key := gm.keyPrefix + DriverLocationKey

	// Add to geospatial index
	err := gm.client.GeoAdd(ctx, key, &redis.GeoLocation{
		Name:      driverID,
		Longitude: lng,
		Latitude:  lat,
	}).Err()
	if err != nil {
		return fmt.Errorf("failed to add driver location to geo index: %w", err)
	}

	// Store metadata
	if metadata != nil {
		metadataKey := gm.keyPrefix + DriverMetadataKey + ":" + driverID
		err = gm.client.HSet(ctx, metadataKey, metadata).Err()
		if err != nil {
			log.Printf("Warning: failed to store driver metadata for %s: %v", driverID, err)
		}
	}

	// Update last seen timestamp
	lastSeenKey := gm.keyPrefix + DriverLastSeenKey
	err = gm.client.HSet(ctx, lastSeenKey, driverID, time.Now().Unix()).Err()
	if err != nil {
		log.Printf("Warning: failed to update last seen for driver %s: %v", driverID, err)
	}

	return nil
}

// UpdateDriverLocation updates only the coordinates of a driver
func (gm *GeoLocationManager) UpdateDriverLocation(ctx context.Context, driverID string, lat, lng float64) error {
	return gm.AddDriverLocation(ctx, driverID, lat, lng, nil)
}

// RemoveDriverLocation removes a driver from the geospatial index
func (gm *GeoLocationManager) RemoveDriverLocation(ctx context.Context, driverID string) error {
	key := gm.keyPrefix + DriverLocationKey

	err := gm.client.ZRem(ctx, key, driverID).Err()
	if err != nil {
		return fmt.Errorf("failed to remove driver location: %w", err)
	}

	// Clean up metadata
	metadataKey := gm.keyPrefix + DriverMetadataKey + ":" + driverID
	gm.client.Del(ctx, metadataKey)

	// Clean up last seen
	lastSeenKey := gm.keyPrefix + DriverLastSeenKey
	gm.client.HDel(ctx, lastSeenKey, driverID)

	return nil
}

// GetDriverLocation retrieves a specific driver's location
func (gm *GeoLocationManager) GetDriverLocation(ctx context.Context, driverID string) (*DriverLocation, error) {
	key := gm.keyPrefix + DriverLocationKey

	// Get coordinates
	positions, err := gm.client.GeoPos(ctx, key, driverID).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get driver position: %w", err)
	}

	if len(positions) == 0 || positions[0] == nil {
		return nil, fmt.Errorf("driver %s not found", driverID)
	}

	pos := positions[0]

	// Get status
	status, err := gm.GetDriverStatus(ctx, driverID)
	if err != nil {
		status = "unknown"
	}

	// Get metadata
	metadataKey := gm.keyPrefix + DriverMetadataKey + ":" + driverID
	metadata, err := gm.client.HGetAll(ctx, metadataKey).Result()
	if err != nil {
		metadata = make(map[string]string)
	}

	// Get last seen
	lastSeenKey := gm.keyPrefix + DriverLastSeenKey
	lastSeenStr, err := gm.client.HGet(ctx, lastSeenKey, driverID).Result()
	var lastSeen time.Time
	if err == nil {
		if timestamp, err := strconv.ParseInt(lastSeenStr, 10, 64); err == nil {
			lastSeen = time.Unix(timestamp, 0)
		}
	}

	// Parse metadata
	vehicleType := metadata["vehicleType"]
	rating := 0.0
	if ratingStr, exists := metadata["rating"]; exists {
		if r, err := strconv.ParseFloat(ratingStr, 64); err == nil {
			rating = r
		}
	}

	return &DriverLocation{
		DriverID:    driverID,
		Latitude:    pos.Latitude,
		Longitude:   pos.Longitude,
		Status:      status,
		LastSeen:    lastSeen,
		VehicleType: vehicleType,
		Rating:      rating,
	}, nil
}

// FindNearbyDrivers finds drivers within a specified radius
func (gm *GeoLocationManager) FindNearbyDrivers(ctx context.Context, lat, lng float64, radius float64, limit int) ([]DriverLocation, error) {
	key := gm.keyPrefix + DriverLocationKey

	// Search for nearby drivers
	results, err := gm.client.GeoRadius(ctx, key, lng, lat, &redis.GeoRadiusQuery{
		Radius:    radius,
		Unit:      "km",
		WithDist:  true,
		WithCoord: true,
		Count:     limit,
		Sort:      "ASC", // Sort by distance
	}).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to search nearby drivers: %w", err)
	}

	var drivers []DriverLocation
	for _, result := range results {
		driverID := result.Name

		// Get additional driver info
		status, _ := gm.GetDriverStatus(ctx, driverID)
		metadataKey := gm.keyPrefix + DriverMetadataKey + ":" + driverID
		metadata, _ := gm.client.HGetAll(ctx, metadataKey).Result()

		// Parse metadata
		vehicleType := metadata["vehicleType"]
		rating := 0.0
		if ratingStr, exists := metadata["rating"]; exists {
			if r, err := strconv.ParseFloat(ratingStr, 64); err == nil {
				rating = r
			}
		}

		drivers = append(drivers, DriverLocation{
			DriverID:    driverID,
			Latitude:    result.Latitude,
			Longitude:   result.Longitude,
			Status:      status,
			Distance:    result.Dist,
			VehicleType: vehicleType,
			Rating:      rating,
		})
	}

	return drivers, nil
}

// FindAvailableDrivers finds only available drivers within radius
func (gm *GeoLocationManager) FindAvailableDrivers(ctx context.Context, lat, lng float64, radius float64, limit int) ([]DriverLocation, error) {
	// First get all nearby drivers
	allDrivers, err := gm.FindNearbyDrivers(ctx, lat, lng, radius, limit*2) // Get more to filter
	if err != nil {
		return nil, err
	}

	var availableDrivers []DriverLocation
	for _, driver := range allDrivers {
		if driver.Status == "available" {
			availableDrivers = append(availableDrivers, driver)
			if len(availableDrivers) >= limit {
				break
			}
		}
	}

	return availableDrivers, nil
}

// FindDriversByVehicleType finds drivers of a specific vehicle type within radius
func (gm *GeoLocationManager) FindDriversByVehicleType(ctx context.Context, lat, lng float64, radius float64, vehicleType string, limit int) ([]DriverLocation, error) {
	// Get all nearby drivers
	allDrivers, err := gm.FindNearbyDrivers(ctx, lat, lng, radius, limit*2)
	if err != nil {
		return nil, err
	}

	var filteredDrivers []DriverLocation
	for _, driver := range allDrivers {
		if driver.VehicleType == vehicleType && driver.Status == "available" {
			filteredDrivers = append(filteredDrivers, driver)
			if len(filteredDrivers) >= limit {
				break
			}
		}
	}

	return filteredDrivers, nil
}

// SetDriverStatus sets a driver's availability status
func (gm *GeoLocationManager) SetDriverStatus(ctx context.Context, driverID string, status string) error {
	key := gm.keyPrefix + DriverStatusKey

	err := gm.client.HSet(ctx, key, driverID, status).Err()
	if err != nil {
		return fmt.Errorf("failed to set driver status: %w", err)
	}

	// Update metadata with status
	metadataKey := gm.keyPrefix + DriverMetadataKey + ":" + driverID
	gm.client.HSet(ctx, metadataKey, "status", status)

	return nil
}

// GetDriverStatus retrieves a driver's current status
func (gm *GeoLocationManager) GetDriverStatus(ctx context.Context, driverID string) (string, error) {
	key := gm.keyPrefix + DriverStatusKey

	status, err := gm.client.HGet(ctx, key, driverID).Result()
	if err != nil {
		if err == redis.Nil {
			return "offline", nil
		}
		return "", fmt.Errorf("failed to get driver status: %w", err)
	}

	return status, nil
}

// GetAvailableDriversCount returns the count of available drivers
func (gm *GeoLocationManager) GetAvailableDriversCount(ctx context.Context) (int64, error) {
	key := gm.keyPrefix + DriverStatusKey

	// Get all driver statuses
	statuses, err := gm.client.HGetAll(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get driver statuses: %w", err)
	}

	count := int64(0)
	for _, status := range statuses {
		if status == "available" {
			count++
		}
	}

	return count, nil
}

// AddMultipleDriverLocations adds multiple driver locations in batch
func (gm *GeoLocationManager) AddMultipleDriverLocations(ctx context.Context, locations []GeoLocation) error {
	key := gm.keyPrefix + DriverLocationKey

	// Prepare geo locations for batch add
	var geoLocations []*redis.GeoLocation
	for _, loc := range locations {
		geoLocations = append(geoLocations, &redis.GeoLocation{
			Name:      loc.Member,
			Longitude: loc.Longitude,
			Latitude:  loc.Latitude,
		})
	}

	// Batch add to geospatial index
	err := gm.client.GeoAdd(ctx, key, geoLocations...).Err()
	if err != nil {
		return fmt.Errorf("failed to add multiple driver locations: %w", err)
	}

	// Store metadata for each driver
	for _, loc := range locations {
		if loc.Metadata != nil {
			metadataKey := gm.keyPrefix + DriverMetadataKey + ":" + loc.Member
			gm.client.HSet(ctx, metadataKey, loc.Metadata)
		}

		// Update last seen
		lastSeenKey := gm.keyPrefix + DriverLastSeenKey
		gm.client.HSet(ctx, lastSeenKey, loc.Member, time.Now().Unix())
	}

	return nil
}

// GetMultipleDriverLocations retrieves locations for multiple drivers
func (gm *GeoLocationManager) GetMultipleDriverLocations(ctx context.Context, driverIDs []string) ([]DriverLocation, error) {
	var locations []DriverLocation

	for _, driverID := range driverIDs {
		location, err := gm.GetDriverLocation(ctx, driverID)
		if err != nil {
			log.Printf("Warning: failed to get location for driver %s: %v", driverID, err)
			continue
		}
		locations = append(locations, *location)
	}

	return locations, nil
}

// Ping tests the Redis connection
func (gm *GeoLocationManager) Ping(ctx context.Context) error {
	return gm.client.Ping(ctx).Err()
}

// GetStats returns Redis geolocation statistics
func (gm *GeoLocationManager) GetStats(ctx context.Context) (map[string]interface{}, error) {
	key := gm.keyPrefix + DriverLocationKey

	// Get total drivers
	totalDrivers, err := gm.client.ZCard(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get total drivers count: %w", err)
	}

	// Get available drivers count
	availableCount, err := gm.GetAvailableDriversCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get available drivers count: %w", err)
	}

	// Get Redis memory usage
	memoryUsage, err := gm.client.MemoryUsage(ctx, key).Result()
	if err != nil {
		memoryUsage = 0
	}

	return map[string]interface{}{
		"totalDrivers":     totalDrivers,
		"availableDrivers": availableCount,
		"memoryUsage":      memoryUsage,
		"timestamp":        time.Now().Unix(),
	}, nil
}

// Helper functions for driver matching

// CalculateDriverScore calculates a score for driver matching
func CalculateDriverScore(driver DriverLocation, riderLat, riderLng float64, preferences map[string]interface{}) float64 {
	score := 100.0 // Base score

	// Distance factor (closer is better)
	if driver.Distance > 0 {
		// Reduce score based on distance (max 50 points reduction)
		distancePenalty := driver.Distance * 2
		if distancePenalty > 50 {
			distancePenalty = 50
		}
		score -= distancePenalty
	}

	// Rating factor (higher rating is better)
	if driver.Rating > 0 {
		score += driver.Rating * 10 // Max 50 points for 5-star rating
	}

	// Vehicle type preference
	if preferredType, exists := preferences["vehicleType"]; exists {
		if driver.VehicleType == preferredType {
			score += 20
		}
	}

	// Status factor
	if driver.Status == "available" {
		score += 10
	}

	return score
}

// SortDriversByScore sorts drivers by their matching score
func SortDriversByScore(drivers []DriverLocation, riderLat, riderLng float64, preferences map[string]interface{}) []DriverLocation {
	// Calculate scores for all drivers
	type DriverWithScore struct {
		Driver DriverLocation
		Score  float64
	}

	var driversWithScores []DriverWithScore
	for _, driver := range drivers {
		score := CalculateDriverScore(driver, riderLat, riderLng, preferences)
		driversWithScores = append(driversWithScores, DriverWithScore{
			Driver: driver,
			Score:  score,
		})
	}

	// Sort by score (descending)
	for i := 0; i < len(driversWithScores)-1; i++ {
		for j := i + 1; j < len(driversWithScores); j++ {
			if driversWithScores[i].Score < driversWithScores[j].Score {
				driversWithScores[i], driversWithScores[j] = driversWithScores[j], driversWithScores[i]
			}
		}
	}

	// Extract sorted drivers
	var sortedDrivers []DriverLocation
	for _, dws := range driversWithScores {
		sortedDrivers = append(sortedDrivers, dws.Driver)
	}

	return sortedDrivers
}
