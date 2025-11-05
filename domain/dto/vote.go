package dto

import (
	"time"

	"github.com/google/uuid"
)

// VoteRequest - Request for voting (upvote/downvote)
type VoteRequest struct {
	TargetID   uuid.UUID `json:"targetId" validate:"required,uuid"`
	TargetType string    `json:"targetType" validate:"required,oneof=post comment"`
	VoteType   string    `json:"voteType" validate:"required,oneof=up down"`
}

// UnvoteRequest - Request for removing vote
type UnvoteRequest struct {
	TargetID   uuid.UUID `json:"targetId" validate:"required,uuid"`
	TargetType string    `json:"targetType" validate:"required,oneof=post comment"`
}

// VoteResponse - Response for vote status
type VoteResponse struct {
	TargetID   uuid.UUID `json:"targetId"`
	TargetType string    `json:"targetType"`
	VoteType   string    `json:"voteType"` // "up" or "down"
	CreatedAt  time.Time `json:"createdAt"`
}

// VoteCountResponse - Response for vote counts
type VoteCountResponse struct {
	TargetID   uuid.UUID `json:"targetId"`
	TargetType string    `json:"targetType"`
	Upvotes    int64     `json:"upvotes"`
	Downvotes  int64     `json:"downvotes"`
	Total      int       `json:"total"` // upvotes - downvotes
}
