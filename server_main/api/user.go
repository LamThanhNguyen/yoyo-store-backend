package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/LamThanhNguyen/yoyo-store-backend/internal/encryption"
	"github.com/LamThanhNguyen/yoyo-store-backend/internal/models"
	"github.com/LamThanhNguyen/yoyo-store-backend/internal/urlsigner"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

// CreateAuthToken creates and sends an auth token, if user supplies valid information
func (server *Server) CreateAuthToken(w http.ResponseWriter, r *http.Request) {
	var userInput struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := server.readJSON(w, r, &userInput)
	if err != nil {
		server.badRequest(w, r, err)
		return
	}

	// get the user from the database by email; send error if invalid email
	user, err := server.DB.GetUserByEmail(userInput.Email)
	if err != nil {
		server.invalidCredentials(w)
		return
	}

	// validate the password; send error if invalid password
	validPassword, err := server.passwordMatches(user.Password, userInput.Password)
	if err != nil {
		server.invalidCredentials(w)
		return
	}

	if !validPassword {
		server.invalidCredentials(w)
		return
	}

	// generate the token
	token, err := models.GenerateToken(user.ID, 24*time.Hour, models.ScopeAuthentication)
	if err != nil {
		server.badRequest(w, r, err)
		return
	}

	// save to database
	err = server.DB.InsertToken(token, user)
	if err != nil {
		server.badRequest(w, r, err)
		return
	}

	// send response

	var payload struct {
		Error   bool          `json:"error"`
		Message string        `json:"message"`
		Token   *models.Token `json:"authentication_token"`
	}
	payload.Error = false
	payload.Message = fmt.Sprintf("token for %s created", userInput.Email)
	payload.Token = token

	_ = server.writeJSON(w, http.StatusOK, payload)
}

// CheckAuthentication checks auth status
func (server *Server) CheckAuthentication(w http.ResponseWriter, r *http.Request) {
	// validate the token, and get associated user
	user, err := server.authenticateToken(r)
	if err != nil {
		server.invalidCredentials(w)
		return
	}

	// valid user
	var payload struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}

	payload.Error = false
	payload.Message = fmt.Sprintf("authenticated user %s", user.Email)
	server.writeJSON(w, http.StatusOK, payload)
}

// SendPasswordResetEmail sends an email with a signed url to allow user to reset password
func (server *Server) SendPasswordResetEmail(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Email string `json:"email"`
	}

	err := server.readJSON(w, r, &payload)
	if err != nil {
		server.badRequest(w, r, err)
		return
	}

	_, err = server.DB.GetUserByEmail(payload.Email)

	if err != nil {
		var resp struct {
			Error   bool   `json:"error"`
			Message string `json:"message"`
		}
		resp.Error = true
		resp.Message = "No matching email found on our system"
		server.writeJSON(w, http.StatusAccepted, resp)
		return
	}

	link := fmt.Sprintf("%s/reset-password?email=%s", server.config.FrontendDomain, payload.Email)

	sign := urlsigner.Signer{
		Secret: []byte(server.config.TokenSymmetricKey),
	}

	signedLink := sign.GenerateTokenFromString(link)

	var data struct {
		Link string
	}

	data.Link = signedLink

	// send mail
	err = server.SendMail("info@yoyo.com", payload.Email, "Password Reset Request", "password-reset", data)
	if err != nil {
		log.Error().Err(err)
		server.badRequest(w, r, err)
		return
	}

	var resp struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}

	resp.Error = false

	server.writeJSON(w, http.StatusCreated, resp)
}

// ResetPassword resets a user's password in the database
func (server *Server) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := server.readJSON(w, r, &payload)
	if err != nil {
		server.badRequest(w, r, err)
		return
	}

	encyrptor := encryption.Encryption{
		Key: []byte(server.config.TokenSymmetricKey),
	}

	realEmail, err := encyrptor.Decrypt(payload.Email)
	if err != nil {
		server.badRequest(w, r, err)
		return
	}

	user, err := server.DB.GetUserByEmail(realEmail)
	if err != nil {
		server.badRequest(w, r, err)
		return
	}

	newHash, err := bcrypt.GenerateFromPassword([]byte(payload.Password), 12)
	if err != nil {
		server.badRequest(w, r, err)
		return
	}

	err = server.DB.UpdatePasswordForUser(user, string(newHash))
	if err != nil {
		server.badRequest(w, r, err)
		return
	}

	var resp struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}
	resp.Error = false
	resp.Message = "password changed"

	server.writeJSON(w, http.StatusCreated, resp)
}
