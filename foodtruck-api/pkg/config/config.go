package config

import (
	debug "foodtruck/pkg/logger"
	"os"
)

var logger = debug.Logger

func GetEnv(key string, defaultVal string) string {
	var val = os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val
}
