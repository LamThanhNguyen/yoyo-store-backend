package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/LamThanhNguyen/yoyo-store-backend/internal/encryption"
	"github.com/LamThanhNguyen/yoyo-store-backend/internal/models"
	"github.com/LamThanhNguyen/yoyo-store-backend/internal/urlsigner"
	"github.com/go-chi/chi/v5"
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
		_ = server.badRequest(w, r, err)
		return
	}

	// get the user from the database by email; send error if invalid email
	user, err := server.DB.GetUserByEmail(userInput.Email)
	if err != nil {
		_ = server.invalidCredentials(w)
		return
	}

	// validate the password; send error if invalid password
	validPassword, err := server.passwordMatches(user.Password, userInput.Password)
	if err != nil {
		_ = server.invalidCredentials(w)
		return
	}

	if !validPassword {
		_ = server.invalidCredentials(w)
		return
	}

	// generate the token
	token, err := models.GenerateToken(user.ID, 24*time.Hour, models.ScopeAuthentication)
	if err != nil {
		_ = server.badRequest(w, r, err)
		return
	}

	// save to database
	err = server.DB.InsertToken(token, user)
	if err != nil {
		_ = server.badRequest(w, r, err)
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
		_ = server.invalidCredentials(w)
		return
	}

	// valid user
	var payload struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}

	payload.Error = false
	payload.Message = fmt.Sprintf("authenticated user %s", user.Email)
	_ = server.writeJSON(w, http.StatusOK, payload)
}

// SendPasswordResetEmail sends an email with a signed url to allow user to reset password
func (server *Server) SendPasswordResetEmail(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Email string `json:"email"`
	}

	err := server.readJSON(w, r, &payload)
	if err != nil {
		_ = server.badRequest(w, r, err)
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
		_ = server.writeJSON(w, http.StatusAccepted, resp)
		return
	}

	link := fmt.Sprintf("%s/reset-password?email=%s", server.config.FrontendAddr, payload.Email)

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
		log.Error().Err(err).Msg("SendPasswordResetEmail")
		_ = server.badRequest(w, r, err)
		return
	}

	var resp struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}

	resp.Error = false

	_ = server.writeJSON(w, http.StatusCreated, resp)
}

// ResetPassword resets a user's password in the database
func (server *Server) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := server.readJSON(w, r, &payload)
	if err != nil {
		_ = server.badRequest(w, r, err)
		return
	}

	encyrptor := encryption.Encryption{
		Key: []byte(server.config.TokenSymmetricKey),
	}

	realEmail, err := encyrptor.Decrypt(payload.Email)
	if err != nil {
		_ = server.badRequest(w, r, err)
		return
	}

	user, err := server.DB.GetUserByEmail(realEmail)
	if err != nil {
		_ = server.badRequest(w, r, err)
		return
	}

	newHash, err := bcrypt.GenerateFromPassword([]byte(payload.Password), 12)
	if err != nil {
		_ = server.badRequest(w, r, err)
		return
	}

	err = server.DB.UpdatePasswordForUser(user, string(newHash))
	if err != nil {
		_ = server.badRequest(w, r, err)
		return
	}

	var resp struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}
	resp.Error = false
	resp.Message = "password changed"
	_ = server.writeJSON(w, http.StatusCreated, resp)
}

// AllUsers returns a JSON file listing all admin users
func (server *Server) AllUsers(w http.ResponseWriter, r *http.Request) {
	pageSize := 10   // default
	currentPage := 1 // default

	// Parse query params
	if val := r.URL.Query().Get("page_size"); val != "" {
		if ps, err := strconv.Atoi(val); err == nil && ps > 0 {
			pageSize = ps
		} else {
			_ = server.badRequest(w, r, errors.New("invalid page_size"))
			return
		}
	}
	// Parse query params
	if val := r.URL.Query().Get("page"); val != "" {
		if cp, err := strconv.Atoi(val); err == nil && cp > 0 {
			currentPage = cp
		} else {
			_ = server.badRequest(w, r, errors.New("invalid page"))
			return
		}
	}
	allUsers, lastPage, totalRecords, err := server.DB.GetAllUsersPaginated(pageSize, currentPage)
	if err != nil {
		_ = server.badRequest(w, r, err)
		return
	}

	var resp struct {
		CurrentPage  int            `json:"current_page"`
		PageSize     int            `json:"page_size"`
		LastPage     int            `json:"last_page"`
		TotalRecords int            `json:"total_records"`
		Users        []*models.User `json:"users"`
	}

	resp.CurrentPage = currentPage
	resp.PageSize = pageSize
	resp.LastPage = lastPage
	resp.TotalRecords = totalRecords
	resp.Users = allUsers

	_ = server.writeJSON(w, http.StatusOK, resp)
}

// OneUser gets one user by id (from the url) and returns it as JSON
func (server *Server) OneUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userID, _ := strconv.Atoi(id)

	user, err := server.DB.GetOneUser(userID)
	if err != nil {
		_ = server.badRequest(w, r, err)
		return
	}

	_ = server.writeJSON(w, http.StatusOK, user)
}

func (server *Server) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User

	err := server.readJSON(w, r, &user)
	if err != nil {
		_ = server.badRequest(w, r, err)
		return
	}
	newHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		_ = server.badRequest(w, r, err)
		return
	}
	err = server.DB.AddUser(user, string(newHash))
	if err != nil {
		_ = server.badRequest(w, r, err)
		return
	}

	var resp struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}

	resp.Error = false
	_ = server.writeJSON(w, http.StatusOK, resp)
}

// EditUser is the handler for adding or editing an existing user
func (server *Server) EditUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userID, _ := strconv.Atoi(id)

	var user models.User

	err := server.readJSON(w, r, &user)
	if err != nil {
		_ = server.badRequest(w, r, err)
		return
	}

	if userID > 0 {
		err = server.DB.EditUser(user)
		if err != nil {
			_ = server.badRequest(w, r, err)
			return
		}

		if user.Password != "" {
			newHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
			if err != nil {
				_ = server.badRequest(w, r, err)
				return
			}

			err = server.DB.UpdatePasswordForUser(user, string(newHash))
			if err != nil {
				_ = server.badRequest(w, r, err)
				return
			}
		}
	} else {
		newHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
		if err != nil {
			_ = server.badRequest(w, r, err)
			return
		}
		err = server.DB.AddUser(user, string(newHash))
		if err != nil {
			_ = server.badRequest(w, r, err)
			return
		}
	}

	var resp struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}

	resp.Error = false
	_ = server.writeJSON(w, http.StatusOK, resp)
}

// DeleteUser deletes a user, and all associated tokens, from the database
func (server *Server) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userID, _ := strconv.Atoi(id)

	err := server.DB.DeleteUser(userID)
	if err != nil {
		_ = server.badRequest(w, r, err)
		return
	}

	var resp struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}

	resp.Error = false
	_ = server.writeJSON(w, http.StatusOK, resp)
}
