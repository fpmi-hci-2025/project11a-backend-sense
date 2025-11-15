package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"sense-backend/internal/domain"
	"sense-backend/internal/infrastructure/jwt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// UseCase handles authentication use cases
type UseCase struct {
	userRepo domain.UserRepository
	tokenSvc *jwt.TokenService
}

// NewUseCase creates a new auth use case
func NewUseCase(userRepo domain.UserRepository, tokenSvc *jwt.TokenService) *UseCase {
	return &UseCase{
		userRepo: userRepo,
		tokenSvc: tokenSvc,
	}
}

// LoginRequest represents login request
type LoginRequest struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// RegisterRequest represents registration request
type RegisterRequest struct {
	Username    string  `json:"username" validate:"required,min=3,max=30"`
	Email       string  `json:"email" validate:"required,email"`
	Password    string  `json:"password" validate:"required,min=8"`
	Phone       *string `json:"phone,omitempty"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=500"`
}

// SessionResponse represents session response
type SessionResponse struct {
	AccessToken string       `json:"access_token"`
	TokenType   string       `json:"token_type"`
	ExpiresIn   int          `json:"expires_in"`
	User        *domain.User `json:"user"`
}

// Login authenticates user and returns session
func (uc *UseCase) Login(ctx context.Context, req *LoginRequest) (*SessionResponse, error) {
	// Get user by login (username or email)
	user, err := uc.userRepo.GetByLogin(ctx, req.Login)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Generate token
	token, err := uc.tokenSvc.GenerateToken(user.ID, user.Username, string(user.Role))
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Remove password hash from response
	user.PasswordHash = ""

	return &SessionResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   86400, // 1 day in seconds
		User:        user,
	}, nil
}

// Register creates a new user and returns session
func (uc *UseCase) Register(ctx context.Context, req *RegisterRequest) (*SessionResponse, error) {
	// Check if username exists
	_, err := uc.userRepo.GetByUsername(ctx, req.Username)
	if err == nil {
		return nil, errors.New("username already exists")
	}

	// Check if email exists
	_, err = uc.userRepo.GetByEmail(ctx, req.Email)
	if err == nil {
		return nil, errors.New("email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &domain.User{
		ID:           uuid.New().String(),
		Username:     req.Username,
		Email:        &req.Email,
		Phone:        req.Phone,
		Description:  req.Description,
		Role:         domain.UserRoleUser,
		PasswordHash: string(hashedPassword),
		RegisteredAt: time.Now(),
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate token
	token, err := uc.tokenSvc.GenerateToken(user.ID, user.Username, string(user.Role))
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Remove password hash from response
	user.PasswordHash = ""

	return &SessionResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   86400, // 1 day in seconds
		User:        user,
	}, nil
}

// CheckToken validates token and returns user
func (uc *UseCase) CheckToken(ctx context.Context, tokenString string) (*domain.User, error) {
	claims, err := uc.tokenSvc.ValidateToken(tokenString)
	if err != nil {
		return nil, errors.New("invalid token")
	}

	user, err := uc.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Remove password hash from response
	user.PasswordHash = ""

	return user, nil
}
