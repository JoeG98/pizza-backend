package models

import (
	"github.com/teris-io/shortid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID       string `gorm:"primaryKey;size:14" json:"id"`
	Username string `gorm:"uniqueIndex;not null" json:"username"`
	Password string `gorm:"not null" json:"-"`
	Role     string `gorm:"not null;default:customer" json:"role"`
}

// Before Saving new user -> Hash Passowrd

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	// generate ID
	if u.ID == "" {
		u.ID = shortid.MustGenerate()
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	u.Password = string(hashed)
	return nil
}
