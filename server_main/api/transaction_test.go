package api

import (
	"errors"
	"testing"

	"github.com/LamThanhNguyen/yoyo-store-backend/internal/models"
	"go.uber.org/mock/gomock"
)

func TestSaveTransaction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := NewMocktransactionInserter(ctrl)

	txn := models.Transaction{Amount: 100}
	mockDB.EXPECT().InsertTransaction(txn).Return(5, nil)

	id, err := saveTransaction(mockDB, txn)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if id != 5 {
		t.Fatalf("expected id 5, got %d", id)
	}
}

func TestSaveTransactionError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := NewMocktransactionInserter(ctrl)

	txn := models.Transaction{Amount: 200}
	mockErr := errors.New("insert failed")
	mockDB.EXPECT().InsertTransaction(txn).Return(0, mockErr)

	id, err := saveTransaction(mockDB, txn)
	if err == nil {
		t.Fatalf("expected error")
	}
	if id != 0 {
		t.Fatalf("expected id 0, got %d", id)
	}
}
