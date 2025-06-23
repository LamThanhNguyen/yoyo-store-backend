package api

import "github.com/LamThanhNguyen/yoyo-store-backend/internal/models"

// customerInserter defines the behaviour required to insert a customer.
// It allows us to generate mocks for testing with gomock.
type customerInserter interface {
	InsertCustomer(models.Customer) (int, error)
}

// saveCustomer inserts a new customer using the provided database interface.
// This helper exists so that we can easily mock database interactions in tests.
func saveCustomer(db customerInserter, firstName, lastName, email string) (int, error) {
	customer := models.Customer{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
	}

	return db.InsertCustomer(customer)
}

func (server *Server) SaveCustomer(
	firstName,
	lastName,
	email string,
) (int, error) {
	return saveCustomer(server.DB, firstName, lastName, email)
}
