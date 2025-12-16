package publication

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"sense-backend/internal/domain"
)

// UseCase handles publication use cases
type UseCase struct {
	publicationRepo domain.PublicationRepository
	userRepo        domain.UserRepository
	mediaRepo       domain.MediaRepository
}

// NewUseCase creates a new publication use case
func NewUseCase(
	publicationRepo domain.PublicationRepository,
	userRepo domain.UserRepository,
	mediaRepo domain.MediaRepository,
) *UseCase {
	return &UseCase{
		publicationRepo: publicationRepo,
		userRepo:        userRepo,
		mediaRepo:       mediaRepo,
	}
}

// CreateRequest represents create publication request
type CreateRequest struct {
	Type       domain.PublicationType `json:"type" validate:"required"`
	Title      string                 `json:"title" validate:"required,max=500"`
	Content    *string                `json:"content,omitempty" validate:"omitempty,max=10000"`
	Source     *string                `json:"source,omitempty" validate:"omitempty,max=200"`
	Visibility domain.VisibilityType  `json:"visibility" validate:"required"`
	MediaIDs   []string               `json:"media_ids,omitempty"`
}

// UpdateRequest represents update publication request
type UpdateRequest struct {
	Title      *string                `json:"title,omitempty" validate:"omitempty,max=500"`
	Content    *string                `json:"content,omitempty" validate:"max=10000"`
	Source     *string                `json:"source,omitempty" validate:"max=200"`
	Visibility *domain.VisibilityType `json:"visibility,omitempty"`
	MediaIDs   []string               `json:"media_ids,omitempty"`
}

// Create creates a new publication
func (uc *UseCase) Create(ctx context.Context, authorID string, req *CreateRequest) (*domain.Publication, error) {
	// Validate media ownership
	for _, mediaID := range req.MediaIDs {
		owned, err := uc.mediaRepo.CheckOwnership(ctx, mediaID, authorID)
		if err != nil || !owned {
			return nil, errors.New("media not found or not owned")
		}
	}

	publication := &domain.Publication{
		ID:              uuid.New().String(),
		AuthorID:        authorID,
		Type:            req.Type,
		Title:           req.Title,
		Content:         req.Content,
		Source:          req.Source,
		PublicationDate: time.Now(),
		Visibility:      req.Visibility,
		LikesCount:      0,
		CommentsCount:   0,
		SavedCount:      0,
	}

	if err := uc.publicationRepo.Create(ctx, publication, req.MediaIDs); err != nil {
		return nil, fmt.Errorf("failed to create publication: %w", err)
	}

	return publication, nil
}

// Get retrieves publication by ID with like status for viewer
func (uc *UseCase) Get(ctx context.Context, id string, viewerUserID *string) (*domain.PublicationWithLikeStatus, error) {
	return uc.publicationRepo.GetByIDWithLikeStatus(ctx, id, viewerUserID)
}

// Update updates publication
func (uc *UseCase) Update(ctx context.Context, id, userID string, req *UpdateRequest) (*domain.Publication, error) {
	publication, err := uc.publicationRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if publication.AuthorID != userID {
		return nil, errors.New("forbidden: not the author")
	}

	if req.Title != nil {
		publication.Title = *req.Title
	}
	if req.Content != nil {
		publication.Content = req.Content
	}
	if req.Source != nil {
		publication.Source = req.Source
	}
	if req.Visibility != nil {
		publication.Visibility = *req.Visibility
	}

	mediaIDs := req.MediaIDs
	if mediaIDs == nil {
		mediaIDs, _ = uc.publicationRepo.GetMediaIDs(ctx, id)
	}

	// Validate media ownership
	for _, mediaID := range mediaIDs {
		owned, err := uc.mediaRepo.CheckOwnership(ctx, mediaID, userID)
		if err != nil || !owned {
			return nil, errors.New("media not found or not owned")
		}
	}

	if err := uc.publicationRepo.Update(ctx, publication, mediaIDs); err != nil {
		return nil, fmt.Errorf("failed to update publication: %w", err)
	}

	return publication, nil
}

// Delete deletes publication
func (uc *UseCase) Delete(ctx context.Context, id, userID string) error {
	publication, err := uc.publicationRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if publication.AuthorID != userID {
		return errors.New("forbidden: not the author")
	}

	return uc.publicationRepo.Delete(ctx, id)
}

// Like toggles like on publication
func (uc *UseCase) Like(ctx context.Context, publicationID, userID string) (bool, int, error) {
	liked, err := uc.publicationRepo.Like(ctx, userID, publicationID)
	if err != nil {
		return false, 0, err
	}

	count, err := uc.publicationRepo.GetLikesCount(ctx, publicationID)
	if err != nil {
		return false, 0, err
	}

	return liked, count, nil
}

// Save saves publication for user
func (uc *UseCase) Save(ctx context.Context, publicationID, userID string, note *string) error {
	return uc.publicationRepo.Save(ctx, userID, publicationID, note)
}

// Unsave removes saved publication
func (uc *UseCase) Unsave(ctx context.Context, publicationID, userID string) error {
	return uc.publicationRepo.Unsave(ctx, userID, publicationID)
}

// GetLikedUsers returns users who liked publication
func (uc *UseCase) GetLikedUsers(ctx context.Context, publicationID string, limit, offset int) ([]*domain.User, int, error) {
	return uc.publicationRepo.GetLikedUsers(ctx, publicationID, limit, offset)
}

