package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"sense-backend/internal/domain"
)

type mediaRepository struct {
	pool *pgxpool.Pool
}

// NewMediaRepository creates a new media repository
func NewMediaRepository(pool *pgxpool.Pool) domain.MediaRepository {
	return &mediaRepository{pool: pool}
}

func (r *mediaRepository) Create(ctx context.Context, media *domain.MediaAsset) error {
	query := `
		INSERT INTO media_assets (id, owner_id, data, filename, mime, width, height, exif, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	
	var exifJSON interface{}
	if media.EXIF != nil {
		exifJSON = media.EXIF
	}
	
	_, err := r.pool.Exec(ctx, query,
		media.ID, media.OwnerID, media.Data, media.Filename, media.MIME,
		media.Width, media.Height, exifJSON, media.CreatedAt,
	)
	return err
}

func (r *mediaRepository) GetByID(ctx context.Context, id string) (*domain.MediaAsset, error) {
	query := `
		SELECT id, owner_id, filename, mime, width, height, exif, created_at, data
		FROM media_assets
		WHERE id = $1
	`
	
	var media domain.MediaAsset
	var exifJSON interface{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&media.ID, &media.OwnerID, &media.Filename, &media.MIME,
		&media.Width, &media.Height, &exifJSON, &media.CreatedAt, &media.Data,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("media not found")
	}
	if err != nil {
		return nil, err
	}
	
	// Parse EXIF if present
	if exifJSON != nil {
		// EXIF is stored as JSONB, will be handled by pgx
		media.EXIF = &domain.EXIFData{}
		// Note: Full EXIF parsing would require JSON unmarshaling
	}
	
	return &media, nil
}

func (r *mediaRepository) GetByOwner(ctx context.Context, ownerID string, limit, offset int) ([]*domain.MediaAsset, int, error) {
	// Get total
	var total int
	err := r.pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM media_assets WHERE owner_id = $1
	`, ownerID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get media (without data for list)
	rows, err := r.pool.Query(ctx, `
		SELECT id, owner_id, filename, mime, width, height, exif, created_at
		FROM media_assets
		WHERE owner_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`, ownerID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var mediaList []*domain.MediaAsset
	for rows.Next() {
		var media domain.MediaAsset
		var exifJSON interface{}
		err := rows.Scan(
			&media.ID, &media.OwnerID, &media.Filename, &media.MIME,
			&media.Width, &media.Height, &exifJSON, &media.CreatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		if exifJSON != nil {
			media.EXIF = &domain.EXIFData{}
		}
		mediaList = append(mediaList, &media)
	}

	return mediaList, total, rows.Err()
}

func (r *mediaRepository) Delete(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM media_assets WHERE id = $1`, id)
	return err
}

func (r *mediaRepository) CheckOwnership(ctx context.Context, mediaID, userID string) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx, `
		SELECT EXISTS(SELECT 1 FROM media_assets WHERE id = $1 AND owner_id = $2)
	`, mediaID, userID).Scan(&exists)
	return exists, err
}

