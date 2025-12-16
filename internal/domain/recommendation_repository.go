package domain

import "context"

// RecommendationRepository defines interface for recommendation data operations
type RecommendationRepository interface {
	// Create creates a new recommendation
	Create(ctx context.Context, recommendation *Recommendation) error
	
	// GetByUser retrieves recommendations for user
	GetByUser(ctx context.Context, userID string, limit, offset int) ([]*Recommendation, int, error)
	
	// Hide hides a recommendation
	Hide(ctx context.Context, recommendationID string) error
	
	// GetPublicationIDs retrieves publication IDs from recommendations
	GetPublicationIDs(ctx context.Context, userID string, limit int) ([]string, error)
}

