package handlers

import (
	"net/http"

	"sense-backend/internal/delivery/http/middleware"
	aiUsecase "sense-backend/internal/usecase/ai"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

// AIHandler handles AI endpoints
type AIHandler struct {
	aiUC      *aiUsecase.UseCase
	validator *validator.Validate
}

// NewAIHandler creates a new AI handler
func NewAIHandler(aiUC *aiUsecase.UseCase, validator *validator.Validate) *AIHandler {
	return &AIHandler{
		aiUC:      aiUC,
		validator: validator,
	}
}

// RegisterRoutes registers AI routes
func (h *AIHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/recommendations", h.GetRecommendations).Methods("POST")
	r.HandleFunc("/recommendations/feed", h.GetRecommendationsFeed).Methods("GET")
	r.HandleFunc("/recommendations/{id}/hide", h.HideRecommendation).Methods("POST")
	r.HandleFunc("/purify", h.PurifyText).Methods("POST")
}

// GetRecommendations handles POST /recommendations
func (h *AIHandler) GetRecommendations(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		WriteError(w, http.StatusUnauthorized, "unauthorized", "Требуется аутентификация", nil)
		return
	}

	var req struct {
		Limit     *int    `json:"limit"`
		Algorithm *string `json:"algorithm"`
	}
	if err := ParseJSON(r, &req); err != nil {
		// Request body is optional, so we continue with defaults
	}

	limit := 20
	if req.Limit != nil {
		if *req.Limit >= 1 && *req.Limit <= 50 {
			limit = *req.Limit
		}
	}

	publications, err := h.aiUC.GetRecommendations(r.Context(), userID, limit, req.Algorithm)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "validation_error", err.Error(), nil)
		return
	}

	// Build response with recommendation metadata
	items := make([]map[string]interface{}, 0, len(publications))
	algorithm := "collaborative_filtering"
	if req.Algorithm != nil {
		algorithm = *req.Algorithm
	}

	for _, pub := range publications {
		item := map[string]interface{}{
			"id":               pub.ID,
			"author_id":        pub.AuthorID,
			"type":             pub.Type,
			"content":          pub.Content,
			"source":           pub.Source,
			"publication_date": pub.PublicationDate.Format("2006-01-02T15:04:05Z07:00"),
			"visibility":       pub.Visibility,
			"likes_count":      pub.LikesCount,
			"comments_count":   pub.CommentsCount,
			"saved_count":      pub.SavedCount,
			"recommendation": map[string]interface{}{
				"algorithm": algorithm,
				"reason":    "Recommended by " + algorithm,
			},
		}

		items = append(items, item)
	}

	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"items": items,
		"total": len(items),
	})
}

// GetRecommendationsFeed handles GET /recommendations/feed
func (h *AIHandler) GetRecommendationsFeed(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		WriteError(w, http.StatusUnauthorized, "unauthorized", "Требуется аутентификация", nil)
		return
	}

	limit, offset := getPagination(r)

	algorithm := r.URL.Query().Get("algorithm")
	var algorithmPtr *string
	if algorithm != "" {
		algorithmPtr = &algorithm
	}

	publications, total, err := h.aiUC.GetRecommendationsFeed(r.Context(), userID, limit, offset, algorithmPtr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "validation_error", err.Error(), nil)
		return
	}

	// Build response with recommendation metadata
	items := make([]map[string]interface{}, 0, len(publications))
	alg := "collaborative_filtering"
	if algorithmPtr != nil {
		alg = *algorithmPtr
	}

	for _, pub := range publications {
		item := map[string]interface{}{
			"id":               pub.ID,
			"author_id":        pub.AuthorID,
			"type":             pub.Type,
			"content":          pub.Content,
			"source":           pub.Source,
			"publication_date": pub.PublicationDate.Format("2006-01-02T15:04:05Z07:00"),
			"visibility":       pub.Visibility,
			"likes_count":      pub.LikesCount,
			"comments_count":   pub.CommentsCount,
			"saved_count":      pub.SavedCount,
			"recommendation": map[string]interface{}{
				"algorithm": alg,
				"reason":    "Recommended by " + alg,
			},
		}

		items = append(items, item)
	}

	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"items":  items,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

// HideRecommendation handles POST /recommendations/{id}/hide
func (h *AIHandler) HideRecommendation(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		WriteError(w, http.StatusUnauthorized, "unauthorized", "Требуется аутентификация", nil)
		return
	}

	vars := mux.Vars(r)
	recommendationID := vars["id"]

	if err := h.aiUC.HideRecommendation(r.Context(), recommendationID); err != nil {
		WriteError(w, http.StatusNotFound, "not_found", "Рекомендация не найдена", nil)
		return
	}

	WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Рекомендация скрыта",
	})
}

// PurifyText handles POST /purify
func (h *AIHandler) PurifyText(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		WriteError(w, http.StatusUnauthorized, "unauthorized", "Требуется аутентификация", nil)
		return
	}

	var req aiUsecase.PurifyRequest
	if err := ParseJSON(r, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "validation_error", "Неверные данные в запросе", nil)
		return
	}

	if err := ValidateRequest(h.validator, &req); err != nil {
		errMsg := err.Error()
		WriteError(w, http.StatusBadRequest, "validation_error", "Неверные данные в запросе", &errMsg)
		return
	}

	resp, err := h.aiUC.PurifyText(r.Context(), &req)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "validation_error", err.Error(), nil)
		return
	}

	WriteJSON(w, http.StatusOK, resp)
}
