package handler

import (
	"net/http"

	"github.com/rs/zerolog/log"
)

// Home displays the home page
func (server *Server) Home(w http.ResponseWriter, r *http.Request) {
	if err := server.renderTemplate(w, r, "home", &templateData{}); err != nil {
		log.Error().Err(err).Msg("Home")
	}
}

// VirtualTerminal displays the virtual terminal page
func (server *Server) VirtualTerminal(w http.ResponseWriter, r *http.Request) {
	if err := server.renderTemplate(w, r, "terminal", &templateData{}); err != nil {
		log.Error().Err(err).Msg("VirtualTerminal")
	}
}
