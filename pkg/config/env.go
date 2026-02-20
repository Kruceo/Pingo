package config

import (
	"os"
	"time"
)

type ClickHouseConfig struct {
	DSN      string
	Database string
	Username string
	Password string
}

func LoadClickHouseConfig() *ClickHouseConfig {
	return &ClickHouseConfig{
		DSN:      getEnv("CLICKHOUSE_DSN", "localhost:9000"),
		Database: getEnv("CLICKHOUSE_DATABASE", "default"),
		Username: getEnv("CLICKHOUSE_USERNAME", "default"),
		Password: getEnv("CLICKHOUSE_PASSWORD", ""),
	}
}

func GetPingInterval(defaultInterval int) time.Duration {
	intervalStr := getEnv("PING_INTERVAL", "")
	if intervalStr != "" {
		duration, err := time.ParseDuration(intervalStr)
		if err != nil {
			return time.Duration(defaultInterval) * time.Second
		}
		return duration
	}
	return time.Duration(defaultInterval) * time.Second
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
