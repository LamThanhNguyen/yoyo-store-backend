package main

import (
	"net/http"

	db "github.com/LamThanhNguyen/yoyo-store-backend/db/sqlc"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

type Server struct {
	config RuntimeConfig
	store  db.Store
	router http.Handler
}

func NewServer(
	config RuntimeConfig,
	store db.Store,
) (*Server, error) {

	return &Server{
		config: config,
		store:  store,
	}, nil
}

func (server *Server) SetupRouter() {
	mux := chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	server.router = mux
}

func (server *Server) Router() http.Handler {
	return server.router
}
