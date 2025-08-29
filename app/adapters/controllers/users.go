package controllers

import (
	"errors"
	"log"
	"net/http"

	"github.com/Leonardo-Henrique/decoreagora/app/core/models"
	"github.com/Leonardo-Henrique/decoreagora/app/core/usecases"
	"github.com/Leonardo-Henrique/decoreagora/app/core/utils"
	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	userUC usecases.UserUsecase
}

func NewUserController(uc usecases.UserUsecase) *UserController {
	return &UserController{
		userUC: uc,
	}
}

func (u *UserController) NewUser(c *fiber.Ctx) error {
	log.Println("Received request to create a new user")

	var user models.User

	if err := c.BodyParser(&user); err != nil {
		log.Println("Error parsing request body for new user", err)
		return utils.ErrorResponse(c, fiber.StatusUnprocessableEntity, errors.New("invalid request body"))
	}

	createdUser, err := u.userUC.CreateUser(user)
	if err != nil {
		log.Println("Usecase failed to create user", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, errors.New("could not create user"))
	}

	return c.Status(http.StatusCreated).JSON(createdUser)
}
