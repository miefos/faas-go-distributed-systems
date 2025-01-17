package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"registry-service/models"
	"registry-service/utils"
)

func (h *Handler) DeleteFunction(w http.ResponseWriter, r *http.Request) {
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

	// Optional: Validate UUID matches userID
	if UUID != userID {
		log.Printf("UUID does not match userID. UUID: %s, userID: %s", UUID, userID)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err := h.KVStore.DeleteFunctionMetadata(UUID, functionName); err != nil {
		http.Error(w, "Failed to delete function metadata", http.StatusInternalServerError)
		return
	}

	log.Printf("Function %s deleted by user %s", functionName, UUID)
	w.WriteHeader(http.StatusNoContent)
}
