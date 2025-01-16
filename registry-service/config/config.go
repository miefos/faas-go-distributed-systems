package config

import (
	"os"
)

type Config struct {
	ServerAddress string
	NATS1Url       string
	NATS2Url       string
	BucketName    string
}

func LoadConfig() *Config {
	serverAddress := getEnv("SERVER_ADDRESS", ":8081")
	nats1Url := getEnv("NATS1_URL", "nats://localhost:4222")
	nats2Url := getEnv("NATS2_URL", "")
	bucketName := getEnv("BUCKET_NAME", "functions")

	return &Config{
		ServerAddress: serverAddress,
		NATS1Url:       nats1Url,
		NATS2Url:       nats2Url,
		BucketName:    bucketName,
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
