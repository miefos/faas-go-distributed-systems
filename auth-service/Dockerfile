# Use the official Golang image as a base
FROM golang:1.23

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker's caching
COPY go.mod go.sum ./

# Download Go module dependencies
RUN go mod download

# Ensure the directory exists and generate a JWT secret key, then write it to secret.txt
RUN mkdir -p ./tmp && echo $(openssl rand -base64 32) > /app/tmp/secret.txt
RUN echo "Secret key has been created in container: /app/tmp/secret.txt"

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o main .

# Expose the port your app runs on
EXPOSE 8080

# Command to run the application
CMD ["./main"]
