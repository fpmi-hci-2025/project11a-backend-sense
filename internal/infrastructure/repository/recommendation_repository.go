package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"sense-backend/internal/domain"
)

type recommendationRepository struct {
	pool *pgxpool.Pool
}

// NewRecommendationRepository creates a new recommendation repository
func NewRecommendationRepository(pool *pgxpool.Pool) domain.RecommendationRepository {
	return &recommendationRepository{pool: pool}
}

func (r *recommendationRepository) Create(ctx context.Context, recommendation *domain.Recommendation) error {
	query := `
		INSERT INTO recommendations (id, user_id, publication_id, algorithm, reason, rank, created_at, hidden)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.pool.Exec(ctx, query,
		recommendation.ID, recommendation.UserID, recommendation.PublicationID,
		recommendation.Algorithm, recommendation.Reason, recommendation.Rank,
		recommendation.CreatedAt, recommendation.Hidden,
	)
	return err
}

func (r *recommendationRepository) GetByUser(ctx context.Context, userID string, limit, offset int) ([]*domain.Recommendation, int, error) {
	// Get total
	var total int
	err := r.pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM recommendations WHERE user_id = $1 AND hidden = false
	`, userID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get recommendations
	rows, err := r.pool.Query(ctx, `
		SELECT id, user_id, publication_id, algorithm, reason, rank, created_at, hidden
		FROM recommendations
		WHERE user_id = $1 AND hidden = false
		ORDER BY rank ASC, created_at DESC
		LIMIT $2 OFFSET $3
	`, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var recommendations []*domain.Recommendation
	for rows.Next() {
		var rec domain.Recommendation
		err := rows.Scan(
			&rec.ID, &rec.UserID, &rec.PublicationID, &rec.Algorithm,
			&rec.Reason, &rec.Rank, &rec.CreatedAt, &rec.Hidden,
		)
		if err != nil {
			return nil, 0, err
		}
		recommendations = append(recommendations, &rec)
	}

	return recommendations, total, rows.Err()
}

func (r *recommendationRepository) Hide(ctx context.Context, recommendationID string) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE recommendations SET hidden = true WHERE id = $1
	`, recommendationID)
	return err
}

func (r *recommendationRepository) GetPublicationIDs(ctx context.Context, userID string, limit int) ([]string, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT publication_id
		FROM recommendations
		WHERE user_id = $1 AND hidden = false
		ORDER BY rank ASC, created_at DESC
		LIMIT $2
	`, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var publicationIDs []string
	for rows.Next() {
		var publicationID string
		if err := rows.Scan(&publicationID); err != nil {
			return nil, err
		}
		publicationIDs = append(publicationIDs, publicationID)
	}

	return publicationIDs, rows.Err()
}

