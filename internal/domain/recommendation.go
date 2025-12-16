package domain

import "time"

// Recommendation represents an AI-generated recommendation
type Recommendation struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	PublicationID string   `json:"publication_id"`
	Algorithm    string    `json:"algorithm"`
	Reason       string     `json:"reason"`
	Rank         int        `json:"rank"`
	CreatedAt    time.Time `json:"created_at"`
	Hidden       bool      `json:"hidden"`
}

