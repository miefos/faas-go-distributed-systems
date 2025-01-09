package handlers

import (
	"auth-service/models"
	"auth-service/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// RegisterHandler handles user registration
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("RegisterHandler called")

	// Decode the user data from the request body
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	log.Println(user)

	// Validate that the username is not taken by checking the NATS KV store
	existingUser, err := utils.KvStore.Get(user.Username)
	if err == nil && existingUser != nil {
		http.Error(w, "Username already taken", http.StatusConflict)
		return
	}

	// Hash the password
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	// Store the hashed password in the NATS KV store
	result, err := utils.KvStore.Put(user.Username, []byte(hashedPassword))
	if err != nil {
		http.Error(w, "Error storing user data", http.StatusInternalServerError)
		return
	}

	// Send a success response
	log.Printf("User %s registered successfully", user.Username)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("User %s (%+v) registered successfully", user.Username, result)))
}
