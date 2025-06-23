package api

import "net/http"

func (server *Server) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := server.authenticateToken(r)
		if err != nil {
			_ = server.invalidCredentials(w)
			return
		}
		next.ServeHTTP(w, r)
	})
}
