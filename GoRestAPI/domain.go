package main

// Snippet from DB
type Snippet struct {
	ID    string `json:"id"`
	Lang  string `json:"lang"`
	About string `json:"about"`
	Code  string `json:"code"`
}

// User from DB
type User struct {
	Name     string    `json:"name"`
	Password string    `json:"password"`
	Snippets []Snippet `json:"-"`
}

// LoginCredentials for DB
type LoginCredentials struct {
	Name     string `json:"username,omitempty" bson:"username,omitempty"`
	Password string `json:"password,omitempty" bson:"password,omitempty"`
}
