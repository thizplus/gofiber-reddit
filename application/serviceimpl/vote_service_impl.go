package serviceimpl

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gofiber-template/domain/dto"
	"gofiber-template/domain/models"
	"gofiber-template/domain/repositories"
	"gofiber-template/domain/services"
)

type VoteServiceImpl struct {
	voteRepo     repositories.VoteRepository
	postRepo     repositories.PostRepository
	commentRepo  repositories.CommentRepository
	userRepo     repositories.UserRepository
	notifService services.NotificationService
}

func NewVoteService(
	voteRepo repositories.VoteRepository,
	postRepo repositories.PostRepository,
	commentRepo repositories.CommentRepository,
	userRepo repositories.UserRepository,
	notifService services.NotificationService,
) services.VoteService {
	return &VoteServiceImpl{
		voteRepo:     voteRepo,
		postRepo:     postRepo,
		commentRepo:  commentRepo,
		userRepo:     userRepo,
		notifService: notifService,
	}
}

func (s *VoteServiceImpl) Vote(ctx context.Context, userID uuid.UUID, req *dto.VoteRequest) (*dto.VoteResponse, error) {
	// Check if user already voted
	existingVote, _ := s.voteRepo.GetVote(ctx, userID, req.TargetID, req.TargetType)

	vote := &models.Vote{
		UserID:     userID,
		TargetID:   req.TargetID,
		TargetType: req.TargetType,
		VoteType:   req.VoteType,
		CreatedAt:  time.Now(),
	}

	// Calculate vote change
	voteChange := 0
	if existingVote == nil {
		// New vote
		if req.VoteType == "up" {
			voteChange = 1
		} else {
			voteChange = -1
		}
	} else if existingVote.VoteType != req.VoteType {
		// Changing vote
		if req.VoteType == "up" {
			voteChange = 2 // from -1 to +1
		} else {
			voteChange = -2 // from +1 to -1
		}
	}

	// Save vote (upsert)
	err := s.voteRepo.Vote(ctx, vote)
	if err != nil {
		return nil, err
	}

	// Update vote count on target
	if voteChange != 0 {
		if req.TargetType == "post" {
			_ = s.postRepo.UpdateVoteCount(ctx, req.TargetID, voteChange)

			// Send notification to post author (only for upvotes, and only if new vote)
			if req.VoteType == "up" && existingVote == nil {
				post, _ := s.postRepo.GetByID(ctx, req.TargetID)
				if post != nil && post.AuthorID != userID {
					_ = s.notifService.CreateNotification(
						ctx,
						post.AuthorID,
						userID,
						"vote",
						"ถูกใจโพสต์ของคุณ",
						&req.TargetID,
						nil,
					)
				}
			}
		} else if req.TargetType == "comment" {
			_ = s.commentRepo.UpdateVoteCount(ctx, req.TargetID, voteChange)

			// Send notification to comment author (only for upvotes, and only if new vote)
			if req.VoteType == "up" && existingVote == nil {
				comment, _ := s.commentRepo.GetByID(ctx, req.TargetID)
				if comment != nil && comment.AuthorID != userID {
					_ = s.notifService.CreateNotification(
						ctx,
						comment.AuthorID,
						userID,
						"vote",
						"ถูกใจความคิดเห็นของคุณ",
						&comment.PostID,
						&req.TargetID,
					)
				}
			}
		}
	}

	return &dto.VoteResponse{
		TargetID:   req.TargetID,
		TargetType: req.TargetType,
		VoteType:   req.VoteType,
		CreatedAt:  vote.CreatedAt,
	}, nil
}

func (s *VoteServiceImpl) Unvote(ctx context.Context, userID uuid.UUID, req *dto.UnvoteRequest) error {
	// Get existing vote
	existingVote, err := s.voteRepo.GetVote(ctx, userID, req.TargetID, req.TargetType)
	if err != nil || existingVote == nil {
		return nil // Already not voted
	}

	// Calculate vote change
	voteChange := 0
	if existingVote.VoteType == "up" {
		voteChange = -1
	} else {
		voteChange = 1
	}

	// Remove vote
	err = s.voteRepo.Unvote(ctx, userID, req.TargetID, req.TargetType)
	if err != nil {
		return err
	}

	// Update vote count on target
	if req.TargetType == "post" {
		_ = s.postRepo.UpdateVoteCount(ctx, req.TargetID, voteChange)
	} else if req.TargetType == "comment" {
		_ = s.commentRepo.UpdateVoteCount(ctx, req.TargetID, voteChange)
	}

	return nil
}

func (s *VoteServiceImpl) GetVote(ctx context.Context, userID uuid.UUID, targetID uuid.UUID, targetType string) (*dto.VoteResponse, error) {
	vote, err := s.voteRepo.GetVote(ctx, userID, targetID, targetType)
	if err != nil {
		return nil, err
	}

	if vote == nil {
		return nil, nil
	}

	return dto.VoteToVoteResponse(vote), nil
}

func (s *VoteServiceImpl) GetVoteCount(ctx context.Context, targetID uuid.UUID, targetType string) (*dto.VoteCountResponse, error) {
	upvotes, downvotes, err := s.voteRepo.GetVoteCount(ctx, targetID, targetType)
	if err != nil {
		return nil, err
	}

	return &dto.VoteCountResponse{
		TargetID:   targetID,
		TargetType: targetType,
		Upvotes:    upvotes,
		Downvotes:  downvotes,
		Total:      int(upvotes - downvotes),
	}, nil
}

func (s *VoteServiceImpl) GetUserVotes(ctx context.Context, userID uuid.UUID, targetType string, offset, limit int) ([]*dto.VoteResponse, error) {
	votes, err := s.voteRepo.ListByUser(ctx, userID, targetType, offset, limit)
	if err != nil {
		return nil, err
	}

	responses := make([]*dto.VoteResponse, len(votes))
	for i, vote := range votes {
		responses[i] = dto.VoteToVoteResponse(vote)
	}

	return responses, nil
}

var _ services.VoteService = (*VoteServiceImpl)(nil)
