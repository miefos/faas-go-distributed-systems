package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	NATS1_URL      string
	NATS2_URL      string
	SERVER_ADDRESS string
	MessageQueue   string
	Timeout        int //seconds
}

func LoadConfig() *Config {
	nats1_url := getEnv("NATS1_URL", "nats://localhost:4222")
	nats2_url := getEnv("NATS2_URL", "nats://localhost:4223")
	server_address := getEnv("SERVER_ADDRESS", ":8083")
	messageQueue := getEnv("MESSAGE_QUEUE", "functions.execute")
	timeout, err := strconv.Atoi(getEnv("TIMEOUT", "30"))
	if err != nil {
		log.Fatalf("Error parsing TIMEOUT: %v", err)
	}
	return &Config{
		NATS1_URL:      nats1_url,
		NATS2_URL:      nats2_url,
		SERVER_ADDRESS: server_address,
		MessageQueue:   messageQueue,
		Timeout:        timeout,
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
