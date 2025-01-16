package main

import (
	"auth-service/handlers"
	"auth-service/utils"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

// NotFoundHandler handles unmatched routes
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("NotFoundHandler called")

	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintln(w, "404 - Page not found!!")
}

func main() {
	// Initialize NATS connection and KV store
	natsURL := os.Getenv("NATS_URL")
	myPort := os.Getenv("SERVER_ADDRESS")
	if natsURL == "" {
		natsURL = "nats://localhost:4222"
	}

	utils.InitNATSConnection(natsURL)
	utils.InitKVStore("users")

	// Create a new router
	r := mux.NewRouter()

	// Routes
	r.HandleFunc("/auth/register", handlers.RegisterHandler).Methods("POST")
	r.HandleFunc("/auth/login", handlers.LoginHandler).Methods("POST")
	r.HandleFunc("/auth/validate", handlers.ValidateHandler).Methods("GET")
	r.NotFoundHandler = http.HandlerFunc(NotFoundHandler)

	// Start server
	log.Println("Auth Service is running on port", myPort)
	log.Fatal(http.ListenAndServe(myPort, r))
}
