package feed

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

func createTestPublication() *domain.Publication {
	content := "Test publication"
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

func TestGetFeed_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	publicationRepo := mocks.NewMockPublicationRepository(ctrl)
	uc := NewUseCase(publicationRepo)

	userID := "user-123"
	filters := &domain.FeedFilters{
		Type:       &[]domain.PublicationType{domain.PublicationTypePost}[0],
		Visibility: &[]domain.VisibilityType{domain.VisibilityTypePublic}[0],
	}

	publications := []*domain.Publication{
		createTestPublication(),
	}

	publicationRepo.EXPECT().
		GetFeed(gomock.Any(), &userID, filters, 10, 0).
		Return(publications, 1, nil)

	result, total, err := uc.GetFeed(context.Background(), &userID, filters, 10, 0)

	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, 1, total)
}

func TestGetFeed_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	publicationRepo := mocks.NewMockPublicationRepository(ctrl)
	uc := NewUseCase(publicationRepo)

	userID := "user-123"
	filters := &domain.FeedFilters{}

	publicationRepo.EXPECT().
		GetFeed(gomock.Any(), &userID, filters, 10, 0).
		Return([]*domain.Publication{}, 0, nil)

	result, total, err := uc.GetFeed(context.Background(), &userID, filters, 10, 0)

	require.NoError(t, err)
	assert.Empty(t, result)
	assert.Equal(t, 0, total)
}

func TestGetUserFeed_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	publicationRepo := mocks.NewMockPublicationRepository(ctrl)
	uc := NewUseCase(publicationRepo)

	filters := &domain.PublicationFilters{
		Type: &[]domain.PublicationType{domain.PublicationTypePost}[0],
	}

	publications := []*domain.Publication{
		createTestPublication(),
	}

	publicationRepo.EXPECT().
		GetByAuthor(gomock.Any(), "user-123", filters, 10, 0).
		Return(publications, 1, nil)

	result, total, err := uc.GetUserFeed(context.Background(), "user-123", filters, 10, 0)

	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, 1, total)
}

func TestGetSavedFeed_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	publicationRepo := mocks.NewMockPublicationRepository(ctrl)
	uc := NewUseCase(publicationRepo)

	filters := &domain.PublicationFilters{}

	savedPublications := []*domain.SavedPublication{
		{
			Publication: *createTestPublication(),
			SavedAt:     time.Now(),
		},
	}

	publicationRepo.EXPECT().
		GetSaved(gomock.Any(), "user-123", filters, 10, 0).
		Return(savedPublications, 1, nil)

	result, total, err := uc.GetSavedFeed(context.Background(), "user-123", filters, 10, 0)

	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, 1, total)
}

func TestGetFeed_WithFilters(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	publicationRepo := mocks.NewMockPublicationRepository(ctrl)
	uc := NewUseCase(publicationRepo)

	userID := "user-123"
	dateFrom := time.Now().Add(-24 * time.Hour)
	dateTo := time.Now()
	pubType := domain.PublicationTypeArticle
	visibility := domain.VisibilityTypePublic

	filters := &domain.FeedFilters{
		Type:       &pubType,
		Visibility: &visibility,
		DateFrom:   &dateFrom,
		DateTo:     &dateTo,
	}

	publications := []*domain.Publication{
		createTestPublication(),
	}

	publicationRepo.EXPECT().
		GetFeed(gomock.Any(), &userID, filters, 20, 10).
		Return(publications, 1, nil)

	result, total, err := uc.GetFeed(context.Background(), &userID, filters, 20, 10)

	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, 1, total)
}

