package services

import (
	"context"
	"gofiber-template/domain/dto"
	"github.com/google/uuid"
)

type VoteService interface {
	// Vote/Unvote
	Vote(ctx context.Context, userID uuid.UUID, req *dto.VoteRequest) (*dto.VoteResponse, error)
	Unvote(ctx context.Context, userID uuid.UUID, req *dto.UnvoteRequest) error

	// Get vote status
	GetVote(ctx context.Context, userID uuid.UUID, targetID uuid.UUID, targetType string) (*dto.VoteResponse, error)
	GetVoteCount(ctx context.Context, targetID uuid.UUID, targetType string) (*dto.VoteCountResponse, error)

	// Get user votes
	GetUserVotes(ctx context.Context, userID uuid.UUID, targetType string, offset, limit int) ([]*dto.VoteResponse, error)
}
