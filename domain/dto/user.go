package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateUserRequest struct {
	Email       string `json:"email" validate:"required,email,max=255"`
	Username    string `json:"username" validate:"required,min=3,max=20"`
	Password    string `json:"password" validate:"required,min=8,max=72"`
	DisplayName string `json:"displayName" validate:"required,min=1,max=100"`
}

type UpdateUserRequest struct {
	DisplayName string `json:"displayName" validate:"omitempty,min=1,max=100"`
	Bio         string `json:"bio" validate:"omitempty,max=500"`
	Location    string `json:"location" validate:"omitempty,max=100"`
	Website     string `json:"website" validate:"omitempty,url,max=255"`
	Avatar      string `json:"avatar" validate:"omitempty,max=500"`
}

type UserResponse struct {
	ID              uuid.UUID `json:"id"`
	Email           string    `json:"email,omitempty"` // Only for owner
	Username        string    `json:"username"`
	DisplayName     string    `json:"displayName"`
	Avatar          string    `json:"avatar,omitempty"`
	Bio             string    `json:"bio,omitempty"`
	Location        string    `json:"location,omitempty"`
	Website         string    `json:"website,omitempty"`
	Karma           int       `json:"karma"`
	FollowersCount  int       `json:"followersCount"`
	FollowingCount  int       `json:"followingCount"`
	Role            string    `json:"role,omitempty"`
	IsActive        bool      `json:"isActive"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
	IsFollowing     *bool     `json:"isFollowing,omitempty"` // Only when authenticated
}

type UserListResponse struct {
	Users []UserResponse `json:"users"`
	Meta  PaginationMeta `json:"meta"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" validate:"required"`
	NewPassword     string `json:"newPassword" validate:"required,min=8,max=72"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,eqfield=NewPassword"`
}
