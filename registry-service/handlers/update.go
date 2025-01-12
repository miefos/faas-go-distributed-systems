package handlers

import (
	"encoding/json"
	"net/http"
	"log"

	"registry-service/models"
)

func (h *Handler) UpdateFunction(w http.ResponseWriter, r *http.Request) {
	var metadata models.FunctionMetadata
	if err := json.NewDecoder(r.Body).Decode(&metadata); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	if metadata.UUID == "" {
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
