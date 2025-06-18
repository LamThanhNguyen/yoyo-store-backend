package api

import (
	"encoding/json"
	"net/http"

	"github.com/LamThanhNguyen/yoyo-store-backend/internal/models"
	"github.com/LamThanhNguyen/yoyo-store-backend/server_main/util"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

type Server struct {
	config util.RuntimeConfig
	DB     models.DBModel
	router http.Handler
}

func NewServer(
	config util.RuntimeConfig,
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

// writeJSON writes aribtrary data out as JSON
func (server *Server) writeJSON(
	w http.ResponseWriter,
	status int,
	data interface{},
	headers ...http.Header,
) error {
	out, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for k, v := range headers[0] {
			w.Header()[k] = v
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(out)

	return nil
}

func (server *Server) failedValidation(
	w http.ResponseWriter,
	_ *http.Request,
	errors map[string]string,
) {
	var payload struct {
		Error   bool              `json:"error"`
		Message string            `json:"message"`
		Errors  map[string]string `json:"errors"`
	}

	payload.Error = true
	payload.Message = "failed validation"
	payload.Errors = errors
	server.writeJSON(w, http.StatusUnprocessableEntity, payload)
}
