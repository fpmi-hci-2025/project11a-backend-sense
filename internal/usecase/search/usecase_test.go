package search

import (
	"context"
	"testing"
	"time"

	"sense-backend/internal/domain"
	"sense-backend/internal/usecase/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const (
	testQuery = "test"
)

func createTestPublication() *domain.Publication {
	content := "Test publication"
	return &domain.Publication{
		ID:              "pub-123",
		AuthorID:        "user-123",
		Type:            domain.PublicationTypePost,
		Title:           "Test Title",
		Content:         &content,
		PublicationDate: time.Now(),
		Visibility:      domain.VisibilityTypePublic,
		LikesCount:      0,
		CommentsCount:   0,
		SavedCount:      0,
	}
}

func createTestPublicationWithLikeStatus() *domain.PublicationWithLikeStatus {
	return &domain.PublicationWithLikeStatus{
		Publication: *createTestPublication(),
		IsLiked:     false,
	}
}

func createTestUser() *domain.User {
	email := "test@example.com"
	return &domain.User{
		ID:           "user-123",
		Username:     "testuser",
		Email:        &email,
		Role:         domain.UserRoleUser,
		RegisteredAt: time.Now(),
	}
}

func createTestTag() *domain.Tag {
	return &domain.Tag{
		ID:         "tag-123",
		Name:       "testtag",
		UsageCount: 10,
		CreatedAt:  time.Now(),
	}
}

func TestSearchPublications_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	publicationRepo := mocks.NewMockPublicationRepository(ctrl)
	userRepo := mocks.NewMockUserRepository(ctrl)
	tagRepo := mocks.NewMockTagRepository(ctrl)
	uc := NewUseCase(publicationRepo, userRepo, tagRepo)

	const testUserID = "user-123"

	query := "test query"
	viewerUserID := testUserID
	filters := &domain.SearchFilters{}

	publications := []*domain.PublicationWithLikeStatus{
		createTestPublicationWithLikeStatus(),
	}

	publicationRepo.EXPECT().
		Search(gomock.Any(), query, &viewerUserID, filters, 10, 0).
		Return(publications, 1, nil)

	result, total, err := uc.SearchPublications(context.Background(), query, &viewerUserID, filters, 10, 0)

	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, 1, total)
}

func TestSearchPublications_WithFilters(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	publicationRepo := mocks.NewMockPublicationRepository(ctrl)
	userRepo := mocks.NewMockUserRepository(ctrl)
	tagRepo := mocks.NewMockTagRepository(ctrl)
	uc := NewUseCase(publicationRepo, userRepo, tagRepo)

	const testUserID = "user-123"

	query := "test query"
	viewerUserID := testUserID
	pubType := domain.PublicationTypeArticle
	visibility := domain.VisibilityTypePublic
	authorID := testUserID

	filters := &domain.SearchFilters{
		Type:       &pubType,
		Visibility: &visibility,
		AuthorID:   &authorID,
	}

	publications := []*domain.PublicationWithLikeStatus{
		createTestPublicationWithLikeStatus(),
	}

	publicationRepo.EXPECT().
		Search(gomock.Any(), query, &viewerUserID, filters, 20, 10).
		Return(publications, 1, nil)

	result, total, err := uc.SearchPublications(context.Background(), query, &viewerUserID, filters, 20, 10)

	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, 1, total)
}

func TestSearchUsers_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	publicationRepo := mocks.NewMockPublicationRepository(ctrl)
	userRepo := mocks.NewMockUserRepository(ctrl)
	tagRepo := mocks.NewMockTagRepository(ctrl)
	uc := NewUseCase(publicationRepo, userRepo, tagRepo)

	query := testQuery

	users := []*domain.User{
		createTestUser(),
	}

	userRepo.EXPECT().
		Search(gomock.Any(), query, nil, 10, 0).
		Return(users, 1, nil)

	result, total, err := uc.SearchUsers(context.Background(), query, nil, 10, 0)

	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, 1, total)
}

func TestSearchUsers_WithRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	publicationRepo := mocks.NewMockPublicationRepository(ctrl)
	userRepo := mocks.NewMockUserRepository(ctrl)
	tagRepo := mocks.NewMockTagRepository(ctrl)
	uc := NewUseCase(publicationRepo, userRepo, tagRepo)

	query := testQuery
	role := domain.UserRoleCreator

	users := []*domain.User{
		createTestUser(),
	}

	userRepo.EXPECT().
		Search(gomock.Any(), query, &role, 10, 0).
		Return(users, 1, nil)

	result, total, err := uc.SearchUsers(context.Background(), query, &role, 10, 0)

	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, 1, total)
}

func TestGetTags_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	publicationRepo := mocks.NewMockPublicationRepository(ctrl)
	userRepo := mocks.NewMockUserRepository(ctrl)
	tagRepo := mocks.NewMockTagRepository(ctrl)
	uc := NewUseCase(publicationRepo, userRepo, tagRepo)

	tags := []*domain.Tag{
		createTestTag(),
	}

	tagRepo.EXPECT().
		GetPopular(gomock.Any(), 10, nil).
		Return(tags, 1, nil)

	result, total, err := uc.GetTags(context.Background(), 10, nil)

	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, 1, total)
}

func TestGetTags_WithSearch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	publicationRepo := mocks.NewMockPublicationRepository(ctrl)
	userRepo := mocks.NewMockUserRepository(ctrl)
	tagRepo := mocks.NewMockTagRepository(ctrl)
	uc := NewUseCase(publicationRepo, userRepo, tagRepo)

	search := testQuery
	tags := []*domain.Tag{
		createTestTag(),
	}

	tagRepo.EXPECT().
		GetPopular(gomock.Any(), 10, &search).
		Return(tags, 1, nil)

	result, total, err := uc.GetTags(context.Background(), 10, &search)

	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, 1, total)
}

