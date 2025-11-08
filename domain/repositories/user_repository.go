package repositories

import (
	"context"
	"gofiber-template/domain/models"
	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	GetByOAuth(ctx context.Context, provider, oauthID string) (*models.User, error)
	Update(ctx context.Context, id uuid.UUID, user *models.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, offset, limit int) ([]*models.User, error)
	Count(ctx context.Context) (int64, error)
	SearchForChat(ctx context.Context, currentUserID uuid.UUID, query string, limit int) ([]*models.User, error)
	GetSuggestedForChat(ctx context.Context, currentUserID uuid.UUID, limit int) ([]*models.User, error)
}