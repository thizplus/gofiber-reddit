package models

import (
	"time"

	"github.com/google/uuid"
)

type Media struct {
	ID     uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID uuid.UUID `gorm:"not null;index"`
	User   User      `gorm:"foreignKey:UserID"`

	// File info
	Type     string `gorm:"not null;index"` // image, video
	FileName string `gorm:"not null"`
	MimeType string `gorm:"not null"`
	Size     int64  `gorm:"not null"` // bytes

	// URLs (Bunny CDN)
	URL       string `gorm:"not null"` // Full CDN URL
	Thumbnail string                   // Thumbnail URL

	// Dimensions
	Width  int
	Height int

	// Video specific
	Duration float64 // seconds (for videos)

	// Usage tracking
	Posts      []Post `gorm:"many2many:post_media;"`
	UsageCount int    `gorm:"default:0"`

	// Timestamps
	CreatedAt time.Time `gorm:"index"`
}

func (Media) TableName() string {
	return "media"
}
