package config

import (
	"os"
	"strconv"
	"strings"
)

// Get an environment variable, specifying a default value if its not set
func GetEnvDefault(key string, defaultValue string) string {
	val, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return val
}

// Get an environment variable as an int64, specifying a default value if its
// not set or can't be parsed properly into an int64
func GetEnvInt64(key string, defaultValue int64) int64 {
	val, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	intVal, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return defaultValue
	}
	return intVal

}

// Get an environment variable as a boolean, specifying a default value if its
// not set or can't be parsed properly into a bool
func GetEnvBool(key string, defaultValue bool) bool {
	val, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	truthy := map[string]bool{
		"true": true, "t": true, "yes": true, "y": true, "on": true, "1": true,
		"enable": true, "enabled": true, "active": true, "affirmative": true,
	}

	falsy := map[string]bool{
		"false": false, "f": false, "no": false, "n": false, "off": false, "0": false,
		"disable": false, "disabled": false, "inactive": false, "negative": false,
	}

	normalized := strings.TrimSpace(strings.ToLower(val))

	if val, ok := truthy[normalized]; ok {
		return val
	}
	if val, ok := falsy[normalized]; ok {
		return val
	}

	return defaultValue
}
