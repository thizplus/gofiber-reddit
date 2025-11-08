package models

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	ID     uuid.UUID `gorm:"primaryKey;type:uuid"`
	PostID uuid.UUID `gorm:"not null;index"`
	Post   Post      `gorm:"foreignKey:PostID"`

	AuthorID uuid.UUID `gorm:"not null;index"`
	Author   User      `gorm:"foreignKey:AuthorID"`

	Content string `gorm:"not null;type:text"`
	Votes   int    `gorm:"default:0;index"`

	// Nested replies
	ParentID *uuid.UUID `gorm:"index"`
	Parent   *Comment   `gorm:"foreignKey:ParentID"`
	Replies  []Comment  `gorm:"foreignKey:ParentID"`
	Depth    int        `gorm:"default:0;index"` // 0 = top-level, max 10

	// Status
	IsDeleted bool `gorm:"default:false"`

	// Timestamps
	CreatedAt time.Time `gorm:"index"`
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func (Comment) TableName() string {
	return "comments"
}
