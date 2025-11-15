package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"sense-backend/internal/domain"
)

type tagRepository struct {
	pool *pgxpool.Pool
}

// NewTagRepository creates a new tag repository
func NewTagRepository(pool *pgxpool.Pool) domain.TagRepository {
	return &tagRepository{pool: pool}
}

func (r *tagRepository) Create(ctx context.Context, tag *domain.Tag) error {
	query := `
		INSERT INTO tags (id, name, description, usage_count, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.pool.Exec(ctx, query,
		tag.ID, tag.Name, tag.Description, tag.UsageCount, tag.CreatedAt,
	)
	return err
}

func (r *tagRepository) GetByID(ctx context.Context, id string) (*domain.Tag, error) {
	query := `
		SELECT id, name, description, usage_count, created_at
		FROM tags
		WHERE id = $1
	`
	
	var tag domain.Tag
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&tag.ID, &tag.Name, &tag.Description, &tag.UsageCount, &tag.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("tag not found")
	}
	return &tag, err
}

func (r *tagRepository) GetByName(ctx context.Context, name string) (*domain.Tag, error) {
	query := `
		SELECT id, name, description, usage_count, created_at
		FROM tags
		WHERE name = $1
	`
	
	var tag domain.Tag
	err := r.pool.QueryRow(ctx, query, name).Scan(
		&tag.ID, &tag.Name, &tag.Description, &tag.UsageCount, &tag.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("tag not found")
	}
	return &tag, err
}

func (r *tagRepository) GetPopular(ctx context.Context, limit int, search *string) ([]*domain.Tag, int, error) {
	where := "1=1"
	args := []interface{}{}
	argIndex := 1

	if search != nil {
		where = fmt.Sprintf("name ILIKE $%d", argIndex)
		args = append(args, "%"+*search+"%")
		argIndex++
	}

	// Get total
	var total int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM tags WHERE %s", where)
	err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get tags
	query := fmt.Sprintf(`
		SELECT id, name, description, usage_count, created_at
		FROM tags
		WHERE %s
		ORDER BY usage_count DESC, name ASC
		LIMIT $%d
	`, where, argIndex)
	args = append(args, limit)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var tags []*domain.Tag
	for rows.Next() {
		var tag domain.Tag
		err := rows.Scan(
			&tag.ID, &tag.Name, &tag.Description, &tag.UsageCount, &tag.CreatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		tags = append(tags, &tag)
	}

	return tags, total, rows.Err()
}

func (r *tagRepository) AttachToPublication(ctx context.Context, publicationID, tagID string) error {
	query := `
		INSERT INTO publication_tags (publication_id, tag_id, created_at)
		VALUES ($1, $2, NOW())
		ON CONFLICT (publication_id, tag_id) DO NOTHING
	`
	_, err := r.pool.Exec(ctx, query, publicationID, tagID)
	return err
}

func (r *tagRepository) DetachFromPublication(ctx context.Context, publicationID, tagID string) error {
	_, err := r.pool.Exec(ctx, `
		DELETE FROM publication_tags WHERE publication_id = $1 AND tag_id = $2
	`, publicationID, tagID)
	return err
}

func (r *tagRepository) GetByPublication(ctx context.Context, publicationID string) ([]*domain.Tag, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT t.id, t.name, t.description, t.usage_count, t.created_at
		FROM tags t
		INNER JOIN publication_tags pt ON t.id = pt.tag_id
		WHERE pt.publication_id = $1
		ORDER BY t.name
	`, publicationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []*domain.Tag
	for rows.Next() {
		var tag domain.Tag
		err := rows.Scan(
			&tag.ID, &tag.Name, &tag.Description, &tag.UsageCount, &tag.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		tags = append(tags, &tag)
	}

	return tags, rows.Err()
}

