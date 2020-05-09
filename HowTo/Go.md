# Go Steps

#### Needed software
https://golang.org/

#### golang-chi
https://github.com/go-chi/chi

Install:
> go get -u github.com/go-chi/chi

Usage:
```golang
import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})
	http.ListenAndServe(":3000", r)
}
```