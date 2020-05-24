package main

import (
	"database/sql"
)

// Snippet from DB
type Snippet struct {
	ID    string `json:"id"`
	Lang  string `json:"lang"`
	About string `json:"about"`
	Code  string `json:"code"`
}

// User from DB
type User struct {
	ID       string    `json:"user_id"`
	Name     string    `json:"name"`
	Password string    `json:"password"`
	Snippets []Snippet `json:"-"`
}

// LoginCredentials for DB
type LoginCredentials struct {
	Name     string `json:"username,omitempty" bson:"username,omitempty"`
	Password string `json:"password,omitempty" bson:"password,omitempty"`
}

func (s *Snippet) getSnippet(db *sql.DB) error {
	return db.QueryRow("SELECT id, lang, about, code FROM Snipppets WHERE id=$1",
		s.ID).Scan(&s.ID, &s.Lang, &s.About, &s.Code)
}

func getSnippets(db *sql.DB, user User) ([]Snippet, error) {
	rows, err := db.Query(
		"SELECT id, lang, about, code FROM snippets snippets WHERE user_id=$1",
		user.ID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	snippets := []Snippet{}

	for rows.Next() {
		var s Snippet
		if err := rows.Scan(&s.ID, &s.Lang, &s.About, &s.Code); err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}

	return snippets, nil
}
