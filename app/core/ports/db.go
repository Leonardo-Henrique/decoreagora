package ports

import "github.com/Leonardo-Henrique/decoreagora/app/core/models"

type Database interface {

	/*
		USER
	*/
	GetUserByEmail(email string) (models.User, error)
	CreateUser(user models.User) (int64, error)
	CreateUserCredit(userID int64) error
	CheckIfEmailIsRegistered(email string) (bool, error)
	GetAccessCodeByUserID(userID int) (models.AccessCode, error)

	CreateAccessCode(accessCode models.AccessCode) error
	DeleteAccessCode(userID int, code string) error
}
