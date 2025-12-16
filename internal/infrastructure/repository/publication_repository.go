package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"sense-backend/internal/domain"
)

type publicationRepository struct {
	pool *pgxpool.Pool
}

// NewPublicationRepository creates a new publication repository
func NewPublicationRepository(pool *pgxpool.Pool) domain.PublicationRepository {
	return &publicationRepository{pool: pool}
}

func (r *publicationRepository) Create(ctx context.Context, publication *domain.Publication, mediaIDs []string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	// Insert publication
	query := `
		INSERT INTO publications (id, author_id, type, content, source, publication_date, visibility)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err = tx.Exec(ctx, query,
		publication.ID, publication.AuthorID, publication.Type, publication.Content,
		publication.Source, publication.PublicationDate, publication.Visibility,
	)
	if err != nil {
		return err
	}

	// Link media
	for i, mediaID := range mediaIDs {
		_, err = tx.Exec(ctx, `
			INSERT INTO publication_media (publication_id, media_id, ord)
			VALUES ($1, $2, $3)
		`, publication.ID, mediaID, i)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *publicationRepository) GetByID(ctx context.Context, id string) (*domain.Publication, error) {
	query := `
		SELECT id, author_id, type, content, source, publication_date, visibility,
		       likes_count, comments_count, saved_count
		FROM publications
		WHERE id = $1
	`
	
	var pub domain.Publication
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&pub.ID, &pub.AuthorID, &pub.Type, &pub.Content, &pub.Source,
		&pub.PublicationDate, &pub.Visibility, &pub.LikesCount,
		&pub.CommentsCount, &pub.SavedCount,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("publication not found")
	}
	return &pub, err
}

func (r *publicationRepository) Update(ctx context.Context, publication *domain.Publication, mediaIDs []string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	// Update publication
	query := `
		UPDATE publications
		SET content = $2, source = $3, visibility = $4
		WHERE id = $1
	`
	_, err = tx.Exec(ctx, query,
		publication.ID, publication.Content, publication.Source, publication.Visibility,
	)
	if err != nil {
		return err
	}

	// Remove old media links
	_, err = tx.Exec(ctx, `DELETE FROM publication_media WHERE publication_id = $1`, publication.ID)
	if err != nil {
		return err
	}

	// Add new media links
	for i, mediaID := range mediaIDs {
		_, err = tx.Exec(ctx, `
			INSERT INTO publication_media (publication_id, media_id, ord)
			VALUES ($1, $2, $3)
		`, publication.ID, mediaID, i)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *publicationRepository) Delete(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM publications WHERE id = $1`, id)
	return err
}

func (r *publicationRepository) GetFeed(ctx context.Context, userID *string, filters *domain.FeedFilters, limit, offset int) ([]*domain.Publication, int, error) {
	where := []string{"1=1"}
	args := []interface{}{}
	argIndex := 1

	if filters != nil {
		if filters.Type != nil {
			where = append(where, fmt.Sprintf("type = $%d", argIndex))
			args = append(args, *filters.Type)
			argIndex++
		}
		if filters.Visibility != nil {
			where = append(where, fmt.Sprintf("visibility = $%d", argIndex))
			args = append(args, *filters.Visibility)
			argIndex++
		}
		if filters.AuthorID != nil {
			where = append(where, fmt.Sprintf("author_id = $%d", argIndex))
			args = append(args, *filters.AuthorID)
			argIndex++
		}
		if filters.DateFrom != nil {
			where = append(where, fmt.Sprintf("publication_date >= $%d", argIndex))
			args = append(args, *filters.DateFrom)
			argIndex++
		}
		if filters.DateTo != nil {
			where = append(where, fmt.Sprintf("publication_date <= $%d", argIndex))
			args = append(args, *filters.DateTo)
			argIndex++
		}
	}

	// If userID provided, filter by visibility (public or community for logged in users)
	if userID != nil {
		where = append(where, fmt.Sprintf("(visibility = 'public' OR visibility = 'community' OR author_id = $%d)", argIndex))
		args = append(args, *userID)
		argIndex++
	} else {
		where = append(where, "visibility = 'public'")
	}

	whereClause := strings.Join(where, " AND ")

	// Get total
	var total int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM publications WHERE %s", whereClause)
	err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get publications
	query := fmt.Sprintf(`
		SELECT id, author_id, type, content, source, publication_date, visibility,
		       likes_count, comments_count, saved_count
		FROM publications
		WHERE %s
		ORDER BY publication_date DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)
	args = append(args, limit, offset)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var publications []*domain.Publication
	for rows.Next() {
		var pub domain.Publication
		err := rows.Scan(
			&pub.ID, &pub.AuthorID, &pub.Type, &pub.Content, &pub.Source,
			&pub.PublicationDate, &pub.Visibility, &pub.LikesCount,
			&pub.CommentsCount, &pub.SavedCount,
		)
		if err != nil {
			return nil, 0, err
		}
		publications = append(publications, &pub)
	}

	return publications, total, rows.Err()
}

func (r *publicationRepository) GetByAuthor(ctx context.Context, authorID string, filters *domain.PublicationFilters, limit, offset int) ([]*domain.Publication, int, error) {
	where := []string{"author_id = $1"}
	args := []interface{}{authorID}
	argIndex := 2

	if filters != nil {
		if filters.Type != nil {
			where = append(where, fmt.Sprintf("type = $%d", argIndex))
			args = append(args, *filters.Type)
			argIndex++
		}
		if filters.Visibility != nil {
			where = append(where, fmt.Sprintf("visibility = $%d", argIndex))
			args = append(args, *filters.Visibility)
			argIndex++
		}
	}

	whereClause := strings.Join(where, " AND ")

	// Get total
	var total int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM publications WHERE %s", whereClause)
	err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get publications
	query := fmt.Sprintf(`
		SELECT id, author_id, type, content, source, publication_date, visibility,
		       likes_count, comments_count, saved_count
		FROM publications
		WHERE %s
		ORDER BY publication_date DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)
	args = append(args, limit, offset)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var publications []*domain.Publication
	for rows.Next() {
		var pub domain.Publication
		err := rows.Scan(
			&pub.ID, &pub.AuthorID, &pub.Type, &pub.Content, &pub.Source,
			&pub.PublicationDate, &pub.Visibility, &pub.LikesCount,
			&pub.CommentsCount, &pub.SavedCount,
		)
		if err != nil {
			return nil, 0, err
		}
		publications = append(publications, &pub)
	}

	return publications, total, rows.Err()
}

func (r *publicationRepository) Like(ctx context.Context, userID, publicationID string) (bool, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return false, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	// Check if already liked
	var exists bool
	err = tx.QueryRow(ctx, `
		SELECT EXISTS(SELECT 1 FROM publication_likes WHERE user_id = $1 AND publication_id = $2)
	`, userID, publicationID).Scan(&exists)
	if err != nil {
		return false, err
	}

	if exists {
		// Unlike
		_, err = tx.Exec(ctx, `
			DELETE FROM publication_likes WHERE user_id = $1 AND publication_id = $2
		`, userID, publicationID)
		if err != nil {
			return false, err
		}
		_, err = tx.Exec(ctx, `
			UPDATE publications SET likes_count = likes_count - 1 WHERE id = $1
		`, publicationID)
		if err != nil {
			return false, err
		}
		return false, tx.Commit(ctx)
	} else {
		// Like
		_, err = tx.Exec(ctx, `
			INSERT INTO publication_likes (user_id, publication_id) VALUES ($1, $2)
		`, userID, publicationID)
		if err != nil {
			return false, err
		}
		_, err = tx.Exec(ctx, `
			UPDATE publications SET likes_count = likes_count + 1 WHERE id = $1
		`, publicationID)
		if err != nil {
			return false, err
		}
		return true, tx.Commit(ctx)
	}
}

func (r *publicationRepository) IsLiked(ctx context.Context, userID, publicationID string) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx, `
		SELECT EXISTS(SELECT 1 FROM publication_likes WHERE user_id = $1 AND publication_id = $2)
	`, userID, publicationID).Scan(&exists)
	return exists, err
}

func (r *publicationRepository) GetLikesCount(ctx context.Context, publicationID string) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx, `
		SELECT likes_count FROM publications WHERE id = $1
	`, publicationID).Scan(&count)
	return count, err
}

func (r *publicationRepository) GetLikedUsers(ctx context.Context, publicationID string, limit, offset int) ([]*domain.User, int, error) {
	// Get total
	var total int
	err := r.pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM publication_likes WHERE publication_id = $1
	`, publicationID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get users
	rows, err := r.pool.Query(ctx, `
		SELECT u.id, u.username, u.email, u.phone, u.icon_url, u.description, u.role, u.registered_at
		FROM users u
		INNER JOIN publication_likes pl ON u.id = pl.user_id
		WHERE pl.publication_id = $1
		ORDER BY pl.created_at DESC
		LIMIT $2 OFFSET $3
	`, publicationID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var user domain.User
		err := rows.Scan(
			&user.ID, &user.Username, &user.Email, &user.Phone, &user.IconURL,
			&user.Description, &user.Role, &user.RegisteredAt,
		)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, &user)
	}

	return users, total, rows.Err()
}

func (r *publicationRepository) Save(ctx context.Context, userID, publicationID string, note *string) error {
	query := `
		INSERT INTO saved_items (user_id, publication_id, note)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, publication_id) DO UPDATE SET note = $3
	`
	_, err := r.pool.Exec(ctx, query, userID, publicationID, note)
	if err != nil {
		return err
	}
	_, err = r.pool.Exec(ctx, `
		UPDATE publications SET saved_count = saved_count + 1 WHERE id = $1
	`, publicationID)
	return err
}

func (r *publicationRepository) Unsave(ctx context.Context, userID, publicationID string) error {
	_, err := r.pool.Exec(ctx, `
		DELETE FROM saved_items WHERE user_id = $1 AND publication_id = $2
	`, userID, publicationID)
	if err != nil {
		return err
	}
	_, err = r.pool.Exec(ctx, `
		UPDATE publications SET saved_count = saved_count - 1 WHERE id = $1
	`, publicationID)
	return err
}

func (r *publicationRepository) IsSaved(ctx context.Context, userID, publicationID string) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx, `
		SELECT EXISTS(SELECT 1 FROM saved_items WHERE user_id = $1 AND publication_id = $2)
	`, userID, publicationID).Scan(&exists)
	return exists, err
}

func (r *publicationRepository) GetSaved(ctx context.Context, userID string, filters *domain.PublicationFilters, limit, offset int) ([]*domain.SavedPublication, int, error) {
	where := []string{"si.user_id = $1"}
	args := []interface{}{userID}
	argIndex := 2

	if filters != nil {
		if filters.Type != nil {
			where = append(where, fmt.Sprintf("p.type = $%d", argIndex))
			args = append(args, *filters.Type)
			argIndex++
		}
		if filters.Visibility != nil {
			where = append(where, fmt.Sprintf("p.visibility = $%d", argIndex))
			args = append(args, *filters.Visibility)
			argIndex++
		}
	}

	whereClause := strings.Join(where, " AND ")

	// Get total
	var total int
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) FROM saved_items si
		INNER JOIN publications p ON si.publication_id = p.id
		WHERE %s
	`, whereClause)
	err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get saved publications
	query := fmt.Sprintf(`
		SELECT p.id, p.author_id, p.type, p.content, p.source, p.publication_date, p.visibility,
		       p.likes_count, p.comments_count, p.saved_count,
		       si.note, si.added_at
		FROM saved_items si
		INNER JOIN publications p ON si.publication_id = p.id
		WHERE %s
		ORDER BY si.added_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)
	args = append(args, limit, offset)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var saved []*domain.SavedPublication
	for rows.Next() {
		var sp domain.SavedPublication
		err := rows.Scan(
			&sp.ID, &sp.AuthorID, &sp.Type, &sp.Content, &sp.Source,
			&sp.PublicationDate, &sp.Visibility, &sp.LikesCount,
			&sp.CommentsCount, &sp.SavedCount, &sp.SavedNote, &sp.SavedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		saved = append(saved, &sp)
	}

	return saved, total, rows.Err()
}

func (r *publicationRepository) Search(ctx context.Context, query string, filters *domain.SearchFilters, limit, offset int) ([]*domain.Publication, int, error) {
	searchQuery := `%` + query + `%`
	where := []string{"content ILIKE $1"}
	args := []interface{}{searchQuery}
	argIndex := 2

	if filters != nil {
		if filters.Type != nil {
			where = append(where, fmt.Sprintf("type = $%d", argIndex))
			args = append(args, *filters.Type)
			argIndex++
		}
		if filters.Visibility != nil {
			where = append(where, fmt.Sprintf("visibility = $%d", argIndex))
			args = append(args, *filters.Visibility)
			argIndex++
		}
		if filters.AuthorID != nil {
			where = append(where, fmt.Sprintf("author_id = $%d", argIndex))
			args = append(args, *filters.AuthorID)
			argIndex++
		}
	}

	whereClause := strings.Join(where, " AND ")

	// Get total
	var total int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM publications WHERE %s", whereClause)
	err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get publications
	queryStr := fmt.Sprintf(`
		SELECT id, author_id, type, content, source, publication_date, visibility,
		       likes_count, comments_count, saved_count
		FROM publications
		WHERE %s
		ORDER BY publication_date DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)
	args = append(args, limit, offset)

	rows, err := r.pool.Query(ctx, queryStr, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var publications []*domain.Publication
	for rows.Next() {
		var pub domain.Publication
		err := rows.Scan(
			&pub.ID, &pub.AuthorID, &pub.Type, &pub.Content, &pub.Source,
			&pub.PublicationDate, &pub.Visibility, &pub.LikesCount,
			&pub.CommentsCount, &pub.SavedCount,
		)
		if err != nil {
			return nil, 0, err
		}
		publications = append(publications, &pub)
	}

	return publications, total, rows.Err()
}

func (r *publicationRepository) GetMediaIDs(ctx context.Context, publicationID string) ([]string, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT media_id FROM publication_media
		WHERE publication_id = $1
		ORDER BY ord
	`, publicationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mediaIDs []string
	for rows.Next() {
		var mediaID string
		if err := rows.Scan(&mediaID); err != nil {
			return nil, err
		}
		mediaIDs = append(mediaIDs, mediaID)
	}

	return mediaIDs, rows.Err()
}

