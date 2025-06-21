package handler

import (
	"net/http"

	"github.com/LamThanhNguyen/yoyo-store-backend/internal/models"
	"github.com/rs/zerolog/log"
)

// SaveOrder saves a order and returns id
func (server *Server) SaveOrder(order models.Order) (int, error) {
	id, err := server.DB.InsertOrder(order)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// AllSales shows the all sales page
func (server *Server) AllSales(w http.ResponseWriter, r *http.Request) {
	if err := server.renderTemplate(w, r, "all-sales", &templateData{}); err != nil {
		log.Error().Err(err)
	}
}

// AllSubscriptions shows all subscription page
func (server *Server) AllSubscriptions(w http.ResponseWriter, r *http.Request) {
	if err := server.renderTemplate(w, r, "all-subscriptions", &templateData{}); err != nil {
		log.Error().Err(err)
	}
}

// ShowSale shows one sale page
func (server *Server) ShowSale(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["title"] = "Sale"
	stringMap["cancel"] = "/admin/all-sales"
	stringMap["refund-url"] = "/api/admin/refund"
	stringMap["refund-btn"] = "Refund Order"
	stringMap["refunded-badge"] = "Refunded"
	stringMap["refunded-msg"] = "Charge refunded"

	if err := server.renderTemplate(w, r, "sale", &templateData{
		StringMap: stringMap,
	}); err != nil {
		log.Error().Err(err)
	}
}

// ShowSubscription shows one subscription page
func (server *Server) ShowSubscription(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["title"] = "Subscription"
	stringMap["cancel"] = "/admin/all-subscriptions"
	stringMap["refund-url"] = "/api/admin/cancel-subscription"
	stringMap["refund-btn"] = "Cancel Subscription"
	stringMap["refunded-badge"] = "Cancelled"
	stringMap["refunded-msg"] = "Subscription cancelled"

	if err := server.renderTemplate(w, r, "sale", &templateData{
		StringMap: stringMap,
	}); err != nil {
		log.Error().Err(err)
	}
}
