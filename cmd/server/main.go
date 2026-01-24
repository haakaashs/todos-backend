package main

import (
	"log"
	"net/http"

	"connectrpc.com/connect"
	"connectrpc.com/validate"
	gen "github.com/haakaashs/todos-backend/gen/protos/todos/v1/todosv1connect"
	handler "github.com/haakaashs/todos-backend/internal/api/v1/handler"
	"github.com/haakaashs/todos-backend/internal/db"
	"github.com/haakaashs/todos-backend/internal/repository"
	"github.com/haakaashs/todos-backend/internal/service"
	"github.com/rs/cors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	// Initialize DB
	repo, err := repository.NewRepository(db.InitializeDB())
	if err != nil {
		log.Fatal("Failed to initialize repository:", err)
	}

	// Create service handler
	handler := handler.NewTodosServiceHandler(service.NewTodosService(repo))

	// Get Connect handler
	path, h := gen.NewTodosServiceHandler(handler, connect.WithInterceptors(validate.NewInterceptor()))

	// Register HTTP handlers
	mux := http.NewServeMux()
	mux.Handle(path, h)

	// CORS (optional, for browser JSON requests)
	handlerWithCORS := cors.AllowAll().Handler(mux)

	// Wrap with HTTP/2 cleartext (h2c)
	h2Server := &http2.Server{}
	server := &http.Server{
		Addr:    ":8080",
		Handler: h2c.NewHandler(handlerWithCORS, h2Server),
	}

	log.Println("Backend listening on :8080 with HTTP/2 (h2c) support")
	log.Fatal(server.ListenAndServe())
}
