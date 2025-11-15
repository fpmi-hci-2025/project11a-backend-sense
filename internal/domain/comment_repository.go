package domain

import "context"

// CommentRepository defines interface for comment data operations
type CommentRepository interface {
	// Create creates a new comment
	Create(ctx context.Context, comment *Comment) error
	
	// GetByID retrieves comment by ID
	GetByID(ctx context.Context, id string) (*Comment, error)
	
	// GetByPublication retrieves comments for publication
	GetByPublication(ctx context.Context, publicationID string, limit, offset int) ([]*Comment, int, error)
	
	// Update updates comment
	Update(ctx context.Context, comment *Comment) error
	
	// Delete deletes comment
	Delete(ctx context.Context, id string) error
	
	// Like toggles like on comment
	Like(ctx context.Context, userID, commentID string) (bool, error) // returns true if liked, false if unliked
	
	// IsLiked checks if user liked comment
	IsLiked(ctx context.Context, userID, commentID string) (bool, error)
	
	// GetLikesCount returns number of likes
	GetLikesCount(ctx context.Context, commentID string) (int, error)
}

