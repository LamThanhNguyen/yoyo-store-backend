package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/LamThanhNguyen/yoyo-store-backend/internal/cards"
	"github.com/LamThanhNguyen/yoyo-store-backend/internal/models"
	"github.com/LamThanhNguyen/yoyo-store-backend/internal/validator"
	"github.com/rs/zerolog/log"
	"github.com/stripe/stripe-go/v82"
)

type stripePayload struct {
	Currency      string `json:"currency"`
	Amount        string `json:"amount"`
	PaymentMethod string `json:"payment_method"`
	Email         string `json:"email"`
	CardBrand     string `json:"card_brand"`
	ExpiryMonth   int    `json:"exp_month"`
	ExpiryYear    int    `json:"exp_year"`
	LastFour      string `json:"last_four"`
	Plan          string `json:"plan"`
	ProductID     string `json:"product_id"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
}

func (server *Server) GetPaymentIntent(w http.ResponseWriter, r *http.Request) {
	var payload stripePayload

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		log.Error().Err(err)
		return
	}

	amount, err := strconv.Atoi(payload.Amount)
	if err != nil {
		log.Error().Err(err)
	}

	card := cards.Card{
		Secret:   server.config.StripeSecret,
		Key:      server.config.StripeKey,
		Currency: payload.Currency,
	}

	okay := true

	pi, msg, err := card.Charge(payload.Currency, amount)
	if err != nil {
		okay = false
	}

	if okay {
		out, err := json.MarshalIndent(pi, "", "   ")
		if err != nil {
			log.Error().Err(err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
	} else {
		j := jsonResponse{
			OK:      false,
			Message: msg,
			Content: "",
		}

		out, err := json.MarshalIndent(j, "", "   ")
		if err != nil {
			log.Error().Err(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
	}
}

// CreateCustomerAndSubscribeToPlan is the handler for subscribing to the bronze plan
func (server *Server) CreateCustomerAndSubscribeToPlan(w http.ResponseWriter, r *http.Request) {
	var data stripePayload

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Error().Err(err)
		return
	}

	// validate data
	v := validator.New()
	v.Check(len(data.FirstName) > 1, "first_name", "must be at least 2 characters")

	if !v.Valid() {
		server.failedValidation(w, r, v.Errors)
		return
	}

	card := cards.Card{
		Secret:   server.config.StripeSecret,
		Key:      server.config.StripeKey,
		Currency: data.Currency,
	}

	okay := true
	var subscription *stripe.Subscription
	txnMsg := "Transaction successful"

	stripeCustomer, msg, err := card.CreateCustomer(data.PaymentMethod, data.Email)
	if err != nil {
		log.Error().Err(err)
		okay = false
		txnMsg = msg
	}

	if okay {
		subscription, err = card.SubscribeToPlan(stripeCustomer, data.Plan, data.Email, data.LastFour, "")
		if err != nil {
			log.Error().Err(err)
			okay = false
			txnMsg = "Error subscribing customer"
		}
	}

	if okay {
		productID, _ := strconv.Atoi(data.ProductID)
		customerID, err := server.SaveCustomer(data.FirstName, data.LastName, data.Email)
		if err != nil {
			log.Error().Err(err)
			return
		}

		// create a new txn
		amount, _ := strconv.Atoi(data.Amount)

		txn := models.Transaction{
			Amount:              amount,
			Currency:            "usd",
			LastFour:            data.LastFour,
			ExpiryMonth:         data.ExpiryMonth,
			ExpiryYear:          data.ExpiryYear,
			TransactionStatusID: 2,
			PaymentIntent:       subscription.ID,
			PaymentMethod:       data.PaymentMethod,
		}

		txnID, err := server.SaveTransaction(txn)
		if err != nil {
			log.Error().Err(err)
			return
		}

		// create order
		order := models.Order{
			WidgetID:      productID,
			TransactionID: txnID,
			CustomerID:    customerID,
			StatusID:      1,
			Quantity:      1,
			Amount:        amount,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		orderID, err := server.SaveOrder(order)
		if err != nil {
			log.Error().Err(err)
			return
		}

		inv := Invoice{
			ID:        orderID,
			Amount:    2000,
			Product:   "Bronze Plan monthly subscription",
			Quantity:  order.Quantity,
			FirstName: data.FirstName,
			LastName:  data.LastName,
			Email:     data.Email,
			CreatedAt: time.Now(),
		}

		err = server.callInvoiceMicro(inv)
		if err != nil {
			log.Error().Err(err)
			return
		}
	}

	resp := jsonResponse{
		OK:      false,
		Message: txnMsg,
	}

	out, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		log.Error().Err(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// VirtualTerminalPaymentSucceeded displays a page with receipt information
func (server *Server) VirtualTerminalPaymentSucceeded(w http.ResponseWriter, r *http.Request) {
	var txnData struct {
		PaymentAmount   int    `json:"amount"`
		PaymentCurrency string `json:"currency"`
		FirstName       string `json:"first_name"`
		LastName        string `json:"last_name"`
		Email           string `json:"email"`
		PaymentIntent   string `json:"payment_intent"`
		PaymentMethod   string `json:"payment_method"`
		BankReturnCode  string `json:"bank_return_code"`
		ExpiryMonth     int    `json:"expiry_month"`
		ExpiryYear      int    `json:"expiry_year"`
		LastFour        string `json:"last_four"`
	}

	err := server.readJSON(w, r, &txnData)
	if err != nil {
		server.badRequest(w, r, err)
		return
	}

	card := cards.Card{
		Secret: server.config.StripeSecret,
		Key:    server.config.StripeKey,
	}

	pi, err := card.RetrievePaymentIntent(txnData.PaymentIntent)
	if err != nil {
		server.badRequest(w, r, err)
		return
	}

	pm, err := card.GetPaymentMethod(txnData.PaymentMethod)
	if err != nil {
		server.badRequest(w, r, err)
		return
	}

	txnData.LastFour = pm.Card.Last4
	txnData.ExpiryMonth = int(pm.Card.ExpMonth)
	txnData.ExpiryYear = int(pm.Card.ExpYear)

	txn := models.Transaction{
		Amount:              txnData.PaymentAmount,
		Currency:            txnData.PaymentCurrency,
		LastFour:            txnData.LastFour,
		ExpiryMonth:         txnData.ExpiryMonth,
		ExpiryYear:          txnData.ExpiryYear,
		PaymentIntent:       txnData.PaymentIntent,
		PaymentMethod:       txnData.PaymentMethod,
		BankReturnCode:      pi.LatestCharge.ID,
		TransactionStatusID: 2,
	}

	_, err = server.SaveTransaction(txn)
	if err != nil {
		server.badRequest(w, r, err)
		return
	}

	server.writeJSON(w, http.StatusOK, txn)
}
