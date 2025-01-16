package models

type FunctionMetadata struct {
	UUID        string `json:"uuid"`        // User ID
	Name        string `json:"name"`        // Function name
	Argument string `json:"argument"` // String argument to pass to the function
}

type RegistryFunction struct {
	UUID 	  string `json:"uuid"`        // User ID
	Name      string `json:"name"`        // Function name
	Description string `json:"description"` // Function description
	Payload string `json:"payload"` // The image reference
}
