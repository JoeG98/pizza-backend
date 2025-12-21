package auth

import (
	"log"
	"strings"

	"github.com/JoeG98/pizza-backend/internal/database"
	"github.com/JoeG98/pizza-backend/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var DB *database.Database

func JWTMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")

	// must start with Bearer

	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "missing or invalid authorization header",
		})
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// Parse and Verify

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		// must be RSA

		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fiber.ErrUnauthorized
		}

		return PublicKey, nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid or expired token",
		})
	}

	// Extract user id

	claims := token.Claims.(jwt.MapClaims)
	userID := claims["sub"]

	var user models.User

	if err := DB.DB.First(&user, "id = ?", userID).Error; err != nil {
		log.Println(err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "user no longer exists",
		})
	}

	c.Locals("user", user)

	return c.Next()
}
