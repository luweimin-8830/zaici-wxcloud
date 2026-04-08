package utils

import "fmt"

// GetString extracts a string value from a map.
func GetString(m map[string]interface{}, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}

// ParseInt extracts an integer value from a map, handling float64 as common in JSON decoding.
func ParseInt(m map[string]interface{}, key string) int {
	if v, ok := m[key].(float64); ok {
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
