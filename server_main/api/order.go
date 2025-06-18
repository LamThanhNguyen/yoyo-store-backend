package api

import "github.com/LamThanhNguyen/yoyo-store-backend/internal/models"

// SaveOrder saves a order and returns id
func (server *Server) SaveOrder(order models.Order) (int, error) {
	id, err := server.DB.InsertOrder(order)
	if err != nil {
		return 0, err
	}
	return id, nil
}
