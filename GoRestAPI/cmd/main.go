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
	Snippets []Snippet `json:"snippets"`
}

var users []User

// Creates Mockdata which in future we will get from the database
func createMockData() {
	snippet1 := Snippet{"snippet1", "java", "stdout", "System.out.println(\"Hello World\");"}
	snippet2 := Snippet{"snippet2", "go", "structure", "type User struct {\n    name string\n    lastname string\n    age int\n} "}
	snippet3 := Snippet{"snippet3", "python", "stdout", "print(\"Hello World\")"}
	snippet4 := Snippet{"snippet4", "javascript", "objects", "var car = {\n    type:\"Fiat\",\n    model:\"500\",\n     color:\"white\"\n};"}

	user1Snippets := []Snippet{snippet1, snippet2}
	user2Snippets := []Snippet{snippet3, snippet4}

	users = append(users, User{"Andreas", "asdf", user1Snippets})
	users = append(users, User{"Markus", "1234", user2Snippets})
}

func getWelcomeMessage(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode([]string{"Welcome to SNIPPETS!"})
}

// Outputs the requested user as json
func getUser(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userID")

	for _, user := range users {
		if user.Name == userId {
			user_json, _ := json.Marshal(users)

			json.NewEncoder(w).Encode(string(user_json))
		}
	}
}

func registerRoutes() http.Handler {
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Get("/", getWelcomeMessage)
	})

	// route: /user/{userName}
	// name equals id
	r.Route("/user", func(r chi.Router) {
		r.Route("/{userID}", func(r chi.Router) {
			r.Get("/", getUser)
		})
		r.Get("/", getUser)
	})
	return r
}

func main() {
	createMockData()

	fmt.Println("REST API started!")
	r := registerRoutes()
	log.Fatal(http.ListenAndServe(":8080", r))
}
