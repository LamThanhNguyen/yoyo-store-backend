package handler

import (
	"fmt"
	"net/http"

	"github.com/LamThanhNguyen/yoyo-store-backend/internal/encryption"
	"github.com/LamThanhNguyen/yoyo-store-backend/internal/urlsigner"
	"github.com/rs/zerolog/log"
)

// AllUsers shows the all users page
func (server *Server) AllUsers(w http.ResponseWriter, r *http.Request) {
	if err := server.renderTemplate(w, r, "all-users", &templateData{}); err != nil {
		log.Error().Err(err)
	}
}

// OneUser shows one admin user for add/edit/delete
func (server *Server) OneUser(w http.ResponseWriter, r *http.Request) {
	if err := server.renderTemplate(w, r, "one-user", &templateData{}); err != nil {
		log.Error().Err(err)
	}
}

// LoginPage displays the login page
func (server *Server) LoginPage(w http.ResponseWriter, r *http.Request) {
	if err := server.renderTemplate(w, r, "login", &templateData{}); err != nil {
		log.Error().Err(err)
	}
}

// PostLoginPage handles the posted login form
func (server *Server) PostLoginPage(w http.ResponseWriter, r *http.Request) {
	server.Session.RenewToken(r.Context())

	err := r.ParseForm()
	if err != nil {
		log.Error().Err(err)
		return
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	id, err := server.DB.Authenticate(email, password)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	server.Session.Put(r.Context(), "userID", id)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Logout logs the user out
func (server *Server) Logout(w http.ResponseWriter, r *http.Request) {
	server.Session.Destroy(r.Context())
	server.Session.RenewToken(r.Context())

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// ForgotPassword shows the forgot password page
func (server *Server) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	if err := server.renderTemplate(w, r, "forgot-password", &templateData{}); err != nil {
		log.Error().Err(err)
	}
}

// ShowResetPassword shows the reset password page (and validates url integrity)
func (server *Server) ShowResetPassword(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	theURL := r.RequestURI
	testURL := fmt.Sprintf("%s%s", server.config.FrontendAddr, theURL)

	signer := urlsigner.Signer{
		Secret: []byte(server.config.TokenSymmetricKey),
	}

	valid := signer.VerifyToken(testURL)

	if !valid {
		log.Error().Msg("Invalid url - tampering detected")
		return
	}

	// make sure not expired
	expired := signer.Expired(testURL, 60)
	if expired {
		log.Error().Msg("Link expired")
		return
	}

	encyrptor := encryption.Encryption{
		Key: []byte(server.config.TokenSymmetricKey),
	}

	encryptedEmail, err := encyrptor.Encrypt(email)
	if err != nil {
		log.Error().Msg("Encryption failed")
		return
	}

	data := make(map[string]interface{})
	data["email"] = encryptedEmail

	if err := server.renderTemplate(w, r, "reset-password", &templateData{
		Data: data,
	}); err != nil {
		log.Error().Err(err)
	}
}
