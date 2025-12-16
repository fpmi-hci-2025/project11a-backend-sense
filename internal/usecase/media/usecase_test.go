package media

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

const (
	testFilename = "test.jpg"
)

func createTestMedia() *domain.MediaAsset {
	filename := testFilename
	width := 1920
	height := 1080
	return &domain.MediaAsset{
		ID:        "media-123",
		OwnerID:   "user-123",
		Filename:  &filename,
		MIME:      "image/jpeg",
		Width:     &width,
		Height:    &height,
		Data:      []byte("fake image data"),
		CreatedAt: time.Now(),
	}
}

func TestUpload_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mediaRepo := mocks.NewMockMediaRepository(ctrl)
	uc := NewUseCase(mediaRepo)

	filename := testFilename
	req := &UploadRequest{
		Data:     []byte("fake image data"),
		Filename: &filename,
		MIME:     "image/jpeg",
		Width:    intPtr(1920),
		Height:   intPtr(1080),
	}

	mediaRepo.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, media *domain.MediaAsset) error {
			assert.Equal(t, "user-123", media.OwnerID)
			assert.Equal(t, "image/jpeg", media.MIME)
			return nil
		})

	media, err := uc.Upload(context.Background(), "user-123", req)

	require.NoError(t, err)
	assert.NotNil(t, media)
	assert.Equal(t, "image/jpeg", media.MIME)
}

func TestUpload_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mediaRepo := mocks.NewMockMediaRepository(ctrl)
	uc := NewUseCase(mediaRepo)

	filename := testFilename
	req := &UploadRequest{
		Data:     []byte("fake image data"),
		Filename: &filename,
		MIME:     "image/jpeg",
	}

	mediaRepo.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(errors.New("database error"))

	media, err := uc.Upload(context.Background(), "user-123", req)

	assert.Error(t, err)
	assert.Nil(t, media)
	assert.Contains(t, err.Error(), "failed to upload media")
}

func TestGet_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mediaRepo := mocks.NewMockMediaRepository(ctrl)
	uc := NewUseCase(mediaRepo)

	media := createTestMedia()

	mediaRepo.EXPECT().
		GetByID(gomock.Any(), "media-123").
		Return(media, nil)

	result, err := uc.Get(context.Background(), "media-123")

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "media-123", result.ID)
}

func TestGet_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mediaRepo := mocks.NewMockMediaRepository(ctrl)
	uc := NewUseCase(mediaRepo)

	mediaRepo.EXPECT().
		GetByID(gomock.Any(), "nonexistent").
		Return(nil, errors.New("not found"))

	result, err := uc.Get(context.Background(), "nonexistent")

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestDelete_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mediaRepo := mocks.NewMockMediaRepository(ctrl)
	uc := NewUseCase(mediaRepo)

	mediaRepo.EXPECT().
		CheckOwnership(gomock.Any(), "media-123", "user-123").
		Return(true, nil)

	mediaRepo.EXPECT().
		Delete(gomock.Any(), "media-123").
		Return(nil)

	err := uc.Delete(context.Background(), "media-123", "user-123")

	require.NoError(t, err)
}

func TestDelete_NotOwner(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mediaRepo := mocks.NewMockMediaRepository(ctrl)
	uc := NewUseCase(mediaRepo)

	mediaRepo.EXPECT().
		CheckOwnership(gomock.Any(), "media-123", "other-user").
		Return(false, nil)

	err := uc.Delete(context.Background(), "media-123", "other-user")

	assert.Error(t, err)
	assert.Equal(t, "forbidden: not the owner", err.Error())
}

func TestDelete_MediaNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mediaRepo := mocks.NewMockMediaRepository(ctrl)
	uc := NewUseCase(mediaRepo)

	mediaRepo.EXPECT().
		CheckOwnership(gomock.Any(), "nonexistent", "user-123").
		Return(false, errors.New("not found"))

	err := uc.Delete(context.Background(), "nonexistent", "user-123")

	assert.Error(t, err)
}

func intPtr(i int) *int {
	return &i
}

