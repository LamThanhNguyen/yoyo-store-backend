package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

// GetItemByID gets one item by id and returns as JSON
func (server *Server) GetItemByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	itemID, _ := strconv.Atoi(id)

	item, err := server.DB.GetItem(itemID)
	if err != nil {
		log.Error().Err(err)
		return
	}

	out, err := json.MarshalIndent(item, "", "   ")
	if err != nil {
		log.Error().Err(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}
