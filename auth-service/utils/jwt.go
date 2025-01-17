package utils

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var jwtKey []byte

func init() {
	// Read the JWT secret key from the secret.txt file
	key, err := os.ReadFile("secret.txt")
	if err != nil {
		fmt.Println("Error reading the secret file:", err)
		os.Exit(1)
	}
	jwtKey = key
}

func GenerateJWT(id string, username string) (string, error) {
	// Define claims similar to the Python structure
	claims := jwt.MapClaims{
		"id":       id,                                    // Same as Python
		"username": username,                              // username
		"iss":      "faas-app",                            // Issuer
		"key":      "faas-app-key",                        // Custom key
		"iat":      time.Now().Unix(),                     // Explicit issued-at time (current time)
		"exp":      time.Now().Add(time.Hour * 72).Unix(), // 3 day expiration
	}

	// Create the token with HS256 method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtKey)
}

func ValidateJWT(tokenStr string) (map[string]interface{}, error) {
	// Parse the JWT token
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	// If there's an error or the token is not valid, return an error
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	// Extract claims if token is valid
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	}

	log.Printf("Token claims: %+v\n", token.Claims)

	return nil, fmt.Errorf("invalid token claims")
}
