package configs

import (
	"os"
	"strconv"
)

func New() (config *Configuration, err error) {
	config = &Configuration{
		API: &APIConfig{
			Host:      getEnvString("PTC_API_HOST", "127.0.0.1"),
			Port:      getEnvUint32("PTC_API_PORT", 5001),
			Env:       getEnvString("PTC_API_ENV", "production"),
			LogPath:   getEnvString("PTC_API_LOG_PATH", "./logs"),
			LogPrefix: getEnvString("PTC_API_LOG_PREFIX", "api"),
			AccessKey: getEnvString("PTC_API_ACCESS_KEY", "port-traffic-control"),
			PIDPath:   getEnvString("PTC_API_PID_PATH", "./pid"),
			PIDFile:   getEnvString("PTC_API_PID_FILE", "api.pid"),
		},
		Database: &DatabaseConfig{
			Path: getEnvString("PTC_DATABASE_PATH", "./data"),
			File: getEnvString("PTC_DATABASE_FILE", "./data.db"),
		},
		TC: &TCConfig{
			InterfaceName:  getEnvString("PTC_TC_INTERFACE_NAME", "eth0"),
			HTBVersion:     getEnvUint32("PTC_TC_HTB_VERSION", 3),
			Rate2Quantum:   getEnvUint32("PTC_TC_RATE2QUANTUM", 10),
			DefaultClassID: getEnvUint32("PTC_TC_DEFAULT_CLASS_ID", 1),
			EnableLogging:  getEnvBool("PTC_TC_ENABLE_LOGGING", false),
			LogLevel:       getEnvString("PTC_TC_LOG_LEVEL", "info"),
		},
	}
	return
}

func getEnvString(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvUint32(key string, defaultValue uint32) uint32 {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseUint(value, 10, 32); err == nil {
			return uint32(parsed)
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}
