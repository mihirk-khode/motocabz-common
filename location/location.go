package location

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
)

// Location represents a geographical location
type Location struct {
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
	Address    string  `json:"address,omitempty"`
	City       string  `json:"city,omitempty"`
	State      string  `json:"state,omitempty"`
	Country    string  `json:"country,omitempty"`
	PostalCode string  `json:"postalCode,omitempty"`
}

// LocationBounds represents geographical bounds
type LocationBounds struct {
	NorthEast Location `json:"northEast"`
	SouthWest Location `json:"southWest"`
}

// DistanceUnit represents the unit for distance calculation
type DistanceUnit string

const (
	DistanceUnitKilometers DistanceUnit = "km"
	DistanceUnitMiles      DistanceUnit = "miles"
	DistanceUnitMeters     DistanceUnit = "meters"
)

// ExtractLocationData extracts latitude and longitude from JSON location data
func ExtractLocationData(locationJSON json.RawMessage) (lat, lng float64, err error) {
	if len(locationJSON) == 0 {
		return 0.0, 0.0, errors.New("empty location JSON")
	}

	var location map[string]interface{}
	if err := json.Unmarshal(locationJSON, &location); err != nil {
		return 0.0, 0.0, err
	}

	latVal, latOk := location["latitude"].(float64)
	lngVal, lngOk := location["longitude"].(float64)

	if !latOk || !lngOk {
		return 0.0, 0.0, errors.New("invalid latitude or longitude format")
	}

	return latVal, lngVal, nil
}

// CreateLocationJSON creates JSON from latitude and longitude
func CreateLocationJSON(lat, lng float64) (json.RawMessage, error) {
	location := map[string]interface{}{
		"latitude":  lat,
		"longitude": lng,
	}

	return json.Marshal(location)
}

// CreateLocationJSONWithAddress creates JSON from location data including address
func CreateLocationJSONWithAddress(lat, lng float64, address string) (json.RawMessage, error) {
	location := map[string]interface{}{
		"latitude":  lat,
		"longitude": lng,
		"address":   address,
	}

	return json.Marshal(location)
}

// IsValidLocation checks if latitude and longitude are within valid ranges
func IsValidLocation(lat, lng float64) bool {
	return lat >= -90 && lat <= 90 && lng >= -180 && lng <= 180
}

// CalculateDistance calculates the distance between two locations using Haversine formula
func CalculateDistance(loc1, loc2 Location, unit DistanceUnit) float64 {
	const earthRadiusKm = 6371.0
	const earthRadiusMiles = 3959.0

	lat1Rad := loc1.Latitude * math.Pi / 180
	lat2Rad := loc2.Latitude * math.Pi / 180
	deltaLatRad := (loc2.Latitude - loc1.Latitude) * math.Pi / 180
	deltaLngRad := (loc2.Longitude - loc1.Longitude) * math.Pi / 180

	a := math.Sin(deltaLatRad/2)*math.Sin(deltaLatRad/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(deltaLngRad/2)*math.Sin(deltaLngRad/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	distanceKm := earthRadiusKm * c

	switch unit {
	case DistanceUnitMiles:
		return distanceKm * 0.621371
	case DistanceUnitMeters:
		return distanceKm * 1000
	default:
		return distanceKm
	}
}

// CalculateBearing calculates the bearing from one location to another
func CalculateBearing(from, to Location) float64 {
	lat1Rad := from.Latitude * math.Pi / 180
	lat2Rad := to.Latitude * math.Pi / 180
	deltaLngRad := (to.Longitude - from.Longitude) * math.Pi / 180

	y := math.Sin(deltaLngRad) * math.Cos(lat2Rad)
	x := math.Cos(lat1Rad)*math.Sin(lat2Rad) - math.Sin(lat1Rad)*math.Cos(lat2Rad)*math.Cos(deltaLngRad)

	bearing := math.Atan2(y, x) * 180 / math.Pi

	// Normalize to 0-360 degrees
	if bearing < 0 {
		bearing += 360
	}

	return bearing
}

// IsLocationWithinBounds checks if a location is within given bounds
func IsLocationWithinBounds(location Location, bounds LocationBounds) bool {
	return location.Latitude >= bounds.SouthWest.Latitude &&
		location.Latitude <= bounds.NorthEast.Latitude &&
		location.Longitude >= bounds.SouthWest.Longitude &&
		location.Longitude <= bounds.NorthEast.Longitude
}

// CreateBoundsFromCenter creates bounds around a center location
func CreateBoundsFromCenter(center Location, radiusKm float64) LocationBounds {
	const earthRadiusKm = 6371.0

	// Convert radius from km to degrees (approximate)
	latDelta := radiusKm / earthRadiusKm * 180 / math.Pi
	lngDelta := radiusKm / (earthRadiusKm * math.Cos(center.Latitude*math.Pi/180)) * 180 / math.Pi

	return LocationBounds{
		NorthEast: Location{
			Latitude:  center.Latitude + latDelta,
			Longitude: center.Longitude + lngDelta,
		},
		SouthWest: Location{
			Latitude:  center.Latitude - latDelta,
			Longitude: center.Longitude - lngDelta,
		},
	}
}

// FindNearestLocation finds the nearest location from a list of locations
func FindNearestLocation(target Location, locations []Location) (Location, float64) {
	if len(locations) == 0 {
		return Location{}, 0
	}

	nearest := locations[0]
	minDistance := CalculateDistance(target, nearest, DistanceUnitKilometers)

	for _, location := range locations[1:] {
		distance := CalculateDistance(target, location, DistanceUnitKilometers)
		if distance < minDistance {
			minDistance = distance
			nearest = location
		}
	}

	return nearest, minDistance
}

// SortLocationsByDistance sorts locations by distance from a target location
func SortLocationsByDistance(target Location, locations []Location) []Location {
	if len(locations) <= 1 {
		return locations
	}

	// Create a slice of location-distance pairs
	type locationDistance struct {
		location Location
		distance float64
	}

	var pairs []locationDistance
	for _, location := range locations {
		distance := CalculateDistance(target, location, DistanceUnitKilometers)
		pairs = append(pairs, locationDistance{location, distance})
	}

	// Sort by distance (simple bubble sort for small lists)
	for i := 0; i < len(pairs)-1; i++ {
		for j := 0; j < len(pairs)-i-1; j++ {
			if pairs[j].distance > pairs[j+1].distance {
				pairs[j], pairs[j+1] = pairs[j+1], pairs[j]
			}
		}
	}

	// Extract sorted locations
	var sorted []Location
	for _, pair := range pairs {
		sorted = append(sorted, pair.location)
	}

	return sorted
}

// ValidateLocationJSON validates that JSON contains valid location data
func ValidateLocationJSON(locationJSON json.RawMessage) error {
	lat, lng, err := ExtractLocationData(locationJSON)
	if err != nil {
		return err
	}

	if !IsValidLocation(lat, lng) {
		return errors.New("invalid latitude or longitude values")
	}

	return nil
}

// ParseLocationFromMap parses location data from a map
func ParseLocationFromMap(data map[string]interface{}) (Location, error) {
	var location Location

	if lat, ok := data["latitude"].(float64); ok {
		location.Latitude = lat
	} else {
		return location, errors.New("latitude not found or invalid type")
	}

	if lng, ok := data["longitude"].(float64); ok {
		location.Longitude = lng
	} else {
		return location, errors.New("longitude not found or invalid type")
	}

	if address, ok := data["address"].(string); ok {
		location.Address = address
	}

	if city, ok := data["city"].(string); ok {
		location.City = city
	}

	if state, ok := data["state"].(string); ok {
		location.State = state
	}

	if country, ok := data["country"].(string); ok {
		location.Country = country
	}

	if postalCode, ok := data["postalCode"].(string); ok {
		location.PostalCode = postalCode
	}

	return location, nil
}

// LocationToMap converts a Location struct to a map
func LocationToMap(location Location) map[string]interface{} {
	result := map[string]interface{}{
		"latitude":  location.Latitude,
		"longitude": location.Longitude,
	}

	if location.Address != "" {
		result["address"] = location.Address
	}

	if location.City != "" {
		result["city"] = location.City
	}

	if location.State != "" {
		result["state"] = location.State
	}

	if location.Country != "" {
		result["country"] = location.Country
	}

	if location.PostalCode != "" {
		result["postalCode"] = location.PostalCode
	}

	return result
}

// GetLocationString returns a formatted string representation of the location
func GetLocationString(location Location) string {
	if location.Address != "" {
		return location.Address
	}

	return formatCoordinate(location.Latitude) + ", " + formatCoordinate(location.Longitude)
}

// formatCoordinate formats a coordinate to a readable string
func formatCoordinate(coord float64) string {
	direction := "N"
	if coord < 0 {
		direction = "S"
		coord = -coord
	}

	degrees := int(coord)
	minutes := int((coord - float64(degrees)) * 60)
	seconds := (coord - float64(degrees) - float64(minutes)/60) * 3600

	return formatCoordinateString(degrees, minutes, seconds, direction)
}

// formatCoordinateString formats coordinate components
func formatCoordinateString(degrees, minutes int, seconds float64, direction string) string {
	if minutes == 0 && seconds < 1 {
		return formatFloat(float64(degrees), 0) + "°" + direction
	}

	if seconds < 1 {
		return formatFloat(float64(degrees), 0) + "°" + formatFloat(float64(minutes), 0) + "'" + direction
	}

	return formatFloat(float64(degrees), 0) + "°" + formatFloat(float64(minutes), 0) + "'" + formatFloat(seconds, 1) + "\"" + direction
}

// formatFloat formats a float to a string with specified decimal places
func formatFloat(f float64, decimals int) string {
	if decimals == 0 {
		return fmt.Sprintf("%.0f", f)
	}

	multiplier := math.Pow(10, float64(decimals))
	return fmt.Sprintf("%."+fmt.Sprintf("%d", decimals)+"f", math.Round(f*multiplier)/multiplier)
}
