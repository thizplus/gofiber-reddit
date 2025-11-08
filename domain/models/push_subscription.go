package models

import (
	"time"

	"github.com/google/uuid"
)

type PushSubscription struct {
	ID             uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	UserID         uuid.UUID  `gorm:"type:uuid;not null;uniqueIndex:idx_user_endpoint" json:"userId"`
	User           User       `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
	Endpoint       string     `gorm:"type:text;not null;uniqueIndex:idx_user_endpoint" json:"endpoint"`
	P256dh         string     `gorm:"type:text;not null" json:"p256dh"`
	Auth           string     `gorm:"type:text;not null" json:"auth"`
	ExpirationTime *int64     `gorm:"type:bigint" json:"expirationTime,omitempty"`
	CreatedAt      time.Time  `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt      time.Time  `gorm:"autoUpdateTime" json:"updatedAt"`
}

// TableName specifies the table name for PushSubscription
func (PushSubscription) TableName() string {
	return "push_subscriptions"
}
