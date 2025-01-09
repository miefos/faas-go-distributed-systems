package handlers

import (
	"auth-service/utils" // Ensure this matches your module name
	"encoding/json"
	"net/http"

	"auth-service/models"
)

// LoginHandler handles user authentication
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// TODO: Fetch the hashed password for the username from NATS KV store
	storedHashedPassword := "" // Replace with fetched value

	// Check if the password matches
	if !utils.CheckPasswordHash(user.Password, storedHashedPassword) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate JWT
	token, err := utils.GenerateJWT(user.Username)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(token))
}
