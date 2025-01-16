package config

import (
	"os"
)

type Config struct {
	NATS1_URL       string
	NATS2_URL       string
	SERVER_ADDRESS string
}

func LoadConfig() *Config {
	nats1_url := getEnv("NATS1_URL", "nats://localhost:4222")
	nats2_url := getEnv("NATS2_URL", "nats://localhost:4223")
	server_address := getEnv("SERVER_ADDRESS", ":8083")

	return &Config{
		NATS1_URL:       nats1_url,
		NATS2_URL:       nats2_url,
		SERVER_ADDRESS: server_address,
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
