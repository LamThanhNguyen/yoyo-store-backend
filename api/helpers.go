package main

import (
	"encoding/json"
	"net/http"
)

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
