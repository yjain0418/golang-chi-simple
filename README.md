# golang-chi-simple

A minimal HTTP server built with [Chi](https://github.com/go-chi/chi), demonstrating how the router's core features fit together.

## Running the server

Requires Go 1.23 or newer.

```bash
go run .
```

The server listens on `http://localhost:3000`.

```bash
curl http://localhost:3000/
```

## Chi features implemented

### Router creation

`chi.NewRouter()` creates the root `*chi.Mux` used to register routes and attach middleware.

```go
r := chi.NewRouter()
```

The returned router is passed directly to `http.ListenAndServe`, so Chi handles each incoming `*http.Request` and writes to the matching `http.ResponseWriter`.

### Not-found handler

`r.NotFound(handleNotFound)` registers the handler used when no route matches the request URL.

- Returns HTTP status `404`
- Response body: `This route does not exist.`

```bash
curl -i http://localhost:3000/missing
```

### Method-not-allowed handler

`r.MethodNotAllowed(handleMethodNotAllowed)` registers the handler used when the URL exists but the HTTP method is not supported by that route.

- Returns HTTP status `405`
- Response body: `This method is not allowed.`

For example, `/test` supports `GET` but not `POST`:

```bash
curl -i -X POST http://localhost:3000/test
```

### Groups

`r.Group` creates a group of routes that share the same configuration. A group does not create a URL prefix by itself.

The first group contains the application routes and enables request logging. The second group contains the health route and has no middleware attached.

```go
r.Group(func(r chi.Router) {
	r.Use(middleware.Logger)
	// Logged application routes
})
```

### Logger middleware

`r.Use(middleware.Logger)` adds Chi's logger middleware to the first route group. The middleware prints request information, such as method and path, to standard output for every request handled by that group.

The `/health` route is outside this group, so it is not logged by `middleware.Logger`.

### HTTP method handlers

Chi provides typed route registration methods for common HTTP verbs:

| Method | Route | Handler | Response |
|---|---|---|---|
| `GET` | `/` | `handleHelloWorldRequest` | `Hello, World!` |
| `GET` | `/test` | `handleTest` | `This is a test.` |
| `GET` | `/health` | inline handler | `OK` |
| `GET` | `/authors/` | `handleGetAuthors` | `Get list of authors.` |
| `POST` | `/authors/` | `handleCreateAuthor` | `Create a new author.` |
| `GET` | `/authors/{id}/` | `handleGetAuthorById` | JSON object containing the author ID |
| `DELETE` | `/authors/{id}/` | `handleDeleteAuthorById` | `Delete author by ID: <id>` |

Example requests:

```bash
curl http://localhost:3000/
curl http://localhost:3000/test
curl http://localhost:3000/health
curl http://localhost:3000/authors/
curl -X POST http://localhost:3000/authors/
curl http://localhost:3000/authors/42/
curl -X DELETE http://localhost:3000/authors/42/
```

### URL parameters

`{bookId}` and `{id}` declare named path parameters:

```go
r.Get("/books/{bookId}", handleGetBookById)
r.Route("/{id}", func(r chi.Router) {
	r.Get("/", handleGetAuthorById)
})
```

Handlers read those values with `chi.URLParam`:

```go
bookId := chi.URLParam(r, "bookId")
authorId := chi.URLParam(r, "id")
```

Try a book or author ID:

```bash
curl http://localhost:3000/books/123
curl http://localhost:3000/authors/42/
```

`/books/123` returns `123`. The author handler returns JSON for the supplied ID, for example:

```json
{"message":"Get author by ID: 42"}
```

### Nested route registration

`r.Route` mounts a group of routes under a shared URL prefix. In this project, the nested `/authors` routes are registered as:

```text
/authors/
/authors/{id}/
```

This structure keeps collection operations (`GET` and `POST`) separate from item operations (`GET` and `DELETE`) without repeating `/authors` in every registration.

### Inline handlers

Routes can accept a function directly instead of a named handler:

```go
r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
})
```

This is useful for small, self-contained handlers such as `/health`.

## Route overview

```text
GET    /                 -> Hello, World!
GET    /test              -> Test message
GET    /books/{bookId}    -> Echo the book ID

GET    /authors/          -> Get the author list
POST   /authors/          -> Create an author
GET    /authors/{id}/     -> Get one author by ID
DELETE /authors/{id}/     -> Delete one author by ID

GET    /health            -> OK
```

## HTTP server

`http.ListenAndServe(":3000", r)` starts a standard Go HTTP server on port `3000` and uses the Chi router as its `http.Handler`.

The server currently does not handle startup errors explicitly, and it has no graceful shutdown or TLS configuration.