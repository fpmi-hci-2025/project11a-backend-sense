package domain

import "context"

// TagRepository defines interface for tag data operations
type TagRepository interface {
	// Create creates a new tag
	Create(ctx context.Context, tag *Tag) error
	
	// GetByID retrieves tag by ID
	GetByID(ctx context.Context, id string) (*Tag, error)
	
	// GetByName retrieves tag by name
	GetByName(ctx context.Context, name string) (*Tag, error)
	
	// GetPopular retrieves popular tags
	GetPopular(ctx context.Context, limit int, search *string) ([]*Tag, int, error)
	
	// AttachToPublication attaches tag to publication
	AttachToPublication(ctx context.Context, publicationID, tagID string) error
	
	// DetachFromPublication detaches tag from publication
	DetachFromPublication(ctx context.Context, publicationID, tagID string) error
	
	// GetByPublication retrieves tags for publication
	GetByPublication(ctx context.Context, publicationID string) ([]*Tag, error)
}

