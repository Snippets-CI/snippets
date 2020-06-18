package main

import (
	"database/sql"
	"log"

	"golang.org/x/crypto/bcrypt"
)

const extensionQueryUUID = `
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
`

const tableCreationQueryUsers = `
CREATE TABLE IF NOT EXISTS users (
	user_id uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
	username VARCHAR(25) NOT NULL,
    mail VARCHAR(50) NOT NULL,
	password VARCHAR(60) NOT NULL,
    UNIQUE(mail)
)`

const tableCreationQuerySnippets = `
CREATE TABLE IF NOT EXISTS snippets (
	snippet_id uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
	owner UUID REFERENCES users(user_id) NOT NULL,
	language VARCHAR(20) NULL,
	title VARCHAR(30) NULL,
	category VARCHAR(50) NULL,
	code VARCHAR NULL
)`

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

func ensureTablesExist(db *sql.DB) {
	if _, err := db.Exec(tableCreationQueryUsers); err != nil {
		log.Fatal(err)
	}

	if _, err := db.Exec(tableCreationQuerySnippets); err != nil {
		log.Fatal(err)
	}
}

func ensureExtensionExists(db *sql.DB) {
	if _, err := db.Exec(extensionQueryUUID); err != nil {
		log.Fatal(err)
	}
}

func (user *User) getUser(db *sql.DB) error {
	var dbUser User

	err := db.QueryRow(`SELECT user_id, mail, username, password FROM "users" WHERE mail=$1`, user.Mail).Scan(&dbUser.ID, &dbUser.Mail, &dbUser.Name, &dbUser.Password)

	// If Query found user with matching email, check the hashed password
	if err == nil {
		err = comparePasswords(dbUser.Password, user.Password)
		user.ID = dbUser.ID
		user.Mail = dbUser.Mail
		user.Name = dbUser.Name
		user.Password = ""
	}

	if err != nil {
		err = user.getSnippets(db)
	}

	return err
}

func (user *User) createUser(db *sql.DB) error {
	pwd := hashAndSalt(user.Password)

	err := db.QueryRow(`INSERT INTO "users" (mail, username, password) VALUES ($1, $2, $3) RETURNING user_id`,
		user.Mail, user.Name, pwd).Scan(&user.ID)

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

func hashAndSalt(password string) string {
	bPwd := []byte(password)
	hash, _ := bcrypt.GenerateFromPassword(bPwd, bcrypt.MinCost)

	return string(hash)
}

func comparePasswords(hashedPwd string, plainPwd string) error {
	bPlainPwd := []byte(plainPwd)
	bHashedPwd := []byte(hashedPwd)

	err := bcrypt.CompareHashAndPassword(bHashedPwd, bPlainPwd)
	return err
}
