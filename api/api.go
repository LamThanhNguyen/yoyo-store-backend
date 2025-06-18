package main

import (
	"net/http"

	"github.com/LamThanhNguyen/yoyo-store-backend/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

type Server struct {
	config RuntimeConfig
	DB     models.DBModel
	router http.Handler
}

func NewServer(
	config RuntimeConfig,
	db models.DBModel,
) (*Server, error) {

	return &Server{
		config: config,
		DB:     db,
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

	mux.Post("/api/v1/payment-intent", server.GetPaymentIntent)
	mux.Get("/api/v1/items/{id}", server.GetItemByID)
	mux.Post("/api/create-customer-and-subscribe-to-plan", server.CreateCustomerAndSubscribeToPlan)

	server.router = mux
}

func (server *Server) Router() http.Handler {
	return server.router
}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message,omitempty"`
	Content string `json:"content,omitempty"`
	ID      int    `json:"id,omitempty"`
}
