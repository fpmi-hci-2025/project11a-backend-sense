package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	publicationUsecase "sense-backend/internal/usecase/publication"
	"sense-backend/internal/delivery/http/middleware"
)

// PublicationHandler handles publication endpoints
type PublicationHandler struct {
	publicationUC *publicationUsecase.UseCase
	validator     *validator.Validate
}

// NewPublicationHandler creates a new publication handler
func NewPublicationHandler(publicationUC *publicationUsecase.UseCase, validator *validator.Validate) *PublicationHandler {
	return &PublicationHandler{
		publicationUC: publicationUC,
		validator:     validator,
	}
}

// RegisterRoutes registers publication routes
func (h *PublicationHandler) RegisterRoutes(r *mux.Router, commentHandler *CommentHandler) {
	r.HandleFunc("/create", h.Create).Methods("POST")
	r.HandleFunc("/{id}", h.Get).Methods("GET")
	r.HandleFunc("/{id}", h.Update).Methods("PUT")
	r.HandleFunc("/{id}", h.Delete).Methods("DELETE")
	r.HandleFunc("/{id}/like", h.Like).Methods("POST")
	r.HandleFunc("/{id}/likes", h.GetLikes).Methods("GET")
	r.HandleFunc("/{id}/save", h.Save).Methods("POST")
	r.HandleFunc("/{id}/save", h.Unsave).Methods("DELETE")
	r.HandleFunc("/{id}/comments", commentHandler.GetByPublication).Methods("GET")
	r.HandleFunc("/{id}/comments", commentHandler.Create).Methods("POST")
}

// Create handles POST /publication/create
func (h *PublicationHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		WriteError(w, http.StatusUnauthorized, "unauthorized", "Требуется аутентификация", nil)
		return
	}

	var req publicationUsecase.CreateRequest
	if err := ParseJSON(r, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "validation_error", "Неверные данные в запросе", nil)
		return
	}

	if err := ValidateRequest(h.validator, &req); err != nil {
		errMsg := err.Error()
		WriteError(w, http.StatusBadRequest, "validation_error", "Неверные данные в запросе", &errMsg)
		return
	}

	publication, err := h.publicationUC.Create(r.Context(), userID, &req)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "validation_error", err.Error(), nil)
		return
	}

	WriteJSON(w, http.StatusCreated, publication)
}

// Get handles GET /publication/{id}
func (h *PublicationHandler) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	publication, err := h.publicationUC.Get(r.Context(), id)
	if err != nil {
		WriteError(w, http.StatusNotFound, "not_found", "Публикация не найдена", nil)
		return
	}

	WriteJSON(w, http.StatusOK, publication)
}

// Update handles PUT /publication/{id}
func (h *PublicationHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		WriteError(w, http.StatusUnauthorized, "unauthorized", "Требуется аутентификация", nil)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	var req publicationUsecase.UpdateRequest
	if err := ParseJSON(r, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "validation_error", "Неверные данные в запросе", nil)
		return
	}

	publication, err := h.publicationUC.Update(r.Context(), id, userID, &req)
	if err != nil {
		if err.Error() == "forbidden: not the author" {
			WriteError(w, http.StatusForbidden, "forbidden", "Недостаточно прав", nil)
			return
		}
		WriteError(w, http.StatusNotFound, "not_found", "Публикация не найдена", nil)
		return
	}

	WriteJSON(w, http.StatusOK, publication)
}

// Delete handles DELETE /publication/{id}
func (h *PublicationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		WriteError(w, http.StatusUnauthorized, "unauthorized", "Требуется аутентификация", nil)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.publicationUC.Delete(r.Context(), id, userID); err != nil {
		if err.Error() == "forbidden: not the author" {
			WriteError(w, http.StatusForbidden, "forbidden", "Недостаточно прав", nil)
			return
		}
		WriteError(w, http.StatusNotFound, "not_found", "Публикация не найдена", nil)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Like handles POST /publication/{id}/like
func (h *PublicationHandler) Like(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		WriteError(w, http.StatusUnauthorized, "unauthorized", "Требуется аутентификация", nil)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	liked, count, err := h.publicationUC.Like(r.Context(), id, userID)
	if err != nil {
		WriteError(w, http.StatusNotFound, "not_found", "Публикация не найдена", nil)
		return
	}

	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"liked":      liked,
		"likes_count": count,
	})
}

// GetLikes handles GET /publication/{id}/likes
func (h *PublicationHandler) GetLikes(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	limit := 20
	offset := 0
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}
	if o := r.URL.Query().Get("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil {
			offset = parsed
		}
	}

	users, total, err := h.publicationUC.GetLikedUsers(r.Context(), id, limit, offset)
	if err != nil {
		WriteError(w, http.StatusNotFound, "not_found", "Публикация не найдена", nil)
		return
	}

	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"items": users,
		"total": total,
		"limit": limit,
		"offset": offset,
	})
}

// Save handles POST /publication/{id}/save
func (h *PublicationHandler) Save(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		WriteError(w, http.StatusUnauthorized, "unauthorized", "Требуется аутентификация", nil)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	var body struct {
		Note *string `json:"note"`
	}
	ParseJSON(r, &body) // Ignore errors, note is optional

	if err := h.publicationUC.Save(r.Context(), id, userID, body.Note); err != nil {
		WriteError(w, http.StatusNotFound, "not_found", "Публикация не найдена", nil)
		return
	}

	WriteJSON(w, http.StatusCreated, map[string]string{"message": "Публикация добавлена в сохраненные"})
}

// Unsave handles DELETE /publication/{id}/save
func (h *PublicationHandler) Unsave(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		WriteError(w, http.StatusUnauthorized, "unauthorized", "Требуется аутентификация", nil)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.publicationUC.Unsave(r.Context(), id, userID); err != nil {
		WriteError(w, http.StatusNotFound, "not_found", "Публикация не найдена", nil)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

