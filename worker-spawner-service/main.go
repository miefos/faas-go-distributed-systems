package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"sync"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/nats-io/nats.go"
)

// Constants
const (
	maxContainers = 5
	natsURL       = "nats://localhost:4222"
	functionTopic = "functions.execute"
	resultTopic   = "functions.results"
)

// Global variables
var (
	activeContainers = 0
	mu               sync.Mutex // Thread sage access to activeContainers var
)

type FunctionRequest struct {
	Function string `json:"image_reference"`
	Body     string `json:"parameter"`
}

func main() {
	//NATS
	nc := connectToNATS(natsURL)
	defer nc.Close()

	//Creation clien docker
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Fatalf("Error creating Docker client: %v", err)
	}
	defer cli.Close()

	log.Println("Client Docker created with success!")

	// Subscring worker queue
	_, err = nc.Subscribe(functionTopic, func(msg *nats.Msg) {
		var req FunctionRequest
		if err := json.Unmarshal(msg.Data, &req); err != nil {
			log.Printf("Error: message parsing: %v", err)
			nc.Publish(msg.Reply, []byte("Error: message not valid!"))
			return
		}

		// Check maxContainers
		mu.Lock()
		if activeContainers >= maxContainers {
			mu.Unlock()
			nc.Publish(msg.Reply, []byte("Error: containers are more than 5"))
			return
		}
		activeContainers++
		mu.Unlock()

		// Go routine, it will manage container life cycle
		go func() {
			defer func() {
				mu.Lock()
				activeContainers--
				mu.Unlock()
			}()

			result, err := executeFunction(cli, req.Function, req.Body)
			if err != nil {
				log.Printf("Errore executing function: %v", err)
				nc.Publish(msg.Reply, []byte(fmt.Sprintf("Error: %v", err)))
				return
			}

			// Publish result on the queue
			nc.Publish(resultTopic, []byte(result))
		}()
	})
	if err != nil {
		log.Fatalf("Error subscribing topic: %v", err)
	}

	select {}
}

func connectToNATS(natsURL string) *nats.Conn {
	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Fatalf("Error NATS connection: %v", err)
	}
	log.Println("Connection to NATS with success:", natsURL)
	return nc
}

func executeFunction(cli *client.Client, functionName, functionBody string) (string, error) {
	ctx := context.Background()

	/*
			// Path relativo per il file function.py
			functionFilePath := "./scripts/function.py"
			log.Printf("Writing function.py in %s", functionFilePath)

			// Create dynamic file `function.py`
			err := os.WriteFile(functionFilePath, []byte(fmt.Sprintf(`def %s():
		    %s
		`, functionName, functionBody)), 0644)
			if err != nil {
				return "", fmt.Errorf("Error creating function.py: %w", err)
			}
			log.Println("File function.py created successfully!")
	*/

	// Create Docker container
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: functionName, // Contains execute.py
		Cmd: []string{
			"python3",
			"/app/scripts/execute.py", // Container executes the script
			functionName,
			functionBody,
		},
		Tty: false,
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

	// Inspect initial container state
	containerJSON, err := cli.ContainerInspect(ctx, resp.ID)
	if err != nil {
		return "", fmt.Errorf("Error inspecting container: %w", err)
	}
	log.Printf("Initial container state: %s", containerJSON.State.Status)

	// Start container
	log.Printf("Starting container %s...", resp.ID)
	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return "", fmt.Errorf("Error starting container: %w", err)
	}

	// Re-inspect container state after starting
	containerJSON, err = cli.ContainerInspect(ctx, resp.ID)
	if err != nil {
		return "", fmt.Errorf("Error re-inspecting container: %w", err)
	}
	log.Printf("Container state after starting: %s", containerJSON.State.Status)

	// Ensure the container is running before copying the file
	if containerJSON.State.Status != "running" {
		return "", fmt.Errorf("Container is not running: current state '%s'", containerJSON.State.Status)
	}

	/*
		// Open function.py for reading
		functionFile, err := os.Open(functionFilePath)
		if err != nil {
			return "", fmt.Errorf("Error opening function file: %w", err)
		}
		defer functionFile.Close()

		// Copy function.py into the container
		log.Println("Copying function.py into container...")
		err = cli.CopyToContainer(ctx, resp.ID, "/app/scripts", functionFile, container.CopyToContainerOptions{})
		if err != nil {
			return "", fmt.Errorf("Error copying file to container: %w", err)
		}
		log.Println("File function.py copied successfully!")
	*/

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
	fmt.Println("output container:")
	fmt.Println(string(output))
	return string(output), nil
}
