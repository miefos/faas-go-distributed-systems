# Use the official Go 1.23 image as the base image
FROM golang:1.23

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files to the container
COPY go.mod go.sum ./

# Download the dependencies defined in go.mod
RUN go mod tidy

# Copy the entire source code into the container
COPY . .

# Build the Go application
RUN go build -o main .

# Define the command to run the application
CMD ["./main"]
