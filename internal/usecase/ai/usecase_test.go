package ai

import (
	"context"
	"errors"
	"testing"
	"time"

	"sense-backend/internal/domain"
	aiClient "sense-backend/internal/infrastructure/ai"
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

func TestGetRecommendations_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	aiClientMock := mocks.NewMockClientInterface(ctrl)
	recommendationRepo := mocks.NewMockRecommendationRepository(ctrl)
	publicationRepo := mocks.NewMockPublicationRepository(ctrl)
	uc := NewUseCase(aiClientMock, recommendationRepo, publicationRepo)

	recommendResp := &aiClient.RecommendResponse{
		Publications: []string{"pub-1", "pub-2", "pub-3"},
	}

	pub1 := createTestPublication()
	pub1.ID = "pub-1"
	pub2 := createTestPublication()
	pub2.ID = "pub-2"

	aiClientMock.EXPECT().
		Recommend(gomock.Any(), "user-123").
		Return(recommendResp, nil)

	publicationRepo.EXPECT().
		GetByID(gomock.Any(), "pub-1").
		Return(pub1, nil)

	publicationRepo.EXPECT().
		GetByID(gomock.Any(), "pub-2").
		Return(pub2, nil)

	publicationRepo.EXPECT().
		GetByID(gomock.Any(), "pub-3").
		Return(nil, errors.New("not found"))

	recommendationRepo.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(nil).
		Times(2)

	result, err := uc.GetRecommendations(context.Background(), "user-123", 10, nil)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestGetRecommendations_AIError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	aiClientMock := mocks.NewMockClientInterface(ctrl)
	recommendationRepo := mocks.NewMockRecommendationRepository(ctrl)
	publicationRepo := mocks.NewMockPublicationRepository(ctrl)
	uc := NewUseCase(aiClientMock, recommendationRepo, publicationRepo)

	aiClientMock.EXPECT().
		Recommend(gomock.Any(), "user-123").
		Return(nil, errors.New("AI service error"))

	result, err := uc.GetRecommendations(context.Background(), "user-123", 10, nil)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get AI recommendations")
}

func TestGetRecommendations_PublicationNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	aiClientMock := mocks.NewMockClientInterface(ctrl)
	recommendationRepo := mocks.NewMockRecommendationRepository(ctrl)
	publicationRepo := mocks.NewMockPublicationRepository(ctrl)
	uc := NewUseCase(aiClientMock, recommendationRepo, publicationRepo)

	recommendResp := &aiClient.RecommendResponse{
		Publications: []string{"pub-1", "pub-2"},
	}

	aiClientMock.EXPECT().
		Recommend(gomock.Any(), "user-123").
		Return(recommendResp, nil)

	publicationRepo.EXPECT().
		GetByID(gomock.Any(), "pub-1").
		Return(nil, errors.New("not found"))

	publicationRepo.EXPECT().
		GetByID(gomock.Any(), "pub-2").
		Return(nil, errors.New("not found"))

	result, err := uc.GetRecommendations(context.Background(), "user-123", 10, nil)

	require.NoError(t, err)
	assert.Empty(t, result)
}

func TestGetRecommendationsFeed_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	aiClientMock := mocks.NewMockClientInterface(ctrl)
	recommendationRepo := mocks.NewMockRecommendationRepository(ctrl)
	publicationRepo := mocks.NewMockPublicationRepository(ctrl)
	uc := NewUseCase(aiClientMock, recommendationRepo, publicationRepo)

	publicationIDs := []string{"pub-1", "pub-2", "pub-3"}

	recommendationRepo.EXPECT().
		GetPublicationIDs(gomock.Any(), "user-123", 10).
		Return(publicationIDs, nil)

	pub1 := createTestPublication()
	pub1.ID = "pub-1"
	pub2 := createTestPublication()
	pub2.ID = "pub-2"

	publicationRepo.EXPECT().
		GetByID(gomock.Any(), "pub-1").
		Return(pub1, nil)

	publicationRepo.EXPECT().
		GetByID(gomock.Any(), "pub-2").
		Return(pub2, nil)

	publicationRepo.EXPECT().
		GetByID(gomock.Any(), "pub-3").
		Return(nil, errors.New("not found"))

	result, total, err := uc.GetRecommendationsFeed(context.Background(), "user-123", 10, 0, nil)

	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, 2, total)
}

func TestGetRecommendationsFeed_WithOffset(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	aiClientMock := mocks.NewMockClientInterface(ctrl)
	recommendationRepo := mocks.NewMockRecommendationRepository(ctrl)
	publicationRepo := mocks.NewMockPublicationRepository(ctrl)
	uc := NewUseCase(aiClientMock, recommendationRepo, publicationRepo)

	publicationIDs := []string{"pub-1", "pub-2", "pub-3", "pub-4", "pub-5"}

	// limit + offset = 2 + 2 = 4
	recommendationRepo.EXPECT().
		GetPublicationIDs(gomock.Any(), "user-123", 4).
		Return(publicationIDs, nil)

	pub3 := createTestPublication()
	pub3.ID = "pub-3"
	pub4 := createTestPublication()
	pub4.ID = "pub-4"

	publicationRepo.EXPECT().
		GetByID(gomock.Any(), "pub-3").
		Return(pub3, nil)

	publicationRepo.EXPECT().
		GetByID(gomock.Any(), "pub-4").
		Return(pub4, nil)

	result, total, err := uc.GetRecommendationsFeed(context.Background(), "user-123", 2, 2, nil)

	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, 2, total)
}

func TestHideRecommendation_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	aiClientMock := mocks.NewMockClientInterface(ctrl)
	recommendationRepo := mocks.NewMockRecommendationRepository(ctrl)
	publicationRepo := mocks.NewMockPublicationRepository(ctrl)
	uc := NewUseCase(aiClientMock, recommendationRepo, publicationRepo)

	recommendationRepo.EXPECT().
		Hide(gomock.Any(), "rec-123").
		Return(nil)

	err := uc.HideRecommendation(context.Background(), "rec-123")

	require.NoError(t, err)
}

func TestPurifyText_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	aiClientMock := mocks.NewMockClientInterface(ctrl)
	recommendationRepo := mocks.NewMockRecommendationRepository(ctrl)
	publicationRepo := mocks.NewMockPublicationRepository(ctrl)
	uc := NewUseCase(aiClientMock, recommendationRepo, publicationRepo)

	req := &PurifyRequest{
		Text: "Original text",
	}

	composeResp := &aiClient.ComposeResponse{
		Text: "Original text",
	}

	aiClientMock.EXPECT().
		Compose(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, composeReq *aiClient.ComposeRequest) (*aiClient.ComposeResponse, error) {
			assert.Contains(t, composeReq.Query, "Original text")
			return composeResp, nil
		})

	result, err := uc.PurifyText(context.Background(), req)

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Original text", result.CleanedText)
	assert.True(t, result.IsClean)
	assert.Equal(t, 0.95, result.Confidence)
}

func TestPurifyText_AIError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	aiClientMock := mocks.NewMockClientInterface(ctrl)
	recommendationRepo := mocks.NewMockRecommendationRepository(ctrl)
	publicationRepo := mocks.NewMockPublicationRepository(ctrl)
	uc := NewUseCase(aiClientMock, recommendationRepo, publicationRepo)

	req := &PurifyRequest{
		Text: "Original text",
	}

	aiClientMock.EXPECT().
		Compose(gomock.Any(), gomock.Any()).
		Return(nil, errors.New("AI service error"))

	result, err := uc.PurifyText(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to purify text")
}

func TestPurifyText_Modified(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	aiClientMock := mocks.NewMockClientInterface(ctrl)
	recommendationRepo := mocks.NewMockRecommendationRepository(ctrl)
	publicationRepo := mocks.NewMockPublicationRepository(ctrl)
	uc := NewUseCase(aiClientMock, recommendationRepo, publicationRepo)

	req := &PurifyRequest{
		Text: "Original text with bad words",
	}

	composeResp := &aiClient.ComposeResponse{
		Text: "Cleaned text",
	}

	aiClientMock.EXPECT().
		Compose(gomock.Any(), gomock.Any()).
		Return(composeResp, nil)

	result, err := uc.PurifyText(context.Background(), req)

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Cleaned text", result.CleanedText)
	assert.False(t, result.IsClean)
	assert.Equal(t, 0.85, result.Confidence)
}

