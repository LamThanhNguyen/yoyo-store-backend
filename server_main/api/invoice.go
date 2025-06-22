package api

import (
	"context"
	"time"

	"github.com/LamThanhNguyen/yoyo-store-backend/internal/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Invoice describes the JSON payload sent to the microservice
type Invoice struct {
	ID        int       `json:"id"`
	ItemID    int       `json:"item_id"`
	Amount    int       `json:"amount"`
	Product   string    `json:"product"`
	Quantity  int       `json:"quantity"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// callInvoiceMicro calls the invoicing microservice
func (server *Server) callInvoiceMicro(inv Invoice) error {
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientConn, err := grpc.NewClient(
		server.config.InvoiceGrpcAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return err
	}
	defer clientConn.Close()

	client := pb.NewInvoiceServiceClient(clientConn)

	// Separate context for the actual request (best practice)
	reqCtx, reqCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer reqCancel()

	_, err = client.CreateAndSendInvoice(reqCtx, &pb.CreateInvoiceRequest{
		Id:        int32(inv.ID),
		Quantity:  int32(inv.Quantity),
		Amount:    int32(inv.Amount),
		Product:   inv.Product,
		FirstName: inv.FirstName,
		LastName:  inv.LastName,
		Email:     inv.Email,
		CreatedAt: timestamppb.New(inv.CreatedAt),
	})
	return err
}
