package auth

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/JoeG98/pizza-backend/internal/models"
)

func CreateRefreshToken(userID string) (string, error) {

	// random 32 bytes
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	token := base64.StdEncoding.EncodeToString(b)

	rt := models.RefreshToken{
		UserID:    userID,
		Token:     token,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 30), // 30 days
	}

	err = DB.DB.Create(&rt).Error
	if err != nil {
		return "", err
	}

	return token, nil
}
