package models

type FunctionMetadata struct {
	ID          string `json:"id"`
	UserID      string `json:"user_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Payload     string `json:"payload"` // Function code
}
