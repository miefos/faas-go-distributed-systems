package config

import (
	"os"
)

type Config struct {
	ServerAddress string
	NATSUrl       string
	BucketName    string
}

func LoadConfig() *Config {
	serverAddress := getEnv("SERVER_ADDRESS", ":8081")
	natsUrl := getEnv("NATS_URL", "nats://localhost:4222")
	bucketName := getEnv("BUCKET_NAME", "functions")

	return &Config{
		ServerAddress: serverAddress,
		NATSUrl:       natsUrl,
		BucketName:    bucketName,
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
