package cards

import (
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/customer"
	"github.com/stripe/stripe-go/v82/paymentintent"
	"github.com/stripe/stripe-go/v82/paymentmethod"
	"github.com/stripe/stripe-go/v82/refund"
	"github.com/stripe/stripe-go/v82/subscription"
)

// Card holds the information needed by this package
type Card struct {
	Secret   string
	Key      string
	Currency string
}

// Transaction is the type to store information for a given transaction
type Transaction struct {
	TransactionStatusID int
	Amount              int
	Currency            string
	LastFour            string
	BankReturnCode      string
}

// Charge is an alias to CreatePaymentIntent
func (c *Card) Charge(currency string, amount int) (*stripe.PaymentIntent, string, error) {
	return c.CreatePaymentIntent(currency, amount)
}

// CreatePaymentIntent attempts to get a payment intent object from Stripe
func (c *Card) CreatePaymentIntent(currency string, amount int) (*stripe.PaymentIntent, string, error) {
	stripe.Key = c.Secret

	// create a payment intent
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(int64(amount)),
		Currency: stripe.String(currency),
	}

	//params.AddMetadata("key", "value")

	pi, err := paymentintent.New(params)
	if err != nil {
		msg := ""
		if stripeErr, ok := err.(*stripe.Error); ok {
			msg = cardErrorMessage(stripeErr.Code)
		}
		return nil, msg, err
	}
	return pi, "", nil
}

// GetPaymentMethod gets the payment method by payment intend id
func (c *Card) GetPaymentMethod(s string) (*stripe.PaymentMethod, error) {
	stripe.Key = c.Secret

	pm, err := paymentmethod.Get(s, nil)
	if err != nil {
		return nil, err
	}
	return pm, nil
}

// RetrievePaymentIntent gets an existing payment intent by id
func (c *Card) RetrievePaymentIntent(id string) (*stripe.PaymentIntent, error) {
	stripe.Key = c.Secret

	pi, err := paymentintent.Get(id, nil)
	if err != nil {
		return nil, err
	}
	return pi, nil
}

// SubscribeToPlan subscribes a stripe customer to a stripe plan
func (c *Card) SubscribeToPlan(cust *stripe.Customer, priceID, email, last4, cardType string) (*stripe.Subscription, error) {
	stripeCustomerID := cust.ID
	items := []*stripe.SubscriptionItemsParams{
		{Price: stripe.String(priceID)},
	}

	params := &stripe.SubscriptionParams{
		Customer: stripe.String(stripeCustomerID),
		Items:    items,
	}

	params.AddMetadata("last_four", last4)
	params.AddMetadata("card_type", cardType)
	params.AddExpand("latest_invoice.payment_intent")
	subscription, err := subscription.New(params)
	if err != nil {
		return nil, err
	}
	return subscription, nil
}

// CreateCustomer creates a stripe customer
func (c *Card) CreateCustomer(pm, email string) (*stripe.Customer, string, error) {
	stripe.Key = c.Secret
	customerParams := &stripe.CustomerParams{
		PaymentMethod: stripe.String(pm),
		Email:         stripe.String(email),
		InvoiceSettings: &stripe.CustomerInvoiceSettingsParams{
			DefaultPaymentMethod: stripe.String(pm),
		},
	}

	cust, err := customer.New(customerParams)
	if err != nil {
		msg := ""
		if stripeErr, ok := err.(*stripe.Error); ok {
			msg = cardErrorMessage(stripeErr.Code)
		}
		return nil, msg, err
	}
	return cust, "", nil
}

// Refund refunds an amount for a paymentIntent
func (c *Card) Refund(pi string, amount int) error {
	stripe.Key = c.Secret
	amountToRefund := int64(amount)

	refundParams := &stripe.RefundParams{
		Amount:        &amountToRefund,
		PaymentIntent: &pi,
	}

	_, err := refund.New(refundParams)
	if err != nil {
		return err
	}

	return nil
}

// CancelSubscription cancels a subscription, by subscription id
func (c *Card) CancelSubscription(subID string) error {
	stripe.Key = c.Secret

	params := &stripe.SubscriptionParams{
		CancelAtPeriodEnd: stripe.Bool(true),
	}

	_, err := subscription.Update(subID, params)
	if err != nil {
		return err
	}
	return nil
}

// cardErrorMessage returns human readable versions of card error messages
func cardErrorMessage(code stripe.ErrorCode) string {
	switch code {
	case stripe.ErrorCodeCardDeclined:
		return "Your card was declined"
	case stripe.ErrorCodeExpiredCard:
		return "Your card is expired"
	case stripe.ErrorCodeIncorrectCVC:
		return "Incorrect CVC code"
	case stripe.ErrorCodeIncorrectZip:
		return "Incorrect zip/postal code"
	case stripe.ErrorCodeAmountTooLarge:
		return "The amount is too large to charge to your card"
	case stripe.ErrorCodeAmountTooSmall:
		return "The amount is too small to charge to your card"
	case stripe.ErrorCodeBalanceInsufficient:
		return "Insufficient balance"
	case stripe.ErrorCodePostalCodeInvalid:
		return "Your postal code is invalid"
	default:
		return "Your card was declined"
	}
}
