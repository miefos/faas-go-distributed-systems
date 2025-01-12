package handlers

import (
	"encoding/json"
	"net/http"
	"log"
	"registry-service/models"
)

func (h *Handler) ListFunctions(w http.ResponseWriter, r *http.Request) {
	var metadata models.FunctionMetadata
	if err := json.NewDecoder(r.Body).Decode(&metadata); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if metadata.UUID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	metadataList, err := h.KVStore.ListFunctions(metadata.UUID)
	if err != nil {
		http.Error(w, "Failed to retrieve functions", http.StatusInternalServerError)
		return
	}

	log.Printf("Functions listed for user %s", metadata.UUID)
	json.NewEncoder(w).Encode(metadataList)
}