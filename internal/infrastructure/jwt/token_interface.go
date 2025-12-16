package jwt

// TokenServiceInterface defines interface for token operations
// This interface is used for mocking in tests
type TokenServiceInterface interface {
	GenerateToken(userID, username, role string) (string, error)
	ValidateToken(tokenString string) (*Claims, error)
}

// Ensure TokenService implements TokenServiceInterface
var _ TokenServiceInterface = (*TokenService)(nil)

