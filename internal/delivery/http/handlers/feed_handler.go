package handlers

import (
	"net/http"
	"strconv"
	"time"

	"sense-backend/internal/delivery/http/middleware"
	"sense-backend/internal/domain"
	feedUsecase "sense-backend/internal/usecase/feed"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

// FeedHandler handles feed endpoints
type FeedHandler struct {
	feedUC    *feedUsecase.UseCase
	validator *validator.Validate
}

// NewFeedHandler creates a new feed handler
func NewFeedHandler(feedUC *feedUsecase.UseCase, validator *validator.Validate) *FeedHandler {
	return &FeedHandler{
		feedUC:    feedUC,
		validator: validator,
	}
}

// RegisterRoutes registers feed routes
func (h *FeedHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("", h.GetFeed).Methods("GET")
	r.HandleFunc("/me", h.GetMe).Methods("GET")
	r.HandleFunc("/me/saved", h.GetSaved).Methods("GET")
	r.HandleFunc("/user/{id}", h.GetUser).Methods("GET")
}

// GetFeed handles GET /feed
func (h *FeedHandler) GetFeed(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context()) // May be empty
	var userIDPtr *string
	if userID != "" {
		userIDPtr = &userID
	}

	limit, offset := getPagination(r)
	filters := h.parseFeedFilters(r)

	publications, total, err := h.feedUC.GetFeed(r.Context(), userIDPtr, filters, limit, offset)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "validation_error", err.Error(), nil)
		return
	}

	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"items":  publications,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

// GetMe handles GET /feed/me
func (h *FeedHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		WriteError(w, http.StatusUnauthorized, "unauthorized", "Требуется аутентификация", nil)
		return
	}

	limit, offset := getPagination(r)
	filters := h.parsePublicationFilters(r)

	// Pass userID as viewerUserID to get like status
	publications, total, err := h.feedUC.GetUserFeed(r.Context(), userID, &userID, filters, limit, offset)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "validation_error", err.Error(), nil)
		return
	}

	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"items":  publications,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

// GetSaved handles GET /feed/me/saved
func (h *FeedHandler) GetSaved(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		WriteError(w, http.StatusUnauthorized, "unauthorized", "Требуется аутентификация", nil)
		return
	}

	limit, offset := getPagination(r)
	filters := h.parsePublicationFilters(r)

	publications, total, err := h.feedUC.GetSavedFeed(r.Context(), userID, filters, limit, offset)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "validation_error", err.Error(), nil)
		return
	}

	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"items":  publications,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

// GetUser handles GET /feed/user/{id}
func (h *FeedHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	authorID := vars["id"]

	// Get viewer user ID for like status (may be empty if not authenticated)
	viewerUserID := middleware.GetUserID(r.Context())
	var viewerUserIDPtr *string
	if viewerUserID != "" {
		viewerUserIDPtr = &viewerUserID
	}

	limit, offset := getPagination(r)
	filters := h.parsePublicationFilters(r)

	publications, total, err := h.feedUC.GetUserFeed(r.Context(), authorID, viewerUserIDPtr, filters, limit, offset)
	if err != nil {
		WriteError(w, http.StatusNotFound, "not_found", "Пользователь не найден", nil)
		return
	}

	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"items":  publications,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

func (h *FeedHandler) parseFeedFilters(r *http.Request) *domain.FeedFilters {
	filters := &domain.FeedFilters{}

	if typeStr := r.URL.Query().Get("type"); typeStr != "" {
		t := domain.PublicationType(typeStr)
		filters.Type = &t
	}
	if visStr := r.URL.Query().Get("visibility"); visStr != "" {
		v := domain.VisibilityType(visStr)
		filters.Visibility = &v
	}
	if authorID := r.URL.Query().Get("author_id"); authorID != "" {
		filters.AuthorID = &authorID
	}
	if dateFrom := r.URL.Query().Get("date_from"); dateFrom != "" {
		if t, err := time.Parse(time.RFC3339, dateFrom); err == nil {
			filters.DateFrom = &t
		}
	}
	if dateTo := r.URL.Query().Get("date_to"); dateTo != "" {
		if t, err := time.Parse(time.RFC3339, dateTo); err == nil {
			filters.DateTo = &t
		}
	}

	return filters
}

func (h *FeedHandler) parsePublicationFilters(r *http.Request) *domain.PublicationFilters {
	filters := &domain.PublicationFilters{}

	if typeStr := r.URL.Query().Get("type"); typeStr != "" {
		t := domain.PublicationType(typeStr)
		filters.Type = &t
	}
	if visStr := r.URL.Query().Get("visibility"); visStr != "" {
		v := domain.VisibilityType(visStr)
		filters.Visibility = &v
	}

	return filters
}

func getPagination(r *http.Request) (limit, offset int) {
	limit = 20
	offset = 0
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}
	if o := r.URL.Query().Get("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}
	return limit, offset
}
