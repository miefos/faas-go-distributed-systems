package handlers

import (
	"auth-service/models"
	"auth-service/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// LoginHandler handles user authentication
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the request body to get the user credentials
	var userPayload models.User
	err := json.NewDecoder(r.Body).Decode(&userPayload)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	log.Printf("Login attempt: %s", userPayload.Username)

	// Fetch the hashed password from NATS KV store using the username
	entry, err := utils.KvStore.Get(userPayload.Username)
	if err != nil || entry == nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Unmarshal the user data from the NATS KV store
	var userValue models.UserValue
	err = json.Unmarshal(entry.Value(), &userValue)
	if err != nil {
		http.Error(w, "Error unmarshaling user data", http.StatusInternalServerError)
		return
	}

	// Extract the hashed password from the user data
	password := userValue.Password
	ID := userValue.ID

	log.Printf("User ID: %s", ID)

	// Validate the password
	if !utils.CheckPasswordHash(userPayload.Password, password) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	// Generate JWT token for the user after successful authentication
	token, err := utils.GenerateJWT(ID, userPayload.Username)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	// Send the token back in the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(fmt.Sprintf(`{"token": "%s"}`, token)))
	if err != nil {
		return
	}
}
