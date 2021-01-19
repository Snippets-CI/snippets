package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
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
func (a *App) Initialize(user, password, dbname string, dbhost string, middlewareEnabled bool) {
	fmt.Println("[*] Waiting for db to settle...")
	time.Sleep(3 * time.Second)

	fmt.Println("[*] Initializing...")

	connectionString := fmt.Sprintf("%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbhost, 5432, user, password, dbname)

	fmt.Println("Connection string: " + connectionString)

	var err error
	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		fmt.Println("[!] Error while opening sql connection")
		log.Fatal(err)
	}
	fmt.Println("[*] Finished connecting")

	fmt.Println("[*] Setting up chi router")
	a.Router = chi.NewRouter()
	fmt.Println("[*] Chi router setup finished")

	if middlewareEnabled {
		fmt.Println("[*] Middleware enabled")
		a.Router.Use(middleware.RequestID)
		a.Router.Use(middleware.RealIP)
		a.Router.Use(middleware.Logger)
		a.Router.Use(middleware.Recoverer)
		fmt.Println("[*] Middleware setup finished")
	}

	a.Router.Use(middleware.Timeout(60 * time.Second))

	a.initializeRoutes()

	fmt.Println("[*] Ensuring extension exists")
	// ensureExtensionExists(a.DB)

	fmt.Println("[*] Ensuring db exists")
	ensureTablesExist(a.DB)

	fmt.Println("[*] Finished Initializing")
}

// Run http listen and serve
func (a *App) Run(addr string) {
	fmt.Println("[*] ListenAndServe on " + addr)
	log.Fatal(http.ListenAndServe(addr, a.Router))
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
		r.Post("/", a.login)
	})

	// route: /users
	a.Router.Route("/user", func(r chi.Router) {
		r.Post("/", a.createUser)
		r.Put("/", a.updateUser)

		// route: /users/{userID}
		r.Route("/{userID}", func(r chi.Router) {

			// route: /users/{userID}/snippets
			r.Route("/snippets", func(r chi.Router) {
				r.Get("/", a.getSnippets)
				r.Post("/", a.createSnippet)

				// route: /users/{userID}/snippets/{snippetID}
				r.Route("/{snippetID}", func(r chi.Router) {
					r.Get("/", a.getSnippet)
					r.Put("/", a.updateSnippet)
					r.Delete("/", a.deleteSnippet)
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

// ---------------
// ---------------
// ---------------
// Rest API functions
// ---------------
// ---------------
// ---------------

func getWelcomeMessage(w http.ResponseWriter, r *http.Request) {
	// swagger:operation GET /
	//
	// Returns a welcome
	// ---
	// produces:
	// - application/json
	// responses:
	//   '200':
	//     description: welcome message :)

	respondWithJSON(w, http.StatusOK, WelcomeMessage)
}

func (a *App) login(w http.ResponseWriter, r *http.Request) {
	// swagger:operation POST /login login
	//
	// Returns a user for provided email and jwt
	// ---
	// produces:
	// - application/json
	// parameters:
	// - mail: userMail
	//   in: body
	//   description: user id for user selection
	//   required: true
	//   type: string
	// - authToken: jwt
	//   in: header
	//   description: web token from user authentication
	//   required: true
	//   type: string
	// responses:
	//   '200':
	//     schema:
	//     "$ref": "#/responses/User"
	//     type: json
	//   '404':
	//      description: no user with that mail found
	//   '401':
	//      description: authorization failed

	user := User{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := user.getUser(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, err.Error())
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	// Create JWT and send it as a response
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Name,
	})

	// !!!!!!!!!!
	// TODO: create secret key and save in env
	// !!!!!!!!!!!!
	// Sign the token with the secrete key
	tokenString, _ := token.SignedString([]byte("secret"))

	respondWithJSON(w, http.StatusOK, tokenString)
}

func (a *App) createUser(w http.ResponseWriter, r *http.Request) {
	// swagger:operation POST /user
	//
	// Create and returns user
	// ---
	// produces:
	// - application/json
	// parameters:
	// - mail: userMail
	//   in: body
	//   description: new mail for user
	//   required: true
	//   type: string
	// - name: userName
	//   in: body
	//   description: new name for user
	//   required: true
	//   type: string
	// - password: userPassword
	//   in: body
	//   description: new password for user
	//   required: true
	//   type: string
	// responses:
	//   '200':
	//      description: user created

	user := User{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := user.createUser(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Create JWT and send it as a response
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Name,
	})

	tokenString, _ := token.SignedString([]byte("secret"))

	respondWithJSON(w, http.StatusCreated, tokenString)
}

func (a *App) updateUser(w http.ResponseWriter, r *http.Request) {
	// swagger:operation PUT /user
	//
	// Updates and returns user
	// ---
	// produces:
	// - application/json
	// parameters:
	// - userID: userId
	//   in: path
	//   description: user id for user selection
	//   required: true
	//   type: string
	// - mail: userMail
	//   in: body
	//   description: new mail for user
	//   required: true
	//   type: string
	// - name: userName
	//   in: body
	//   description: new name for user
	//   required: true
	//   type: string
	// - authToken: jwt
	//   in: body
	//   description: authorization token
	//   required: true
	//   type: string
	// responses:
	//   '200':
	//      description: user updated
	//   '404':
	//      description: no user with that id found
	//   '401':
	//      description: authorization failed

	user := User{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	tokenString := r.Header.Get("Authorization")
	userID, err := validateJwtToken(tokenString)

	if err != nil {
		respondWithError(w, http.StatusForbidden, err.Error())

	} else if userID != user.ID {
		respondWithError(w, http.StatusUnauthorized, err.Error())
	}

	if err := user.updateUser(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, user)
}

func (a *App) getSnippet(w http.ResponseWriter, r *http.Request) {
	// swagger:operation GET /user/{userID}/snippets/{snippetID} getSnippet
	//
	// Returns a snippet for a given user and snippet id
	// ---
	// produces:
	// - application/json
	// parameters:
	// - userID: userId
	//   in: path
	//   description: user id for snippet selection
	//   required: true
	//   type: string
	// - snippetID: snippetId
	//   in: path
	//   description: snippet which should be retrieved
	//   required: true
	//   type: string
	// - authToken: jwt
	//   in: header
	//   description: authorization token
	//   required: true
	//   type: string
	// responses:
	//   '200':
	//      description: a snippet to be returned
	//      schema:
	//      $ref: '#/responses/Snippet'
	//   '404':
	//      description: no snippet with that id found
	//   '401':
	//      description: authorization failed

	id := chi.URLParam(r, "userID")
	val := chi.URLParam(r, "snippetID")

	tokenString := r.Header.Get("Authorization")
	userID, err := validateJwtToken(tokenString)

	if err != nil {
		respondWithError(w, http.StatusForbidden, err.Error())

	} else if userID != id {
		respondWithError(w, http.StatusUnauthorized, err.Error())
	}

	p := Snippet{ID: val, Owner: userID}
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

func (a *App) createSnippet(w http.ResponseWriter, r *http.Request) {
	// swagger:operation POST /user/{userID}/snippets/ createSnippet
	//
	// Create a snippet with a given id
	// ---
	// produces:
	// - application/json
	// parameters:
	// - userID: userId
	//   in: path
	//   description: user id for snippet selection
	//   required: true
	//   type: string
	// - authToken: jwt
	//   in: header
	//   description: authorization token
	//   required: true
	//   type: string
	// - id: snippetID
	//   in: body
	//   description: snippet id
	//   type: string
	// - name: name
	//   in: body
	//   description: snippet name
	//   type: string
	// - lang: snippetLang
	//   in: body
	//   description: snippet language
	//   type: string
	// - about: snippetAbout
	//   in: body
	//   description: snippet description
	//   type: string
	// - code: snippetCode
	//   in: body
	//   description: snippet code text
	//   type: string
	// responses:
	//   '200':
	//      description: an updated snippet to be returned
	//      schema:
	//      $ref: '#/responses/Snippet'
	//   '404':
	//      description: no snippet with that id found
	//   '401':
	//      description: authorization failed

	id := chi.URLParam(r, "userID")

	tokenString := r.Header.Get("Authorization")
	userID, err := validateJwtToken(tokenString)

	if err != nil {
		respondWithError(w, http.StatusForbidden, err.Error())

	} else if userID != id {
		respondWithError(w, http.StatusUnauthorized, err.Error())
	}

	snippet := Snippet{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&snippet); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload "+err.Error())
		return
	}
	defer r.Body.Close()

	if err := snippet.createSnippet(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Snippet not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusCreated, snippet)
}

func (a *App) updateSnippet(w http.ResponseWriter, r *http.Request) {
	// swagger:operation PUT /user/{userID}/snippets/{snippetID} updateSnippet
	//
	// Update a snippet with a given id
	// ---
	// produces:
	// - application/json
	// parameters:
	// - userID: userId
	//   in: path
	//   description: user id for snippet selection
	//   required: true
	//   type: string
	// - id: snippetID
	//   in: path
	//   description: snippet id
	//   type: string
	// - authToken: jwt
	//   in: header
	//   description: authorization token
	//   required: true
	//   type: string
	// - name: name
	//   in: body
	//   description: snippet name
	//   type: string
	// - lang: snippetLang
	//   in: body
	//   description: snippet language
	//   type: string
	// - about: snippetAbout
	//   in: body
	//   description: snippet description
	//   type: string
	// - code: snippetCode
	//   in: body
	//   description: snippet code text
	//   type: string
	// responses:
	//   '200':
	//      description: an updated snippet to be returned
	//      schema:
	//      $ref: '#/responses/Snippet'
	//   '404':
	//      description: no snippet with that id found
	//   '401':
	//      description: authorization failed

	id := chi.URLParam(r, "userID")
	snippetID := chi.URLParam(r, "snippetID")

	tokenString := r.Header.Get("Authorization")
	userID, err := validateJwtToken(tokenString)

	if err != nil {
		respondWithError(w, http.StatusForbidden, err.Error())

	} else if userID != id {
		respondWithError(w, http.StatusUnauthorized, err.Error())
	}

	snippet := Snippet{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&snippet); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	snippet.ID = snippetID

	if err := snippet.updateSnippet(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Snippet not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, snippet)
}

func (a *App) deleteSnippet(w http.ResponseWriter, r *http.Request) {
	// swagger:operation DELETE /user/{userID}/snippets/{snippetID} deleteSnippet
	//
	// Deletes a snippet for a given id
	// ---
	// produces:
	// - application/json
	// parameters:
	// - userID: userId
	//   in: path
	//   description: user id for snippet selection
	//   required: true
	//   type: string
	// - snippetID: snippetId
	//   in: path
	//   description: snippet which should be deleted
	//   required: true
	//   type: string
	// - authToken: jwt
	//   in: header
	//   description: authorization token
	//   required: true
	//   type: string
	// responses:
	//   '200':
	//      description: result success
	//   '404':
	//      description: no snippet with that id found
	//   '401':
	//      description: authorization failed

	id := chi.URLParam(r, "userID")
	val := chi.URLParam(r, "snippetID")

	tokenString := r.Header.Get("Authorization")
	userID, err := validateJwtToken(tokenString)

	if err != nil {
		respondWithError(w, http.StatusForbidden, err.Error())

	} else if userID != id {
		respondWithError(w, http.StatusUnauthorized, err.Error())
	}

	p := Snippet{ID: val}
	if err := p.deleteSnippet(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) getSnippets(w http.ResponseWriter, r *http.Request) {
	// swagger:operation GET /user/{userID}/snippets/ getSnippets
	//
	// Returns all snippets for a user
	// ---
	// produces:
	// - application/json
	// parameters:
	// - userID: userId
	//   in: path
	//   description: user id for snippet selection
	//   required: true
	//   type: string
	// - authToken: jwt
	//   in: header
	//   description: authorization token
	//   required: true
	//   type: string
	// responses:
	//   '200':
	//     description: HTTP status code 200 and Snippets
	//     schema:
	//       type: array
	//       items:
	//         "$ref": "#/responses/Snippet"
	//   '404':
	//     description: no user(user) with that id found
	//   '401':
	//     description: authorization failed

	id := chi.URLParam(r, "userID")

	user := User{ID: id}

	tokenString := r.Header.Get("Authorization")
	userID, err := validateJwtToken(tokenString)

	if err != nil {
		respondWithError(w, http.StatusForbidden, err.Error())

	} else if userID != id {
		respondWithError(w, http.StatusUnauthorized, err.Error())
	}

	if err := user.getSnippets(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, user.Snippets)
}

// Validates the token and returns the userid or an error in case of an error
func validateJwtToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte("secret"), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["user_id"].(string), err

	}

	return "", err
}
