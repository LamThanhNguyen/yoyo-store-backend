package handler

import "github.com/LamThanhNguyen/yoyo-store-backend/internal/models"

// SaveCustomer saves a customer and returns id
func (server *Server) SaveCustomer(firstName, lastName, email string) (int, error) {
	customer := models.Customer{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
	}

	id, err := server.DB.InsertCustomer(customer)
	if err != nil {
		return 0, err
	}
	return id, nil
}
