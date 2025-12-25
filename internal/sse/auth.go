package sse

import (
	"errors"
	"strings"

	"github.com/JoeG98/pizza-backend/internal/auth"
	"github.com/JoeG98/pizza-backend/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// authenticateAdmin validates JWT and enforces admin role
func authenticateAdmin(c *fiber.Ctx) (*models.User, error) {
	user, err := authenticateUser(c)
	if err != nil {
		return nil, err
	}

	if user.Role != "admin" {
		return nil, errors.New("admin access required")
	}

	return user, nil
}

// authenticateUser validates JWT and returns the user
func authenticateUser(c *fiber.Ctx) (*models.User, error) {
	authHeader := c.Get("Authorization")

	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return nil, errors.New("missing or invalid authorization header")
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// Parse & verify JWT (same logic as JWTMiddleware)
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("invalid signing method")
		}
		return auth.PublicKey, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid or expired token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	userID, ok := claims["sub"]
	if !ok {
		return nil, errors.New("missing subject claim")
	}

	var user models.User
	if err := auth.DB.DB.First(&user, "id = ?", userID).Error; err != nil {
		return nil, errors.New("user not found")
	}

	return &user, nil
}
