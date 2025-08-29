package usecases

import "github.com/Leonardo-Henrique/decoreagora/app/core/models"

type UserUsecase struct{}

func NewUserUsecase() *UserUsecase {
	return &UserUsecase{}
}

func (u *UserUsecase) CreateUser(user models.User) (models.User, error) {
	return models.User{}, nil
}
