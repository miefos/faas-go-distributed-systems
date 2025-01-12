# Use the latest Go image with Alpine as the base image
FROM golang:1.20-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files to download dependencies
COPY go.mod go.sum ./

# Download the dependencies
RUN go mod download

# Copy the entire source code into the working directory
COPY . .

# Build the application and create the binary
RUN go build -o worker-spawner .

# Expose the required port (if applicable)
EXPOSE 8082

# Run the application
CMD ["./worker-spawner"]
