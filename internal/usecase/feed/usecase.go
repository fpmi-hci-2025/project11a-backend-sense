package feed

import (
	"context"
	"sense-backend/internal/domain"
)

// UseCase handles feed use cases
type UseCase struct {
	publicationRepo domain.PublicationRepository
}

// NewUseCase creates a new feed use case
func NewUseCase(publicationRepo domain.PublicationRepository) *UseCase {
	return &UseCase{publicationRepo: publicationRepo}
}

// GetFeed retrieves feed with filters
func (uc *UseCase) GetFeed(ctx context.Context, userID *string, filters *domain.FeedFilters, limit, offset int) ([]*domain.Publication, int, error) {
	return uc.publicationRepo.GetFeed(ctx, userID, filters, limit, offset)
}

// GetUserFeed retrieves publications by user
func (uc *UseCase) GetUserFeed(ctx context.Context, authorID string, filters *domain.PublicationFilters, limit, offset int) ([]*domain.Publication, int, error) {
	return uc.publicationRepo.GetByAuthor(ctx, authorID, filters, limit, offset)
}

// GetSavedFeed retrieves saved publications for user
func (uc *UseCase) GetSavedFeed(ctx context.Context, userID string, filters *domain.PublicationFilters, limit, offset int) ([]*domain.SavedPublication, int, error) {
	return uc.publicationRepo.GetSaved(ctx, userID, filters, limit, offset)
}

