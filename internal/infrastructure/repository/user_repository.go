package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"sense-backend/internal/domain"
)

type userRepository struct {
	pool *pgxpool.Pool
}

// NewUserRepository creates a new user repository
func NewUserRepository(pool *pgxpool.Pool) domain.UserRepository {
	return &userRepository{pool: pool}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (id, username, email, phone, icon_url, description, role, password_hash, registered_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	
	_, err := r.pool.Exec(ctx, query,
		user.ID, user.Username, user.Email, user.Phone, user.IconURL, user.Description,
		user.Role, user.PasswordHash, user.RegisteredAt,
	)
	return err
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	query := `
		SELECT id, username, email, phone, icon_url, description, role, registered_at, password_hash
		FROM users
		WHERE id = $1
	`
	
	var user domain.User
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.Phone, &user.IconURL,
		&user.Description, &user.Role, &user.RegisteredAt, &user.PasswordHash,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	return &user, err
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	query := `
		SELECT id, username, email, phone, icon_url, description, role, registered_at, password_hash
		FROM users
		WHERE username = $1
	`
	
	var user domain.User
	err := r.pool.QueryRow(ctx, query, username).Scan(
		&user.ID, &user.Username, &user.Email, &user.Phone, &user.IconURL,
		&user.Description, &user.Role, &user.RegisteredAt, &user.PasswordHash,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	return &user, err
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, username, email, phone, icon_url, description, role, registered_at, password_hash
		FROM users
		WHERE email = $1
	`
	
	var user domain.User
	err := r.pool.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.Phone, &user.IconURL,
		&user.Description, &user.Role, &user.RegisteredAt, &user.PasswordHash,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	return &user, err
}

func (r *userRepository) GetByLogin(ctx context.Context, login string) (*domain.User, error) {
	query := `
		SELECT id, username, email, phone, icon_url, description, role, registered_at, password_hash
		FROM users
		WHERE username = $1 OR email = $1
	`
	
	var user domain.User
	err := r.pool.QueryRow(ctx, query, login).Scan(
		&user.ID, &user.Username, &user.Email, &user.Phone, &user.IconURL,
		&user.Description, &user.Role, &user.RegisteredAt, &user.PasswordHash,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	return &user, err
}

func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users
		SET username = $2, email = $3, phone = $4, icon_url = $5, description = $6, role = $7
		WHERE id = $1
	`
	
	_, err := r.pool.Exec(ctx, query,
		user.ID, user.Username, user.Email, user.Phone, user.IconURL,
		user.Description, user.Role,
	)
	return err
}

func (r *userRepository) GetStats(ctx context.Context, userID string) (*domain.UserStatistic, error) {
	stats := &domain.UserStatistic{}
	
	// Get publications count
	err := r.pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM publications WHERE author_id = $1
	`, userID).Scan(&stats.PublicationsCount)
	if err != nil {
		return nil, err
	}
	
	// Get followers and following from users table (updated by triggers)
	err = r.pool.QueryRow(ctx, `
		SELECT followers_count, following_count FROM users WHERE id = $1
	`, userID).Scan(&stats.FollowersCount, &stats.FollowingCount)
	if err != nil {
		return nil, err
	}
	
	// Get likes received
	err = r.pool.QueryRow(ctx, `
		SELECT COALESCE(SUM(likes_count), 0) FROM publications WHERE author_id = $1
	`, userID).Scan(&stats.LikesReceived)
	if err != nil {
		return nil, err
	}
	
	// Get comments received
	err = r.pool.QueryRow(ctx, `
		SELECT COALESCE(SUM(comments_count), 0) FROM publications WHERE author_id = $1
	`, userID).Scan(&stats.CommentsReceived)
	if err != nil {
		return nil, err
	}
	
	// Get saved count
	err = r.pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM saved_items WHERE publication_id IN (
			SELECT id FROM publications WHERE author_id = $1
		)
	`, userID).Scan(&stats.SavedCount)
	if err != nil {
		return nil, err
	}
	
	return stats, nil
}

func (r *userRepository) GetFollowersCount(ctx context.Context, userID string) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx, `
		SELECT followers_count FROM users WHERE id = $1
	`, userID).Scan(&count)
	return count, err
}

func (r *userRepository) GetFollowingCount(ctx context.Context, userID string) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx, `
		SELECT following_count FROM users WHERE id = $1
	`, userID).Scan(&count)
	return count, err
}

func (r *userRepository) IsFollowing(ctx context.Context, followerID, followingID string) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx, `
		SELECT EXISTS(SELECT 1 FROM user_follows WHERE follower_id = $1 AND following_id = $2)
	`, followerID, followingID).Scan(&exists)
	return exists, err
}

func (r *userRepository) Follow(ctx context.Context, followerID, followingID string) error {
	query := `
		INSERT INTO user_follows (follower_id, following_id)
		VALUES ($1, $2)
		ON CONFLICT (follower_id, following_id) DO NOTHING
	`
	_, err := r.pool.Exec(ctx, query, followerID, followingID)
	return err
}

func (r *userRepository) Unfollow(ctx context.Context, followerID, followingID string) error {
	query := `
		DELETE FROM user_follows
		WHERE follower_id = $1 AND following_id = $2
	`
	_, err := r.pool.Exec(ctx, query, followerID, followingID)
	return err
}

func (r *userRepository) Search(ctx context.Context, query string, role *domain.UserRole, limit, offset int) ([]*domain.User, int, error) {
	searchQuery := `%` + query + `%`
	baseQuery := `
		SELECT id, username, email, phone, icon_url, description, role, registered_at
		FROM users
		WHERE (username ILIKE $1 OR description ILIKE $1)
	`
	args := []interface{}{searchQuery}
	argIndex := 2
	
	if role != nil {
		baseQuery += fmt.Sprintf(" AND role = $%d", argIndex)
		args = append(args, *role)
		argIndex++
	}
	
	// Get total count
	var total int
	countQuery := "SELECT COUNT(*) FROM users WHERE (username ILIKE $1 OR description ILIKE $1)"
	if role != nil {
		countQuery += " AND role = $2"
	}
	err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}
	
	// Get users
	baseQuery += fmt.Sprintf(" ORDER BY username LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, limit, offset)
	
	rows, err := r.pool.Query(ctx, baseQuery, args...)
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

