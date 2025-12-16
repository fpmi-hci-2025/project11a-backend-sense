package domain

import (
	"context"
	"time"
)

// PublicationRepository defines interface for publication data operations
type PublicationRepository interface {
	// Create creates a new publication
	Create(ctx context.Context, publication *Publication, mediaIDs []string) error
	
	// GetByID retrieves publication by ID
	GetByID(ctx context.Context, id string) (*Publication, error)
	
	// Update updates publication
	Update(ctx context.Context, publication *Publication, mediaIDs []string) error
	
	// Delete deletes publication
	Delete(ctx context.Context, id string) error
	
	// GetFeed retrieves feed with filters
	GetFeed(ctx context.Context, userID *string, filters *FeedFilters, limit, offset int) ([]*Publication, int, error)
	
	// GetByAuthor retrieves publications by author
	GetByAuthor(ctx context.Context, authorID string, filters *PublicationFilters, limit, offset int) ([]*Publication, int, error)
	
	// Like toggles like on publication
	Like(ctx context.Context, userID, publicationID string) (bool, error) // returns true if liked, false if unliked
	
	// IsLiked checks if user liked publication
	IsLiked(ctx context.Context, userID, publicationID string) (bool, error)
	
	// GetLikesCount returns number of likes
	GetLikesCount(ctx context.Context, publicationID string) (int, error)
	
	// GetLikedUsers returns users who liked publication
	GetLikedUsers(ctx context.Context, publicationID string, limit, offset int) ([]*User, int, error)
	
	// Save saves publication for user
	Save(ctx context.Context, userID, publicationID string, note *string) error
	
	// Unsave removes saved publication
	Unsave(ctx context.Context, userID, publicationID string) error
	
	// IsSaved checks if publication is saved by user
	IsSaved(ctx context.Context, userID, publicationID string) (bool, error)
	
	// GetSaved retrieves saved publications for user
	GetSaved(ctx context.Context, userID string, filters *PublicationFilters, limit, offset int) ([]*SavedPublication, int, error)
	
	// Search searches publications by query
	Search(ctx context.Context, query string, filters *SearchFilters, limit, offset int) ([]*Publication, int, error)
	
	// GetMediaIDs retrieves media IDs for publication
	GetMediaIDs(ctx context.Context, publicationID string) ([]string, error)
}

// FeedFilters represents filters for feed
type FeedFilters struct {
	Type       *PublicationType
	Visibility *VisibilityType
	AuthorID   *string
	DateFrom   *time.Time
	DateTo     *time.Time
}

// PublicationFilters represents filters for publications
type PublicationFilters struct {
	Type       *PublicationType
	Visibility *VisibilityType
}

// SearchFilters represents filters for search
type SearchFilters struct {
	Type       *PublicationType
	Visibility *VisibilityType
	AuthorID   *string
}

// SavedPublication represents publication with saved metadata
type SavedPublication struct {
	Publication
	SavedNote *string    `json:"saved_note,omitempty"`
	SavedAt   time.Time  `json:"saved_at"`
}

