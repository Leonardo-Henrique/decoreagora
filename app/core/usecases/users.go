package usecases

import (
	"github.com/Leonardo-Henrique/decoreagora/app/core/models"
	"github.com/Leonardo-Henrique/decoreagora/app/core/ports"
)

type UserUsecase struct {
	db ports.Database
}

func NewUserUsecase(db ports.Database) *UserUsecase {
	return &UserUsecase{
		db: db,
	}
}

func (u *UserUsecase) CreateUser(user models.User) (models.User, error) {
	if err := user.ValidateRequiredField(); err != nil {
		return models.User{}, err
	}
	user.InitializeFreshUser()
	userID, err := u.db.CreateUser(user)
	if err != nil {
		return models.User{}, err
	}
	if err := u.db.CreateUserCredit(userID); err != nil {
		return models.User{}, err
	}
	user.ID = int(userID)
	return user, nil
}

func (u *UserUsecase) GetMe(userID int) (models.UserInfoMe, error) {
	return u.db.GetUserResume(userID)
}
