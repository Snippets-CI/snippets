package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	_ "github.com/lib/pq"
)

// App struct with router and db specification
type App struct {
	Router *chi.Mux
	DB     *sql.DB
}

// Initialize app and connect to db
func (a *App) Initialize(user, password, dbname string) {
	fmt.Println("[*] Initialize...")
	createMockData()

	connectionString :=
		fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, dbname)

	var err error
	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	a.Router = chi.NewRouter()
	a.Router.Use(middleware.Logger)

	a.initializeRoutes()
}

// Run http listen and serve
func (a *App) Run(addr string) {
	fmt.Println("[*] ListenAndServe...")
	log.Fatal(http.ListenAndServe(":8010", a.Router))
}

func (a *App) initializeRoutes() {
	a.Router.Use(render.SetContentType(render.ContentTypeJSON))
	a.Router.Use(cors.Handler(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	a.Router.Route("/", func(r chi.Router) {
		r.Get("/", getWelcomeMessage)
	})

	a.Router.Route("/login", func(r chi.Router) {
		r.Post("/", loginUser)
	})

	// route: /users
	a.Router.Route("/users", func(r chi.Router) {
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
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (a *App) getSnippet(w http.ResponseWriter, r *http.Request) {
	val := chi.URLParam(r, "snippetID")

	_, err := strconv.Atoi(val)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	p := Snippet{ID: val}
	if err := p.getSnippet(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Snippet not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, p)
}

// ---------------
// ---------------
// ---------------
// MOCKUP
// ---------------
// ---------------
// ---------------

var users = []User{}

// Creates Mockdata which in future we will get from the database
func createMockData() {
	snippet1 := Snippet{`snippet1`, `java`, `stdout`, `System.out.println("Hello World");`}
	snippet2 := Snippet{`snippet2`, `go`, `structure`, `type User struct {\n    name string\n    lastname string\n    age int\n} `}
	snippet3 := Snippet{`snippet3`, `python`, `stdout`, `print("Hello World")`}
	snippet4 := Snippet{`snippet4`, `javascript`, `objects`, `var car = {\n    type:"Fiat",\n    model:"500",\n     color:"white"\n};`}

	user1Snippets := []Snippet{snippet1, snippet2}
	user2Snippets := []Snippet{snippet3, snippet4}

	users = append(users, User{"0", "Andreas", "asdf", user1Snippets})
	users = append(users, User{"1", "Markus", "1234", user2Snippets})
}

func getUser(userID string) User {
	for _, user := range users {
		if user.ID == userID {
			return user
		}
	}
	return User{}
}

func getUserSnippet(user User, snippetID string) Snippet {
	for _, snippet := range user.Snippets {
		if snippet.ID == snippetID {
			return snippet
		}
	}
	return Snippet{}
}

// GET - Request
// Outputs Welcome Msg
func getWelcomeMessage(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(WelcomeMessage)
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
	userID := chi.URLParam(r, "userID")

	user := getUser(userID)

	userJSON, _ := json.Marshal(user)
	json.NewEncoder(w).Encode(string(userJSON))

}

// GET - Request
// Outputs all saved Snippets of a user
func getUserSnippets(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")

	user := getUser(userID)

	snippetsJSON, _ := json.Marshal(user.Snippets)
	json.NewEncoder(w).Encode(string(snippetsJSON))
}

// GET - Request
// Outputs the requested snipped of a user
func getUserSnippetDetails(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	snippetID := chi.URLParam(r, "snippetID")

	user := getUser(userID)
	snippet := getUserSnippet(user, snippetID)

	snippetJSON, _ := json.Marshal(snippet)
	fmt.Println(snippet)
	json.NewEncoder(w).Encode(string(snippetJSON))
}

// PUT - Request
// Updates a snippet of a user
func putUserSnippetDetails(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	snippetID := chi.URLParam(r, "snippetID")
	snippetString := chi.URLParam(r, "snippet")
	snippet := Snippet{}
	err := json.Unmarshal([]byte(snippetString), &snippet)

	fmt.Println(r.Context().Value("snippet"))
	fmt.Println(userID + " - " + snippetID + " - " + snippetString)

	if err != nil {
		json.NewEncoder(w).Encode(ErrorMessage)
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
		json.NewEncoder(w).Encode(ErrorMessageMalformedJSON)
	}

	fmt.Println("########")
	fmt.Println(dbUserName)
	fmt.Println(dbPassword)
}

func getUserCredentialsDB() (string, string) {
	// Quick and Dirty
	// Should be refactored
	db, err := sql.Open("postgres",
		"postgresql://restapi@localhost:5432/snippet?sslmode=disable")
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
