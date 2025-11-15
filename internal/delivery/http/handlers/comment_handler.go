package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	commentUsecase "sense-backend/internal/usecase/comment"
	"sense-backend/internal/delivery/http/middleware"
)

// CommentHandler handles comment endpoints
type CommentHandler struct {
	commentUC *commentUsecase.UseCase
	validator *validator.Validate
}

// NewCommentHandler creates a new comment handler
func NewCommentHandler(commentUC *commentUsecase.UseCase, validator *validator.Validate) *CommentHandler {
	return &CommentHandler{
		commentUC: commentUC,
		validator: validator,
	}
}

// RegisterRoutes registers comment routes
func (h *CommentHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/{id}", h.Get).Methods("GET")
	r.HandleFunc("/{id}", h.Update).Methods("PUT")
	r.HandleFunc("/{id}", h.Delete).Methods("DELETE")
	r.HandleFunc("/{id}/reply", h.Reply).Methods("POST")
	r.HandleFunc("/{id}/like", h.Like).Methods("POST")
}

// GetByPublication handles GET /publication/{id}/comments
func (h *CommentHandler) GetByPublication(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	publicationID := vars["id"]

	limit, offset := getPagination(r)

	comments, total, err := h.commentUC.GetByPublication(r.Context(), publicationID, limit, offset)
	if err != nil {
		WriteError(w, http.StatusNotFound, "not_found", "Публикация не найдена", nil)
		return
	}

	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"items": comments,
		"total": total,
		"limit": limit,
		"offset": offset,
	})
}

// Create handles POST /publication/{id}/comments
func (h *CommentHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		WriteError(w, http.StatusUnauthorized, "unauthorized", "Требуется аутентификация", nil)
		return
	}

	vars := mux.Vars(r)
	publicationID := vars["id"]

	var req commentUsecase.CreateRequest
	if err := ParseJSON(r, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "validation_error", "Неверные данные в запросе", nil)
		return
	}

	if err := ValidateRequest(h.validator, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "validation_error", "Неверные данные в запросе", nil)
		return
	}

	comment, err := h.commentUC.Create(r.Context(), publicationID, userID, &req)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "validation_error", err.Error(), nil)
		return
	}

	WriteJSON(w, http.StatusCreated, comment)
}

// Get handles GET /comment/{id}
func (h *CommentHandler) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	comment, err := h.commentUC.Get(r.Context(), id)
	if err != nil {
		WriteError(w, http.StatusNotFound, "not_found", "Комментарий не найден", nil)
		return
	}

	WriteJSON(w, http.StatusOK, comment)
}

// Update handles PUT /comment/{id}
func (h *CommentHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		WriteError(w, http.StatusUnauthorized, "unauthorized", "Требуется аутентификация", nil)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	var body struct {
		Text string `json:"text" validate:"required,max=2000"`
	}
	if err := ParseJSON(r, &body); err != nil {
		WriteError(w, http.StatusBadRequest, "validation_error", "Неверные данные в запросе", nil)
		return
	}

	if err := ValidateRequest(h.validator, &body); err != nil {
		WriteError(w, http.StatusBadRequest, "validation_error", "Неверные данные в запросе", nil)
		return
	}

	comment, err := h.commentUC.Update(r.Context(), id, userID, body.Text)
	if err != nil {
		if err.Error() == "forbidden: not the author" {
			WriteError(w, http.StatusForbidden, "forbidden", "Недостаточно прав", nil)
			return
		}
		WriteError(w, http.StatusNotFound, "not_found", "Комментарий не найден", nil)
		return
	}

	WriteJSON(w, http.StatusOK, comment)
}

// Delete handles DELETE /comment/{id}
func (h *CommentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		WriteError(w, http.StatusUnauthorized, "unauthorized", "Требуется аутентификация", nil)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.commentUC.Delete(r.Context(), id, userID); err != nil {
		if err.Error() == "forbidden: not the author" {
			WriteError(w, http.StatusForbidden, "forbidden", "Недостаточно прав", nil)
			return
		}
		WriteError(w, http.StatusNotFound, "not_found", "Комментарий не найден", nil)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Reply handles POST /comment/{id}/reply
func (h *CommentHandler) Reply(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		WriteError(w, http.StatusUnauthorized, "unauthorized", "Требуется аутентификация", nil)
		return
	}

	vars := mux.Vars(r)
	parentID := vars["id"]

	var body struct {
		Text string `json:"text" validate:"required,max=2000"`
	}
	if err := ParseJSON(r, &body); err != nil {
		WriteError(w, http.StatusBadRequest, "validation_error", "Неверные данные в запросе", nil)
		return
	}

	if err := ValidateRequest(h.validator, &body); err != nil {
		WriteError(w, http.StatusBadRequest, "validation_error", "Неверные данные в запросе", nil)
		return
	}

	// Get parent comment to get publication_id
	parentComment, err := h.commentUC.Get(r.Context(), parentID)
	if err != nil {
		WriteError(w, http.StatusNotFound, "not_found", "Комментарий не найден", nil)
		return
	}

	req := commentUsecase.CreateRequest{
		Text:     body.Text,
		ParentID: &parentID,
	}

	comment, err := h.commentUC.Create(r.Context(), parentComment.PublicationID, userID, &req)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "validation_error", err.Error(), nil)
		return
	}

	WriteJSON(w, http.StatusCreated, comment)
}

// Like handles POST /comment/{id}/like
func (h *CommentHandler) Like(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		WriteError(w, http.StatusUnauthorized, "unauthorized", "Требуется аутентификация", nil)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	liked, count, err := h.commentUC.Like(r.Context(), id, userID)
	if err != nil {
		WriteError(w, http.StatusNotFound, "not_found", "Комментарий не найден", nil)
		return
	}

	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"liked":      liked,
		"likes_count": count,
	})
}

func getPaginationFromRequest(r *http.Request) (limit, offset int) {
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

