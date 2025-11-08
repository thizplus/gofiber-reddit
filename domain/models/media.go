package models

import (
	"time"

	"github.com/google/uuid"
)

type Media struct {
	ID     uuid.UUID `gorm:"primaryKey;type:uuid"`
	UserID uuid.UUID `gorm:"not null;index"`
	User   User      `gorm:"foreignKey:UserID"`

	// File info
	Type      string `gorm:"not null;index"` // image, video, file
	FileName  string `gorm:"not null"`
	Extension string `gorm:"type:varchar(10);index"` // pdf, doc, zip, etc. (without dot)
	MimeType  string `gorm:"not null"`
	Size      int64  `gorm:"not null"` // bytes

	// URLs (Bunny CDN)
	URL       string `gorm:"not null"` // Full CDN URL
	Thumbnail string                   // Thumbnail URL

	// Dimensions
	Width  int
	Height int

	// Video specific
	Duration float64 // seconds (for videos)

	// Video streaming (Bunny Stream)
	VideoID          string `gorm:"type:varchar(255);index"`                      // Bunny Stream video ID
	HLSURL           string `gorm:"type:text"`                                    // HLS playlist URL (m3u8)
	EncodingStatus   string `gorm:"type:varchar(20);default:'pending';index"`    // pending, processing, completed, failed
	EncodingProgress int    `gorm:"default:0"`                                    // 0-100 percentage

	// Usage tracking
	Posts      []Post `gorm:"many2many:post_media;"`
	UsageCount int    `gorm:"default:0"`

	// Timestamps
	CreatedAt time.Time `gorm:"index"`
}

func (Media) TableName() string {
	return "media"
}
