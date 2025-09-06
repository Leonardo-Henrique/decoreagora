package ports

import (
	"time"

	"github.com/Leonardo-Henrique/decoreagora/app/core/models"
)

type Database interface {

	/*
		USER
	*/
	GetUserByEmail(email string) (models.User, error)
	CreateUser(user models.User) (int64, error)
	CreateUserCredit(userID int64) error
	CheckIfEmailIsRegistered(email string) (bool, error)
	GetAccessCodeByUserID(userID int) ([]models.AccessCode, error)
	GetUserCredits(userID int) (int, error)
	GetUserResume(userID int) (models.UserInfoMe, error)

	CreateAccessCode(accessCode models.AccessCode) error
	DeleteAccessCode(userID int, code string) error

	GetUserImages(userID int) ([]models.EditedImageResponse, error)
	CreateNewImageEntry(publicID, imageKey, prompt_description string, userID int, date time.Time) (int, error)
	FinishImageEdition(generatedFileBucketKey, originalFilePublicKey string) error

	GetCurrentCredits(userID int) (int, error)
	AtomicDecrementCredit(userID int) (bool, error)

	CreateNewSubscription(userID int, tier string, isActive bool, email string) error
	UpdateUserCustomerID(userID int, customerID string) error
	GetSubscriptionByEmail(email string) (models.Subscription, error)
	IncrementUserCreditsByCustomerID(userID, qtd int) error
	GetSubscriptionByCustomerID(customerID string) (models.Subscription, error)
	UpdateUserTier(customerID, tier string) error
	CreatePaymentHistoryEntry(paymentEntry models.PaymentHistory) error
}
