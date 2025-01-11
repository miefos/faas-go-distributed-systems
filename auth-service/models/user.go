package models

type User struct {
	ID       string `json:"id"` // UUID for the user
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserValue struct {
	ID       string `json:"id"`       // UUID for the user
	Password string `json:"password"` // Hashed password
}
