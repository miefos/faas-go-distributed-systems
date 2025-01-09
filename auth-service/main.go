package main

import (
	"auth-service/handlers"
	"auth-service/utils"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	// Initialize NATS connection and KV store
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = "nats://localhost:4222"
	}

	utils.InitNATSConnection(natsURL)
	utils.InitKVStore("users")

	// Create a new router
	r := mux.NewRouter()

	// Routes
	r.HandleFunc("/register", handlers.RegisterHandler).Methods("POST")
	r.HandleFunc("/login", handlers.LoginHandler).Methods("POST")
	r.HandleFunc("/validate", handlers.ValidateHandler).Methods("GET")

	// Start server
	log.Println("Auth Service is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
