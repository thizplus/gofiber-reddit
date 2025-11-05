package repositories

import (
	"context"
	"gofiber-template/domain/models"
	"github.com/google/uuid"
)

type VoteRepository interface {
	// Create or update vote (upsert)
	Vote(ctx context.Context, vote *models.Vote) error

	// Remove vote
	Unvote(ctx context.Context, userID uuid.UUID, targetID uuid.UUID, targetType string) error

	// Get user's vote on a target
	GetVote(ctx context.Context, userID uuid.UUID, targetID uuid.UUID, targetType string) (*models.Vote, error)

	// Check if user voted
	HasVoted(ctx context.Context, userID uuid.UUID, targetID uuid.UUID, targetType string) (bool, error)

	// Get vote counts
	GetVoteCount(ctx context.Context, targetID uuid.UUID, targetType string) (upvotes int64, downvotes int64, err error)

	// Get user's votes
	ListByUser(ctx context.Context, userID uuid.UUID, targetType string, offset, limit int) ([]*models.Vote, error)

	// Get votes on target
	ListByTarget(ctx context.Context, targetID uuid.UUID, targetType string, offset, limit int) ([]*models.Vote, error)

	// Batch get user votes (for checking multiple posts/comments at once)
	GetUserVotesForTargets(ctx context.Context, userID uuid.UUID, targetIDs []uuid.UUID, targetType string) (map[uuid.UUID]*models.Vote, error)
}
