package domain

import (
	"time"
)

// MediaAsset represents a media file in the system
type MediaAsset struct {
	ID        string    `json:"id"`
	OwnerID   string    `json:"owner_id"`
	Filename  *string   `json:"filename,omitempty"`
	MIME      string    `json:"mime"`
	Width     *int      `json:"width,omitempty"`
	Height    *int      `json:"height,omitempty"`
	EXIF      *EXIFData `json:"exif,omitempty"`
	Data      []byte    `json:"-"` // Binary data, not exposed in JSON
	CreatedAt time.Time `json:"created_at"`
}

// EXIFData represents EXIF metadata for images
type EXIFData struct {
	Camera       *string `json:"camera,omitempty"`
	Lens         *string `json:"lens,omitempty"`
	ISO          *int    `json:"iso,omitempty"`
	Aperture     *string `json:"aperture,omitempty"`
	ShutterSpeed *string `json:"shutter_speed,omitempty"`
}
