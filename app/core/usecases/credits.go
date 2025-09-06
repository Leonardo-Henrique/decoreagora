package usecases

import (
	"errors"
	"log"

	"github.com/Leonardo-Henrique/decoreagora/app/core/ports"
)

type CreditsUsecase struct {
	db ports.Database
}

func NewCreditsUsecase(db ports.Database) *CreditsUsecase {
	return &CreditsUsecase{
		db: db,
	}
}

func (c *CreditsUsecase) DecrementCredit(userID int) error {

	ok, err := c.db.AtomicDecrementCredit(userID)
	if err != nil {
		log.Println(err)
		return err
	}

	if !ok {
		// good to throw an alert
		return errors.New("user doesnt have necessary funds")
	}

	return nil

}

func (c *CreditsUsecase) GetUserCredits(userID int) (int, error) {
	return c.db.GetUserCredits(userID)
}
