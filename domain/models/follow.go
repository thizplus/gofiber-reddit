package models

import (
	"time"

	"github.com/google/uuid"
)

type Follow struct {
	FollowerID uuid.UUID `gorm:"primaryKey;index"`
	Follower   User      `gorm:"foreignKey:FollowerID"`

	FollowingID uuid.UUID `gorm:"primaryKey;index"`
	Following   User      `gorm:"foreignKey:FollowingID"`

	CreatedAt time.Time `gorm:"index"`
}

func (Follow) TableName() string {
	return "follows"
}
