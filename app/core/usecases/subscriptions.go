package usecases

import (
	"errors"

	"github.com/Leonardo-Henrique/decoreagora/app/core/ports"
)

type SubscriptionUsecase struct {
	db ports.Database
}

func NewSubscriptionUsecase(db ports.Database) *SubscriptionUsecase {
	return &SubscriptionUsecase{
		db: db,
	}
}

func (s *SubscriptionUsecase) CreateNewSubscription(userID int, tier string, isActive bool, email string) error {
	if userID == 0 || tier == "" {
		return errors.New("required info for subs not found")
	}
	return s.db.CreateNewSubscription(userID, tier, isActive, email)
}
