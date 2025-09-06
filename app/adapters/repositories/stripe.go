package repositories

import (
	"github.com/Leonardo-Henrique/decoreagora/app/core/config"
	"github.com/Leonardo-Henrique/decoreagora/app/core/ports"
	"github.com/stripe/stripe-go/v82"
	bp "github.com/stripe/stripe-go/v82/billingportal/session"
	"github.com/stripe/stripe-go/v82/checkout/session"
	"github.com/stripe/stripe-go/v82/customer"
)

type Stripe struct {
}

func NewStripeClient() *Stripe {
	stripe.Key = config.C.STRIPE_SECRET_KEY
	return &Stripe{}
}

func (s *Stripe) CreateCustomer(email string, name string) (ports.Customer, error) {
	params := &stripe.CustomerParams{
		Email: stripe.String(email),
		Name:  stripe.String(name),
	}
	customer, err := customer.New(params)
	if err != nil {
		return ports.Customer{}, err
	}
	return ports.Customer{
		ID:    customer.ID,
		Email: customer.Email,
	}, nil
}

func (s *Stripe) CreatePaymentSession(customerID, priceID, successURL, cancelURL string) (ports.CheckoutSession, error) {
	params := &stripe.CheckoutSessionParams{
		Customer: stripe.String(customerID),
		Mode:     stripe.String(string(stripe.CheckoutSessionModePayment)), // Changed from subscription to payment
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(priceID),
				Quantity: stripe.Int64(1),
			},
		},
		SuccessURL: stripe.String(successURL),
		CancelURL:  stripe.String(cancelURL),
		Metadata: map[string]string{
			"price_id": priceID,
		},
	}
	session, err := session.New(params)
	if err != nil {
		return ports.CheckoutSession{}, err
	}
	return ports.CheckoutSession{
		ID:  session.ID,
		URL: session.URL,
	}, nil
}

func (s *Stripe) StartSubscription(customerID, priceID, successURL, cancelURL string) (ports.CheckoutSession, error) {
	params := &stripe.CheckoutSessionParams{
		Customer: stripe.String(customerID),
		Mode:     stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(priceID),
				Quantity: stripe.Int64(1),
			},
		},
		SuccessURL: stripe.String(successURL),
		CancelURL:  stripe.String(cancelURL),
	}
	session, err := session.New(params)
	if err != nil {
		return ports.CheckoutSession{}, err
	}

	return ports.CheckoutSession{
		ID:  session.ID,
		URL: session.URL,
	}, nil
}

func (s *Stripe) CreateBillingPortalSession(customerID, returnURL string) (string, error) {
	params := &stripe.BillingPortalSessionParams{
		Customer:  stripe.String(customerID),
		ReturnURL: stripe.String(returnURL),
	}
	session, err := bp.New(params)
	if err != nil {
		return "", err
	}
	return session.URL, nil
}
