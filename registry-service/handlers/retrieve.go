package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"registry-service/models"
	"registry-service/utils"
)

func (h *Handler) RetrieveFunction(w http.ResponseWriter, r *http.Request) {
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

	UUID := metadata.UUID
	functionName := metadata.Name

	// Optional: Additional validation if UUID must match userID
	if UUID != userID {
		log.Printf("UUID does not match userID. UUID: %s, userID: %s", UUID, userID)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	res, err := h.KVStore.GetFunctionMetadata(UUID, functionName)
	if err != nil {
		http.Error(w, "Function not found", http.StatusNotFound)
		return
	}

	log.Printf("Function %s retrieved by user %s", res.Name, res.UUID)
	json.NewEncoder(w).Encode(res)
}
