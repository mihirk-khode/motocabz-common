package util

import (
	"encoding/json"
	"log"
)

// RawMessageToMap converts json.RawMessage to map[string]interface{}
func RawMessageToMap(raw json.RawMessage) map[string]interface{} {
	if len(raw) == 0 {
		return nil
	}
	var m map[string]interface{}
	if err := json.Unmarshal(raw, &m); err != nil {
		log.Printf("Failed to unmarshal json.RawMessage: %v", err)
		return nil
	}
	return m
}

// ExtractCoordinates extracts latitude and longitude from a location object
// Handles both json.RawMessage and map[string]interface{}
func ExtractCoordinates(location interface{}) (lat, lng float64, ok bool) {
	var locationMap map[string]interface{}

	// Handle json.RawMessage
	if rawMsg, ok := location.(json.RawMessage); ok {
		if len(rawMsg) == 0 {
			return 0, 0, false
		}
		if err := json.Unmarshal(rawMsg, &locationMap); err != nil {
			return 0, 0, false
		}
	} else if m, ok := location.(map[string]interface{}); ok {
		locationMap = m
	} else {
		return 0, 0, false
	}

	// Extract latitude
	if latVal, exists := locationMap["latitude"]; exists {
		switch v := latVal.(type) {
		case float64:
			lat = v
		case float32:
			lat = float64(v)
		case int:
			lat = float64(v)
		case int64:
			lat = float64(v)
		default:
			return 0, 0, false
		}
	} else {
		return 0, 0, false
	}

	// Extract longitude
	if lngVal, exists := locationMap["longitude"]; exists {
		switch v := lngVal.(type) {
		case float64:
			lng = v
		case float32:
			lng = float64(v)
		case int:
			lng = float64(v)
		case int64:
			lng = float64(v)
		default:
			return 0, 0, false
		}
	} else {
		return 0, 0, false
	}

	return lat, lng, true
}

// HandleRawMessage unmarshals and re-marshals json.RawMessage
func HandleRawMessage(raw json.RawMessage) json.RawMessage {
	if len(raw) == 0 {
		return nil
	}

	var data map[string]interface{}
	if err := json.Unmarshal(raw, &data); err != nil {
		if err.Error() != "unexpected end of JSON input" {
			log.Printf("Failed to unmarshal %v", err)
		}
		return nil
	}
	bytes, err := json.Marshal(data)
	if err != nil {
		log.Printf("Failed to marshal %v", err)
		return nil
	}
	return bytes
}

// GetStringFromMap extracts a string value from a map
func GetStringFromMap(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

// GetFloat64FromMap extracts a float64 value from a map
func GetFloat64FromMap(m map[string]interface{}, key string) float64 {
	if val, ok := m[key]; ok {
		switch v := val.(type) {
		case float64:
			return v
		case float32:
			return float64(v)
		case int:
			return float64(v)
		case int64:
			return float64(v)
		}
	}
	return 0.0
}
