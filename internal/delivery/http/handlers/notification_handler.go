package handlers

import (
	"net/http"
	"strconv"

	"sense-backend/internal/delivery/http/middleware"
	notificationUsecase "sense-backend/internal/usecase/notification"

	"github.com/go-playground/validator/v10"
)

// NotificationHandler handles notification endpoints
type NotificationHandler struct {
	notificationUC *notificationUsecase.UseCase
	validator      *validator.Validate
}

// NewNotificationHandler creates a new notification handler
func NewNotificationHandler(notificationUC *notificationUsecase.UseCase, validator *validator.Validate) *NotificationHandler {
	return &NotificationHandler{
		notificationUC: notificationUC,
		validator:      validator,
	}
}

// GetNotifications handles GET /notifications
func (h *NotificationHandler) GetNotifications(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		WriteError(w, http.StatusUnauthorized, "unauthorized", "Требуется аутентификация", nil)
		return
	}

	limit, offset := getPagination(r)

	unreadOnly := false
	if unreadStr := r.URL.Query().Get("unread_only"); unreadStr != "" {
		if parsed, err := strconv.ParseBool(unreadStr); err == nil {
			unreadOnly = parsed
		}
	}

	notifications, total, err := h.notificationUC.GetByUser(r.Context(), userID, unreadOnly, limit, offset)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "validation_error", err.Error(), nil)
		return
	}

	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"items":  notifications,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}
