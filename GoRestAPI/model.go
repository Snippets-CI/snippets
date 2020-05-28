package main

import (
	"database/sql"
)

// Snippet from DB
// HTTP status code 200 and Snippet
// swagger:response Snippet
type Snippet struct {
	ID       string `json:"snippet_id"`
	Owner    string `json:"owner"`
	Title    string `json:"title"`
	Lang     string `json:"language"`
	Category string `json:"category"`
	Code     string `json:"code"`
}

// User from DB
// HTTP status code 200 and User
// swagger:response User
type User struct {
	ID       string    `json:"user_id"`
	Mail     string    `json:"mail"`
	Name     string    `json:"username"`
	Password string    `json:"password"`
	Snippets []Snippet `json:"-"`
}

// LoginCredentials for DB
type LoginCredentials struct {
	Name     string `json:"username,omitempty" bson:"username,omitempty"`
	Password string `json:"password,omitempty" bson:"password,omitempty"`
}

func (user *User) getUser(db *sql.DB) error {
	err := db.QueryRow(`SELECT user_id, mail, username FROM "users" WHERE mail=$1`, user.Mail).Scan(&user.ID, &user.Mail, &user.Name)

	if err != nil {
		err = user.getSnippets(db)
	}

	return err
}

func (user *User) createUser(db *sql.DB) error {
	err := db.QueryRow(`INSERT INTO "users" (mail, username, password) VALUES ($1, $2, $3) RETURNING user_id`,
		user.Mail, user.Name, user.Password).Scan(&user.ID)

	if err != nil {
		return err
	}

	return nil
}

func (user *User) updateUser(db *sql.DB) error {
	_, err := db.Exec(`UPDATE "users" SET mail=$1, username=$2 WHERE user_id=$3`, user.Mail, user.Name, user.ID)

	return err
}

func (s *Snippet) getSnippet(db *sql.DB) error {
	return db.QueryRow(`SELECT language, title, code, category FROM "snippets" WHERE snippet_id=$1 AND owner=$2`,
		s.ID, s.Owner).Scan(&s.Lang, &s.Title, &s.Code, &s.Category)
}

func (s *Snippet) updateSnippet(db *sql.DB) error {
	_, err := db.Exec(`UPDATE "snippets" SET title=$1, language=$2, category=$3, code=$4 WHERE snippet_id=$5`,
		s.Title, s.Lang, s.Category, s.Code, s.ID)

	return err
}

func (s *Snippet) createSnippet(db *sql.DB) error {
	err := db.QueryRow(`INSERT INTO "snippets"(owner, title, language, category, code) VALUES($1, $2, $3, $4, $5) RETURNING snippet_id`,
		s.Owner, s.Title, s.Lang, s.Category, s.Code).Scan(&s.ID)

	if err != nil {
		return err
	}

	return nil
}

func (s *Snippet) deleteSnippet(db *sql.DB) error {
	_, err := db.Exec(`DELETE FROM "snippets" WHERE snippet_id=$1`, s.ID)
	return err
}

func (user *User) getSnippets(db *sql.DB) error {
	var snippets, err = getSnippets(db, user.ID)
	user.Snippets = snippets

	return err
}

func getSnippets(db *sql.DB, userID string) ([]Snippet, error) {
	rows, err := db.Query(
		`SELECT snippet_id, language, title, category, code FROM "snippets" WHERE owner=$1`,
		userID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	snippets := []Snippet{}

	for rows.Next() {
		var s Snippet
		if err := rows.Scan(&s.ID, &s.Lang, &s.Title, &s.Category, &s.Code); err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}

	return snippets, nil
}
