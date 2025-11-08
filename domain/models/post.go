package models

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID      uuid.UUID `gorm:"primaryKey;type:uuid"`
	Title   string    `gorm:"not null;type:varchar(300);index"`
	Content string    `gorm:"not null;type:text"`

	// Author
	AuthorID uuid.UUID `gorm:"not null;index"`
	Author   User      `gorm:"foreignKey:AuthorID"`

	// Stats
	Votes        int `gorm:"default:0;index"`
	CommentCount int `gorm:"default:0"`

	// Crosspost (optional)
	SourcePostID *uuid.UUID `gorm:"index"`
	SourcePost   *Post      `gorm:"foreignKey:SourcePostID"`

	// Media & Tags (relationships)
	Media []Media `gorm:"many2many:post_media;"`
	Tags  []Tag   `gorm:"many2many:post_tags;"`

	// Status
	IsDeleted bool `gorm:"default:false;index"`

	// Timestamps
	CreatedAt time.Time `gorm:"index"`
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"index"`
}

func (Post) TableName() string {
	return "posts"
}
