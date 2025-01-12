# Use an official Go image as the build stage
FROM golang:1.23

# Set the current working directory inside the container
WORKDIR /app

# Copy the Go modules file and download the dependencies
COPY . ./

RUN go mod init registry-service

# Download the Go dependencies
RUN go get github.com/gorilla/mux && go get github.com/nats-io/nats.go

# Compile the Go application when the container is built
RUN go build -o registry-service main.go

# Expose the port your app runs on
EXPOSE 8081

# Set the entry point for the container
CMD ["./registry-service"]
