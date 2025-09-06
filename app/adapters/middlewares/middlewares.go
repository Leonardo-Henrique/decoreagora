package middlewares

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/Leonardo-Henrique/decoreagora/app/core/ports"
	"github.com/Leonardo-Henrique/decoreagora/app/core/utils"
	"github.com/gofiber/fiber/v2"
)

type Middleware struct {
	db ports.Database
	th ports.TokenHandler
}

func NewMiddleware(db ports.Database, th ports.TokenHandler) *Middleware {
	return &Middleware{
		db: db,
		th: th,
	}
}

func (m *Middleware) AuthMiddleware(next fiber.Handler) fiber.Handler {
	return func(c *fiber.Ctx) error {

		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "no token was passed",
			})
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(http.StatusForbidden).JSON(fiber.Map{
				"error": "the provided token is in a wrong form",
			})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := m.th.ValidateToken(tokenString)
		if err != nil {
			log.Println(err)
			return c.Status(http.StatusForbidden).JSON(fiber.Map{
				"error": "the token is invalid",
			})
		}

		c.Locals("userClaims", claims)
		c.Locals("userID", claims.UserID)
		return next(c)
	}
}

func (m *Middleware) CreditsMiddleware(next fiber.Handler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := utils.GetCurrentUserID(c)
		if userID == 0 {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "we couldnt identify the user"})
		}

		fmt.Println("user from midd", userID)

		qtdCredits, err := m.db.GetUserCredits(userID)
		if err != nil {
			log.Println(err)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "we couldnt calculate available credits"})
		}

		if qtdCredits < 1 {
			return c.Status(http.StatusForbidden).JSON(fiber.Map{"error": "user doesnt have available credits"})
		}

		return next(c)
	}
}
