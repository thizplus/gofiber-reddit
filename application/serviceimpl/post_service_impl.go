package serviceimpl

import (
	"context"
	"errors"
	"math"
	"time"

	"github.com/google/uuid"
	"gofiber-template/domain/dto"
	"gofiber-template/domain/models"
	"gofiber-template/domain/repositories"
	"gofiber-template/domain/services"
)

type PostServiceImpl struct {
	postRepo      repositories.PostRepository
	userRepo      repositories.UserRepository
	voteRepo      repositories.VoteRepository
	savedPostRepo repositories.SavedPostRepository
	tagService    services.TagService
	mediaRepo     repositories.MediaRepository
}

func NewPostService(
	postRepo repositories.PostRepository,
	userRepo repositories.UserRepository,
	voteRepo repositories.VoteRepository,
	savedPostRepo repositories.SavedPostRepository,
	tagService services.TagService,
	mediaRepo repositories.MediaRepository,
) services.PostService {
	return &PostServiceImpl{
		postRepo:      postRepo,
		userRepo:      userRepo,
		voteRepo:      voteRepo,
		savedPostRepo: savedPostRepo,
		tagService:    tagService,
		mediaRepo:     mediaRepo,
	}
}

func (s *PostServiceImpl) CreatePost(ctx context.Context, userID uuid.UUID, req *dto.CreatePostRequest) (*dto.PostResponse, error) {
	// Create post
	post := &models.Post{
		ID:           uuid.New(),
		Title:        req.Title,
		Content:      req.Content,
		AuthorID:     userID,
		Votes:        0,
		CommentCount: 0,
		IsDeleted:    false,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Handle source post for crosspost
	if req.SourcePostID != nil {
		post.SourcePostID = req.SourcePostID
	}

	// Create post in database
	err := s.postRepo.Create(ctx, post)
	if err != nil {
		return nil, err
	}

	// Handle tags
	if len(req.Tags) > 0 {
		tagIDs, err := s.tagService.GetOrCreateTags(ctx, req.Tags)
		if err != nil {
			return nil, err
		}
		err = s.postRepo.AttachTags(ctx, post.ID, tagIDs)
		if err != nil {
			return nil, err
		}
	}

	// Handle media attachments
	if len(req.MediaIDs) > 0 {
		err = s.postRepo.AttachMedia(ctx, post.ID, req.MediaIDs)
		if err != nil {
			return nil, err
		}
		// Increment usage count for each media
		for _, mediaID := range req.MediaIDs {
			_ = s.mediaRepo.IncrementUsageCount(ctx, mediaID)
		}
	}

	// Get full post with relations
	return s.GetPost(ctx, post.ID, &userID)
}

func (s *PostServiceImpl) GetPost(ctx context.Context, postID uuid.UUID, userID *uuid.UUID) (*dto.PostResponse, error) {
	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return nil, err
	}

	resp := dto.PostToPostResponse(post)

	// Add user-specific data if authenticated
	if userID != nil {
		// Get user's vote
		vote, _ := s.voteRepo.GetVote(ctx, *userID, postID, "post")
		if vote != nil {
			resp.UserVote = &vote.VoteType
		}

		// Check if saved
		isSaved, _ := s.savedPostRepo.IsSaved(ctx, *userID, postID)
		resp.IsSaved = &isSaved
	}

	return resp, nil
}

func (s *PostServiceImpl) UpdatePost(ctx context.Context, postID uuid.UUID, userID uuid.UUID, req *dto.UpdatePostRequest) (*dto.PostResponse, error) {
	// Get existing post
	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return nil, err
	}

	// Check ownership
	if post.AuthorID != userID {
		return nil, errors.New("unauthorized: not post owner")
	}

	// Update fields
	if req.Title != "" {
		post.Title = req.Title
	}
	if req.Content != "" {
		post.Content = req.Content
	}
	post.UpdatedAt = time.Now()

	err = s.postRepo.Update(ctx, postID, post)
	if err != nil {
		return nil, err
	}

	// Update tags if provided
	if len(req.Tags) > 0 {
		tagIDs, err := s.tagService.GetOrCreateTags(ctx, req.Tags)
		if err != nil {
			return nil, err
		}
		err = s.postRepo.SyncTags(ctx, postID, tagIDs)
		if err != nil {
			return nil, err
		}
	}

	return s.GetPost(ctx, postID, &userID)
}

func (s *PostServiceImpl) DeletePost(ctx context.Context, postID uuid.UUID, userID uuid.UUID) error {
	// Get post
	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return err
	}

	// Check ownership
	if post.AuthorID != userID {
		return errors.New("unauthorized: not post owner")
	}

	// Soft delete
	return s.postRepo.Delete(ctx, postID)
}

func (s *PostServiceImpl) ListPosts(ctx context.Context, offset, limit int, sortBy repositories.PostSortBy, userID *uuid.UUID) (*dto.PostListResponse, error) {
	posts, err := s.postRepo.List(ctx, offset, limit, sortBy)
	if err != nil {
		return nil, err
	}

	count, err := s.postRepo.Count(ctx)
	if err != nil {
		return nil, err
	}

	return s.buildPostListResponse(ctx, posts, count, offset, limit, userID)
}

func (s *PostServiceImpl) ListPostsByAuthor(ctx context.Context, authorID uuid.UUID, offset, limit int, userID *uuid.UUID) (*dto.PostListResponse, error) {
	posts, err := s.postRepo.ListByAuthor(ctx, authorID, offset, limit)
	if err != nil {
		return nil, err
	}

	count, err := s.postRepo.CountByAuthor(ctx, authorID)
	if err != nil {
		return nil, err
	}

	return s.buildPostListResponse(ctx, posts, count, offset, limit, userID)
}

func (s *PostServiceImpl) ListPostsByTag(ctx context.Context, tagName string, offset, limit int, sortBy repositories.PostSortBy, userID *uuid.UUID) (*dto.PostListResponse, error) {
	posts, err := s.postRepo.ListByTag(ctx, tagName, offset, limit, sortBy)
	if err != nil {
		return nil, err
	}

	// For tag-filtered lists, we'll use a simplified count
	count := int64(len(posts))
	if len(posts) == limit {
		count = int64(offset + limit + 1) // Approximate
	}

	return s.buildPostListResponse(ctx, posts, count, offset, limit, userID)
}

func (s *PostServiceImpl) ListPostsByTagID(ctx context.Context, tagID uuid.UUID, offset, limit int, sortBy repositories.PostSortBy, userID *uuid.UUID) (*dto.PostListResponse, error) {
	posts, err := s.postRepo.ListByTagID(ctx, tagID, offset, limit, sortBy)
	if err != nil {
		return nil, err
	}

	// For tag-filtered lists, we'll use a simplified count
	count := int64(len(posts))
	if len(posts) == limit {
		count = int64(offset + limit + 1) // Approximate
	}

	return s.buildPostListResponse(ctx, posts, count, offset, limit, userID)
}

func (s *PostServiceImpl) SearchPosts(ctx context.Context, query string, offset, limit int, userID *uuid.UUID) (*dto.PostListResponse, error) {
	posts, err := s.postRepo.Search(ctx, query, offset, limit)
	if err != nil {
		return nil, err
	}

	count := int64(len(posts))
	if len(posts) == limit {
		count = int64(offset + limit + 1)
	}

	return s.buildPostListResponse(ctx, posts, count, offset, limit, userID)
}

func (s *PostServiceImpl) CreateCrosspost(ctx context.Context, userID uuid.UUID, sourcePostID uuid.UUID, req *dto.CreatePostRequest) (*dto.PostResponse, error) {
	// Verify source post exists
	_, err := s.postRepo.GetByID(ctx, sourcePostID)
	if err != nil {
		return nil, errors.New("source post not found")
	}

	// Set source post ID
	req.SourcePostID = &sourcePostID

	return s.CreatePost(ctx, userID, req)
}

func (s *PostServiceImpl) GetCrossposts(ctx context.Context, postID uuid.UUID, offset, limit int, userID *uuid.UUID) (*dto.PostListResponse, error) {
	posts, err := s.postRepo.GetCrossposts(ctx, postID, offset, limit)
	if err != nil {
		return nil, err
	}

	count := int64(len(posts))
	return s.buildPostListResponse(ctx, posts, count, offset, limit, userID)
}

func (s *PostServiceImpl) GetFeed(ctx context.Context, userID uuid.UUID, offset, limit int, sortBy repositories.PostSortBy) (*dto.PostFeedResponse, error) {
	// For now, just return all posts sorted by the requested method
	// TODO: Implement personalized feed based on followed users
	posts, err := s.postRepo.List(ctx, offset, limit, sortBy)
	if err != nil {
		return nil, err
	}

	count, err := s.postRepo.Count(ctx)
	if err != nil {
		return nil, err
	}

	listResp, err := s.buildPostListResponse(ctx, posts, count, offset, limit, &userID)
	if err != nil {
		return nil, err
	}

	return &dto.PostFeedResponse{
		Posts: listResp.Posts,
		Meta:  listResp.Meta,
	}, nil
}

// Helper function to build post list response with user-specific data
func (s *PostServiceImpl) buildPostListResponse(ctx context.Context, posts []*models.Post, count int64, offset, limit int, userID *uuid.UUID) (*dto.PostListResponse, error) {
	responses := make([]dto.PostResponse, len(posts))

	// Collect all post IDs for batch operations
	postIDs := make([]uuid.UUID, len(posts))
	for i, post := range posts {
		postIDs[i] = post.ID
	}

	// Batch get user votes and saved status if authenticated
	var voteMap map[uuid.UUID]*models.Vote
	var savedMap map[uuid.UUID]bool
	if userID != nil {
		voteMap, _ = s.voteRepo.GetUserVotesForTargets(ctx, *userID, postIDs, "post")
		savedMap, _ = s.savedPostRepo.GetSavedStatus(ctx, *userID, postIDs)
	}

	// Build responses
	for i, post := range posts {
		resp := dto.PostToPostResponse(post)

		// Calculate hot score
		hoursSinceCreation := time.Since(post.CreatedAt).Hours()
		hotScore := float64(post.Votes) / math.Pow(hoursSinceCreation+2, 1.5)
		resp.HotScore = &hotScore

		// Add user-specific data
		if userID != nil {
			if vote, ok := voteMap[post.ID]; ok {
				resp.UserVote = &vote.VoteType
			}
			// Always set isSaved for authenticated users (true or false)
			isSaved := false
			if saved, ok := savedMap[post.ID]; ok {
				isSaved = saved
			}
			resp.IsSaved = &isSaved
		}

		responses[i] = *resp
	}

	return &dto.PostListResponse{
		Posts: responses,
		Meta: dto.PaginationMeta{
			Total:  count,
			Offset: offset,
			Limit:  limit,
		},
	}, nil
}

var _ services.PostService = (*PostServiceImpl)(nil)
