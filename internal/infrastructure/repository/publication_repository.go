package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"sense-backend/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
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
		INSERT INTO publications (id, author_id, type, title, content, source, publication_date, visibility)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err = tx.Exec(ctx, query,
		publication.ID, publication.AuthorID, publication.Type, publication.Title, publication.Content,
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
		SELECT p.id, p.author_id, p.type, p.title, p.content, p.source, p.publication_date, p.visibility,
		       COALESCE(likes.count, 0) as likes_count,
		       COALESCE(comments.count, 0) as comments_count,
		       COALESCE(saved.count, 0) as saved_count
		FROM publications p
		LEFT JOIN (SELECT publication_id, COUNT(*) as count FROM publication_likes GROUP BY publication_id) likes 
		  ON p.id = likes.publication_id
		LEFT JOIN (SELECT publication_id, COUNT(*) as count FROM comments GROUP BY publication_id) comments 
		  ON p.id = comments.publication_id
		LEFT JOIN (SELECT publication_id, COUNT(*) as count FROM saved_items GROUP BY publication_id) saved 
		  ON p.id = saved.publication_id
		WHERE p.id = $1
	`

	var pub domain.Publication
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&pub.ID, &pub.AuthorID, &pub.Type, &pub.Title, &pub.Content, &pub.Source,
		&pub.PublicationDate, &pub.Visibility, &pub.LikesCount,
		&pub.CommentsCount, &pub.SavedCount,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("publication not found")
	}
	return &pub, err
}

func (r *publicationRepository) GetByIDWithLikeStatus(ctx context.Context, id string, viewerUserID *string) (*domain.PublicationWithLikeStatus, error) {
	pub, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	result := &domain.PublicationWithLikeStatus{
		Publication: *pub,
		IsLiked:     false,
	}

	if viewerUserID != nil {
		isLiked, err := r.IsLiked(ctx, *viewerUserID, id)
		if err == nil {
			result.IsLiked = isLiked
		}
	}

	return result, nil
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
		SET title = $2, content = $3, source = $4, visibility = $5
		WHERE id = $1
	`
	_, err = tx.Exec(ctx, query,
		publication.ID, publication.Title, publication.Content, publication.Source, publication.Visibility,
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

func (r *publicationRepository) GetFeed(ctx context.Context, userID *string, filters *domain.FeedFilters, limit, offset int) ([]*domain.PublicationWithLikeStatus, int, error) {
	where := []string{"1=1"}
	args := []interface{}{}
	argIndex := 1

	// userID parameter index for LEFT JOIN (will be set if userID is provided)
	var userIDArgIndex int

	if userID != nil {
		userIDArgIndex = argIndex
		args = append(args, *userID)
		argIndex++
	}

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
		if filters.AuthorID != nil {
			where = append(where, fmt.Sprintf("p.author_id = $%d", argIndex))
			args = append(args, *filters.AuthorID)
			argIndex++
		}
		if filters.DateFrom != nil {
			where = append(where, fmt.Sprintf("p.publication_date >= $%d", argIndex))
			args = append(args, *filters.DateFrom)
			argIndex++
		}
		if filters.DateTo != nil {
			where = append(where, fmt.Sprintf("p.publication_date <= $%d", argIndex))
			args = append(args, *filters.DateTo)
			argIndex++
		}
	}

	// If userID provided, filter by visibility (public or community for logged in users)
	if userID != nil {
		where = append(where, fmt.Sprintf("(p.visibility = 'public' OR p.visibility = 'community' OR p.author_id = $%d)", userIDArgIndex))
	} else {
		where = append(where, "p.visibility = 'public'")
	}

	whereClause := strings.Join(where, " AND ")

	// Get total
	var total int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM publications p WHERE %s", whereClause)
	err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Build query with LEFT JOIN for like status and dynamic counts
	var query string
	if userID != nil {
		query = fmt.Sprintf(`
			SELECT p.id, p.author_id, p.type, p.title, p.content, p.source, p.publication_date, p.visibility,
			       COALESCE(likes.count, 0) as likes_count,
			       COALESCE(comments.count, 0) as comments_count,
			       COALESCE(saved.count, 0) as saved_count,
			       CASE WHEN pl.user_id IS NOT NULL THEN true ELSE false END as is_liked
			FROM publications p
			LEFT JOIN publication_likes pl ON p.id = pl.publication_id AND pl.user_id = $%d
			LEFT JOIN (SELECT publication_id, COUNT(*) as count FROM publication_likes GROUP BY publication_id) likes 
			  ON p.id = likes.publication_id
			LEFT JOIN (SELECT publication_id, COUNT(*) as count FROM comments GROUP BY publication_id) comments 
			  ON p.id = comments.publication_id
			LEFT JOIN (SELECT publication_id, COUNT(*) as count FROM saved_items GROUP BY publication_id) saved 
			  ON p.id = saved.publication_id
			WHERE %s
			ORDER BY p.publication_date DESC
			LIMIT $%d OFFSET $%d
		`, userIDArgIndex, whereClause, argIndex, argIndex+1)
	} else {
		query = fmt.Sprintf(`
			SELECT p.id, p.author_id, p.type, p.title, p.content, p.source, p.publication_date, p.visibility,
			       COALESCE(likes.count, 0) as likes_count,
			       COALESCE(comments.count, 0) as comments_count,
			       COALESCE(saved.count, 0) as saved_count,
			       false as is_liked
			FROM publications p
			LEFT JOIN (SELECT publication_id, COUNT(*) as count FROM publication_likes GROUP BY publication_id) likes 
			  ON p.id = likes.publication_id
			LEFT JOIN (SELECT publication_id, COUNT(*) as count FROM comments GROUP BY publication_id) comments 
			  ON p.id = comments.publication_id
			LEFT JOIN (SELECT publication_id, COUNT(*) as count FROM saved_items GROUP BY publication_id) saved 
			  ON p.id = saved.publication_id
			WHERE %s
			ORDER BY p.publication_date DESC
			LIMIT $%d OFFSET $%d
		`, whereClause, argIndex, argIndex+1)
	}
	args = append(args, limit, offset)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var publications []*domain.PublicationWithLikeStatus
	for rows.Next() {
		var pub domain.PublicationWithLikeStatus
		err := rows.Scan(
			&pub.ID, &pub.AuthorID, &pub.Type, &pub.Title, &pub.Content, &pub.Source,
			&pub.PublicationDate, &pub.Visibility, &pub.LikesCount,
			&pub.CommentsCount, &pub.SavedCount, &pub.IsLiked,
		)
		if err != nil {
			return nil, 0, err
		}
		publications = append(publications, &pub)
	}

	return publications, total, rows.Err()
}

func (r *publicationRepository) GetByAuthor(ctx context.Context, authorID string, viewerUserID *string, filters *domain.PublicationFilters, limit, offset int) ([]*domain.PublicationWithLikeStatus, int, error) {
	// Build WHERE for count query (author + filters; no viewer to keep placeholders dense)
	countWhere := []string{"p.author_id = $1"}
	countArgs := []interface{}{authorID}
	countIdx := 2
	if filters != nil {
		if filters.Type != nil {
			countWhere = append(countWhere, fmt.Sprintf("p.type = $%d", countIdx))
			countArgs = append(countArgs, *filters.Type)
			countIdx++
		}
		if filters.Visibility != nil {
			countWhere = append(countWhere, fmt.Sprintf("p.visibility = $%d", countIdx))
			countArgs = append(countArgs, *filters.Visibility)
			countIdx++
		}
	}
	countWhereClause := strings.Join(countWhere, " AND ")

	var total int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM publications p WHERE %s", countWhereClause)
	if err := r.pool.QueryRow(ctx, countQuery, countArgs...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Build WHERE for main query (author + optional viewer + filters)
	queryWhere := []string{"p.author_id = $1"}
	queryArgs := []interface{}{authorID}
	queryIdx := 2
	viewerUserIDArgIndex := 0
	if viewerUserID != nil {
		viewerUserIDArgIndex = queryIdx
		queryArgs = append(queryArgs, *viewerUserID)
		queryIdx++
	}
	if filters != nil {
		if filters.Type != nil {
			queryWhere = append(queryWhere, fmt.Sprintf("p.type = $%d", queryIdx))
			queryArgs = append(queryArgs, *filters.Type)
			queryIdx++
		}
		if filters.Visibility != nil {
			queryWhere = append(queryWhere, fmt.Sprintf("p.visibility = $%d", queryIdx))
			queryArgs = append(queryArgs, *filters.Visibility)
			queryIdx++
		}
	}
	queryWhereClause := strings.Join(queryWhere, " AND ")

	limitPlaceholder := queryIdx
	offsetPlaceholder := queryIdx + 1

	var query string
	if viewerUserID != nil {
		query = fmt.Sprintf(`
			SELECT p.id, p.author_id, p.type, p.title, p.content, p.source, p.publication_date, p.visibility,
			       COALESCE(likes.count, 0) as likes_count,
			       COALESCE(comments.count, 0) as comments_count,
			       COALESCE(saved.count, 0) as saved_count,
			       CASE WHEN pl.user_id IS NOT NULL THEN true ELSE false END as is_liked
			FROM publications p
			LEFT JOIN publication_likes pl ON p.id = pl.publication_id AND pl.user_id = $%d
			LEFT JOIN (SELECT publication_id, COUNT(*) as count FROM publication_likes GROUP BY publication_id) likes 
			  ON p.id = likes.publication_id
			LEFT JOIN (SELECT publication_id, COUNT(*) as count FROM comments GROUP BY publication_id) comments 
			  ON p.id = comments.publication_id
			LEFT JOIN (SELECT publication_id, COUNT(*) as count FROM saved_items GROUP BY publication_id) saved 
			  ON p.id = saved.publication_id
			WHERE %s
			ORDER BY p.publication_date DESC
			LIMIT $%d OFFSET $%d
		`, viewerUserIDArgIndex, queryWhereClause, limitPlaceholder, offsetPlaceholder)
	} else {
		query = fmt.Sprintf(`
			SELECT p.id, p.author_id, p.type, p.title, p.content, p.source, p.publication_date, p.visibility,
			       COALESCE(likes.count, 0) as likes_count,
			       COALESCE(comments.count, 0) as comments_count,
			       COALESCE(saved.count, 0) as saved_count,
			       false as is_liked
			FROM publications p
			LEFT JOIN (SELECT publication_id, COUNT(*) as count FROM publication_likes GROUP BY publication_id) likes 
			  ON p.id = likes.publication_id
			LEFT JOIN (SELECT publication_id, COUNT(*) as count FROM comments GROUP BY publication_id) comments 
			  ON p.id = comments.publication_id
			LEFT JOIN (SELECT publication_id, COUNT(*) as count FROM saved_items GROUP BY publication_id) saved 
			  ON p.id = saved.publication_id
			WHERE %s
			ORDER BY p.publication_date DESC
			LIMIT $%d OFFSET $%d
		`, queryWhereClause, limitPlaceholder, offsetPlaceholder)
	}

	queryArgs = append(queryArgs, limit, offset)

	rows, err := r.pool.Query(ctx, query, queryArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var publications []*domain.PublicationWithLikeStatus
	for rows.Next() {
		var pub domain.PublicationWithLikeStatus
		if err := rows.Scan(
			&pub.ID, &pub.AuthorID, &pub.Type, &pub.Title, &pub.Content, &pub.Source,
			&pub.PublicationDate, &pub.Visibility, &pub.LikesCount,
			&pub.CommentsCount, &pub.SavedCount, &pub.IsLiked,
		); err != nil {
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
		return false, tx.Commit(ctx)
	} else {
		// Like
		_, err = tx.Exec(ctx, `
			INSERT INTO publication_likes (user_id, publication_id) VALUES ($1, $2)
		`, userID, publicationID)
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
		SELECT COUNT(*) FROM publication_likes WHERE publication_id = $1
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
	return err
}

func (r *publicationRepository) Unsave(ctx context.Context, userID, publicationID string) error {
	_, err := r.pool.Exec(ctx, `
		DELETE FROM saved_items WHERE user_id = $1 AND publication_id = $2
	`, userID, publicationID)
	return err
}

func (r *publicationRepository) IsSaved(ctx context.Context, userID, publicationID string) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx, `
		SELECT EXISTS(SELECT 1 FROM saved_items WHERE user_id = $1 AND publication_id = $2)
	`, userID, publicationID).Scan(&exists)
	return exists, err
}

func (r *publicationRepository) GetSaved(ctx context.Context, userID string, filters *domain.PublicationFilters, limit, offset int) ([]*domain.SavedPublicationWithLikeStatus, int, error) {
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

	// Get saved publications with like status (userID is the viewer)
	query := fmt.Sprintf(`
		SELECT p.id, p.author_id, p.type, p.title, p.content, p.source, p.publication_date, p.visibility,
		       COALESCE(likes.count, 0) as likes_count,
		       COALESCE(comments.count, 0) as comments_count,
		       COALESCE(saved.count, 0) as saved_count,
		       si.note, si.added_at,
		       CASE WHEN pl.user_id IS NOT NULL THEN true ELSE false END as is_liked
		FROM saved_items si
		INNER JOIN publications p ON si.publication_id = p.id
		LEFT JOIN publication_likes pl ON p.id = pl.publication_id AND pl.user_id = $1
		LEFT JOIN (SELECT publication_id, COUNT(*) as count FROM publication_likes GROUP BY publication_id) likes 
		  ON p.id = likes.publication_id
		LEFT JOIN (SELECT publication_id, COUNT(*) as count FROM comments GROUP BY publication_id) comments 
		  ON p.id = comments.publication_id
		LEFT JOIN (SELECT publication_id, COUNT(*) as count FROM saved_items GROUP BY publication_id) saved 
		  ON p.id = saved.publication_id
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

	var saved []*domain.SavedPublicationWithLikeStatus
	for rows.Next() {
		var sp domain.SavedPublicationWithLikeStatus
		err := rows.Scan(
			&sp.ID, &sp.AuthorID, &sp.Type, &sp.Title, &sp.Content, &sp.Source,
			&sp.PublicationDate, &sp.Visibility, &sp.LikesCount,
			&sp.CommentsCount, &sp.SavedCount, &sp.SavedNote, &sp.SavedAt, &sp.IsLiked,
		)
		if err != nil {
			return nil, 0, err
		}
		saved = append(saved, &sp)
	}

	return saved, total, rows.Err()
}

func (r *publicationRepository) Search(ctx context.Context, query string, viewerUserID *string, filters *domain.SearchFilters, limit, offset int) ([]*domain.PublicationWithLikeStatus, int, error) {
	searchQuery := `%` + query + `%`
	where := []string{"(p.content ILIKE $1 OR p.title ILIKE $1)"}
	filterValues := []interface{}{}

	// viewerUserID placeholder (used only in JOIN, not in WHERE)
	viewerUserIDArgIndex := 0

	// Filters start after search query and optional viewer placeholder
	filterStartIndex := 2
	if viewerUserID != nil {
		viewerUserIDArgIndex = filterStartIndex
		filterStartIndex++
	}

	if filters != nil {
		if filters.Type != nil {
			where = append(where, fmt.Sprintf("p.type = $%d", filterStartIndex))
			filterValues = append(filterValues, *filters.Type)
			filterStartIndex++
		}
		if filters.Visibility != nil {
			where = append(where, fmt.Sprintf("p.visibility = $%d", filterStartIndex))
			filterValues = append(filterValues, *filters.Visibility)
			filterStartIndex++
		}
		if filters.AuthorID != nil {
			where = append(where, fmt.Sprintf("p.author_id = $%d", filterStartIndex))
			filterValues = append(filterValues, *filters.AuthorID)
			filterStartIndex++
		}
	}

	whereClause := strings.Join(where, " AND ")

	// Count query args (exclude viewerUserID to keep placeholders dense)
	countArgs := []interface{}{searchQuery}
	countArgs = append(countArgs, filterValues...)
	var total int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM publications p WHERE %s", whereClause)
	if err := r.pool.QueryRow(ctx, countQuery, countArgs...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Main query args: search query, optional viewer, filters
	args := []interface{}{searchQuery}
	if viewerUserID != nil {
		viewerUserIDArgIndex = len(args) + 1 // 2
		args = append(args, *viewerUserID)
	}
	args = append(args, filterValues...)

	limitPlaceholder := len(args) + 1
	offsetPlaceholder := len(args) + 2

	var queryStr string
	if viewerUserID != nil {
		queryStr = fmt.Sprintf(`
			SELECT p.id, p.author_id, p.type, p.title, p.content, p.source, p.publication_date, p.visibility,
			       COALESCE(likes.count, 0) as likes_count,
			       COALESCE(comments.count, 0) as comments_count,
			       COALESCE(saved.count, 0) as saved_count,
			       CASE WHEN pl.user_id IS NOT NULL THEN true ELSE false END as is_liked
			FROM publications p
			LEFT JOIN publication_likes pl ON p.id = pl.publication_id AND pl.user_id = $%d
			LEFT JOIN (SELECT publication_id, COUNT(*) as count FROM publication_likes GROUP BY publication_id) likes 
			  ON p.id = likes.publication_id
			LEFT JOIN (SELECT publication_id, COUNT(*) as count FROM comments GROUP BY publication_id) comments 
			  ON p.id = comments.publication_id
			LEFT JOIN (SELECT publication_id, COUNT(*) as count FROM saved_items GROUP BY publication_id) saved 
			  ON p.id = saved.publication_id
			WHERE %s
			ORDER BY p.publication_date DESC
			LIMIT $%d OFFSET $%d
		`, viewerUserIDArgIndex, whereClause, limitPlaceholder, offsetPlaceholder)
	} else {
		queryStr = fmt.Sprintf(`
			SELECT p.id, p.author_id, p.type, p.title, p.content, p.source, p.publication_date, p.visibility,
			       COALESCE(likes.count, 0) as likes_count,
			       COALESCE(comments.count, 0) as comments_count,
			       COALESCE(saved.count, 0) as saved_count,
			       false as is_liked
			FROM publications p
			LEFT JOIN (SELECT publication_id, COUNT(*) as count FROM publication_likes GROUP BY publication_id) likes 
			  ON p.id = likes.publication_id
			LEFT JOIN (SELECT publication_id, COUNT(*) as count FROM comments GROUP BY publication_id) comments 
			  ON p.id = comments.publication_id
			LEFT JOIN (SELECT publication_id, COUNT(*) as count FROM saved_items GROUP BY publication_id) saved 
			  ON p.id = saved.publication_id
			WHERE %s
			ORDER BY p.publication_date DESC
			LIMIT $%d OFFSET $%d
		`, whereClause, limitPlaceholder, offsetPlaceholder)
	}

	args = append(args, limit, offset)

	rows, err := r.pool.Query(ctx, queryStr, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var publications []*domain.PublicationWithLikeStatus
	for rows.Next() {
		var pub domain.PublicationWithLikeStatus
		if err := rows.Scan(
			&pub.ID, &pub.AuthorID, &pub.Type, &pub.Title, &pub.Content, &pub.Source,
			&pub.PublicationDate, &pub.Visibility, &pub.LikesCount,
			&pub.CommentsCount, &pub.SavedCount, &pub.IsLiked,
		); err != nil {
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
