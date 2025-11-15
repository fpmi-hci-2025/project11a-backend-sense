package domain

import "context"

// UserRepository defines interface for user data operations
type UserRepository interface {
	// Create creates a new user
	Create(ctx context.Context, user *User) error
	
	// GetByID retrieves user by ID
	GetByID(ctx context.Context, id string) (*User, error)
	
	// GetByUsername retrieves user by username
	GetByUsername(ctx context.Context, username string) (*User, error)
	
	// GetByEmail retrieves user by email
	GetByEmail(ctx context.Context, email string) (*User, error)
	
	// GetByLogin retrieves user by username or email
	GetByLogin(ctx context.Context, login string) (*User, error)
	
	// Update updates user information
	Update(ctx context.Context, user *User) error
	
	// GetStats retrieves user statistics
	GetStats(ctx context.Context, userID string) (*UserStatistic, error)
	
	// GetFollowersCount returns number of followers
	GetFollowersCount(ctx context.Context, userID string) (int, error)
	
	// GetFollowingCount returns number of users being followed
	GetFollowingCount(ctx context.Context, userID string) (int, error)
	
	// IsFollowing checks if followerID follows followingID
	IsFollowing(ctx context.Context, followerID, followingID string) (bool, error)
	
	// Follow creates a follow relationship
	Follow(ctx context.Context, followerID, followingID string) error
	
	// Unfollow removes a follow relationship
	Unfollow(ctx context.Context, followerID, followingID string) error
	
	// Search searches users by query
	Search(ctx context.Context, query string, role *UserRole, limit, offset int) ([]*User, int, error)
}

