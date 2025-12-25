package domain

import "context"

// MediaRepository defines interface for media data operations
type MediaRepository interface {
	// Create creates a new media asset
	Create(ctx context.Context, media *MediaAsset) error

	// GetByID retrieves media by ID
	GetByID(ctx context.Context, id string) (*MediaAsset, error)

	// GetByOwner retrieves media by owner
	GetByOwner(ctx context.Context, ownerID string, limit, offset int) ([]*MediaAsset, int, error)

	// GetByPublicationID retrieves media by publication ID
	GetByPublicationID(ctx context.Context, publicationID string) ([]*MediaAsset, error)

	// GetByPublicationIDs retrieves media for multiple publications (batch)
	GetByPublicationIDs(ctx context.Context, publicationIDs []string) (map[string][]*MediaAsset, error)

	// Delete deletes media asset
	Delete(ctx context.Context, id string) error

	// CheckOwnership checks if user owns media
	CheckOwnership(ctx context.Context, mediaID, userID string) (bool, error)
}
