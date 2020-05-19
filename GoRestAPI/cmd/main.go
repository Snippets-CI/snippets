package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	_ "github.com/lib/pq"
)

type Snippet struct {
	Id    string `json:"id"`
	Lang  string `json:"lang"`
	About string `json:"about"`
	Code  string `json:"code"`
}

type User struct {
	Name     string    `json:"name"`
	Password string    `json:"password"`
	Snippets []Snippet `json:"-"`
}

type LoginCredentials struct {
	Name     string `json:"username,omitempty" bson:"username,omitempty"`
	Password string `json:"password,omitempty" bson:"password,omitempty"`
}

const WELCOME_MESSAGE = "Welcome to SNIPPETS!"
const ERROR_MESSAGE = "Error Resource not found"
const ERROR_MESSAGE_MAILFORMED_JSON = "Invalid Json found"

var users = []User{}

// Creates Mockdata which in future we will get from the database
func createMockData() {
	snippet1 := Snippet{`snippet1`, `java`, `stdout`, `System.out.println("Hello World");`}
	snippet2 := Snippet{`snippet2`, `go`, `structure`, `type User struct {\n    name string\n    lastname string\n    age int\n} `}
	snippet3 := Snippet{`snippet3`, `python`, `stdout`, `print("Hello World")`}
	snippet4 := Snippet{`snippet4`, `javascript`, `objects`, `var car = {\n    type:"Fiat",\n    model:"500",\n     color:"white"\n};`}

	user1Snippets := []Snippet{snippet1, snippet2}
	user2Snippets := []Snippet{snippet3, snippet4}

	users = append(users, User{"Andreas", "asdf", user1Snippets})
	users = append(users, User{"Markus", "1234", user2Snippets})
}

func getUser(userID string) User {
	for _, user := range users {
		if user.Name == userID {
			return user
		}
	}
	return User{}
}

func getUserSnippet(user User, snippetID string) Snippet {
	for _, snippet := range user.Snippets {
		if snippet.Id == snippetID {
			return snippet
		}
	}
	return Snippet{}
}

// GET - Request
// Outputs Welcome Msg
func getWelcomeMessage(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(WELCOME_MESSAGE)
}

// GET - Request
// Outputs all available names of users
func getUsers(w http.ResponseWriter, r *http.Request) {
	availableUsers := []string{}

	for _, user := range users {
		availableUsers = append(availableUsers, user.Name)
	}

	json.NewEncoder(w).Encode(availableUsers)
}

// GET - Request
// Outputs the requested user as json
func getUserDetails(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userID")

	user := getUser(userId)

	user_json, _ := json.Marshal(user)
	json.NewEncoder(w).Encode(string(user_json))

}

// GET - Request
// Outputs all saved Snippets of a user
func getUserSnippets(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userID")

	user := getUser(userId)

	snippets_json, _ := json.Marshal(user.Snippets)
	json.NewEncoder(w).Encode(string(snippets_json))
}

// GET - Request
// Outputs the requested snipped of a user
func getUserSnippetDetails(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userID")
	snippetId := chi.URLParam(r, "snippetID")

	user := getUser(userId)
	snippet := getUserSnippet(user, snippetId)

	snippet_json, _ := json.Marshal(snippet)
	json.NewEncoder(w).Encode(string(snippet_json))
}

// PUT - Request
// Updates a snippet of a user
func putUserSnippetDetails(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userID")
	snippetId := chi.URLParam(r, "snippetID")
	snippetString := chi.URLParam(r, "snippet")
	snippet := Snippet{}
	err := json.Unmarshal([]byte(snippetString), &snippet)

	fmt.Println(r.Context().Value("snippet"))
	fmt.Println(userId + " - " + snippetId + " - " + snippetString)

	if err != nil {
		json.NewEncoder(w).Encode(ERROR_MESSAGE)
	}

	fmt.Println(snippet)
}

// POST - Reqest
// Authenticates user

func loginUser(w http.ResponseWriter, r *http.Request) {

	dbUserName, dbPassword := getUserCredentialsDB()

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	var credentials LoginCredentials
	err := dec.Decode(&credentials)

	if err != nil {
		json.NewEncoder(w).Encode(ERROR_MESSAGE_MAILFORMED_JSON)
	}

	fmt.Println("########")
	fmt.Println(dbUserName)
	fmt.Println(dbPassword)
}

func getUserCredentialsDB() (string, string) {
	// Quick and Dirty
	// Should be refactored
	db, err := sql.Open("postgres",
		"postgresql://restapi@localhost:26257/snippet?sslmode=disable")
	if err != nil {
		log.Fatal("error connecting to the database: ", err)
	}

	//TODO update selcet query and handle response
	rows, err := db.Query("SELECT username, password FROM users")
	if err != nil {
		log.Fatal(err)
	}

	//Close after for loop
	defer rows.Close()

	for rows.Next() {
		var username, password string
		if err := rows.Scan(&username, &password); err != nil {
			log.Fatal(err)
		}
		return username, password
	}

	return "sdf", "asfd"
}

func registerRoutes() http.Handler {
	r := chi.NewRouter()

	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Route("/", func(r chi.Router) {
		r.Get("/", getWelcomeMessage)
	})

	r.Route("/login", func(r chi.Router) {
		r.Post("/", loginUser)
	})

	// route: /users
	r.Route("/users", func(r chi.Router) {
		r.Get("/", getUsers)

		// route: /users/{userID}
		r.Route("/{userID}", func(r chi.Router) {
			r.Get("/", getUserDetails)

			// route: /users/{userID}/snippets
			r.Route("/snippets", func(r chi.Router) {
				r.Get("/", getUserSnippets)

				// route: /users/{userID}/snippets/{snippetID}
				r.Route("/{snippetID}", func(r chi.Router) {
					r.Get("/", getUserSnippetDetails)
					r.Post("/", putUserSnippetDetails)
				})
			})
		})
	})
	return r
}

func main() {
	createMockData()

	fmt.Println("REST API started!")
	r := registerRoutes()
	log.Fatal(http.ListenAndServe(":8080", r))
}
