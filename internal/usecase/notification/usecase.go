package notification

import (
	"context"
	"sense-backend/internal/domain"
)

// UseCase handles notification use cases
type UseCase struct {
	notificationRepo domain.NotificationRepository
}

// NewUseCase creates a new notification use case
func NewUseCase(notificationRepo domain.NotificationRepository) *UseCase {
	return &UseCase{
		notificationRepo: notificationRepo,
	}
}

// GetByUser retrieves notifications for a user
func (uc *UseCase) GetByUser(ctx context.Context, userID string, unreadOnly bool, limit, offset int) ([]*domain.Notification, int, error) {
	return uc.notificationRepo.GetByUser(ctx, userID, unreadOnly, limit, offset)
}

// MarkAsRead marks a notification as read
func (uc *UseCase) MarkAsRead(ctx context.Context, notificationID string) error {
	return uc.notificationRepo.MarkAsRead(ctx, notificationID)
}

// Create creates a new notification
func (uc *UseCase) Create(ctx context.Context, notification *domain.Notification) error {
	return uc.notificationRepo.Create(ctx, notification)
}

