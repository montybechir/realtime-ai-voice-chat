package routes

import (
	"interviews-ai/internal/ai"
	"interviews-ai/internal/common/handlers"
	"interviews-ai/internal/common/middleware"
	"interviews-ai/internal/services"
	"net/http"

	"interviews-ai/internal/common/config"
)

func NewServer(cfg config.Config) (http.Handler, error) {

	userService := services.NewUserService(cfg.DatabaseURL)

	userHandler := handlers.NewUserHandler(userService)

	hub := ai.NewHub()
	go hub.Run()

	wsHandler := handlers.NewWsHandler(hub, &cfg)

	logger := middleware.NewLogger(cfg.LogLevel)

	mux := http.NewServeMux()
	// User routes with method handling
	mux.HandleFunc("/api/v1/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			userHandler.ListUsers().ServeHTTP(w, r)
		case http.MethodPost:
			userHandler.CreateUser().ServeHTTP(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Single user routes
	mux.HandleFunc("/api/v1/users/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			userHandler.GetUser().ServeHTTP(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.Handle("/api/v1/ws", logger(wsHandler))

	handler := middleware.Chain(
		mux,
		logger,
	)

	return handler, nil

}
