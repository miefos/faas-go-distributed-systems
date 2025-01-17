package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"publisher-service/models"
	"time"

	"github.com/nats-io/nats.go"
)

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

	// Request body coming from the outside
	var incomingRequest models.FunctionMetadata

	if err := json.NewDecoder(r.Body).Decode(&incomingRequest); err != nil {
		http.Error(w, "Error Deserealization incomingRequest: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Print debug informations
	log.Printf("UUID utente: %s", incomingRequest.UUID)
	log.Printf("Name: %s", incomingRequest.Name)
	log.Printf("Arguments: %s", incomingRequest.Argument)

	functionArgument := incomingRequest.Argument

	// Obtain all information from registry
	retrievedCompleteImage, err := getFunction(&incomingRequest)
	if err != nil {
		http.Error(w, "Error retrieve function: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert returned object to a json of tyoe RegistryFunction
	var jsonCompleteImage models.RegistryFunction
	err = json.Unmarshal(retrievedCompleteImage, &jsonCompleteImage)
	if err != nil {
		http.Error(w, "Error unmarshal json: "+err.Error(), http.StatusInternalServerError)
		return
	}

	imageReference := jsonCompleteImage.Payload

	// Build payload to send to spawner as a json
	payload := []byte(fmt.Sprintf(`{"image_reference": "%s", "parameter": "%s"}`, imageReference, functionArgument))

	// Public response on NAT Worker Queue
	msg, err := h.NATSConn.Request(h.RequestTopic, payload, time.Duration(h.ReplyTimeout)*time.Second)
	if err != nil {
		http.Error(w, "Error submission worker "+err.Error(), http.StatusInternalServerError)
		return
	}

	//return worker data
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(msg.Data)
}

func getFunction(function *models.FunctionMetadata) ([]byte, error) {
	// Prepare JSON body for querying the registry
	queryBody, err := json.Marshal(map[string]string{
		"uuid": function.UUID,
		"name": function.Name,
	})
	if err != nil {
		return nil, fmt.Errorf("error marshaling JSON: %w", err)
	}

	// Create URL
	baseURL := "http://registry-service:8082/registry/retrieve"

	// Create a new request with JSON body
	req, err := http.NewRequest("GET", baseURL, bytes.NewBuffer(queryBody))
	if err != nil {
		return nil, fmt.Errorf("error creating new request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error during GET request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response error register-service: %s", resp.Status)
	}

	// Deserialize JSON body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading body JSON: %w", err)
	}

	return body, nil
}
