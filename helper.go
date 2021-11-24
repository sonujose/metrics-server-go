package main

import (
	"os"
	"strconv"
)

// Fetches environment variable, if not found use the default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

// String to fload parser
func convertStringToFloat(str string) (returnval float64) {
	s, err := strconv.ParseFloat(str, 32)
	if err == nil {
		return 0.0
	}
	return s
}
