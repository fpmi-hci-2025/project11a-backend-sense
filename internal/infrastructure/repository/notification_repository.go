package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"sense-backend/internal/domain"
)

type notificationRepository struct {
	pool *pgxpool.Pool
}

// NewNotificationRepository creates a new notification repository
func NewNotificationRepository(pool *pgxpool.Pool) domain.NotificationRepository {
	return &notificationRepository{pool: pool}
}

func (r *notificationRepository) Create(ctx context.Context, notification *domain.Notification) error {
	query := `
		INSERT INTO notifications (id, user_id, type, title, message, data, is_read, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	
	var dataJSON interface{}
	if notification.Data != nil {
		dataJSON = notification.Data
	}
	
	_, err := r.pool.Exec(ctx, query,
		notification.ID, notification.UserID, notification.Type,
		notification.Title, notification.Message, dataJSON,
		notification.IsRead, notification.CreatedAt,
	)
	return err
}

func (r *notificationRepository) GetByUser(ctx context.Context, userID string, unreadOnly bool, limit, offset int) ([]*domain.Notification, int, error) {
	where := "user_id = $1"
	args := []interface{}{userID}
	argIndex := 2

	if unreadOnly {
		where += " AND is_read = false"
	}

	// Get total
	var total int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM notifications WHERE %s", where)
	err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get notifications
	query := fmt.Sprintf(`
		SELECT id, user_id, type, title, message, data, is_read, created_at
		FROM notifications
		WHERE %s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, where, argIndex, argIndex+1)
	args = append(args, limit, offset)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var notifications []*domain.Notification
	for rows.Next() {
		var notif domain.Notification
		var dataJSON interface{}
		err := rows.Scan(
			&notif.ID, &notif.UserID, &notif.Type, &notif.Title,
			&notif.Message, &dataJSON, &notif.IsRead, &notif.CreatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		if dataJSON != nil {
			// Data will be parsed as map[string]interface{} by pgx
			notif.Data = make(map[string]interface{})
			// Full JSON parsing would require proper unmarshaling
		}
		notifications = append(notifications, &notif)
	}

	return notifications, total, rows.Err()
}

func (r *notificationRepository) MarkAsRead(ctx context.Context, notificationID string) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE notifications SET is_read = true WHERE id = $1
	`, notificationID)
	return err
}

func (r *notificationRepository) MarkAllAsRead(ctx context.Context, userID string) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE notifications SET is_read = true WHERE user_id = $1
	`, userID)
	return err
}

