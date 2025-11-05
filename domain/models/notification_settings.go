package models

import (
	"time"

	"github.com/google/uuid"
)

type NotificationSettings struct {
	UserID uuid.UUID `gorm:"primaryKey"`
	User   User      `gorm:"foreignKey:UserID"`

	Replies            bool `gorm:"default:true"`
	Mentions           bool `gorm:"default:true"`
	Votes              bool `gorm:"default:false"`
	Follows            bool `gorm:"default:true"`
	EmailNotifications bool `gorm:"default:false"`

	UpdatedAt time.Time
}

func (NotificationSettings) TableName() string {
	return "notification_settings"
}
