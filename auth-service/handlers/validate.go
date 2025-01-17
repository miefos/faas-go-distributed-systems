package handlers

import (
	"auth-service/utils" // Ensure this matches your module name
	"net/http"
	"strings"
)

func ValidateHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the token from the Authorization header
	tokenStr := r.Header.Get("Authorization")

	// Check if the token is missing or doesn't have the 'Bearer ' prefix
	if tokenStr == "" || !strings.HasPrefix(tokenStr, "Bearer ") {
		http.Error(w, "Missing or invalid token", http.StatusUnauthorized)
		return
	}

	// Remove the 'Bearer ' prefix to get the actual token
	tokenStr = tokenStr[7:] // Removes "Bearer " (7 characters)

	// Validate JWT
	claims, err := utils.ValidateJWT(tokenStr)
	if err != nil {
		http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
		return
	}

	// Respond with a success message, including username and id from the token
	_, err = w.Write([]byte("Token is valid. Username: " + claims["username"].(string) + ", ID: " + claims["id"].(string)))
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}
