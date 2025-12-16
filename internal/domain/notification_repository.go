package domain

import "context"

// NotificationRepository defines interface for notification data operations
type NotificationRepository interface {
	// Create creates a new notification
	Create(ctx context.Context, notification *Notification) error
	
	// GetByUser retrieves notifications for user
	GetByUser(ctx context.Context, userID string, unreadOnly bool, limit, offset int) ([]*Notification, int, error)
	
	// MarkAsRead marks notification as read
	MarkAsRead(ctx context.Context, notificationID string) error
	
	// MarkAllAsRead marks all user notifications as read
	MarkAllAsRead(ctx context.Context, userID string) error
}

