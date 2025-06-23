package api

import "github.com/LamThanhNguyen/yoyo-store-backend/internal/models"

// transactionInserter allows inserting a transaction into the database.
type transactionInserter interface {
	InsertTransaction(models.Transaction) (int, error)
}

// saveTransaction inserts a transaction using the provided database interface.
func saveTransaction(db transactionInserter, txn models.Transaction) (int, error) {
	return db.InsertTransaction(txn)
}

func (server *Server) SaveTransaction(txn models.Transaction) (int, error) {
	return saveTransaction(server.DB, txn)
}
