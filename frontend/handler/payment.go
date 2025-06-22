package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/LamThanhNguyen/yoyo-store-backend/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

// ChargeOnce displays the page to buy one yoyo
func (server *Server) ChargeOnce(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	itemID, _ := strconv.Atoi(id)

	yoyo, err := server.DB.GetItem(itemID)
	if err != nil {
		log.Error().Err(err).Msg("ChargeOnce")
		return
	}

	data := make(map[string]interface{})
	data["item"] = yoyo

	if err := server.renderTemplate(w, r, "buy-once", &templateData{
		Data: data,
	}, "stripe-js"); err != nil {
		log.Error().Err(err).Msg("ChargeOnce")
	}
}

// BronzePlan displays the bronze plan page
func (server *Server) BronzePlan(w http.ResponseWriter, r *http.Request) {
	item, err := server.DB.GetItem(2)
	if err != nil {
		log.Error().Err(err).Msg("BronzePlan")
		return
	}

	data := make(map[string]interface{})
	data["item"] = item

	if err := server.renderTemplate(w, r, "bronze-plan", &templateData{
		Data: data,
	}); err != nil {
		log.Error().Err(err).Msg("BronzePlan")
	}
}

// BronzePlanReceipt displays the receipt for bronze plans
func (server *Server) BronzePlanReceipt(w http.ResponseWriter, r *http.Request) {
	if err := server.renderTemplate(w, r, "receipt-plan", &templateData{}); err != nil {
		log.Error().Err(err).Msg("BronzePlanReceipt")
	}
}

// PaymentSucceeded displays the receipt page
func (server *Server) PaymentSucceeded(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Error().Err(err).Msg("PaymentSucceeded")
		return
	}

	// read posted data
	yoyoID, _ := strconv.Atoi(r.Form.Get("product_id"))

	txnData, err := server.GetTransactionData(r)
	if err != nil {
		log.Error().Err(err).Msg("PaymentSucceeded")
		return
	}

	// create a new customer
	customerID, err := server.SaveCustomer(txnData.FirstName, txnData.LastName, txnData.Email)
	if err != nil {
		log.Error().Err(err).Msg("PaymentSucceeded")
		return
	}

	// create a new transaction
	txn := models.Transaction{
		Amount:              txnData.PaymentAmount,
		Currency:            txnData.PaymentCurrency,
		LastFour:            txnData.LastFour,
		ExpiryMonth:         txnData.ExpiryMonth,
		ExpiryYear:          txnData.ExpiryYear,
		BankReturnCode:      txnData.BankReturnCode,
		PaymentIntent:       txnData.PaymentIntentID,
		PaymentMethod:       txnData.PaymentMethodID,
		TransactionStatusID: 2,
	}

	txnID, err := server.SaveTransaction(txn)
	if err != nil {
		log.Error().Err(err).Msg("PaymentSucceeded")
		return
	}

	// create a new order
	order := models.Order{
		ItemID:        yoyoID,
		TransactionID: txnID,
		CustomerID:    customerID,
		StatusID:      1,
		Quantity:      1,
		Amount:        txnData.PaymentAmount,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	orderID, err := server.SaveOrder(order)
	if err != nil {
		log.Error().Err(err).Msg("PaymentSucceeded")
		return
	}

	// call microservice
	inv := Invoice{
		ID:        orderID,
		Amount:    order.Amount,
		Product:   "Yoyo",
		Quantity:  order.Quantity,
		FirstName: txnData.FirstName,
		LastName:  txnData.LastName,
		Email:     txnData.Email,
		CreatedAt: time.Now(),
	}

	err = server.callInvoiceMicro(inv)
	if err != nil {
		log.Error().Err(err).Msg("PaymentSucceeded")
	}

	// write this data to session, and then redirect user to new page
	server.Session.Put(r.Context(), "receipt", txnData)
	http.Redirect(w, r, "/receipt", http.StatusSeeOther)
}
