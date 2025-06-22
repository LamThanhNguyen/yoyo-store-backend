package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/LamThanhNguyen/yoyo-store-backend/internal/pb"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Invoice struct {
	ID        int       `json:"id"`
	Quantity  int       `json:"quantity"`
	Amount    int       `json:"amount"`
	Product   string    `json:"product"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// Receipt displays a receipt
func (server *Server) Receipt(w http.ResponseWriter, r *http.Request) {
	txn := server.Session.Get(r.Context(), "receipt").(TransactionData)
	data := make(map[string]interface{})
	data["txn"] = txn
	server.Session.Remove(r.Context(), "receipt")
	if err := server.renderTemplate(w, r, "receipt", &templateData{
		Data: data,
	}); err != nil {
		log.Error().Err(err).Msg("Receipt")
	}
}

// callInvoiceMicro calls the invoicing microservice
func (server *Server) callInvoiceMicro(inv Invoice) error {
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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
	reqCtx, reqCancel := context.WithTimeout(context.Background(), 5*time.Second)
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
