package api

import (
	"net/http"
	"strconv"

	"github.com/LamThanhNguyen/yoyo-store-backend/internal/models"
	"github.com/go-chi/chi/v5"
)

// SaveOrder saves a order and returns id
func (server *Server) SaveOrder(order models.Order) (int, error) {
	id, err := server.DB.InsertOrder(order)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// AllSales returns all sales as a slice
func (server *Server) AllSales(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		PageSize    int `json:"page_size"`
		CurrentPage int `json:"page"`
	}

	err := server.readJSON(w, r, &payload)
	if err != nil {
		server.badRequest(w, r, err)
		return
	}

	allSales, lastPage, totalRecords, err := server.DB.GetAllOrdersPaginated(payload.PageSize, payload.CurrentPage)
	if err != nil {
		server.badRequest(w, r, err)
		return
	}

	var resp struct {
		CurrentPage  int             `json:"current_page"`
		PageSize     int             `json:"page_size"`
		LastPage     int             `json:"last_page"`
		TotalRecords int             `json:"total_records"`
		Orders       []*models.Order `json:"orders"`
	}

	resp.CurrentPage = payload.CurrentPage
	resp.PageSize = payload.PageSize
	resp.LastPage = lastPage
	resp.TotalRecords = totalRecords
	resp.Orders = allSales

	server.writeJSON(w, http.StatusOK, resp)
}

// AllSubscriptions returns all subscriptions as a slice
func (server *Server) AllSubscriptions(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		PageSize    int `json:"page_size"`
		CurrentPage int `json:"page"`
	}

	err := server.readJSON(w, r, &payload)
	if err != nil {
		server.badRequest(w, r, err)
		return
	}

	allSales, lastPage, totalRecords, err := server.DB.GetAllSubscriptionsPaginated(payload.PageSize, payload.CurrentPage)
	if err != nil {
		server.badRequest(w, r, err)
		return
	}

	var resp struct {
		CurrentPage  int             `json:"current_page"`
		PageSize     int             `json:"page_size"`
		LastPage     int             `json:"last_page"`
		TotalRecords int             `json:"total_records"`
		Orders       []*models.Order `json:"orders"`
	}

	resp.CurrentPage = payload.CurrentPage
	resp.PageSize = payload.PageSize
	resp.LastPage = lastPage
	resp.TotalRecords = totalRecords
	resp.Orders = allSales

	server.writeJSON(w, http.StatusOK, resp)
}

// GetSale returns one sale as json, by id
func (server *Server) GetSale(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	orderID, _ := strconv.Atoi(id)

	order, err := server.DB.GetOrderByID(orderID)
	if err != nil {
		server.badRequest(w, r, err)
		return
	}

	server.writeJSON(w, http.StatusOK, order)
}
