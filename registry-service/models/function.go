package models

type FunctionMetadata struct {
	UUID      string `json:"uuid"` // User ID
	Name        string `json:"name"` // Function name
	Dscription string `json:"description"` // Function description
	Payload     string `json:"payload"` // Function reference to the docker image
}
