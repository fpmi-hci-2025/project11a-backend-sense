package media

import (
	"context"
	"errors"
	"fmt"
	"time"

	"sense-backend/internal/domain"

	"github.com/google/uuid"
)

// UseCase handles media use cases
type UseCase struct {
	mediaRepo domain.MediaRepository
}

// NewUseCase creates a new media use case
func NewUseCase(mediaRepo domain.MediaRepository) *UseCase {
	return &UseCase{mediaRepo: mediaRepo}
}

// UploadRequest represents upload media request
type UploadRequest struct {
	Data     []byte
	Filename *string
	MIME     string
	Width    *int
	Height   *int
	EXIF     *domain.EXIFData
}

// Upload uploads a media file
func (uc *UseCase) Upload(ctx context.Context, ownerID string, req *UploadRequest) (*domain.MediaAsset, error) {
	media := &domain.MediaAsset{
		ID:        uuid.New().String(),
		OwnerID:   ownerID,
		Filename:  req.Filename,
		MIME:      req.MIME,
		Width:     req.Width,
		Height:    req.Height,
		EXIF:      req.EXIF,
		Data:      req.Data,
		CreatedAt: time.Now(),
	}

	if err := uc.mediaRepo.Create(ctx, media); err != nil {
		return nil, fmt.Errorf("failed to upload media: %w", err)
	}

	return media, nil
}

// Get retrieves media by ID
func (uc *UseCase) Get(ctx context.Context, id string) (*domain.MediaAsset, error) {
	return uc.mediaRepo.GetByID(ctx, id)
}

// Delete deletes media
func (uc *UseCase) Delete(ctx context.Context, id, userID string) error {
	owned, err := uc.mediaRepo.CheckOwnership(ctx, id, userID)
	if err != nil {
		return err
	}
	if !owned {
		return errors.New("forbidden: not the owner")
	}

	return uc.mediaRepo.Delete(ctx, id)
}
