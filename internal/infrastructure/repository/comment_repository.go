package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"sense-backend/internal/domain"
)

type commentRepository struct {
	pool *pgxpool.Pool
}

// NewCommentRepository creates a new comment repository
func NewCommentRepository(pool *pgxpool.Pool) domain.CommentRepository {
	return &commentRepository{pool: pool}
}

func (r *commentRepository) Create(ctx context.Context, comment *domain.Comment) error {
	query := `
		INSERT INTO comments (id, publication_id, parent_id, author_id, text, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.pool.Exec(ctx, query,
		comment.ID, comment.PublicationID, comment.ParentID, comment.AuthorID,
		comment.Text, comment.CreatedAt,
	)
	if err != nil {
		return err
	}

	// Update comments count
	_, err = r.pool.Exec(ctx, `
		UPDATE publications SET comments_count = comments_count + 1 WHERE id = $1
	`, comment.PublicationID)
	return err
}

func (r *commentRepository) GetByID(ctx context.Context, id string) (*domain.Comment, error) {
	query := `
		SELECT id, publication_id, parent_id, author_id, text, created_at, likes_count
		FROM comments
		WHERE id = $1
	`
	
	var comment domain.Comment
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&comment.ID, &comment.PublicationID, &comment.ParentID, &comment.AuthorID,
		&comment.Text, &comment.CreatedAt, &comment.LikesCount,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("comment not found")
	}
	return &comment, err
}

func (r *commentRepository) GetByPublication(ctx context.Context, publicationID string, limit, offset int) ([]*domain.Comment, int, error) {
	// Get total
	var total int
	err := r.pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM comments WHERE publication_id = $1
	`, publicationID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get comments
	rows, err := r.pool.Query(ctx, `
		SELECT id, publication_id, parent_id, author_id, text, created_at, likes_count
		FROM comments
		WHERE publication_id = $1
		ORDER BY created_at ASC
		LIMIT $2 OFFSET $3
	`, publicationID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var comments []*domain.Comment
	for rows.Next() {
		var comment domain.Comment
		err := rows.Scan(
			&comment.ID, &comment.PublicationID, &comment.ParentID, &comment.AuthorID,
			&comment.Text, &comment.CreatedAt, &comment.LikesCount,
		)
		if err != nil {
			return nil, 0, err
		}
		comments = append(comments, &comment)
	}

	return comments, total, rows.Err()
}

func (r *commentRepository) Update(ctx context.Context, comment *domain.Comment) error {
	query := `
		UPDATE comments
		SET text = $2
		WHERE id = $1
	`
	_, err := r.pool.Exec(ctx, query, comment.ID, comment.Text)
	return err
}

func (r *commentRepository) Delete(ctx context.Context, id string) error {
	// Get publication_id before deleting
	var publicationID string
	err := r.pool.QueryRow(ctx, `
		SELECT publication_id FROM comments WHERE id = $1
	`, id).Scan(&publicationID)
	if err != nil {
		return err
	}

	// Delete comment
	_, err = r.pool.Exec(ctx, `DELETE FROM comments WHERE id = $1`, id)
	if err != nil {
		return err
	}

	// Update comments count
	_, err = r.pool.Exec(ctx, `
		UPDATE publications SET comments_count = comments_count - 1 WHERE id = $1
	`, publicationID)
	return err
}

func (r *commentRepository) Like(ctx context.Context, userID, commentID string) (bool, error) {
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
		SELECT EXISTS(SELECT 1 FROM comment_likes WHERE user_id = $1 AND comment_id = $2)
	`, userID, commentID).Scan(&exists)
	if err != nil {
		return false, err
	}

	if exists {
		// Unlike
		_, err = tx.Exec(ctx, `
			DELETE FROM comment_likes WHERE user_id = $1 AND comment_id = $2
		`, userID, commentID)
		if err != nil {
			return false, err
		}
		_, err = tx.Exec(ctx, `
			UPDATE comments SET likes_count = likes_count - 1 WHERE id = $1
		`, commentID)
		if err != nil {
			return false, err
		}
		return false, tx.Commit(ctx)
	} else {
		// Like
		_, err = tx.Exec(ctx, `
			INSERT INTO comment_likes (user_id, comment_id) VALUES ($1, $2)
		`, userID, commentID)
		if err != nil {
			return false, err
		}
		_, err = tx.Exec(ctx, `
			UPDATE comments SET likes_count = likes_count + 1 WHERE id = $1
		`, commentID)
		if err != nil {
			return false, err
		}
		return true, tx.Commit(ctx)
	}
}

func (r *commentRepository) IsLiked(ctx context.Context, userID, commentID string) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx, `
		SELECT EXISTS(SELECT 1 FROM comment_likes WHERE user_id = $1 AND comment_id = $2)
	`, userID, commentID).Scan(&exists)
	return exists, err
}

func (r *commentRepository) GetLikesCount(ctx context.Context, commentID string) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx, `
		SELECT likes_count FROM comments WHERE id = $1
	`, commentID).Scan(&count)
	return count, err
}

