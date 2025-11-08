package models

import (
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	ID uuid.UUID `gorm:"primaryKey;type:uuid"`

	UserID uuid.UUID `gorm:"not null;index"` // Recipient
	User   User      `gorm:"foreignKey:UserID"`

	SenderID uuid.UUID `gorm:"not null"` // Who triggered notification
	Sender   User      `gorm:"foreignKey:SenderID"`

	Type    string `gorm:"not null;index"` // reply, vote, mention, follow
	Message string `gorm:"not null"`

	// Optional references
	PostID    *uuid.UUID
	Post      *Post `gorm:"foreignKey:PostID"`
	CommentID *uuid.UUID
	Comment   *Comment `gorm:"foreignKey:CommentID"`

	IsRead    bool      `gorm:"default:false;index"`
	CreatedAt time.Time `gorm:"index"`
}

func (Notification) TableName() string {
	return "notifications"
}
