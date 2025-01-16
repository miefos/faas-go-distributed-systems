package main

import (
	"log"
	"net/http"
	"publisher-service/config"
	"publisher-service/handlers"
	"os"

	"github.com/gorilla/mux"
	"github.com/nats-io/nats.go"
)

func main() {
	log.Println("Starting publisher-service...")

	// Carica la configurazione
	cfg := config.LoadConfig()

	// Connetti a NATS
	log.Printf("Connecting to NATS at %s...", cfg.NATS1_URL)
	nc, err := nats.Connect(cfg.NATS1_URL)
	if err != nil {
		log.Fatalf("Error connecting to NATS: %v", err)
		nc, err = nats.Connect(cfg.NATS2_URL)
		if err != nil {
			log.Fatalf("Error connecting to NATS: %v", err)
			os.Exit(-1)
		}
	}
	defer nc.Close()
	log.Println("Connected to NATS successfully.")

	// Inizializza il router
	r := mux.NewRouter()

	// Configura l'handler
	publisherHandler := handlers.NewPublisherHandler(nc, "functions.execute", 30)
	r.HandleFunc("/publish", publisherHandler.PublishHandlerMethod).Methods("POST")

	// Avvia il server HTTP
	log.Printf("Publisher-service is running on port %s", cfg.SERVER_ADDRESS)
	if err := http.ListenAndServe(cfg.SERVER_ADDRESS, r); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
