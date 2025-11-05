package models

import (
	"time"

	"github.com/google/uuid"
)

type Vote struct {
	UserID uuid.UUID `gorm:"primaryKey"`
	User   User      `gorm:"foreignKey:UserID"`

	TargetID   uuid.UUID `gorm:"primaryKey;index"` // post_id or comment_id
	TargetType string    `gorm:"primaryKey"`        // 'post' or 'comment'

	VoteType string `gorm:"not null"` // 'up' or 'down'

	CreatedAt time.Time
}

func (Vote) TableName() string {
	return "votes"
}
