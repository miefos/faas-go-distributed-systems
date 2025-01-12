package handlers

import (
	"encoding/json"
	"net/http"
	"log"
	"registry-service/models"
)

func (h *Handler) RetrieveFunction(w http.ResponseWriter, r *http.Request) {
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

	res, err := h.KVStore.GetFunctionMetadata(UUID, functionName)
	if err != nil {
		http.Error(w, "Function not found", http.StatusNotFound)
		return
	}

	log.Printf("Function %s retrieved by user %s", res.Name, res.UUID)
	json.NewEncoder(w).Encode(res)
}