package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"sense-backend/pkg/config"
)

// Claims represents JWT claims
type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// TokenService handles JWT token operations
type TokenService struct {
	secret []byte
	expiry time.Duration
}

// NewTokenService creates a new token service
func NewTokenService(cfg *config.JWTConfig) *TokenService {
	return &TokenService{
		secret: []byte(cfg.Secret),
		expiry: time.Duration(cfg.Expiry) * time.Second,
	}
}

// GenerateToken generates a new JWT token for user
func (ts *TokenService) GenerateToken(userID, username, role string) (string, error) {
	claims := &Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ts.expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(ts.secret)
}

// ValidateToken validates and parses JWT token
func (ts *TokenService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return ts.secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

