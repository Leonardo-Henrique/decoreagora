package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

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

	UserInfoMe struct {
		Name             string `json:"name"`
		Email            string `json:"email"`
		AvailableCredits int    `json:"available_credits"`
		Tier             string `json:"plan_tier"`
		IsPlanActive     bool   `json:"is_plan_active"`
	}

	Subscription struct {
		ID                   int    `json:"id"`
		Email                string `json:"email"`
		Tier                 string `json:"tier"`
		StripeCostumerID     string `json:"stripe_costumer_id"`
		StripeSubscriptionID string `json:"stripe_subscription_id"`
		StripePriceID        string `json:"stripe_price_id"`
		IsActive             *bool  `json:"is_active"`
		UserID               int    `json:"user_id"`
	}

	AvailableCredits struct {
		ID     string `json:"id"`
		UserID int    `json:"user_id"`
		Total  int    `json:"total"`
	}
)

func (u *User) ValidateRequiredField() error {
	if u.Name == "" || u.Email == "" {
		return errors.New("user without required infos")
	}
	return nil
}

func (u *User) InitializeFreshUser() {
	u.PublicID = uuid.New().String()
	u.LastLogin = time.Now()
	u.AvailableCredits = AvailableCredits{Total: 0}
}
