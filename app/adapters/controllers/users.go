package controllers

import (
	"errors"
	"log"
	"net/http"

	"github.com/Leonardo-Henrique/decoreagora/app/core/config/logger"
	"github.com/Leonardo-Henrique/decoreagora/app/core/models"
	"github.com/Leonardo-Henrique/decoreagora/app/core/usecases"
	"github.com/Leonardo-Henrique/decoreagora/app/core/utils"
	"github.com/Leonardo-Henrique/decoreagora/app/core/utils/templates"
	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	userUC  usecases.UserUsecase
	loginUC usecases.LoginUsecase
	subsUC  usecases.SubscriptionUsecase
	email   usecases.EmailUseCase
}

func NewUserController(loginUC usecases.LoginUsecase, uc usecases.UserUsecase, subsUC usecases.SubscriptionUsecase, email usecases.EmailUseCase) *UserController {
	return &UserController{
		loginUC: loginUC,
		userUC:  uc,
		subsUC:  subsUC,
		email:   email,
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

	log.Println("created user", createdUser)

	if err := u.subsUC.CreateNewSubscription(createdUser.ID, "free", true, user.Email); err != nil {
		log.Println("Usecase failed to create user subscription", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, errors.New("could not create user sub"))
	}

	code, err := u.loginUC.Login(models.LoginRequest{
		Email: user.Email,
	})
	if err != nil {
		log.Println("Usecase failed to create user", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, errors.New("could not create user"))
	}

	emailTemplate := templates.CodeViaEmail(code)
	if err := u.email.SendEmail(user.Email, "Seu código de login no DecoreAgora", emailTemplate); err != nil {
		logger.Logging.Error("Error when sending auth code to email", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, errors.New("we couldnt send the email"))
	}

	//TODO
	// make a new model without unnecessary infos
	createdUser.ID = 0

	return c.Status(http.StatusCreated).JSON(createdUser)
}

func (u *UserController) GetMe(c *fiber.Ctx) error {
	userID := utils.GetCurrentUserID(c)
	if userID == 0 {
		log.Println("Usecase failed to get /me")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, errors.New("no user to return info"))
	}

	userData, err := u.userUC.GetMe(userID)
	if err != nil {
		log.Println("Usecase failed to get /me", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, errors.New("no user to return info"))
	}

	return c.Status(http.StatusOK).JSON(userData)
}
