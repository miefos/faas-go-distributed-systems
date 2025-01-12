package handlers

import (
	"encoding/json"
	"net/http"
	"log"
	"registry-service/models"
)

func (h *Handler) DeleteFunction(w http.ResponseWriter, r *http.Request) {
	var metadata models.FunctionMetadata
	if err := json.NewDecoder(r.Body).Decode(&metadata); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	UUID := metadata.UUID
	functionName := metadata.Name

	if UUID == "" {
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