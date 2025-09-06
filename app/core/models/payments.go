package models

import (
	"errors"
	"time"
)

type (
	CreateSessionRequest struct {
		Email string `json:"email"`
		Plan  string `json:"plan"`
	}

	CreatePaymentSessionRequest struct {
		StripeCustomerID string
		PriceID          string
		SuccessUrl       string
		CancelUrl        string
	}

	PaymentHistory struct {
		ID              int       `json:"id"`
		PublicID        string    `json:"public_id"`
		CustomerID      string    `json:"customer_id"`
		ProcessedAt     time.Time `json:"processed_at"`
		StripePriceID   string    `json:"stripe_price_id"`
		AmountPaid      int       `json:"amount_paid"`
		CreditsReceived int       `json:"credits_received"`
	}
)

func (c *CreateSessionRequest) ValidateRequiredFields() error {
	if c.Email == "" || c.Plan == "" {
		return errors.New("required fields are missing")
	}
	return nil
}
