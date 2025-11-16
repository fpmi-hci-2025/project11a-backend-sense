package ai

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"sense-backend/internal/domain"
	aiClient "sense-backend/internal/infrastructure/ai"
)

// UseCase handles AI use cases
type UseCase struct {
	aiClient          aiClient.ClientInterface
	recommendationRepo domain.RecommendationRepository
	publicationRepo    domain.PublicationRepository
}

// NewUseCase creates a new AI use case
func NewUseCase(
	aiClient aiClient.ClientInterface,
	recommendationRepo domain.RecommendationRepository,
	publicationRepo domain.PublicationRepository,
) *UseCase {
	return &UseCase{
		aiClient:          aiClient,
		recommendationRepo: recommendationRepo,
		publicationRepo:    publicationRepo,
	}
}

// PurifyRequest represents purify text request
type PurifyRequest struct {
	Text string `json:"text" validate:"required,max=10000"`
}

// PurifyResponse represents purify text response
type PurifyResponse struct {
	CleanedText string  `json:"cleaned_text"`
	IsClean     bool    `json:"is_clean"`
	Confidence  float64 `json:"confidence,omitempty"`
}

// GetRecommendations retrieves recommendations for user
func (uc *UseCase) GetRecommendations(ctx context.Context, userID string, limit int, algorithm *string) ([]*domain.Publication, error) {
	// Get recommendations from AI service
	recommendResp, err := uc.aiClient.Recommend(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get AI recommendations: %w", err)
	}

	// Get publications
	var publications []*domain.Publication
	for i, publicationID := range recommendResp.Publications {
		if i >= limit {
			break
		}

		pub, err := uc.publicationRepo.GetByID(ctx, publicationID)
		if err != nil {
			continue // Skip if publication not found
		}

		// Save recommendation to database
		alg := "collaborative_filtering"
		if algorithm != nil {
			alg = *algorithm
		}

		recommendation := &domain.Recommendation{
			ID:            uuid.New().String(),
			UserID:        userID,
			PublicationID: publicationID,
			Algorithm:     alg,
			Reason:        fmt.Sprintf("Recommended by %s", alg),
			Rank:          i,
			CreatedAt:     time.Now(),
			Hidden:        false,
		}
		_ = uc.recommendationRepo.Create(ctx, recommendation) // Ignore errors

		publications = append(publications, pub)
	}

	return publications, nil
}

// GetRecommendationsFeed retrieves recommendations feed
func (uc *UseCase) GetRecommendationsFeed(ctx context.Context, userID string, limit, offset int, algorithm *string) ([]*domain.Publication, int, error) {
	// Get publication IDs from recommendations
	publicationIDs, err := uc.recommendationRepo.GetPublicationIDs(ctx, userID, limit+offset)
	if err != nil {
		return nil, 0, err
	}

	// Apply offset
	if offset < len(publicationIDs) {
		publicationIDs = publicationIDs[offset:]
	} else {
		publicationIDs = []string{}
	}

	// Limit results
	if limit < len(publicationIDs) {
		publicationIDs = publicationIDs[:limit]
	}

	// Get publications
	var publications []*domain.Publication
	for _, publicationID := range publicationIDs {
		pub, err := uc.publicationRepo.GetByID(ctx, publicationID)
		if err != nil {
			continue
		}
		publications = append(publications, pub)
	}

	return publications, len(publications), nil
}

// HideRecommendation hides a recommendation
func (uc *UseCase) HideRecommendation(ctx context.Context, recommendationID string) error {
	return uc.recommendationRepo.Hide(ctx, recommendationID)
}

// PurifyText purifies text using AI service
func (uc *UseCase) PurifyText(ctx context.Context, req *PurifyRequest) (*PurifyResponse, error) {
	// Use AI compose endpoint for text purification
	composeReq := &aiClient.ComposeRequest{
		Query: fmt.Sprintf("Purify and clean this text: %s", req.Text),
		Metadata: map[string]interface{}{
			"task": "purify",
		},
	}

	composeResp, err := uc.aiClient.Compose(ctx, composeReq)
	if err != nil {
		return nil, fmt.Errorf("failed to purify text: %w", err)
	}

	// Simple check if text was modified
	isClean := composeResp.Text == req.Text
	confidence := 0.95
	if !isClean {
		confidence = 0.85
	}

	return &PurifyResponse{
		CleanedText: composeResp.Text,
		IsClean:     isClean,
		Confidence:  confidence,
	}, nil
}

