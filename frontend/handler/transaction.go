package handler

import (
	"net/http"
	"strconv"

	"github.com/LamThanhNguyen/yoyo-store-backend/internal/cards"
	"github.com/LamThanhNguyen/yoyo-store-backend/internal/models"
	"github.com/rs/zerolog/log"
)

type TransactionData struct {
	FirstName       string
	LastName        string
	Email           string
	PaymentIntentID string
	PaymentMethodID string
	PaymentAmount   int
	PaymentCurrency string
	LastFour        string
	ExpiryMonth     int
	ExpiryYear      int
	BankReturnCode  string
}

// GetTransactionData gets txn data from post and stripe
func (server *Server) GetTransactionData(r *http.Request) (TransactionData, error) {
	var txnData TransactionData
	err := r.ParseForm()
	if err != nil {
		log.Error().Err(err).Msg("GetTransactionData")
		return txnData, err
	}

	firstName := r.Form.Get("first_name")
	lastName := r.Form.Get("last_name")
	email := r.Form.Get("email")
	paymentIntent := r.Form.Get("payment_intent")
	paymentMethod := r.Form.Get("payment_method")
	paymentAmount := r.Form.Get("payment_amount")
	paymentCurrency := r.Form.Get("payment_currency")
	amount, _ := strconv.Atoi(paymentAmount)

	card := cards.Card{
		Secret: server.config.StripeSecret,
		Key:    server.config.StripeKey,
	}

	pi, err := card.RetrievePaymentIntent(paymentIntent)
	if err != nil {
		log.Error().Err(err).Msg("GetTransactionData")
		return txnData, err
	}

	pm, err := card.GetPaymentMethod(paymentMethod)
	if err != nil {
		log.Error().Err(err).Msg("GetTransactionData")
		return txnData, err
	}

	lastFour := pm.Card.Last4
	expiryMonth := pm.Card.ExpMonth
	expiryYear := pm.Card.ExpYear

	txnData = TransactionData{
		FirstName:       firstName,
		LastName:        lastName,
		Email:           email,
		PaymentIntentID: paymentIntent,
		PaymentMethodID: paymentMethod,
		PaymentAmount:   amount,
		PaymentCurrency: paymentCurrency,
		LastFour:        lastFour,
		ExpiryMonth:     int(expiryMonth),
		ExpiryYear:      int(expiryYear),
		BankReturnCode:  pi.LatestCharge.ID,
	}
	return txnData, nil
}

// SaveTransaction saves a txn and returns id
func (server *Server) SaveTransaction(txn models.Transaction) (int, error) {
	id, err := server.DB.InsertTransaction(txn)
	if err != nil {
		return 0, err
	}
	return id, nil
}
