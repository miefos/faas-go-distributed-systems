package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"registry-service/utils"

	"registry-service/models"
)

func (h *Handler) UpdateFunction(w http.ResponseWriter, r *http.Request) {
	var metadata models.FunctionMetadata
	if err := json.NewDecoder(r.Body).Decode(&metadata); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID, err := utils.GetUserIDFromToken(r)
	if err != nil {
		log.Printf("Error extracting user ID: %v", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	log.Printf("ID from JWT: %s", userID)

	// Check and set UUID
	if metadata.UUID == "" {
		log.Printf("UUID is empty, setting UUID to userID: %s", userID)
		metadata.UUID = userID
	}

	// Optional: Additional validation if UUID must match userID
	if metadata.UUID != userID {
		log.Printf("UUID does not match userID. UUID: %s, userID: %s", metadata.UUID, userID)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Delete old metadata before updating
	if err := h.KVStore.DeleteFunctionMetadata(metadata.UUID, metadata.Name); err != nil {
		http.Error(w, "Failed to update function metadata", http.StatusInternalServerError)
		return
	}

	if err := h.KVStore.StoreFunctionMetadata(metadata.UUID, &metadata); err != nil {
		http.Error(w, "Failed to update function metadata", http.StatusInternalServerError)
		return
	}

	log.Printf("Function %s updated by user %s", metadata.Name, metadata.UUID)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(metadata)
}
