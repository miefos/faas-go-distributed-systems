package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"publisher-service/models"
	"time"

	"github.com/nats-io/nats.go"
)

//<Da fuori riceve ID dell'utente
// fa il retrive con quel ID della funzione
// ottiene function ID della funzione

type PublisherHandler struct {
	NATSConn     *nats.Conn
	RequestTopic string
	ReplyTimeout int
}

func NewPublisherHandler(nc *nats.Conn, requestTopic string, replyTimeout int) *PublisherHandler { //constructor PublishHandler, associato nel main
	return &PublisherHandler{
		NATSConn:     nc,
		RequestTopic: requestTopic,
		ReplyTimeout: replyTimeout,
	}
}

func (h *PublisherHandler) PublishHandlerMethod(w http.ResponseWriter, r *http.Request) {
	log.Printf("Publish Handler called")

	var functionRequest models.FunctionMetadata

	// Decodifica il corpo della richiesta in un oggetto FunctionMetadata
	if err := json.NewDecoder(r.Body).Decode(&functionRequest); err != nil {
		http.Error(w, "Errore nella deserializzazione della richiesta: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Stampa informazioni di debug
	log.Printf("UUID utente: %s", functionRequest.UUID)
	log.Printf("Name: %s", functionRequest.Name)
	log.Printf("Description: %s", functionRequest.Description)
	log.Printf("Payload: %s", functionRequest.Payload)

	// Ottieni la funzione completa
	newFunctionObject, err := getFunction(&functionRequest)
	if err != nil {
		http.Error(w, "Error retrieve function: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Converte l'oggetto in JSON per inviarlo al worker
	payload, err := json.Marshal(newFunctionObject)
	if err != nil {
		http.Error(w, "Error deserealization function: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Pubblica sulla coda NATS e ottieni la risposta
	msg, err := h.NATSConn.Request(h.RequestTopic, payload, time.Duration(h.ReplyTimeout)*time.Second)
	if err != nil {
		http.Error(w, "Error submission worker "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Rispondi al client con i dati del worker
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(msg.Data)
}

func getFunction(function *models.FunctionMetadata) (*models.FunctionMetadata, error) {
	// Prepara la query string
	query := url.Values{}
	query.Set("uuid", function.UUID)
	query.Set("name", function.Name)

	// Costruisce l'URL completo
	baseURL := "http://localhost:8082/retrieve"
	fullURL := fmt.Sprintf("%s?%s", baseURL, query.Encode())

	// Effettua la richiesta GET
	log.Printf("Performing GET request to %s", fullURL)
	resp, err := http.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("errore during GET request: %w", err)
	}
	defer resp.Body.Close()

	// Verifica lo status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response error register-service: %s", resp.Status)
	}

	// Legge e deserializza la risposta JSON
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error read body JSON: %w", err)
	}

	var result models.FunctionMetadata //contains functionID, functionName
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("error read body JSON: %w", err)
	}

	log.Printf("UUID funzione: %s", result.UUID)
	log.Printf("Name: %s", result.Name)
	log.Printf("Description: %s", result.Description)
	log.Printf("Payload: %s", result.Payload)

	return &result, nil
}
