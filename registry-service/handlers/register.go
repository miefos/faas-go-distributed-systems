package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"registry-service/models"
	"registry-service/storage"

	"github.com/gorilla/mux"
)

type Handler struct {
	KVStore *storage.KVStore
}

func RegisterRoutes(router *mux.Router, kvStore *storage.KVStore) {
	h := &Handler{KVStore: kvStore}

	router.HandleFunc("/registry/register", h.RegisterFunction).Methods("POST")
	router.HandleFunc("/registry/retrieve", h.RetrieveFunction).Methods("GET")
	router.HandleFunc("/registry/list", h.ListFunctions).Methods("GET")
	router.HandleFunc("/registry/delete", h.DeleteFunction).Methods("DELETE")
	// router.HandleFunc("/update", h.UpdateFunction).Methods("PUT")

	router.HandleFunc("/registry/health", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Health check")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "OK")
	}).Methods("GET")
}

func (h *Handler) RegisterFunction(w http.ResponseWriter, r *http.Request) {
	var metadata models.FunctionMetadata
	if err := json.NewDecoder(r.Body).Decode(&metadata); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get UUID from body
	if metadata.UUID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Check if the function with the same name and same UUID already exists
	if _, err := h.KVStore.GetFunctionMetadata(metadata.UUID, metadata.Name); err == nil {
		http.Error(w, "Function already exists", http.StatusConflict)
		return
	}

	if err := h.KVStore.StoreFunctionMetadata(metadata.UUID, &metadata); err != nil {
		http.Error(w, "Failed to store function metadata", http.StatusInternalServerError)
		return
	}

	log.Printf("Function %s registered by user %s", metadata.Name, metadata.UUID)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(metadata)
}
