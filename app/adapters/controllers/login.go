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

type LoginController struct {
	loginUC usecases.LoginUsecase
}

func NewLoginController(uc usecases.LoginUsecase) *LoginController {
	return &LoginController{
		loginUC: uc,
	}
}

func (lc *LoginController) Login(c *fiber.Ctx) error {
	log.Println("Received request to login user")

	var req models.LoginRequest

	if err := c.BodyParser(&req); err != nil {
		log.Println("Error parsing request body for login", err)
		return utils.ErrorResponse(c, fiber.StatusUnprocessableEntity, errors.New("invalid request body"))
	}

	if err := lc.loginUC.Login(req); err != nil {
		log.Println("Usecase failed to login user", err)
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, errors.New("invalid credentials"))
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

	loginResponse, err := lc.loginUC.AuthenticateCode(req)
	if err != nil {
		log.Println("Usecase failed to authenticate code", err)
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, errors.New("invalid credentials"))
	}

	return c.JSON(loginResponse)
}
