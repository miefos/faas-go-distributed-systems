package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/docker/docker/client"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/container"
	"github.com/nats-io/nats.go"
)

type Config struct {
	ServerAddress string
	NATSUrl       string
	functionTopic string
	resultTopic   string
	maxContainers int
}

// Global variables
var (
	activeContainers = 0
	mu               sync.Mutex // Thread-safe access to activeContainers var
)

type FunctionRequest struct {
	ImageReference string `json:"image_reference"`
	Parameter      string `json:"parameter"`
}

func main() {
	cfg := LoadConfig()

	// Set the Docker API version to a compatible version
	os.Setenv("DOCKER_API_VERSION", "1.44")

	// Connect to NATS
	nc := connectToNATS(cfg.NATSUrl)
	defer nc.Close()

	// Create Docker client
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Fatalf("Error creating Docker client: %v", err)
	}
	defer cli.Close()

	log.Println("Docker client created successfully!")

	// Subscribe to the function execution topic
	_, err = nc.Subscribe(cfg.functionTopic, func(msg *nats.Msg) {
		var req FunctionRequest
		log.Printf("Received message: %s", msg.Data)
		if err := json.Unmarshal(msg.Data, &req); err != nil {
			log.Printf("Error parsing message: %v", err)
			nc.Publish(msg.Reply, []byte("Error: invalid message format"))
			return
		}

		// Check maxContainers
		mu.Lock()
		if activeContainers >= cfg.maxContainers {
			mu.Unlock()
			nc.Publish(msg.Reply, []byte("Error: too many active containers"))
			return
		}
		activeContainers++
		mu.Unlock()

		// Go routine to manage container lifecycle
		go func() {
			defer func() {
				mu.Lock()
				activeContainers--
				mu.Unlock()
			}()

			result, err := executeFunction(cli, req.ImageReference, req.Parameter)
			if err != nil {
				log.Printf("Error executing function: %v", err)
				nc.Publish(msg.Reply, []byte(fmt.Sprintf("Error: %v", err)))
				return
			}

			// Publish result to the result topic
			nc.Publish(cfg.resultTopic, []byte(result))
		}()
	})
	if err != nil {
		log.Fatalf("Error subscribing to topic: %v", err)
	}

	select {}
}

func connectToNATS(natsURL string) *nats.Conn {
	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Fatalf("Error connecting to NATS: %v", err)
	}
	log.Println("Connected to NATS successfully:", natsURL)
	return nc
}

func executeFunction(cli *client.Client, imageReference, parameter string) (string, error) {
	ctx := context.Background()

	// Pull the Docker image
	log.Printf("Pulling image %s...", imageReference)
	reader, err := cli.ImagePull(ctx, imageReference, image.PullOptions{})
	if err != nil {
		return "", fmt.Errorf("Error pulling image: %w", err)
	}
	defer reader.Close()

	// Read the output to ensure the image is pulled
	_, err = io.Copy(os.Stdout, reader)
	if err != nil {
		return "", fmt.Errorf("Error reading image pull response: %w", err)
	}

	// Create Docker container
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: imageReference,
		Cmd:   []string{parameter},
		Tty:   false,
	}, nil, nil, nil, "")
	if err != nil {
		return "", fmt.Errorf("Error creating container: %w", err)
	}

	// Ensure container cleanup
	defer func() {
		if err := cli.ContainerRemove(ctx, resp.ID, container.RemoveOptions{Force: true}); err != nil {
			log.Printf("Error removing container %s: %v", resp.ID, err)
		}
	}()

	// Start container
	log.Printf("Starting container %s...", resp.ID)
	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return "", fmt.Errorf("Error starting container: %w", err)
	}

	// Wait for container to finish execution
	log.Printf("Waiting for container %s to finish execution...", resp.ID)
	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return "", fmt.Errorf("Error waiting for container: %w", err)
		}
	case <-statusCh:
	}

	// Read logs for the result
	log.Println("Reading container logs...")
	out, err := cli.ContainerLogs(ctx, resp.ID, container.LogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		return "", fmt.Errorf("Error reading container logs: %w", err)
	}

	output, err := io.ReadAll(out)
	if err != nil {
		return "", fmt.Errorf("Error reading container output: %w", err)
	}
	return string(output), nil
}

func LoadConfig() *Config {
	serverAddress := getEnv("SERVER_ADDRESS", ":8082")
	natsUrl := getEnv("NATS_URL", "nats://localhost:4222")
	functionTopic := getEnv("FUNCTION_TOPIC", "functions.execute")
	resultTopic := getEnv("RESULT_TOPIC", "functions.results")
	maxContainers, err := strconv.Atoi(getEnv("MAX_CONTAINERS", "10"))
	if err != nil {
		log.Fatalf("Error parsing MAX_CONTAINERS: %v", err)
	}

	return &Config{
		ServerAddress: serverAddress,
		NATSUrl:       natsUrl,
		functionTopic: functionTopic,
		resultTopic:   resultTopic,
		maxContainers: maxContainers,
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}