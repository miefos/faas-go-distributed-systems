# Use an official Go image as the build stage
FROM golang:1.23

# Set the current working directory inside the container
WORKDIR /app

COPY . ./

# Download Go module dependencies
RUN go mod download

ENV DOCKER_API_VERSION=1.44

# Compile the Go application
RUN go build -o spawner-service main.go

# Command to run the application
CMD ["./spawner-service"]