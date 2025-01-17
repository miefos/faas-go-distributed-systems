package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"registry-service/utils"
)

func (h *Handler) ListFunctions(w http.ResponseWriter, r *http.Request) {
	// Extract userID (or UUID) from the token
	userID, err := utils.GetUserIDFromToken(r)
	if err != nil {
		log.Printf("Error extracting user ID: %v", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Retrieve the list of functions for the user
	metadataList, err := h.KVStore.ListFunctions(userID)
	if err != nil {
		log.Printf("Failed to retrieve functions for user %s: %v", userID, err)
		http.Error(w, "Failed to retrieve functions", http.StatusInternalServerError)
		return
	}

	// Log and respond with the list of functions
	log.Printf("Functions listed for user %s", userID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metadataList)
}
