package serviceimpl

import (
	"context"
	"errors"
	"math"
	"time"

	"github.com/google/uuid"
	"gofiber-template/domain/dto"
	"gofiber-template/domain/repositories"
	"gofiber-template/domain/services"
)

type SavedPostServiceImpl struct {
	savedPostRepo repositories.SavedPostRepository
	postRepo      repositories.PostRepository
	voteRepo      repositories.VoteRepository
}

func NewSavedPostService(
	savedPostRepo repositories.SavedPostRepository,
	postRepo repositories.PostRepository,
	voteRepo repositories.VoteRepository,
) services.SavedPostService {
	return &SavedPostServiceImpl{
		savedPostRepo: savedPostRepo,
		postRepo:      postRepo,
		voteRepo:      voteRepo,
	}
}

func (s *SavedPostServiceImpl) SavePost(ctx context.Context, userID uuid.UUID, postID uuid.UUID) (*dto.SavedPostResponse, error) {
	// Check if post exists
	_, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return nil, errors.New("post not found")
	}

	// Check if already saved
	isSaved, _ := s.savedPostRepo.IsSaved(ctx, userID, postID)
	if isSaved {
		return nil, errors.New("post already saved")
	}

	// Save post
	err = s.savedPostRepo.SavePost(ctx, userID, postID)
	if err != nil {
		return nil, err
	}

	return &dto.SavedPostResponse{
		PostID:  postID,
		SavedAt: time.Now(),
	}, nil
}

func (s *SavedPostServiceImpl) UnsavePost(ctx context.Context, userID uuid.UUID, postID uuid.UUID) error {
	// Check if saved
	isSaved, _ := s.savedPostRepo.IsSaved(ctx, userID, postID)
	if !isSaved {
		return errors.New("post not saved")
	}

	// Unsave post
	return s.savedPostRepo.UnsavePost(ctx, userID, postID)
}

func (s *SavedPostServiceImpl) IsSaved(ctx context.Context, userID uuid.UUID, postID uuid.UUID) (*dto.SaveStatusResponse, error) {
	isSaved, err := s.savedPostRepo.IsSaved(ctx, userID, postID)
	if err != nil {
		return nil, err
	}

	return &dto.SaveStatusResponse{
		IsSaved: isSaved,
	}, nil
}

func (s *SavedPostServiceImpl) GetSavedPosts(ctx context.Context, userID uuid.UUID, offset, limit int) (*dto.SavedPostListResponse, error) {
	posts, err := s.savedPostRepo.GetSavedPosts(ctx, userID, offset, limit)
	if err != nil {
		return nil, err
	}

	count, err := s.savedPostRepo.CountSavedPosts(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Build post responses with user-specific data
	postResponses := make([]dto.PostResponse, len(posts))
	postIDs := make([]uuid.UUID, len(posts))
	for i, post := range posts {
		postIDs[i] = post.ID
	}

	// Batch get user votes
	voteMap, _ := s.voteRepo.GetUserVotesForTargets(ctx, userID, postIDs, "post")

	for i, post := range posts {
		resp := dto.PostToPostResponse(post)

		// Calculate hot score
		hoursSinceCreation := time.Since(post.CreatedAt).Hours()
		hotScore := float64(post.Votes) / math.Pow(hoursSinceCreation+2, 1.5)
		resp.HotScore = &hotScore

		// Add user vote
		if vote, ok := voteMap[post.ID]; ok {
			resp.UserVote = &vote.VoteType
		}

		// All posts in saved list are saved
		isSaved := true
		resp.IsSaved = &isSaved

		postResponses[i] = *resp
	}

	return &dto.SavedPostListResponse{
		Posts: postResponses,
		Meta: dto.PaginationMeta{
			Total:  count,
			Offset: offset,
			Limit:  limit,
		},
	}, nil
}

var _ services.SavedPostService = (*SavedPostServiceImpl)(nil)
