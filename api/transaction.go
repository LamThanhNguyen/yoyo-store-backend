package main

import "github.com/LamThanhNguyen/yoyo-store-backend/internal/models"

func (server *Server) SaveTransaction(txn models.Transaction) (int, error) {
	id, err := server.DB.InsertTransaction(txn)
	if err != nil {
		return 0, err
	}
	return id, nil
}
