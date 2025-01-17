package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"net/http"
)

func GetUserIDFromToken(r *http.Request) (string, error) {
	header := r.Header.Get("Authorization")
	if header == "" {
		return "", fmt.Errorf("no authorization header")
	}

	// Extract the token from the Authorization header
	const bearerPrefix = "Bearer "
	if len(header) <= len(bearerPrefix) || header[:len(bearerPrefix)] != bearerPrefix {
		return "", fmt.Errorf("invalid authorization header")
	}

	tokenString := header[len(bearerPrefix):]

	// Parse the token without validating
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		log.Printf("Error parsing token: %v", err)
		return "", err
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("invalid token claims")
	}

	// Safely access the "user_id" claim
	id, exists := claims["id"]
	if !exists {
		return "", fmt.Errorf("id not found in token")
	}

	userID, ok := id.(string)
	if !ok {
		return "", fmt.Errorf("id claim is not a string")
	}

	return userID, nil
}
