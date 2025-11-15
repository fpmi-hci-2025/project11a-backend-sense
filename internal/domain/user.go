package domain

import "time"

// UserRole represents user role in the system
type UserRole string

const (
	UserRoleReader  UserRole = "reader"
	UserRoleUser   UserRole = "user"
	UserRoleCreator UserRole = "creator"
	UserRoleExpert UserRole = "expert"
	UserRoleSuper  UserRole = "super"
)

// User represents a user in the system
type User struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	Email        *string   `json:"email,omitempty"`
	Phone        *string   `json:"phone,omitempty"`
	IconURL      *string   `json:"icon_url,omitempty"`
	Description *string   `json:"description,omitempty"`
	Role         UserRole  `json:"role"`
	RegisteredAt time.Time `json:"registered_at"`
	PasswordHash string    `json:"-"` // Not exposed in JSON
	Statistic    *UserStatistic `json:"statistic,omitempty"`
}

// UserStatistic represents user statistics
type UserStatistic struct {
	PublicationsCount int `json:"publications_count"`
	FollowersCount    int `json:"followers_count"`
	FollowingCount    int `json:"following_count"`
	LikesReceived     int `json:"likes_received"`
	CommentsReceived  int `json:"comments_received"`
	SavedCount        int `json:"saved_count"`
}

