package handlers

import (
	"auth-service/utils" // Ensure this matches your module name
	"net/http"
)

// ValidateHandler verifies JWT token
func ValidateHandler(w http.ResponseWriter, r *http.Request) {
	tokenStr := r.Header.Get("Authorization")
	if tokenStr == "" {
		http.Error(w, "Missing token", http.StatusUnauthorized)
		return
	}

	// Validate JWT
	claims, err := utils.ValidateJWT(tokenStr)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	_, err = w.Write([]byte("Token is valid. Username: " + claims["username"].(string) + ", ID: " + claims["id"].(string)))
	if err != nil {
		return
	}
}
