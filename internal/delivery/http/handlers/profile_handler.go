package handlers

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	profileUsecase "sense-backend/internal/usecase/profile"
	"sense-backend/internal/delivery/http/middleware"
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

	WriteJSON(w, http.StatusOK, profile)
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

