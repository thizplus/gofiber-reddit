package postgres

import (
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gofiber-template/domain/models"
	"gofiber-template/domain/repositories"
)

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repositories.UserRepository {
	return &UserRepositoryImpl{db: db}
}

func (r *UserRepositoryImpl) Create(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *UserRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryImpl) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryImpl) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryImpl) GetByOAuth(ctx context.Context, provider, oauthID string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).
		Where("oauth_provider = ? AND oauth_id = ?", provider, oauthID).
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryImpl) Update(ctx context.Context, id uuid.UUID, user *models.User) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Updates(user).Error
}

func (r *UserRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.User{}).Error
}

func (r *UserRepositoryImpl) List(ctx context.Context, offset, limit int) ([]*models.User, error) {
	var users []*models.User
	err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&users).Error
	return users, err
}

func (r *UserRepositoryImpl) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.User{}).Count(&count).Error
	return count, err
}

// SearchForChat searches users for chat (excludes self, blocked users)
func (r *UserRepositoryImpl) SearchForChat(ctx context.Context, currentUserID uuid.UUID, query string, limit int) ([]*models.User, error) {
	var users []*models.User

	// Build subquery for blocked users
	blockedUsersSubquery := r.db.Table("blocks").
		Select("CASE WHEN blocker_id = ? THEN blocked_id WHEN blocked_id = ? THEN blocker_id END as user_id", currentUserID, currentUserID).
		Where("blocker_id = ? OR blocked_id = ?", currentUserID, currentUserID)

	// Build query to exclude self and blocked users
	err := r.db.WithContext(ctx).
		Where("id != ?", currentUserID).
		Where("is_active = ?", true).
		Where(
			r.db.Where("LOWER(username) LIKE ?", "%"+query+"%").
			Or("LOWER(display_name) LIKE ?", "%"+query+"%"),
		).
		// Exclude blocked users
		Where("id NOT IN (?)", blockedUsersSubquery).
		// Order by: followers first, then following, then by karma
		Order(r.db.Raw("(SELECT COUNT(*) FROM follows WHERE follower_id = users.id AND following_id = ?) DESC", currentUserID)).
		Order(r.db.Raw("(SELECT COUNT(*) FROM follows WHERE follower_id = ? AND following_id = users.id) DESC", currentUserID)).
		Order("karma DESC").
		Limit(limit).
		Find(&users).Error

	return users, err
}

// GetSuggestedForChat gets suggested users for chat (followers, following, popular users)
func (r *UserRepositoryImpl) GetSuggestedForChat(ctx context.Context, currentUserID uuid.UUID, limit int) ([]*models.User, error) {
	var users []*models.User

	// Build subquery for blocked users
	blockedUsersSubquery := r.db.Table("blocks").
		Select("CASE WHEN blocker_id = ? THEN blocked_id WHEN blocked_id = ? THEN blocker_id END as user_id", currentUserID, currentUserID).
		Where("blocker_id = ? OR blocked_id = ?", currentUserID, currentUserID)

	// Get users who follow current user or current user follows
	err := r.db.WithContext(ctx).
		Where("id != ?", currentUserID).
		Where("is_active = ?", true).
		Where(
			r.db.Where("id IN (?)",
				r.db.Table("follows").
					Select("follower_id").
					Where("following_id = ?", currentUserID),
			).Or("id IN (?)",
				r.db.Table("follows").
					Select("following_id").
					Where("follower_id = ?", currentUserID),
			),
		).
		// Exclude blocked users
		Where("id NOT IN (?)", blockedUsersSubquery).
		// Order by recent conversations first (if exists), then by karma
		Order(r.db.Raw("(SELECT MAX(last_message_at) FROM conversations WHERE user1_id = users.id OR user2_id = users.id) DESC NULLS LAST")).
		Order("karma DESC").
		Limit(limit).
		Find(&users).Error

	return users, err
}