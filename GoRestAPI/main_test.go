package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var a App

const tableCreationQueryUsers = `
CREATE TABLE IF NOT EXISTS users (
	user_id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
	username VARCHAR(25) NOT NULL,
	password VARCHAR(60) NOT NULL
)`

const tableCreationQuerySnippets = `
CREATE TABLE IF NOT EXISTS snippets (
        snippet_id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
        owner UUID REFERENCES users(user_id) NOT NULL,
        language VARCHAR(20) NULL,
        title VARCHAR(30) NULL,
        category VARCHAR(50) NULL,
        code VARCHAR NULL
)`

func TestMain(m *testing.M) {
	a.Initialize("admin", "123", "SnippetsTest")
	ensureTablesExist()

	code := m.Run()

	clearTable()
	os.Exit(code)
}

func ensureTablesExist() {
	if _, err := a.DB.Exec(tableCreationQueryUsers); err != nil {
		log.Fatal(err)
	}

	if _, err := a.DB.Exec(tableCreationQuerySnippets); err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	a.DB.Exec("DELETE FROM users")
	a.DB.Exec("DELETE FROM snippets")
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

/* ******************************
 * Testing functions
 *******************************/

func TestEmptyUserTable(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/users", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func TestEmptySnippetTable(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/users", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func TestWelcomeMessage(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body == "" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func TestCreateUser(t *testing.T) {
	var jsonStr = []byte(`{"user_id":"1", "name": "User1", "password": "123"}`)
	req, _ := http.NewRequest("POST", "/user", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["user_id"] != "1" {
		t.Errorf("Expected snippet id to be '1'. Got '%v'", m["user_id"])
	}

	if m["name"] != "User1" {
		t.Errorf("Expected snippet name to be 'User1'. Got '%v'", m["name"])
	}
}

func TestCreateSnippet(t *testing.T) {
	var jsonStr = []byte(`{"id":"1", "name": "snippet1", "lang": "python", "about": "Hello World", "code": "print("Hello World from Go Rest API")"}`)
	req, _ := http.NewRequest("POST", "/product", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["id"] != "1" {
		t.Errorf("Expected snippet id to be '1'. Got '%v'", m["id"])
	}

	if m["name"] != "snippet1" {
		t.Errorf("Expected snippet name to be 'snippet1'. Got '%v'", m["name"])
	}

	if m["lang"] != "python" {
		t.Errorf("Expected snippet lang to be 'python'. Got '%v'", m["lang"])
	}

	if m["about"] != "Hello World" {
		t.Errorf("Expected snippet about to be 'Hello World'. Got '%v'", m["about"])
	}

	if m["code"] != `print("Hello World from Go Rest API")` {
		t.Errorf(`Expected snippet about to be 'print("Hello World from Go Rest API")'. Got '%v'`, m["code"])
	}
}
