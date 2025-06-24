package api

import (
	"encoding/json"
	"errors"
	"io"
	"maps"
	"net/http"
	"os"

	"github.com/LamThanhNguyen/yoyo-store-backend/server_invoice/util"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/rs/zerolog/log"
)

type Server struct {
	config util.Config
	router http.Handler
}

func NewServer(
	config util.Config,
) (*Server, error) {

	return &Server{
		config: config,
	}, nil
}

func (server *Server) SetupRouter() {
	mux := chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   server.config.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	mux.Get("/invoice/v1/health", server.handleHealthCheck)
	// mux.Post("/invoice/create-and-send", server.CreateAndSendInvoice)

	server.router = mux
}

func (server *Server) Router() http.Handler {
	return server.router
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
	if _, err := w.Write(out); err != nil {
		return err
	}

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
	if _, err := w.Write(out); err != nil {
		return err
	}
	return nil
}

func (server *Server) CreateDirIfNotExist(path string) error {
	const mode = 0755
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, mode)
		if err != nil {
			log.Error().Err(err)
			return err
		}
	}
	return nil
}

func (server *Server) handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	var data = map[string]string{"status": "ok"}
	_ = server.writeJSON(w, http.StatusOK, data)
}
