package handlers

import (
	"net/http"

	"sense-backend/internal/delivery/http/middleware"
	profileUsecase "sense-backend/internal/usecase/profile"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

// ProfileHandler handles profile endpoints
type ProfileHandler struct {
	profileUC *profileUsecase.UseCase
	validator *validator.Validate
}

// NewProfileHandler creates a new profile handler
func NewProfileHandler(profileUC *profileUsecase.UseCase, validator *validator.Validate) *ProfileHandler {
	return &ProfileHandler{
		profileUC: profileUC,
		validator: validator,
	}
}

// RegisterRoutes registers profile routes
func (h *ProfileHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/me", h.GetMe).Methods("GET")
	r.HandleFunc("/me", h.UpdateMe).Methods("POST")
	r.HandleFunc("/{id}", h.Get).Methods("GET")
	r.HandleFunc("/{id}/stats", h.GetStats).Methods("GET")
}

// RegisterFollowRoutes registers follow routes (separate because they don't use /profile prefix)
func (h *ProfileHandler) RegisterFollowRoutes(r *mux.Router) {
	r.HandleFunc("/{id}", h.Follow).Methods("POST")
	r.HandleFunc("/{id}", h.Unfollow).Methods("DELETE")
}

// GetMe handles GET /profile/me
func (h *ProfileHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		WriteError(w, http.StatusUnauthorized, "unauthorized", "Требуется аутентификация", nil)
		return
	}

	profile, err := h.profileUC.GetProfile(r.Context(), userID)
	if err != nil {
		WriteError(w, http.StatusNotFound, "not_found", "Профиль не найден", nil)
		return
	}

	WriteJSON(w, http.StatusOK, profile)
}

// UpdateMe handles POST /profile/me
func (h *ProfileHandler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		WriteError(w, http.StatusUnauthorized, "unauthorized", "Требуется аутентификация", nil)
		return
	}

	var req profileUsecase.UpdateRequest
	if err := ParseJSON(r, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "validation_error", "Неверные данные в запросе", nil)
		return
	}

	if err := ValidateRequest(h.validator, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "validation_error", "Неверные данные в запросе", nil)
		return
	}

	profile, err := h.profileUC.UpdateProfile(r.Context(), userID, &req)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "validation_error", err.Error(), nil)
		return
	}

	WriteJSON(w, http.StatusOK, profile)
}

// Get handles GET /profile/{id}
func (h *ProfileHandler) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	profile, err := h.profileUC.GetProfile(r.Context(), id)
	if err != nil {
		WriteError(w, http.StatusNotFound, "not_found", "Профиль не найден", nil)
		return
	}

	// Check if current user is following this profile
	currentUserID := middleware.GetUserID(r.Context())
	isFollowing := false
	if currentUserID != "" && currentUserID != id {
		isFollowing, _ = h.profileUC.IsFollowing(r.Context(), currentUserID, id)
	}

	// Create response with is_following field
	response := map[string]interface{}{
		"id":            profile.ID,
		"username":      profile.Username,
		"email":         profile.Email,
		"phone":         profile.Phone,
		"icon_url":      profile.IconURL,
		"description":   profile.Description,
		"role":          profile.Role,
		"registered_at": profile.RegisteredAt,
		"is_following":  isFollowing,
	}

	if profile.Statistic != nil {
		response["statistic"] = profile.Statistic
	}

	WriteJSON(w, http.StatusOK, response)
}

// GetStats handles GET /profile/{id}/stats
func (h *ProfileHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	stats, err := h.profileUC.GetStats(r.Context(), id)
	if err != nil {
		WriteError(w, http.StatusNotFound, "not_found", "Профиль не найден", nil)
		return
	}

	WriteJSON(w, http.StatusOK, stats)
}

// Follow handles POST /follow/{id}
func (h *ProfileHandler) Follow(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		WriteError(w, http.StatusUnauthorized, "unauthorized", "Требуется аутентификация", nil)
		return
	}

	vars := mux.Vars(r)
	followingID := vars["id"]

	if err := h.profileUC.Follow(r.Context(), userID, followingID); err != nil {
		if err.Error() == "cannot follow yourself" {
			WriteError(w, http.StatusBadRequest, "validation_error", "Нельзя подписаться на себя", nil)
			return
		}
		WriteError(w, http.StatusNotFound, "not_found", "Пользователь не найден", nil)
		return
	}

	WriteJSON(w, http.StatusCreated, map[string]string{
		"message": "Вы подписались на пользователя",
	})
}

// Unfollow handles DELETE /follow/{id}
func (h *ProfileHandler) Unfollow(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		WriteError(w, http.StatusUnauthorized, "unauthorized", "Требуется аутентификация", nil)
		return
	}

	vars := mux.Vars(r)
	followingID := vars["id"]

	if err := h.profileUC.Unfollow(r.Context(), userID, followingID); err != nil {
		WriteError(w, http.StatusNotFound, "not_found", "Пользователь не найден", nil)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
