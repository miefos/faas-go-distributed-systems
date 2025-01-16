package main

import (
	"auth-service/handlers"
	"auth-service/utils"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize NATS connection and KV store
	nats1URL := os.Getenv("NATS1_URL")
	nats2URL := os.Getenv("NATS2_URL")
	myPort := os.Getenv("SERVER_ADDRESS")
	if nats1URL == "" {
		nats1URL = "nats://localhost:4222"
	}

	if(utils.InitNATSConnection(nats1URL) != 0){
		if(utils.InitNATSConnection(nats2URL) != 0){
			log.Fatalf("Error connecting to all NATS servers")
			os.Exit(-1)
		}
	}

	if(utils.InitKVStore("users") != 0){
		log.Fatalf("Error initializing KV store")
		os.Exit(-1)
	}

	// Create a new router
	r := mux.NewRouter()

	// Routes
	r.HandleFunc("/register", handlers.RegisterHandler).Methods("POST")
	r.HandleFunc("/login", handlers.LoginHandler).Methods("POST")
	r.HandleFunc("/validate", handlers.ValidateHandler).Methods("GET")

	// Start server
	log.Println("Auth Service is running on port", myPort)
	log.Fatal(http.ListenAndServe(myPort, r))
}
