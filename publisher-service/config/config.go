package config

import (
	"os"
)

type Config struct {
	NATS_URL       string
	SERVER_ADDRESS string
}

func LoadConfig() *Config {
	nats_url := getEnv("NATS_URL", "nats://localhost:4222")
	server_address := getEnv("SERVER_ADDRESS", ":8083")

	return &Config{
		NATS_URL:       nats_url,
		SERVER_ADDRESS: server_address,
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
