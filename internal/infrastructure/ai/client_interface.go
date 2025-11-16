package ai

import "context"

// ClientInterface defines interface for AI client operations
// This interface is used for mocking in tests
type ClientInterface interface {
	HealthCheck(ctx context.Context) error
	Compose(ctx context.Context, req *ComposeRequest) (*ComposeResponse, error)
	Recommend(ctx context.Context, userID string) (*RecommendResponse, error)
}

// Ensure Client implements ClientInterface
var _ ClientInterface = (*Client)(nil)

