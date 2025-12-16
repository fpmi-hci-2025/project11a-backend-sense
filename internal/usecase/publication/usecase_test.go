package publication

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

func createTestPublication() *domain.Publication {
	content := "Test publication content"
	return &domain.Publication{
		ID:             "pub-123",
		AuthorID:       "user-123",
		Type:           domain.PublicationTypePost,
		Content:        &content,
		PublicationDate: time.Now(),
		Visibility:     domain.VisibilityTypePublic,
		LikesCount:     0,
		CommentsCount:  0,
		SavedCount:     0,
	}
}

func TestCreate_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	publicationRepo := mocks.NewMockPublicationRepository(ctrl)
	userRepo := mocks.NewMockUserRepository(ctrl)
	mediaRepo := mocks.NewMockMediaRepository(ctrl)
	uc := NewUseCase(publicationRepo, userRepo, mediaRepo)

	req := &CreateRequest{
		Type:       domain.PublicationTypePost,
		Content:    stringPtr("Test content"),
		Visibility: domain.VisibilityTypePublic,
		MediaIDs:   []string{},
	}

	publicationRepo.EXPECT().
		Create(gomock.Any(), gomock.Any(), []string{}).
		DoAndReturn(func(ctx context.Context, pub *domain.Publication, mediaIDs []string) error {
			assert.Equal(t, "user-123", pub.AuthorID)
			assert.Equal(t, domain.PublicationTypePost, pub.Type)
			return nil
		})

	pub, err := uc.Create(context.Background(), "user-123", req)

	require.NoError(t, err)
	assert.NotNil(t, pub)
	assert.Equal(t, domain.PublicationTypePost, pub.Type)
}

func TestCreate_WithMedia(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	publicationRepo := mocks.NewMockPublicationRepository(ctrl)
	userRepo := mocks.NewMockUserRepository(ctrl)
	mediaRepo := mocks.NewMockMediaRepository(ctrl)
	uc := NewUseCase(publicationRepo, userRepo, mediaRepo)

	req := &CreateRequest{
		Type:       domain.PublicationTypePost,
		Content:    stringPtr("Test content"),
		Visibility: domain.VisibilityTypePublic,
		MediaIDs:   []string{"media-1", "media-2"},
	}

	mediaRepo.EXPECT().
		CheckOwnership(gomock.Any(), "media-1", "user-123").
		Return(true, nil)

	mediaRepo.EXPECT().
		CheckOwnership(gomock.Any(), "media-2", "user-123").
		Return(true, nil)

	publicationRepo.EXPECT().
		Create(gomock.Any(), gomock.Any(), []string{"media-1", "media-2"}).
		Return(nil)

	pub, err := uc.Create(context.Background(), "user-123", req)

	require.NoError(t, err)
	assert.NotNil(t, pub)
}

func TestCreate_MediaNotOwned(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	publicationRepo := mocks.NewMockPublicationRepository(ctrl)
	userRepo := mocks.NewMockUserRepository(ctrl)
	mediaRepo := mocks.NewMockMediaRepository(ctrl)
	uc := NewUseCase(publicationRepo, userRepo, mediaRepo)

	req := &CreateRequest{
		Type:       domain.PublicationTypePost,
		Content:    stringPtr("Test content"),
		Visibility: domain.VisibilityTypePublic,
		MediaIDs:   []string{"media-1"},
	}

	mediaRepo.EXPECT().
		CheckOwnership(gomock.Any(), "media-1", "user-123").
		Return(false, nil)

	pub, err := uc.Create(context.Background(), "user-123", req)

	assert.Error(t, err)
	assert.Nil(t, pub)
	assert.Equal(t, "media not found or not owned", err.Error())
}

func TestGet_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	publicationRepo := mocks.NewMockPublicationRepository(ctrl)
	userRepo := mocks.NewMockUserRepository(ctrl)
	mediaRepo := mocks.NewMockMediaRepository(ctrl)
	uc := NewUseCase(publicationRepo, userRepo, mediaRepo)

	pub := createTestPublication()

	publicationRepo.EXPECT().
		GetByID(gomock.Any(), "pub-123").
		Return(pub, nil)

	result, err := uc.Get(context.Background(), "pub-123")

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "pub-123", result.ID)
}

func TestGet_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	publicationRepo := mocks.NewMockPublicationRepository(ctrl)
	userRepo := mocks.NewMockUserRepository(ctrl)
	mediaRepo := mocks.NewMockMediaRepository(ctrl)
	uc := NewUseCase(publicationRepo, userRepo, mediaRepo)

	publicationRepo.EXPECT().
		GetByID(gomock.Any(), "nonexistent").
		Return(nil, errors.New("not found"))

	result, err := uc.Get(context.Background(), "nonexistent")

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestUpdate_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	publicationRepo := mocks.NewMockPublicationRepository(ctrl)
	userRepo := mocks.NewMockUserRepository(ctrl)
	mediaRepo := mocks.NewMockMediaRepository(ctrl)
	uc := NewUseCase(publicationRepo, userRepo, mediaRepo)

	pub := createTestPublication()
	newContent := "Updated content"

	req := &UpdateRequest{
		Content: &newContent,
	}

	publicationRepo.EXPECT().
		GetByID(gomock.Any(), "pub-123").
		Return(pub, nil)

	publicationRepo.EXPECT().
		GetMediaIDs(gomock.Any(), "pub-123").
		Return([]string{}, nil)

	publicationRepo.EXPECT().
		Update(gomock.Any(), gomock.Any(), []string{}).
		DoAndReturn(func(ctx context.Context, pub *domain.Publication, mediaIDs []string) error {
			assert.Equal(t, "Updated content", *pub.Content)
			return nil
		})

	result, err := uc.Update(context.Background(), "pub-123", "user-123", req)

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Updated content", *result.Content)
}

func TestUpdate_NotAuthor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	publicationRepo := mocks.NewMockPublicationRepository(ctrl)
	userRepo := mocks.NewMockUserRepository(ctrl)
	mediaRepo := mocks.NewMockMediaRepository(ctrl)
	uc := NewUseCase(publicationRepo, userRepo, mediaRepo)

	pub := createTestPublication()

	req := &UpdateRequest{
		Content: stringPtr("Updated content"),
	}

	publicationRepo.EXPECT().
		GetByID(gomock.Any(), "pub-123").
		Return(pub, nil)

	result, err := uc.Update(context.Background(), "pub-123", "other-user", req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "forbidden: not the author", err.Error())
}

func TestUpdate_MediaOwnership(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	publicationRepo := mocks.NewMockPublicationRepository(ctrl)
	userRepo := mocks.NewMockUserRepository(ctrl)
	mediaRepo := mocks.NewMockMediaRepository(ctrl)
	uc := NewUseCase(publicationRepo, userRepo, mediaRepo)

	pub := createTestPublication()

	req := &UpdateRequest{
		MediaIDs: []string{"media-1"},
	}

	publicationRepo.EXPECT().
		GetByID(gomock.Any(), "pub-123").
		Return(pub, nil)

	mediaRepo.EXPECT().
		CheckOwnership(gomock.Any(), "media-1", "user-123").
		Return(false, nil)

	result, err := uc.Update(context.Background(), "pub-123", "user-123", req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "media not found or not owned", err.Error())
}

func TestDelete_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	publicationRepo := mocks.NewMockPublicationRepository(ctrl)
	userRepo := mocks.NewMockUserRepository(ctrl)
	mediaRepo := mocks.NewMockMediaRepository(ctrl)
	uc := NewUseCase(publicationRepo, userRepo, mediaRepo)

	pub := createTestPublication()

	publicationRepo.EXPECT().
		GetByID(gomock.Any(), "pub-123").
		Return(pub, nil)

	publicationRepo.EXPECT().
		Delete(gomock.Any(), "pub-123").
		Return(nil)

	err := uc.Delete(context.Background(), "pub-123", "user-123")

	require.NoError(t, err)
}

func TestDelete_NotAuthor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	publicationRepo := mocks.NewMockPublicationRepository(ctrl)
	userRepo := mocks.NewMockUserRepository(ctrl)
	mediaRepo := mocks.NewMockMediaRepository(ctrl)
	uc := NewUseCase(publicationRepo, userRepo, mediaRepo)

	pub := createTestPublication()

	publicationRepo.EXPECT().
		GetByID(gomock.Any(), "pub-123").
		Return(pub, nil)

	err := uc.Delete(context.Background(), "pub-123", "other-user")

	assert.Error(t, err)
	assert.Equal(t, "forbidden: not the author", err.Error())
}

func TestLike_Toggle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	publicationRepo := mocks.NewMockPublicationRepository(ctrl)
	userRepo := mocks.NewMockUserRepository(ctrl)
	mediaRepo := mocks.NewMockMediaRepository(ctrl)
	uc := NewUseCase(publicationRepo, userRepo, mediaRepo)

	publicationRepo.EXPECT().
		Like(gomock.Any(), "user-123", "pub-123").
		Return(true, nil)

	publicationRepo.EXPECT().
		GetLikesCount(gomock.Any(), "pub-123").
		Return(5, nil)

	liked, count, err := uc.Like(context.Background(), "pub-123", "user-123")

	require.NoError(t, err)
	assert.True(t, liked)
	assert.Equal(t, 5, count)
}

func TestSave_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	publicationRepo := mocks.NewMockPublicationRepository(ctrl)
	userRepo := mocks.NewMockUserRepository(ctrl)
	mediaRepo := mocks.NewMockMediaRepository(ctrl)
	uc := NewUseCase(publicationRepo, userRepo, mediaRepo)

	note := "My note"

	publicationRepo.EXPECT().
		Save(gomock.Any(), "user-123", "pub-123", &note).
		Return(nil)

	err := uc.Save(context.Background(), "pub-123", "user-123", &note)

	require.NoError(t, err)
}

func TestUnsave_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	publicationRepo := mocks.NewMockPublicationRepository(ctrl)
	userRepo := mocks.NewMockUserRepository(ctrl)
	mediaRepo := mocks.NewMockMediaRepository(ctrl)
	uc := NewUseCase(publicationRepo, userRepo, mediaRepo)

	publicationRepo.EXPECT().
		Unsave(gomock.Any(), "user-123", "pub-123").
		Return(nil)

	err := uc.Unsave(context.Background(), "pub-123", "user-123")

	require.NoError(t, err)
}

func TestGetLikedUsers_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	publicationRepo := mocks.NewMockPublicationRepository(ctrl)
	userRepo := mocks.NewMockUserRepository(ctrl)
	mediaRepo := mocks.NewMockMediaRepository(ctrl)
	uc := NewUseCase(publicationRepo, userRepo, mediaRepo)

	users := []*domain.User{
		{ID: "user-1", Username: "user1"},
		{ID: "user-2", Username: "user2"},
	}

	publicationRepo.EXPECT().
		GetLikedUsers(gomock.Any(), "pub-123", 10, 0).
		Return(users, 2, nil)

	result, total, err := uc.GetLikedUsers(context.Background(), "pub-123", 10, 0)

	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, 2, total)
}

func stringPtr(s string) *string {
	return &s
}

