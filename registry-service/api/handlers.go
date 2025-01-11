package api

import (
	"encoding/json"
	"net/http"
	"log"

	"registry-service/models"
	"registry-service/storage"

	"github.com/gorilla/mux"
)

type Handler struct {
	KVStore *storage.KVStore
}

func RegisterRoutes(router *mux.Router, kvStore *storage.KVStore) {
	h := &Handler{KVStore: kvStore}

	router.HandleFunc("/register", h.RegisterFunction).Methods("POST")
	router.HandleFunc("/retrieve/{id}", h.RetrieveFunction).Methods("GET")
	router.HandleFunc("/list", h.ListFunctions).Methods("GET")
	router.HandleFunc("/delete/{id}", h.DeleteFunction).Methods("DELETE")
	router.HandleFunc("/update/{id}", h.UpdateFunction).Methods("PUT")
}

func (h *Handler) RegisterFunction(w http.ResponseWriter, r *http.Request) {
	var metadata models.FunctionMetadata
	if err := json.NewDecoder(r.Body).Decode(&metadata); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID := r.Header.Get("UserID") // Example user authentication
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	metadata.UserID = userID

	// Check if the function with the same ID and user ID already exists
	existingMetadata, err := h.KVStore.GetFunctionMetadata(userID, metadata.ID)
	if err == nil && existingMetadata != nil {
		http.Error(w, "Function with the same ID already exists for this user", http.StatusConflict)
		return
	}

	if err := h.KVStore.StoreFunctionMetadata(userID, &metadata); err != nil {
		http.Error(w, "Failed to store function metadata", http.StatusInternalServerError)
		return
	}

	log.Printf("Function %s registered by user %s", metadata.ID, metadata.UserID)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(metadata)
}

func (h *Handler) RetrieveFunction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	functionID := vars["id"]
	userID := r.Header.Get("UserID")

	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	metadata, err := h.KVStore.GetFunctionMetadata(userID, functionID)
	if err != nil {
		http.Error(w, "Function not found", http.StatusNotFound)
		return
	}

	log.Printf("Function %s retrieved by user %s", metadata.ID, metadata.UserID)
	json.NewEncoder(w).Encode(metadata)
}

func (h *Handler) ListFunctions(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("UserID")

	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	metadataList, err := h.KVStore.ListFunctions(userID)
	if err != nil {
		http.Error(w, "Failed to retrieve functions", http.StatusInternalServerError)
		return
	}

	log.Printf("Functions listed for user %s", userID)
	json.NewEncoder(w).Encode(metadataList)
}

func (h *Handler) DeleteFunction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	functionID := vars["id"]
	userID := r.Header.Get("UserID")

	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err := h.KVStore.DeleteFunctionMetadata(userID, functionID); err != nil {
		http.Error(w, "Failed to delete function metadata", http.StatusInternalServerError)
		return
	}

	log.Printf("Function %s deleted by user %s", functionID, userID)
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) UpdateFunction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	functionID := vars["id"]
	userID := r.Header.Get("UserID")

	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var metadata models.FunctionMetadata
	if err := json.NewDecoder(r.Body).Decode(&metadata); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	metadata.ID = functionID
	metadata.UserID = userID

	// Delete old metadata before updating
	if err := h.KVStore.DeleteFunctionMetadata(userID, functionID); err != nil {
		http.Error(w, "Failed to update function metadata", http.StatusInternalServerError)
		return
	}

	if err := h.KVStore.StoreFunctionMetadata(userID, &metadata); err != nil {
		http.Error(w, "Failed to update function metadata", http.StatusInternalServerError)
		return
	}

	log.Printf("Function %s updated by user %s", metadata.ID, metadata.UserID)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(metadata)
}
