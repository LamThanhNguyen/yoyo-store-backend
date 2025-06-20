package api

import (
	"encoding/json"
	"errors"
	"io"
	"maps"
	"net/http"
	"strings"

	"github.com/LamThanhNguyen/yoyo-store-backend/internal/models"
	"github.com/LamThanhNguyen/yoyo-store-backend/server_main/util"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"golang.org/x/crypto/bcrypt"
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
		AllowedOrigins:   server.config.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	mux.Post("/api/v1/payment-intent", server.GetPaymentIntent)
	mux.Get("/api/v1/items/{id}", server.GetItemByID)
	mux.Post("/api/create-customer-and-subscribe-to-plan", server.CreateCustomerAndSubscribeToPlan)

	mux.Post("/api/authenticate", server.CreateAuthToken)
	mux.Post("/api/is-authenticated", server.CheckAuthentication)
	mux.Post("/api/forgot-password", server.SendPasswordResetEmail)
	mux.Post("/api/reset-password", server.ResetPassword)

	mux.Route("/api/admin", func(mux chi.Router) {
		mux.Use(server.Auth)

		mux.Post("/virtual-terminal-succeeded", server.VirtualTerminalPaymentSucceeded)
		mux.Get("/all-sales", server.AllSales)
		mux.Get("/all-subscriptions", server.AllSubscriptions)

		mux.Get("/get-sale/{id}", server.GetSale)

		mux.Post("/refund", server.RefundCharge)
		mux.Post("/cancel-subscription", server.CancelSubscription)

		mux.Get("/all-users", server.AllUsers)
		mux.Get("/all-users/{id}", server.OneUser)
		mux.Patch("/all-users/edit/{id}", server.EditUser)

		mux.Delete("/all-users/delete/{id}", server.DeleteUser)
	})

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
		maps.Copy(w.Header(), headers[0])
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(out)

	return nil
}

// readJSON reads json from request body into data. We only accept a single json value in the body
func (server *Server) readJSON(
	w http.ResponseWriter,
	r *http.Request,
	data interface{},
) error {
	maxBytes := 1048576 // max one megabyte in request body
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		return err
	}

	// we only allow one entry in the json file
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only have a single JSON value")
	}

	return nil
}

// badRequest sends a JSON response with status http.StatusBadRequest, describing the error
func (server *Server) badRequest(
	w http.ResponseWriter,
	_ *http.Request,
	err error,
) error {
	var payload struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}

	payload.Error = true
	payload.Message = err.Error()

	out, err := json.MarshalIndent(payload, "", "\t")
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
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

func (server *Server) invalidCredentials(w http.ResponseWriter) error {
	var payload struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}

	payload.Error = true
	payload.Message = "invalid authentication credentials"

	err := server.writeJSON(w, http.StatusUnauthorized, payload)
	if err != nil {
		return err
	}
	return nil
}

func (server *Server) passwordMatches(hash, password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

// authenticateToken checks an auth token for validity
func (server *Server) authenticateToken(r *http.Request) (*models.User, error) {
	authorizationHeader := r.Header.Get("Authorization")
	if authorizationHeader == "" {
		return nil, errors.New("no authorization header received")
	}

	headerParts := strings.Split(authorizationHeader, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return nil, errors.New("no authorization header received")
	}

	token := headerParts[1]
	if len(token) != 26 {
		return nil, errors.New("authentication token wrong size")
	}

	// get the user from the tokens table
	user, err := server.DB.GetUserForToken(token)
	if err != nil {
		return nil, errors.New("no matching user found")
	}

	return user, nil
}
