package auth

import (
	"errors"

	"github.com/JoeG98/pizza-backend/internal/database"
	"github.com/JoeG98/pizza-backend/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	db *database.Database
}

func AuthService(db *database.Database) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) CreateUser(username string, password string) error {
	// check if username exists
	var existing models.User
	err := s.db.DB.Where("username = ?", username).First(&existing).Error
	if err == nil {
		return errors.New("username already exists")
	}

	// create user (password will be hashed by BeforeCreate hook)
	user := models.User{
		Username: username,
		Password: password,
	}

	return s.db.DB.Create(&user).Error
}

func (s *Service) AuthenticateUser(username string, password string) (*models.User, error) {
	var user models.User

	// Find User

	err := s.db.DB.Where("username = ?", username).First(&user).Error

	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	// compare password

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid username or password")
	}

	return &user, nil
}
