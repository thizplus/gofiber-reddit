package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// MessageType represents the type of message
type MessageType string

const (
	MessageTypeText  MessageType = "text"
	MessageTypeImage MessageType = "image"
	MessageTypeVideo MessageType = "video"
	MessageTypeFile  MessageType = "file"
)

// MessageMedia represents media attached to a message
type MessageMedia struct {
	URL       string  `json:"url"`
	Thumbnail *string `json:"thumbnail,omitempty"`
	Type      string  `json:"type"` // "image", "video", "file"
	Filename  *string `json:"filename,omitempty"`
	MimeType  *string `json:"mimeType,omitempty"`
	Size      *int64  `json:"size,omitempty"`
	Width     *int    `json:"width,omitempty"`
	Height    *int    `json:"height,omitempty"`
	Duration  *int    `json:"duration,omitempty"` // seconds, for videos
}

type Message struct {
	ID             uuid.UUID `gorm:"primaryKey;type:uuid"`
	ConversationID uuid.UUID `gorm:"not null;index:idx_conversation_messages"`
	Conversation   Conversation `gorm:"foreignKey:ConversationID"`

	// Sender & Receiver
	SenderID   uuid.UUID `gorm:"not null;index"`
	Sender     User      `gorm:"foreignKey:SenderID"`
	ReceiverID uuid.UUID `gorm:"not null;index"`
	Receiver   User      `gorm:"foreignKey:ReceiverID"`

	// Message Type (text, image, video, file)
	Type MessageType `gorm:"type:varchar(20);not null;default:'text';index:idx_messages_type"`

	// Content (nullable - for media-only messages)
	Content *string `gorm:"type:text;default:null"`

	// Media (JSONB array of MessageMedia)
	Media datatypes.JSON `gorm:"type:jsonb"`

	// Read Status
	IsRead bool       `gorm:"default:false;index"`
	ReadAt *time.Time

	// Timestamps (for cursor pagination)
	CreatedAt time.Time `gorm:"index:idx_conversation_messages"`
	UpdatedAt time.Time
}

func (Message) TableName() string {
	return "messages"
}

// BeforeCreate hook to generate UUID before creating message
func (m *Message) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}
