package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Block struct {
	ID uuid.UUID `gorm:"primaryKey;type:uuid"`

	// Blocker (user who blocks)
	BlockerID uuid.UUID `gorm:"not null;index:idx_blocker_blocked,priority:1;index"`
	Blocker   User      `gorm:"foreignKey:BlockerID"`

	// Blocked (user being blocked)
	BlockedID uuid.UUID `gorm:"not null;index:idx_blocker_blocked,priority:2;index"`
	Blocked   User      `gorm:"foreignKey:BlockedID"`

	// Timestamps
	CreatedAt time.Time `gorm:"index"`
}

func (Block) TableName() string {
	return "blocks"
}

// BeforeCreate hook to generate UUID before creating block
func (b *Block) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}
