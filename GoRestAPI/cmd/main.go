package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

func getAllContact(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode([]string{"message 1", "message 2", "message 3"})
}

func registerRoutes() http.Handler {
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Get("/", getAllContact)
	})
	return r
}

func main() {
	r := registerRoutes()
	log.Fatal(http.ListenAndServe(":8080", r))
}
