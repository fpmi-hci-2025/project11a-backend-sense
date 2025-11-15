package domain

import "time"

// Comment represents a comment on a publication
type Comment struct {
	ID           string     `json:"id"`
	PublicationID string    `json:"publication_id"`
	ParentID     *string    `json:"parent_id,omitempty"`
	AuthorID     string     `json:"author_id"`
	Text         string     `json:"text"`
	CreatedAt    time.Time  `json:"created_at"`
	LikesCount   int        `json:"likes_count"`
}

