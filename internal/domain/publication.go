package domain

import "time"

// PublicationType represents type of publication
type PublicationType string

const (
	PublicationTypeQuote   PublicationType = "quote"
	PublicationTypePost    PublicationType = "post"
	PublicationTypeArticle  PublicationType = "article"
)

// VisibilityType represents visibility level of publication
type VisibilityType string

const (
	VisibilityTypePublic    VisibilityType = "public"
	VisibilityTypeCommunity VisibilityType = "community"
	VisibilityTypePrivate   VisibilityType = "private"
)

// Publication represents a publication in the system
type Publication struct {
	ID             string          `json:"id"`
	AuthorID       string          `json:"author_id"`
	Type           PublicationType `json:"type"`
	Content        *string         `json:"content,omitempty"`
	Source         *string         `json:"source,omitempty"`
	PublicationDate time.Time      `json:"publication_date"`
	Visibility     VisibilityType  `json:"visibility"`
	LikesCount     int             `json:"likes_count"`
	CommentsCount  int             `json:"comments_count"`
	SavedCount     int             `json:"saved_count"`
}

