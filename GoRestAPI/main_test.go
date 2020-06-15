package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	jwt "github.com/dgrijalva/jwt-go"
)

var a App

var createdUser User
var createdSnippet Snippet
var currentAuthToken string

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

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func TestMain(m *testing.M) {

	a.Initialize(getEnv("POSTGRES_USER", "admin"), getEnv("POSTGRES_PASSWORD", "123"), getEnv("POSTGRES_DB", "postgres"), getEnv("POSTGRES_HOST_NAME", "host=localhost"), false)
	ensureExtensionExists()
	ensureTablesExist()
	clearTable()

	code := m.Run()

	clearTable()
	os.Exit(code)
}

func ensureExtensionExists() {
	if _, err := a.DB.Exec(extensionQueryUUID); err != nil {
		log.Fatal(err)
	}
}

func ensureTablesExist() {
	if _, err := a.DB.Exec(tableCreationQueryUsers); err != nil {
		log.Fatal(err)
	}

	if _, err := a.DB.Exec(tableCreationQuerySnippets); err != nil {
		log.Fatal(err)
	}
}

// Order is important, as tables with fks need to be deleted first
func clearTable() {
	a.DB.Exec("DELETE FROM snippets")
	a.DB.Exec("DELETE FROM users")
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

func TestWelcomeMessage(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body == "" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func TestCreateUser(t *testing.T) {
	var jsonStr = []byte(`{"mail":"test@mail", "username": "User1", "password": "123"}`)
	req, _ := http.NewRequest("POST", "/user", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m string
	json.Unmarshal(response.Body.Bytes(), &m)

	token, _ := jwt.Parse(m, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte("secret"), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if claims["username"].(string) != "User1" {
			t.Errorf("Expected snippet name to be 'User1'. Got '%v'", claims["username"])
		}

		createdUser = User{ID: claims["user_id"].(string), Name: claims["username"].(string)}
		currentAuthToken = token.Raw
	}
}

func TestCreateSnippet(t *testing.T) {
	snippet1 := Snippet{Owner: createdUser.ID, Title: "snippet1", Category: "HelloWorld", Code: "print(\"Hello World from Go Rest API\")", Lang: "python"}

	snippetMarshal, err := json.Marshal(snippet1)

	if err != nil {
		t.Errorf(`Error marshaling snippet`)
	}

	req, _ := http.NewRequest("POST", fmt.Sprintf(`/user/%s/snippets`, createdUser.ID), bytes.NewBuffer(snippetMarshal))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", currentAuthToken)

	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	snippet := Snippet{}
	err = json.Unmarshal(response.Body.Bytes(), &snippet)

	if err != nil {
		t.Errorf("Error while unmarshaling: " + err.Error())
	}

	snippet1.ID = snippet.ID

	if !reflect.DeepEqual(snippet, snippet1) {
		t.Errorf("Error while deepequal, object values are not the same")
	}

	createdSnippet = snippet
}

func TestCreateSnippet2(t *testing.T) {
	snippet2 := Snippet{Owner: createdUser.ID, Title: "snippet2", Category: "HelloWorld", Code: "print(\"Hello World from Go Rest API\")", Lang: "python"}

	snippetMarshal, err := json.Marshal(snippet2)

	if err != nil {
		t.Errorf(`Error marshaling snippet`)
	}

	req, _ := http.NewRequest("POST", fmt.Sprintf(`/user/%s/snippets`, createdUser.ID), bytes.NewBuffer(snippetMarshal))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", currentAuthToken)

	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	snippet := Snippet{}
	err = json.Unmarshal(response.Body.Bytes(), &snippet)

	if err != nil {
		t.Errorf("Error while unmarshaling: " + err.Error())
	}

	snippet2.ID = snippet.ID

	if !reflect.DeepEqual(snippet, snippet2) {
		t.Errorf("Error while deepequal, object values are not the same")
	}
}

func TestGetSnippet(t *testing.T) {
	connectionString := fmt.Sprintf(`/user/%s/snippets/%s`, createdUser.ID, createdSnippet.ID)
	req, _ := http.NewRequest("GET", connectionString, nil)
	req.Header.Set("Authorization", currentAuthToken)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	snippet := Snippet{}
	err := json.Unmarshal(response.Body.Bytes(), &snippet)

	if err != nil {
		t.Errorf("Error while unmarshaling: " + err.Error())
	}

	if !reflect.DeepEqual(snippet, createdSnippet) {
		t.Errorf("Error while deepequal, object values are not the same")
	}
}

func TestUpdateSnippet(t *testing.T) {

	createdSnippet.Lang = "java"
	createdSnippet.Code = "class Changed {}"

	snippetMarshal, err := json.Marshal(createdSnippet)

	if err != nil {
		t.Errorf(`Error marshaling snippet`)
	}

	req, _ := http.NewRequest("PUT", fmt.Sprintf(`/user/%s/snippets/%s`, createdUser.ID, createdSnippet.ID), bytes.NewBuffer(snippetMarshal))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", currentAuthToken)

	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	snippet := Snippet{}
	err = json.Unmarshal(response.Body.Bytes(), &snippet)

	if err != nil {
		t.Errorf("Error while unmarshaling: " + err.Error())
	}

	if !reflect.DeepEqual(snippet, createdSnippet) {
		t.Errorf("Error while deepequal, object values are not the same")
	}
}

func TestDeleteSnippet(t *testing.T) {

	req, _ := http.NewRequest("DELETE", fmt.Sprintf(`/user/%s/snippets/%s`, createdUser.ID, createdSnippet.ID), nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", currentAuthToken)

	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", fmt.Sprintf(`/user/%s/snippets/%s`, createdUser.ID, createdSnippet.ID), nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", currentAuthToken)

	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}

func TestGetSnippets(t *testing.T) {
	connectionString := fmt.Sprintf(`/user/%s/snippets/`, createdUser.ID)
	req, _ := http.NewRequest("GET", connectionString, nil)
	req.Header.Set("Authorization", currentAuthToken)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	snippets := make([]Snippet, 0)
	err := json.Unmarshal(response.Body.Bytes(), &snippets)

	if err != nil {
		t.Errorf("Error while unmarshaling: " + err.Error())
	}

	createdUser.Snippets = []Snippet{}
}
