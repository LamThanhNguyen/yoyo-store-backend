package api

import (
	"errors"
	"testing"

	"github.com/LamThanhNguyen/yoyo-store-backend/internal/models"
	"go.uber.org/mock/gomock"
)

func TestSaveCustomer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := NewMockcustomerInserter(ctrl)

	fname, lname, email := "John", "Doe", "john@example.com"
	expected := models.Customer{FirstName: fname, LastName: lname, Email: email}

	mockDB.EXPECT().InsertCustomer(expected).Return(1, nil)

	id, err := saveCustomer(mockDB, fname, lname, email)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if id != 1 {
		t.Fatalf("expected id 1, got %d", id)
	}
}

func TestSaveCustomerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := NewMockcustomerInserter(ctrl)

	fname, lname, email := "Jane", "Smith", "jane@example.com"
	expected := models.Customer{FirstName: fname, LastName: lname, Email: email}
	mockErr := errors.New("insert failed")
	mockDB.EXPECT().InsertCustomer(expected).Return(0, mockErr)

	id, err := saveCustomer(mockDB, fname, lname, email)
	if err == nil {
		t.Fatalf("expected error")
	}
	if id != 0 {
		t.Fatalf("expected id 0, got %d", id)
	}
}
