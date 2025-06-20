package api

import (
	"context"
	"fmt"

	"github.com/LamThanhNguyen/yoyo-store-backend/internal/pb"
)

// GRPCServer implements pb.InvoiceServiceServer
type GRPCServer struct {
	pb.UnimplementedInvoiceServiceServer
	*Server
}

func NewGRPCServer(s *Server) *GRPCServer {
	return &GRPCServer{Server: s}
}

func (g *GRPCServer) CreateAndSendInvoice(ctx context.Context, req *pb.CreateInvoiceRequest) (*pb.CreateInvoiceResponse, error) {
	order := Order{
		ID:        int(req.Id),
		Quantity:  int(req.Quantity),
		Amount:    int(req.Amount),
		Product:   req.Product,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		CreatedAt: req.CreatedAt.AsTime(),
	}

	if err := g.createInvoicePDF(order); err != nil {
		return nil, err
	}

	attachments := []string{fmt.Sprintf("./invoices/%d.pdf", order.ID)}
	if err := g.SendMail("info@yoyo.com", order.Email, "Your invoice", "invoice", attachments, nil); err != nil {
		return nil, err
	}

	msg := fmt.Sprintf("Invoice %d.pdf created and sent to %s", order.ID, order.Email)
	return &pb.CreateInvoiceResponse{Message: msg}, nil
}
