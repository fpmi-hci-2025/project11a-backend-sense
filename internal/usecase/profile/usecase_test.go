package profile

import (
	"context"
	"errors"
	"testing"
	"time"

	"sense-backend/internal/domain"
	"sense-backend/internal/usecase/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func createTestUser() *domain.User {
	email := "test@example.com"
	return &domain.User{
		ID:           "user-123",
		Username:     "testuser",
		Email:        &email,
		Role:         domain.UserRoleUser,
		RegisteredAt: time.Now(),
		PasswordHash: "hashed",
	}
}

func createTestStats() *domain.UserStatistic {
	return &domain.UserStatistic{
		PublicationsCount: 10,
		FollowersCount:    5,
		FollowingCount:    3,
		LikesReceived:     20,
		CommentsReceived:  15,
		SavedCount:        7,
	}
}

func TestGetProfile_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	uc := NewUseCase(userRepo)

	user := createTestUser()
	stats := createTestStats()

	userRepo.EXPECT().
		GetByID(gomock.Any(), "user-123").
		Return(user, nil)

	userRepo.EXPECT().
		GetStats(gomock.Any(), "user-123").
		Return(stats, nil)

	result, err := uc.GetProfile(context.Background(), "user-123")

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "user-123", result.ID)
	assert.Empty(t, result.PasswordHash)
	assert.NotNil(t, result.Statistic)
	assert.Equal(t, 10, result.Statistic.PublicationsCount)
}

func TestGetProfile_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	uc := NewUseCase(userRepo)

	userRepo.EXPECT().
		GetByID(gomock.Any(), "nonexistent").
		Return(nil, errors.New("not found"))

	result, err := uc.GetProfile(context.Background(), "nonexistent")

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetProfile_StatsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	uc := NewUseCase(userRepo)

	user := createTestUser()

	userRepo.EXPECT().
		GetByID(gomock.Any(), "user-123").
		Return(user, nil)

	userRepo.EXPECT().
		GetStats(gomock.Any(), "user-123").
		Return(nil, errors.New("stats error"))

	result, err := uc.GetProfile(context.Background(), "user-123")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get stats")
}

func TestUpdateProfile_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	uc := NewUseCase(userRepo)

	user := createTestUser()
	description := "Updated description"

	req := &UpdateRequest{
		Description: &description,
	}

	userRepo.EXPECT().
		GetByID(gomock.Any(), "user-123").
		Return(user, nil)

	userRepo.EXPECT().
		Update(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, user *domain.User) error {
			assert.Equal(t, "Updated description", *user.Description)
			return nil
		})

	result, err := uc.UpdateProfile(context.Background(), "user-123", req)

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Updated description", *result.Description)
	assert.Empty(t, result.PasswordHash)
}

func TestUpdateProfile_UpdateIcon(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	uc := NewUseCase(userRepo)

	user := createTestUser()
	iconURL := "https://example.com/icon.jpg"

	req := &UpdateRequest{
		IconURL: &iconURL,
	}

	userRepo.EXPECT().
		GetByID(gomock.Any(), "user-123").
		Return(user, nil)

	userRepo.EXPECT().
		Update(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, user *domain.User) error {
			assert.Equal(t, iconURL, *user.IconURL)
			return nil
		})

	result, err := uc.UpdateProfile(context.Background(), "user-123", req)

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, iconURL, *result.IconURL)
}

func TestUpdateProfile_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	uc := NewUseCase(userRepo)

	description := "Updated description"
	req := &UpdateRequest{
		Description: &description,
	}

	userRepo.EXPECT().
		GetByID(gomock.Any(), "nonexistent").
		Return(nil, errors.New("not found"))

	result, err := uc.UpdateProfile(context.Background(), "nonexistent", req)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetStats_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	uc := NewUseCase(userRepo)

	stats := createTestStats()

	userRepo.EXPECT().
		GetStats(gomock.Any(), "user-123").
		Return(stats, nil)

	result, err := uc.GetStats(context.Background(), "user-123")

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 10, result.PublicationsCount)
}

func TestFollow_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	uc := NewUseCase(userRepo)

	userRepo.EXPECT().
		GetByID(gomock.Any(), "user-456").
		Return(&domain.User{ID: "user-456"}, nil)

	userRepo.EXPECT().
		Follow(gomock.Any(), "user-123", "user-456").
		Return(nil)

	err := uc.Follow(context.Background(), "user-123", "user-456")

	require.NoError(t, err)
}

func TestFollow_Self(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	uc := NewUseCase(userRepo)

	err := uc.Follow(context.Background(), "user-123", "user-123")

	assert.Error(t, err)
	assert.Equal(t, "cannot follow yourself", err.Error())
}

func TestFollow_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	uc := NewUseCase(userRepo)

	userRepo.EXPECT().
		GetByID(gomock.Any(), "nonexistent").
		Return(nil, assert.AnError)

	err := uc.Follow(context.Background(), "user-123", "nonexistent")

	assert.Error(t, err)
}

func TestUnfollow_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	uc := NewUseCase(userRepo)

	userRepo.EXPECT().
		GetByID(gomock.Any(), "user-456").
		Return(&domain.User{ID: "user-456"}, nil)

	userRepo.EXPECT().
		Unfollow(gomock.Any(), "user-123", "user-456").
		Return(nil)

	err := uc.Unfollow(context.Background(), "user-123", "user-456")

	require.NoError(t, err)
}

func TestUnfollow_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	uc := NewUseCase(userRepo)

	userRepo.EXPECT().
		GetByID(gomock.Any(), "nonexistent").
		Return(nil, assert.AnError)

	err := uc.Unfollow(context.Background(), "user-123", "nonexistent")

	assert.Error(t, err)
}

func TestIsFollowing_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	uc := NewUseCase(userRepo)

	userRepo.EXPECT().
		IsFollowing(gomock.Any(), "user-123", "user-456").
		Return(true, nil)

	result, err := uc.IsFollowing(context.Background(), "user-123", "user-456")

	require.NoError(t, err)
	assert.True(t, result)
}
