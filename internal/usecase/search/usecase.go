package search

import (
	"context"
	"sense-backend/internal/domain"
)

// UseCase handles search use cases
type UseCase struct {
	publicationRepo domain.PublicationRepository
	userRepo        domain.UserRepository
	tagRepo         domain.TagRepository
}

// NewUseCase creates a new search use case
func NewUseCase(
	publicationRepo domain.PublicationRepository,
	userRepo domain.UserRepository,
	tagRepo domain.TagRepository,
) *UseCase {
	return &UseCase{
		publicationRepo: publicationRepo,
		userRepo:        userRepo,
		tagRepo:         tagRepo,
	}
}

// SearchPublications searches publications
func (uc *UseCase) SearchPublications(ctx context.Context, query string, filters *domain.SearchFilters, limit, offset int) ([]*domain.Publication, int, error) {
	return uc.publicationRepo.Search(ctx, query, filters, limit, offset)
}

// SearchUsers searches users
func (uc *UseCase) SearchUsers(ctx context.Context, query string, role *domain.UserRole, limit, offset int) ([]*domain.User, int, error) {
	return uc.userRepo.Search(ctx, query, role, limit, offset)
}

// GetTags retrieves popular tags
func (uc *UseCase) GetTags(ctx context.Context, limit int, search *string) ([]*domain.Tag, int, error) {
	return uc.tagRepo.GetPopular(ctx, limit, search)
}

