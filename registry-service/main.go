package main

import (
	"log"
	"net/http"
	"os"

	"registry-service/handlers"
	"registry-service/config"
	"registry-service/storage"

	"github.com/gorilla/mux"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	if(storage.InitNATSConnection(cfg.NATS1Url) != 0){
		if(storage.InitNATSConnection(cfg.NATS2Url) != 0){
			log.Fatalf("Error connecting to all NATS servers")
			os.Exit(-1)
		}
	}

	// Initialize NATS KeyValueStore
	kvStore, err := storage.NewKVStore(cfg.BucketName)
	if err != nil {
		log.Fatalf("Failed to initialize KeyValueStore: %v", err)
		os.Exit(-1)
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
