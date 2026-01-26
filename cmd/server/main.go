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
	todosHandler := handler.NewTodosServiceHandler(service.NewTodosService(repo))

	// Get Connect handler
	path, h := gen.NewTodosServiceHandler(todosHandler, connect.WithInterceptors(validate.NewInterceptor()))

	// Register HTTP handlers
	mux := http.NewServeMux()
	mux.Handle(path, h)

	// Specialized CORS Configuration
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://todos.localhost", "http://localhost:3000"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{
			"Connect-Protocol-Version",
			"Content-Type",
			"Connect-Timeout-Ms",
			"Authorization",
		},
		ExposedHeaders: []string{
			"Grpc-Status",
			"Grpc-Message",
			"Grpc-Status-Details-Bin",
		},
		// Prevents the 404 by returning 200 to OPTIONS requests
		OptionsPassthrough: false,
	})

	// Wrap with CORS
	handlerWithCORS := c.Handler(mux)

	// IMPORTANT: Ensure h2c is the OUTSIDE wrapper
	h2Server := &http2.Server{}
	mainHandler := h2c.NewHandler(handlerWithCORS, h2Server)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mainHandler,
	}

	log.Println("Backend listening on :8080 with HTTP/2 (h2c) and specialized CORS support")
	log.Fatal(server.ListenAndServe())
}
