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
	dbConn := db.InitializeDB()
	repo, err := repository.NewRepository(dbConn)
	if err != nil {
		log.Fatal("Failed to initialize repository:", err)
	}

	// Create service handler
	todosHandler := handler.NewTodosServiceHandler(service.NewTodosService(repo))

	// Get Connect handler
	path, h := gen.NewTodosServiceHandler(todosHandler, connect.WithInterceptors(validate.NewInterceptor()))

	// new mux server
	mux := http.NewServeMux()

	// Register probe health check
	mux.HandleFunc("/health", handler.HealthHandler)

	// Register probe readiness check
	mux.HandleFunc("/ready", handler.ReadyHandler(dbConn))

	// Register connectRPC handler
	mux.Handle(path, h)

	// CORS Configuration helps UI from diff pod access the backend
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
		OptionsPassthrough: false,
	})

	// CORS handler wraper
	handlerWithCORS := c.Handler(mux)

	// h2c wraper for using http2 server helps to serve http2 requests
	h2Server := &http2.Server{}
	mainHandler := h2c.NewHandler(handlerWithCORS, h2Server)

	// Create serve
	server := &http.Server{
		Addr:    ":8080",
		Handler: mainHandler,
	}

	log.Println("Backend listening on :8080 with HTTP/2 (h2c) with CORS support")

	// Listen & serve
	log.Fatal(server.ListenAndServe())
}
