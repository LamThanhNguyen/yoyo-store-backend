package api

import (
	"testing"
	"time"

	pb "github.com/LamThanhNguyen/yoyo-store-backend/internal/pb"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestSendInvoice(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := pb.NewMockInvoiceServiceClient(ctrl)
	inv := Invoice{
		ID:        1,
		Quantity:  2,
		Amount:    100,
		Product:   "test",
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
		CreatedAt: time.Now(),
	}

	mockClient.EXPECT().CreateAndSendInvoice(gomock.Any(), &pb.CreateInvoiceRequest{
		Id:        int32(inv.ID),
		Quantity:  int32(inv.Quantity),
		Amount:    int32(inv.Amount),
		Product:   inv.Product,
		FirstName: inv.FirstName,
		LastName:  inv.LastName,
		Email:     inv.Email,
		CreatedAt: timestamppb.New(inv.CreatedAt),
	}).Return(&pb.CreateInvoiceResponse{}, nil)

	server := &Server{}
	if err := server.sendInvoice(mockClient, inv); err != nil {
		t.Fatalf("sendInvoice returned error: %v", err)
	}
}
