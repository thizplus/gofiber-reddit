package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Conversation struct {
	ID uuid.UUID `gorm:"primaryKey;type:uuid"`

	// Participants (ordered by UUID for consistency)
	User1ID uuid.UUID `gorm:"not null;index:idx_conversation_users"`
	User1   User      `gorm:"foreignKey:User1ID"`

	User2ID uuid.UUID `gorm:"not null;index:idx_conversation_users"`
	User2   User      `gorm:"foreignKey:User2ID"`

	// Last Message (denormalized for performance)
	// Note: No FK constraint to avoid circular dependency with Message table
	// We manually manage this relationship in application layer
	LastMessageID *uuid.UUID `gorm:"index"`
	LastMessage   *Message   `gorm:"-"` // Skip this field during migration
	LastMessageAt time.Time  `gorm:"index"`

	// Unread Counts (denormalized for performance)
	User1UnreadCount int `gorm:"default:0"`
	User2UnreadCount int `gorm:"default:0"`

	// Timestamps
	CreatedAt time.Time `gorm:"index"`
	UpdatedAt time.Time
}

func (Conversation) TableName() string {
	return "conversations"
}

// BeforeCreate hook to generate UUID before creating conversation
func (c *Conversation) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}
