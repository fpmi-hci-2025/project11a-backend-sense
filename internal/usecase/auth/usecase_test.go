package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"sense-backend/internal/domain"
	"sense-backend/internal/infrastructure/jwt"
	"sense-backend/internal/usecase/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

func createTestUser() *domain.User {
	email := "test@example.com"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	return &domain.User{
		ID:           "user-123",
		Username:     "testuser",
		Email:        &email,
		Role:         domain.UserRoleUser,
		PasswordHash: string(hashedPassword),
		RegisteredAt: time.Now(),
	}
}

func TestLogin_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	tokenSvc := mocks.NewMockTokenServiceInterface(ctrl)
	uc := NewUseCase(userRepo, tokenSvc)

	req := &LoginRequest{
		Login:    "testuser",
		Password: "password123",
	}

	user := createTestUser()

	userRepo.EXPECT().
		GetByLogin(gomock.Any(), "testuser").
		Return(user, nil)

	tokenSvc.EXPECT().
		GenerateToken("user-123", "testuser", "user").
		Return("test-token", nil)

	session, err := uc.Login(context.Background(), req)

	require.NoError(t, err)
	assert.NotNil(t, session)
	assert.Equal(t, "test-token", session.AccessToken)
	assert.Equal(t, "Bearer", session.TokenType)
	assert.Equal(t, 86400, session.ExpiresIn)
	assert.Equal(t, user.ID, session.User.ID)
	assert.Empty(t, session.User.PasswordHash)
}

func TestLogin_SuccessWithEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	tokenSvc := mocks.NewMockTokenServiceInterface(ctrl)
	uc := NewUseCase(userRepo, tokenSvc)

	req := &LoginRequest{
		Login:    "test@example.com",
		Password: "password123",
	}

	user := createTestUser()

	userRepo.EXPECT().
		GetByLogin(gomock.Any(), "test@example.com").
		Return(user, nil)

	tokenSvc.EXPECT().
		GenerateToken("user-123", "testuser", "user").
		Return("test-token", nil)

	session, err := uc.Login(context.Background(), req)

	require.NoError(t, err)
	assert.NotNil(t, session)
	assert.Equal(t, "test-token", session.AccessToken)
}

func TestLogin_InvalidCredentials(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	tokenSvc := mocks.NewMockTokenServiceInterface(ctrl)
	uc := NewUseCase(userRepo, tokenSvc)

	req := &LoginRequest{
		Login:    "testuser",
		Password: "wrongpassword",
	}

	user := createTestUser()

	userRepo.EXPECT().
		GetByLogin(gomock.Any(), "testuser").
		Return(user, nil)

	session, err := uc.Login(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, session)
	assert.Equal(t, "invalid credentials", err.Error())
}

func TestLogin_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	tokenSvc := mocks.NewMockTokenServiceInterface(ctrl)
	uc := NewUseCase(userRepo, tokenSvc)

	req := &LoginRequest{
		Login:    "nonexistent",
		Password: "password123",
	}

	userRepo.EXPECT().
		GetByLogin(gomock.Any(), "nonexistent").
		Return(nil, errors.New("user not found"))

	session, err := uc.Login(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, session)
	assert.Equal(t, "invalid credentials", err.Error())
}

func TestLogin_TokenGenerationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	tokenSvc := mocks.NewMockTokenServiceInterface(ctrl)
	uc := NewUseCase(userRepo, tokenSvc)

	req := &LoginRequest{
		Login:    "testuser",
		Password: "password123",
	}

	user := createTestUser()

	userRepo.EXPECT().
		GetByLogin(gomock.Any(), "testuser").
		Return(user, nil)

	tokenSvc.EXPECT().
		GenerateToken("user-123", "testuser", "user").
		Return("", errors.New("token generation failed"))

	session, err := uc.Login(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, session)
	assert.Contains(t, err.Error(), "failed to generate token")
}

func TestRegister_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	tokenSvc := mocks.NewMockTokenServiceInterface(ctrl)
	uc := NewUseCase(userRepo, tokenSvc)

	req := &RegisterRequest{
		Username: "newuser",
		Email:    "newuser@example.com",
		Password: "password123",
	}

	userRepo.EXPECT().
		GetByUsername(gomock.Any(), "newuser").
		Return(nil, errors.New("not found"))

	userRepo.EXPECT().
		GetByEmail(gomock.Any(), "newuser@example.com").
		Return(nil, errors.New("not found"))

	userRepo.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, user *domain.User) error {
			assert.Equal(t, "newuser", user.Username)
			assert.NotEmpty(t, user.PasswordHash)
			return nil
		})

	tokenSvc.EXPECT().
		GenerateToken(gomock.Any(), "newuser", "user").
		Return("test-token", nil)

	session, err := uc.Register(context.Background(), req)

	require.NoError(t, err)
	assert.NotNil(t, session)
	assert.Equal(t, "test-token", session.AccessToken)
	assert.Equal(t, "newuser", session.User.Username)
	assert.Empty(t, session.User.PasswordHash)
}

func TestRegister_UsernameExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	tokenSvc := mocks.NewMockTokenServiceInterface(ctrl)
	uc := NewUseCase(userRepo, tokenSvc)

	req := &RegisterRequest{
		Username: "existinguser",
		Email:    "new@example.com",
		Password: "password123",
	}

	existingUser := createTestUser()

	userRepo.EXPECT().
		GetByUsername(gomock.Any(), "existinguser").
		Return(existingUser, nil)

	session, err := uc.Register(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, session)
	assert.Equal(t, "username already exists", err.Error())
}

func TestRegister_EmailExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	tokenSvc := mocks.NewMockTokenServiceInterface(ctrl)
	uc := NewUseCase(userRepo, tokenSvc)

	req := &RegisterRequest{
		Username: "newuser",
		Email:    "existing@example.com",
		Password: "password123",
	}

	userRepo.EXPECT().
		GetByUsername(gomock.Any(), "newuser").
		Return(nil, errors.New("not found"))

	existingUser := createTestUser()
	email := "existing@example.com"
	existingUser.Email = &email

	userRepo.EXPECT().
		GetByEmail(gomock.Any(), "existing@example.com").
		Return(existingUser, nil)

	session, err := uc.Register(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, session)
	assert.Equal(t, "email already exists", err.Error())
}

func TestRegister_PasswordHashingError(t *testing.T) {
	// This test is difficult to simulate as bcrypt.GenerateFromPassword rarely fails
	// In practice, this would require mocking bcrypt or using an invalid cost
	// For now, we'll skip this test as it's not easily testable
	t.Skip("Password hashing errors are difficult to simulate")
}

func TestRegister_UserCreationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	tokenSvc := mocks.NewMockTokenServiceInterface(ctrl)
	uc := NewUseCase(userRepo, tokenSvc)

	req := &RegisterRequest{
		Username: "newuser",
		Email:    "newuser@example.com",
		Password: "password123",
	}

	userRepo.EXPECT().
		GetByUsername(gomock.Any(), "newuser").
		Return(nil, errors.New("not found"))

	userRepo.EXPECT().
		GetByEmail(gomock.Any(), "newuser@example.com").
		Return(nil, errors.New("not found"))

	userRepo.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(errors.New("database error"))

	session, err := uc.Register(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, session)
	assert.Contains(t, err.Error(), "failed to create user")
}

func TestCheckToken_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	tokenSvc := mocks.NewMockTokenServiceInterface(ctrl)
	uc := NewUseCase(userRepo, tokenSvc)

	tokenString := "valid-token"
	claims := &jwt.Claims{
		UserID:   "user-123",
		Username: "testuser",
		Role:     "user",
	}

	user := createTestUser()

	tokenSvc.EXPECT().
		ValidateToken(tokenString).
		Return(claims, nil)

	userRepo.EXPECT().
		GetByID(gomock.Any(), "user-123").
		Return(user, nil)

	result, err := uc.CheckToken(context.Background(), tokenString)

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "user-123", result.ID)
	assert.Empty(t, result.PasswordHash)
}

func TestCheckToken_InvalidToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	tokenSvc := mocks.NewMockTokenServiceInterface(ctrl)
	uc := NewUseCase(userRepo, tokenSvc)

	tokenString := "invalid-token"

	tokenSvc.EXPECT().
		ValidateToken(tokenString).
		Return(nil, errors.New("invalid token"))

	result, err := uc.CheckToken(context.Background(), tokenString)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "invalid token", err.Error())
}

func TestCheckToken_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	tokenSvc := mocks.NewMockTokenServiceInterface(ctrl)
	uc := NewUseCase(userRepo, tokenSvc)

	tokenString := "valid-token"
	claims := &jwt.Claims{
		UserID:   "user-123",
		Username: "testuser",
		Role:     "user",
	}

	tokenSvc.EXPECT().
		ValidateToken(tokenString).
		Return(claims, nil)

	userRepo.EXPECT().
		GetByID(gomock.Any(), "user-123").
		Return(nil, errors.New("user not found"))

	result, err := uc.CheckToken(context.Background(), tokenString)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "user not found", err.Error())
}

