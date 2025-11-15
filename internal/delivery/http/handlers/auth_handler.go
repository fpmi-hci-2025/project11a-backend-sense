package handlers

import (
	"net/http"

	authUsecase "sense-backend/internal/usecase/auth"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	authUC    *authUsecase.UseCase
	validator *validator.Validate
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authUC *authUsecase.UseCase, validator *validator.Validate) *AuthHandler {
	return &AuthHandler{
		authUC:    authUC,
		validator: validator,
	}
}

// RegisterRoutes registers auth routes
func (h *AuthHandler) RegisterRoutes(r *mux.Router, tokenSvc interface{}) {
	authRouter := r.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/login", h.Login).Methods("POST")
	authRouter.HandleFunc("/register", h.Register).Methods("POST")
	authRouter.HandleFunc("/logout", h.Logout).Methods("POST")
	// Check will be registered with auth middleware in router setup
}

// Login handles POST /auth/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req authUsecase.LoginRequest
	if err := ParseJSON(r, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "validation_error", "Неверные данные в запросе", nil)
		return
	}

	if err := ValidateRequest(h.validator, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "validation_error", "Неверные данные в запросе", nil)
		return
	}

	session, err := h.authUC.Login(r.Context(), &req)
	if err != nil {
		WriteError(w, http.StatusUnauthorized, "unauthorized", "Неверные учетные данные", nil)
		return
	}

	WriteJSON(w, http.StatusOK, session)
}

// Register handles POST /auth/register
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req authUsecase.RegisterRequest
	if err := ParseJSON(r, &req); err != nil {
		errMsg := err.Error()
		WriteError(w, http.StatusBadRequest, "validation_error", "Неверные данные в запросе", &errMsg)
		return
	}

	if err := ValidateRequest(h.validator, &req); err != nil {
		errMsg := err.Error()
		WriteError(w, http.StatusBadRequest, "validation_error", "Неверные данные в запросе", &errMsg)
		return
	}

	session, err := h.authUC.Register(r.Context(), &req)
	if err != nil {
		if err.Error() == "username already exists" || err.Error() == "email already exists" {
			WriteError(w, http.StatusConflict, "user_exists", "Пользователь с таким именем или email уже существует", nil)
			return
		}
		WriteError(w, http.StatusBadRequest, "validation_error", err.Error(), nil)
		return
	}

	WriteJSON(w, http.StatusCreated, session)
}

// Logout handles POST /auth/logout
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	WriteJSON(w, http.StatusOK, map[string]string{"message": "Успешный выход из системы"})
}

// Check handles GET /auth/check
func (h *AuthHandler) Check(w http.ResponseWriter, r *http.Request) {
	// Token already validated by middleware
	authHeader := r.Header.Get("Authorization")
	tokenString := authHeader[7:] // Remove "Bearer " prefix

	user, err := h.authUC.CheckToken(r.Context(), tokenString)
	if err != nil {
		WriteError(w, http.StatusUnauthorized, "unauthorized", "Недействительный токен", nil)
		return
	}

	WriteJSON(w, http.StatusOK, user)
}
