package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/LamThanhNguyen/yoyo-store-backend/internal/models"
	"github.com/go-chi/chi/v5"
)

// orderInserter provides the behaviour required to insert an order.
// Having this interface allows the use of gomock in tests.
type orderInserter interface {
	InsertOrder(models.Order) (int, error)
}

// saveOrder inserts a new order through the provided interface.
func saveOrder(db orderInserter, order models.Order) (int, error) {
	return db.InsertOrder(order)
}

// SaveOrder saves a order and returns id
func (server *Server) SaveOrder(order models.Order) (int, error) {
	return saveOrder(server.DB, order)
}

// AllSales returns all sales as a slice
func (server *Server) AllSales(w http.ResponseWriter, r *http.Request) {
	pageSize := 10   // default
	currentPage := 1 // default

	// Parse query params
	if val := r.URL.Query().Get("page_size"); val != "" {
		if ps, err := strconv.Atoi(val); err == nil && ps > 0 {
			pageSize = ps
		} else {
			server.badRequest(w, r, errors.New("invalid page_size"))
			return
		}
	}
	// Parse query params
	if val := r.URL.Query().Get("page"); val != "" {
		if cp, err := strconv.Atoi(val); err == nil && cp > 0 {
			currentPage = cp
		} else {
			server.badRequest(w, r, errors.New("invalid page"))
			return
		}
	}

	allSales, lastPage, totalRecords, err := server.DB.GetAllOrdersPaginated(pageSize, currentPage)
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

	resp.CurrentPage = currentPage
	resp.PageSize = pageSize
	resp.LastPage = lastPage
	resp.TotalRecords = totalRecords
	resp.Orders = allSales

	server.writeJSON(w, http.StatusOK, resp)
}

// AllSubscriptions returns all subscriptions as a slice
func (server *Server) AllSubscriptions(w http.ResponseWriter, r *http.Request) {
	pageSize := 10   // default
	currentPage := 1 // default

	// Parse query params
	if val := r.URL.Query().Get("page_size"); val != "" {
		if ps, err := strconv.Atoi(val); err == nil && ps > 0 {
			pageSize = ps
		} else {
			server.badRequest(w, r, errors.New("invalid page_size"))
			return
		}
	}
	// Parse query params
	if val := r.URL.Query().Get("page"); val != "" {
		if cp, err := strconv.Atoi(val); err == nil && cp > 0 {
			currentPage = cp
		} else {
			server.badRequest(w, r, errors.New("invalid page"))
			return
		}
	}

	allSales, lastPage, totalRecords, err := server.DB.GetAllSubscriptionsPaginated(pageSize, currentPage)
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

	resp.CurrentPage = currentPage
	resp.PageSize = pageSize
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
