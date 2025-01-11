package handlers

import (
	"auth-service/models"
	"auth-service/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
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

	// Generate a UUID for the user
	user.ID = uuid.NewString()

	// Hash the password
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	// Store the user data in the NATS KV store (with hashed password and ID)
	userData := models.UserValue{
		ID:       user.ID,
		Password: hashedPassword,
	}

	userDataJSON, err := json.Marshal(userData)
	if err != nil {
		http.Error(w, "Error marshaling user data", http.StatusInternalServerError)
		return
	}

	result, err := utils.KvStore.Put(user.Username, userDataJSON)
	if err != nil {
		http.Error(w, "Error storing user data", http.StatusInternalServerError)
		return
	}

	// Send a success response
	log.Printf("User %s registered successfully with ID %s", user.Username, user.ID)
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(fmt.Sprintf("User %s (ID: %s) registered successfully: %+v", user.Username, user.ID, result)))
	if err != nil {
		return
	}
}
