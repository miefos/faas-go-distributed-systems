# Worker Spawner

The **Worker Spawner** is a Go-based microservice designed to manage the execution of containerized functions dynamically. It listens to a message queue (NATS) for function execution requests, spawns containers to execute the functions, and sends back the results. This service does not expose REST APIs; its interactions are entirely based on NATS messaging.

## Disclaimer: library used 
Pay attention on the libraries used:
- "github.com/docker/docker/client"
- "github.com/docker/docker/api/types/image"
- "github.com/docker/docker/api/types/container"
Those libraries change often (definition of the functions and types). If you want to modify this code, check the official documentation :) at link:
[Documentation](https://pkg.go.dev/github.com/docker/docker/client)

---

## Features
- **Dynamic Container Management**: Spawns containers based on received requests and ensures proper cleanup after execution.
- **Concurrency Control**: Limits the number of active containers to a predefined maximum (default: 5)
- **Queue Integration**: Communicates with a NATS message queue for receiving execution requests and publishing results.
- **Flexible Function Execution**: Executes containerized functions based on provided Docker images and parameters.

---

## Environment
Ensure the following environment variables or configurations:
- `NATS_URL`: URL of the NATS server (default: `nats://localhost:4222`)
- `MAX_CONTAINERS`: Maximum number of active containers allowed (default: `5`)

---

## Usage

### NATS Topics

#### `functions.execute`
- **Purpose**: Receives execution requests.
- **Message Format**: JSON object containing the following fields:

```json
{
  "image_reference": "<docker_image>",
  "parameter": "<command_parameter>"
}
```

- **Example**:
```json
{
  "image_reference": "python:3.11-slim",
  "parameter": "python3 -c \"print('Hello, World!')\""
}
```

#### `functions.results`
- **Purpose**: Publishes the results of executed functions.
- **Message Format**: Raw string output from the container.

---

## How It Works
1. **Message Subscription**:
   - The worker subscribes to the `functions.execute` topic on NATS.
   - On receiving a valid request, it checks if the active container limit is reached. Then checks if the image exits in the system, otherwise it will pull from the Docker registry.
2. **Container Execution**:
   - Creates a Docker container using the specified `image_reference`.
   - Passes the `parameter` as the command to be executed inside the container.
   - Waits for the container to complete execution and captures its output.
3. **Result Publication**:
   - Publishes the container's output to the `functions.results` topic on NATS, that is user's queue.
4. **Cleanup**:
   - Ensures proper cleanup by removing the container after execution.

---

## Installation and Setup

### 1. Clone the Repository
```bash
git clone <repository-url>
cd <repository>
```

### 2. Build the Application
```bash
go build -o worker-spawner .
```

### 3. Run the Application
```bash
./worker-spawner
```

---

## Dockerization

### Build Docker Image

Create a `Dockerfile` as follows:

```Dockerfile
FROM golang:1.20-alpine

WORKDIR /app
COPY . .
RUN go build -o worker-spawner .
CMD ["./worker-spawner"]
```

Then build and run the Docker image:

```bash
docker build -t worker-spawner .
docker run --rm -e NATS_URL=nats://<nats-server>:4222 worker-spawner
```

---

## Testing
Use the `nats` CLI to send test requests:

### Send Execution Request
```bash
nats pub functions.execute '{"image_reference": "python:3.11-slim", "parameter": "python3 -c \"print('Hello, World!')\""}'
```

### Listen for Results
```bash
nats sub functions.results
```

---


