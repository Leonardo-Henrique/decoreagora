package controllers

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/Leonardo-Henrique/decoreagora/app/core/config/logger"
	"github.com/Leonardo-Henrique/decoreagora/app/core/models"
	"github.com/Leonardo-Henrique/decoreagora/app/core/usecases"
	"github.com/Leonardo-Henrique/decoreagora/app/core/utils"
	"github.com/Leonardo-Henrique/decoreagora/app/core/utils/templates"
	"github.com/gofiber/fiber/v2"
)

type LoginController struct {
	loginUC usecases.LoginUsecase
	email   usecases.EmailUseCase
}

func NewLoginController(uc usecases.LoginUsecase, email usecases.EmailUseCase) *LoginController {
	return &LoginController{
		loginUC: uc,
		email:   email,
	}
}

func (lc *LoginController) Login(c *fiber.Ctx) error {
	log.Println("Received request to login user")

	var req models.LoginRequest

	if err := c.BodyParser(&req); err != nil {
		log.Println("Error parsing request body for login", err)
		return utils.ErrorResponse(c, fiber.StatusUnprocessableEntity, errors.New("invalid request body"))
	}

	code, err := lc.loginUC.Login(req)
	if err != nil {
		log.Println("Usecase failed to login user", err)
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, errors.New("invalid credentials"))
	}

	emailTemplate := templates.CodeViaEmail(code)
	if err := lc.email.SendEmail(req.Email, "Seu código de login no DecoreAgora", emailTemplate); err != nil {
		logger.Logging.Error("Error when sending auth code to email", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, errors.New("we couldnt send the email"))
	}

	return c.Status(http.StatusOK).JSON(nil)
}

func (lc *LoginController) AuthenticateCode(c *fiber.Ctx) error {
	log.Println("Received request to authenticate user code")

	var req models.AutheticateCodeRequest

	if err := c.BodyParser(&req); err != nil {
		log.Println("Error parsing request body for code authentication", err)
		return utils.ErrorResponse(c, fiber.StatusUnprocessableEntity, errors.New("invalid request body"))
	}
	fmt.Println("code received", req)

	loginResponse, err := lc.loginUC.AuthenticateCode(req)
	if err != nil {
		log.Println("Usecase failed to authenticate code", err)
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, errors.New("invalid credentials"))
	}

	return c.JSON(loginResponse)
}
