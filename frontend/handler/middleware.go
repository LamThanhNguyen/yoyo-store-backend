package handler

import (
	"net/http"
)

// SessionLoad peforms the load and save of a session, per request
func (server *Server) SessionLoad(next http.Handler) http.Handler {
	return server.Session.LoadAndSave(next)
}

// Auth checks for user authentication status by checking for the key
// userID in the session
func (server *Server) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !server.Session.Exists(r.Context(), "userID") {
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}
		next.ServeHTTP(w, r)
	})
}
