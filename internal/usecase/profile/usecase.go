package profile

import (
	"context"
	"errors"
	"fmt"
	"sense-backend/internal/domain"
)

// UseCase handles profile use cases
type UseCase struct {
	userRepo domain.UserRepository
}

// NewUseCase creates a new profile use case
func NewUseCase(userRepo domain.UserRepository) *UseCase {
	return &UseCase{userRepo: userRepo}
}

// UpdateRequest represents update profile request
type UpdateRequest struct {
	Description *string `json:"description,omitempty" validate:"max=500"`
	IconURL     *string `json:"icon_url,omitempty"`
}

// GetProfile retrieves user profile
func (uc *UseCase) GetProfile(ctx context.Context, userID string) (*domain.User, error) {
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	stats, err := uc.userRepo.GetStats(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get stats: %w", err)
	}

	user.Statistic = stats
	user.PasswordHash = ""
	return user, nil
}

// UpdateProfile updates user profile
func (uc *UseCase) UpdateProfile(ctx context.Context, userID string, req *UpdateRequest) (*domain.User, error) {
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if req.Description != nil {
		user.Description = req.Description
	}
	if req.IconURL != nil {
		user.IconURL = req.IconURL
	}

	if err := uc.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}

	user.PasswordHash = ""
	return user, nil
}

// GetStats retrieves user statistics
func (uc *UseCase) GetStats(ctx context.Context, userID string) (*domain.UserStatistic, error) {
	return uc.userRepo.GetStats(ctx, userID)
}

// Follow follows a user
func (uc *UseCase) Follow(ctx context.Context, followerID, followingID string) error {
	if followerID == followingID {
		return errors.New("cannot follow yourself")
	}
	return uc.userRepo.Follow(ctx, followerID, followingID)
}

// Unfollow unfollows a user
func (uc *UseCase) Unfollow(ctx context.Context, followerID, followingID string) error {
	return uc.userRepo.Unfollow(ctx, followerID, followingID)
}

// IsFollowing checks if user is following another user
func (uc *UseCase) IsFollowing(ctx context.Context, followerID, followingID string) (bool, error) {
	return uc.userRepo.IsFollowing(ctx, followerID, followingID)
}

