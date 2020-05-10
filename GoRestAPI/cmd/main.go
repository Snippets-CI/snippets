package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
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

const WELCOME_MESSAGE = "Welcome to SNIPPETS!"
const ERROR_MESSAGE = "Error Resource not found"

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

func registerRoutes() http.Handler {
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Get("/", getWelcomeMessage)
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
