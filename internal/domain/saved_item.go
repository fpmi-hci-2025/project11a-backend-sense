package domain

import "time"

// SavedItem represents a saved publication by a user
type SavedItem struct {
	ID           string     `json:"id"`
	UserID       string     `json:"user_id"`
	PublicationID string    `json:"publication_id"`
	AddedAt      time.Time  `json:"added_at"`
	Note         *string    `json:"note,omitempty"`
}

