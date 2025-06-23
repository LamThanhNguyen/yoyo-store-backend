package api

import (
	"errors"
	"testing"

	"github.com/LamThanhNguyen/yoyo-store-backend/internal/models"
	"go.uber.org/mock/gomock"
)

func TestSaveOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := NewMockorderInserter(ctrl)

	order := models.Order{ID: 1, Quantity: 2}

	mockDB.EXPECT().InsertOrder(order).Return(3, nil)

	id, err := saveOrder(mockDB, order)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if id != 3 {
		t.Fatalf("expected id 3, got %d", id)
	}
}

func TestSaveOrderError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := NewMockorderInserter(ctrl)

	order := models.Order{ID: 2, Quantity: 1}
	mockErr := errors.New("insert failed")
	mockDB.EXPECT().InsertOrder(order).Return(0, mockErr)

	id, err := saveOrder(mockDB, order)
	if err == nil {
		t.Fatalf("expected error")
	}
	if id != 0 {
		t.Fatalf("expected id 0, got %d", id)
	}
}
