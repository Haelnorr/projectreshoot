package server

import (
	"os"
	"strconv"
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
