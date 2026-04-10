package utils

import (
	"encoding/json"
	"fmt"
)

// GetString extracts a string value from a map.
func GetString(m map[string]interface{}, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}

// ParseInt extracts an integer value from a map, handling float64 and json.Number as common in JSON decoding.
func ParseInt(m map[string]interface{}, key string) int {
	if v, ok := m[key].(float64); ok {
		return int(v)
	}
	if v, ok := m[key].(json.Number); ok {
		i, _ := v.Int64()
		return int(i)
	}
	if v, ok := m[key].(int); ok {
		return v
	}
	if v, ok := m[key].(int64); ok {
		return int(v)
	}
	return 0
}

// ParseUint parses a string to a uint.
func ParseUint(s string) uint {
	var u uint
	fmt.Sscanf(s, "%d", &u)
	return u
}
