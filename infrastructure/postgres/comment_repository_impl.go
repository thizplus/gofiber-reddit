package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gofiber-template/domain/models"
	"gofiber-template/domain/repositories"
)

type CommentRepositoryImpl struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) repositories.CommentRepository {
	return &CommentRepositoryImpl{db: db}
}

func (r *CommentRepositoryImpl) Create(ctx context.Context, comment *models.Comment) error {
	return r.db.WithContext(ctx).Create(comment).Error
}

func (r *CommentRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*models.Comment, error) {
	var comment models.Comment
	err := r.db.WithContext(ctx).
		Preload("Author").
		Preload("Post").
		Where("id = ? AND is_deleted = ?", id, false).
		First(&comment).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

func (r *CommentRepositoryImpl) Update(ctx context.Context, id uuid.UUID, comment *models.Comment) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Updates(comment).Error
}

func (r *CommentRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&models.Comment{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_deleted": true,
			"deleted_at": now,
		}).Error
}

func (r *CommentRepositoryImpl) ListByPost(ctx context.Context, postID uuid.UUID, offset, limit int, sortBy repositories.CommentSortBy) ([]*models.Comment, error) {
	var comments []*models.Comment
	query := r.db.WithContext(ctx).
		Preload("Author").
		Where("post_id = ? AND parent_id IS NULL AND is_deleted = ?", postID, false)

	switch sortBy {
	case repositories.CommentSortByHot:
		query = query.Order(r.hotScoreSQL() + " DESC")
	case repositories.CommentSortByNew:
		query = query.Order("created_at DESC")
	case repositories.CommentSortByTop:
		query = query.Order("votes DESC")
	case repositories.CommentSortByOld:
		query = query.Order("created_at ASC")
	default:
		query = query.Order("created_at DESC")
	}

	err := query.Offset(offset).Limit(limit).Find(&comments).Error
	return comments, err
}

func (r *CommentRepositoryImpl) ListByAuthor(ctx context.Context, authorID uuid.UUID, offset, limit int) ([]*models.Comment, error) {
	var comments []*models.Comment
	err := r.db.WithContext(ctx).
		Preload("Author").
		Preload("Post").
		Preload("Post.Author").
		Where("author_id = ? AND is_deleted = ?", authorID, false).
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&comments).Error
	return comments, err
}

func (r *CommentRepositoryImpl) ListReplies(ctx context.Context, parentID uuid.UUID, offset, limit int, sortBy repositories.CommentSortBy) ([]*models.Comment, error) {
	var comments []*models.Comment
	query := r.db.WithContext(ctx).
		Preload("Author").
		Where("parent_id = ? AND is_deleted = ?", parentID, false)

	switch sortBy {
	case repositories.CommentSortByHot:
		query = query.Order(r.hotScoreSQL() + " DESC")
	case repositories.CommentSortByNew:
		query = query.Order("created_at DESC")
	case repositories.CommentSortByTop:
		query = query.Order("votes DESC")
	case repositories.CommentSortByOld:
		query = query.Order("created_at ASC")
	default:
		query = query.Order("created_at DESC")
	}

	err := query.Offset(offset).Limit(limit).Find(&comments).Error
	return comments, err
}

func (r *CommentRepositoryImpl) GetCommentTree(ctx context.Context, postID uuid.UUID, maxDepth int) ([]*models.Comment, error) {
	var comments []*models.Comment
	// Get all comments for the post up to maxDepth
	err := r.db.WithContext(ctx).
		Preload("Author").
		Where("post_id = ? AND is_deleted = ? AND depth <= ?", postID, false, maxDepth).
		Order("depth ASC, created_at ASC").
		Find(&comments).Error
	return comments, err
}

func (r *CommentRepositoryImpl) GetParentChain(ctx context.Context, commentID uuid.UUID) ([]*models.Comment, error) {
	var comments []*models.Comment
	var currentComment *models.Comment

	// Start with the current comment
	err := r.db.WithContext(ctx).
		Preload("Author").
		Where("id = ?", commentID).
		First(&currentComment).Error
	if err != nil {
		return nil, err
	}

	comments = append(comments, currentComment)

	// Traverse up the parent chain
	for currentComment.ParentID != nil {
		var parent models.Comment
		err := r.db.WithContext(ctx).
			Preload("Author").
			Where("id = ?", *currentComment.ParentID).
			First(&parent).Error
		if err != nil {
			break
		}
		comments = append([]*models.Comment{&parent}, comments...) // Prepend parent
		currentComment = &parent
	}

	return comments, nil
}

func (r *CommentRepositoryImpl) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Comment{}).Where("is_deleted = ?", false).Count(&count).Error
	return count, err
}

func (r *CommentRepositoryImpl) CountByPost(ctx context.Context, postID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Comment{}).
		Where("post_id = ? AND parent_id IS NULL AND is_deleted = ?", postID, false).
		Count(&count).Error
	return count, err
}

func (r *CommentRepositoryImpl) CountByAuthor(ctx context.Context, authorID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Comment{}).
		Where("author_id = ? AND is_deleted = ?", authorID, false).
		Count(&count).Error
	return count, err
}

func (r *CommentRepositoryImpl) CountReplies(ctx context.Context, parentID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Comment{}).
		Where("parent_id = ? AND is_deleted = ?", parentID, false).
		Count(&count).Error
	return count, err
}

func (r *CommentRepositoryImpl) UpdateVoteCount(ctx context.Context, commentID uuid.UUID, voteChange int) error {
	return r.db.WithContext(ctx).
		Model(&models.Comment{}).
		Where("id = ?", commentID).
		UpdateColumn("votes", gorm.Expr("votes + ?", voteChange)).Error
}

// hotScoreSQL generates SQL for hot score calculation: votes / (hours + 2)^1.5
func (r *CommentRepositoryImpl) hotScoreSQL() string {
	return fmt.Sprintf(
		"votes / POWER((EXTRACT(EPOCH FROM (NOW() - created_at)) / 3600.0) + 2, %.1f)",
		1.5,
	)
}

var _ repositories.CommentRepository = (*CommentRepositoryImpl)(nil)
