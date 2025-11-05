package models

import (
	"time"

	"github.com/google/uuid"
)

type SavedPost struct {
	UserID uuid.UUID `gorm:"primaryKey"`
	User   User      `gorm:"foreignKey:UserID"`

	PostID uuid.UUID `gorm:"primaryKey"`
	Post   Post      `gorm:"foreignKey:PostID"`

	SavedAt time.Time `gorm:"index"`
}

func (SavedPost) TableName() string {
	return "saved_posts"
}
