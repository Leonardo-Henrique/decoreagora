package ports

import "github.com/Leonardo-Henrique/decoreagora/app/core/models"

type Database interface {

	/*
		USER
	*/
	CreateUser(user models.User) (int64, error)
	CreateUserCredit(userID int64) error
}
