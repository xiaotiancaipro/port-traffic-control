package configs

import (
	"fmt"
	"os"
	"strconv"
)

func New() (config *Configuration, err error) {
	config = &Configuration{
		API: &APIConfig{
			Host:      getEnv("PORT_TC_API_HOST", "127.0.0.1"),
			Env:       getEnv("PORT_TC_API_ENV", "production"),
			LogPath:   getEnv("PORT_TC_API_LOG_PATH", "./logs"),
			LogPrefix: getEnv("PORT_TC_API_LOG_PREFIX", "api"),
		},
		Database: &DatabaseConfig{
			Path: getEnv("PORT_TC_DATABASE_PATH", "./data/data.db"),
		},
	}
	portStr := getEnv("PORT_TC_API_PORT", "6001")
	port, err := strconv.Atoi(portStr)
	if err != nil || port < 1 || port > 65535 {
		err = fmt.Errorf("invalid port %s", portStr)
		return
	}
	config.API.Port = port
	return
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
