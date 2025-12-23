package handlers

import (
	"net/http"
	"strconv"

	"sense-backend/internal/delivery/http/middleware"
	"sense-backend/internal/domain"
	searchUsecase "sense-backend/internal/usecase/search"

	"github.com/go-playground/validator/v10"
)

// SearchHandler handles search endpoints
type SearchHandler struct {
	searchUC  *searchUsecase.UseCase
	validator *validator.Validate
}

// NewSearchHandler creates a new search handler
func NewSearchHandler(searchUC *searchUsecase.UseCase, validator *validator.Validate) *SearchHandler {
	return &SearchHandler{
		searchUC:  searchUC,
		validator: validator,
	}
}

// SearchPublications handles GET /search
func (h *SearchHandler) SearchPublications(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		WriteError(w, http.StatusBadRequest, "validation_error", "Параметр 'q' обязателен", nil)
		return
	}

	// Get viewer user ID for like status (may be empty if not authenticated)
	viewerUserID := middleware.GetUserID(r.Context())
	var viewerUserIDPtr *string
	if viewerUserID != "" {
		viewerUserIDPtr = &viewerUserID
	}

	limit, offset := getPagination(r)
	filters := h.parseSearchFilters(r)

	publications, total, err := h.searchUC.SearchPublications(r.Context(), query, viewerUserIDPtr, filters, limit, offset)
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

// SearchUsers handles GET /search/users
func (h *SearchHandler) SearchUsers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		WriteError(w, http.StatusBadRequest, "validation_error", "Параметр 'q' обязателен", nil)
		return
	}

	limit, offset := getPagination(r)

	var rolePtr *domain.UserRole
	if roleStr := r.URL.Query().Get("role"); roleStr != "" {
		role := domain.UserRole(roleStr)
		rolePtr = &role
	}

	users, total, err := h.searchUC.SearchUsers(r.Context(), query, rolePtr, limit, offset)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "validation_error", err.Error(), nil)
		return
	}

	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"items":  users,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

// WarmupIndex handles POST /search/warmup
func (h *SearchHandler) WarmupIndex(w http.ResponseWriter, r *http.Request) {
	// For now, this is a placeholder that returns success
	// In a real implementation, this would trigger search index rebuilding
	// which could be a background job using the filters provided

	taskID := "123e4567-e89b-12d3-a456-426614174000" // Placeholder task ID

	WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Индексация запущена",
		"task_id": taskID,
	})
}

// GetTags handles GET /tags
func (h *SearchHandler) GetTags(w http.ResponseWriter, r *http.Request) {
	limit := 20
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	var searchPtr *string
	if search := r.URL.Query().Get("search"); search != "" {
		searchPtr = &search
	}

	tags, total, err := h.searchUC.GetTags(r.Context(), limit, searchPtr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "validation_error", err.Error(), nil)
		return
	}

	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"items": tags,
		"total": total,
	})
}

func (h *SearchHandler) parseSearchFilters(r *http.Request) *domain.SearchFilters {
	filters := &domain.SearchFilters{}

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

	return filters
}
