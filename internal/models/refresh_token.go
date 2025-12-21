package models

import (
	"time"

	"github.com/teris-io/shortid"
	"gorm.io/gorm"
)

type RefreshToken struct {
	ID        string    `gorm:"primaryKey;size:14"`
	UserID    string    `gorm:"not null;index"`
	Token     string    `gorm:"not null;uniqueIndex"`
	ExpiresAt time.Time `gorm:"not null"`
}

func (t *RefreshToken) BeforeCreate(tx *gorm.DB) error {
	if t.ID == "" {
		t.ID = shortid.MustGenerate()
	}
	return nil
}
