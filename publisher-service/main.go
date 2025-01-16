package main

import (
	"log"
	"net/http"
	"publisher-service/config"
	"publisher-service/handlers"

	"github.com/gorilla/mux"
	"github.com/nats-io/nats.go"
)

func main() {
	log.Println("Starting publisher-service...")

	// Carica la configurazione
	cfg := config.LoadConfig()

	// Connetti a NATS
	log.Printf("Connecting to NATS at %s...", cfg.NATS_URL)
	nc, err := nats.Connect(cfg.NATS_URL)
	if err != nil {
		log.Fatalf("Error connecting to NATS: %v", err)
	}
	defer nc.Close()
	log.Println("Connected to NATS successfully.")

	// Inizializza il router
	r := mux.NewRouter()

	// Configura l'handler
	publisherHandler := handlers.NewPublisherHandler(nc, "functions.execution", 30)
	r.HandleFunc("/publish", publisherHandler.PublishHandlerMethod).Methods("POST")

	// Avvia il server HTTP
	log.Printf("Publisher-service is running on port %s", cfg.SERVER_ADDRESS)
	if err := http.ListenAndServe(cfg.SERVER_ADDRESS, r); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
