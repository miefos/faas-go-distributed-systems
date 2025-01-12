package main

import (
	"log"
	"net/http"

	"registry-service/handlers"
	"registry-service/config"
	"registry-service/storage"

	"github.com/gorilla/mux"
)

// TODO: Change to use the random UUID as a body parameter

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize NATS KeyValueStore
	kvStore, err := storage.NewKVStore(cfg.NATSUrl, cfg.BucketName)
	if err != nil {
		log.Fatalf("Failed to initialize KeyValueStore: %v", err)
	} else {
		log.Println("KeyValueStore initialized")
	}

	// Set up HTTP router
	router := mux.NewRouter()
	handlers.RegisterRoutes(router, kvStore)

	// Start HTTP server
	log.Printf("Starting Registry Service on %s...", cfg.ServerAddress)
	log.Fatal(http.ListenAndServe(cfg.ServerAddress, router))
}
