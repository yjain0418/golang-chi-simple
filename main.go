package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()

	r.NotFound(handleNotFound)
	r.MethodNotAllowed(handleMethodNotAllowed)

	r.Group(func(r chi.Router) {
		r.Use(middleware.Logger)

		r.Get("/", handleHelloWorldRequest)
		r.Get("/test", handleTest)

		r.Get("/books/{bookId}", handleGetBookById)

		// GET /authors -> get list of authors
		// POST /authors -> insert new author
		// GET /authors/{id} -> get author by id
		// DELETE /authors/{id} -> delete author by id
		r.Route("/authors", func(r chi.Router) {
			r.Get("/", handleGetAuthors)
			r.Post("/", handleCreateAuthor)
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", handleGetAuthorById)
				r.Delete("/", handleDeleteAuthorById)
			})
		})
	})

	r.Group(func(r chi.Router) {
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("OK"))
		})
	})

	http.ListenAndServe(":3000", r)
}

func handleHelloWorldRequest(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!"))
}

func handleTest(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is a test."))
}

func handleGetBookById(w http.ResponseWriter, r *http.Request) {
	bookId := chi.URLParam(r, "bookId")
	w.Write([]byte(bookId))
}

func handleNotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	w.Write([]byte("This route does not exist."))
}

func handleMethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(405)
	w.Write([]byte("This method is not allowed."))
}

func handleGetAuthors(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Get list of authors."))
}

func handleCreateAuthor(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Create a new author."))
}

func handleGetAuthorById(w http.ResponseWriter, r *http.Request) {
	authorId := chi.URLParam(r, "id")
	data := map[string]string{"message": "Get author by ID: " + authorId}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func handleDeleteAuthorById(w http.ResponseWriter, r *http.Request) {
	authorId := chi.URLParam(r, "id")
	w.Write([]byte("Delete author by ID: " + authorId))
}
