package middlewares

import (
	"log"
	"net/http"
	"strings"

	"github.com/Leonardo-Henrique/decoreagora/app/core/ports"
	"github.com/gofiber/fiber/v2"
)

type Middleware struct {
	th ports.TokenHandler
}

func NewMiddleware(th ports.TokenHandler) *Middleware {
	return &Middleware{
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
