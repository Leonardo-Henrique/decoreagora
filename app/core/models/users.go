package models

import "time"

type (
	User struct {
		ID               int              `json:"id"`
		PublicID         string           `json:"public_id"`
		Name             string           `json:"name"`
		Email            string           `json:"email"`
		LastLogin        time.Time        `json:"last_login"`
		AvailableCredits AvailableCredits `json:"available_credits"`
		Subscription     Subscription     `json:"subscription"`
	}

	Subscription struct {
		ID                   int    `json:"id"`
		StripeCostumerID     string `json:"stripe_costumer_id"`
		StripeSubscriptionID string `json:"stripe_subscription_id"`
		StripePriceID        string `json:"stripe_price_id"`
		IsActive             bool   `json:"is_active"`
	}

	AvailableCredits struct {
		ID    string `json:"id"`
		Total int    `json:"total"`
	}
)
