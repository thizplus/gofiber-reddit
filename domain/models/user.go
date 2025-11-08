package models

import (
	"time"
	"github.com/google/uuid"
)

type User struct {
	// Core Fields
	ID       uuid.UUID `gorm:"primaryKey;type:uuid"`
	Email    string    `gorm:"uniqueIndex;not null"`
	Username string    `gorm:"uniqueIndex;not null"`
	Password string    // Optional for OAuth users

	// OAuth Fields
	OAuthProvider string `gorm:"index"` // google, facebook, github, etc.
	OAuthID       string `gorm:"index"` // OAuth provider's user ID
	IsOAuthUser   bool   `gorm:"default:false"`

	// Profile Fields
	DisplayName string `gorm:"not null"`
	Avatar      string
	Bio         string `gorm:"type:text"`
	Location    string
	Website     string

	// Social Stats
	Karma          int `gorm:"default:0;index"`
	FollowersCount int `gorm:"default:0"`
	FollowingCount int `gorm:"default:0"`

	// Status
	Role     string `gorm:"default:'user'"` // user, admin
	IsActive bool   `gorm:"default:true"`

	// Timestamps
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (User) TableName() string {
	return "users"
}