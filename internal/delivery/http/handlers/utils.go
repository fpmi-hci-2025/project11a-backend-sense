package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

// ErrorResponse represents error response
type ErrorResponse struct {
	Error   string  `json:"error"`
	Message string  `json:"message"`
	Details *string `json:"details,omitempty"`
}

// WriteError writes error response
func WriteError(w http.ResponseWriter, statusCode int, errorCode, message string, details *string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error:   errorCode,
		Message: message,
		Details: details,
	})
}

// WriteJSON writes JSON response
func WriteJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// ValidateRequest validates request using validator
func ValidateRequest(v *validator.Validate, req interface{}) error {
	return v.Struct(req)
}

// ParseJSON parses JSON request body
func ParseJSON(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

// ParseMultipartForm parses multipart form data and extracts file and description
func ParseMultipartForm(r *http.Request, maxMemory int64) (file multipart.File, fileHeader *multipart.FileHeader, description *string, err error) {
	if err := r.ParseMultipartForm(maxMemory); err != nil {
		return nil, nil, nil, fmt.Errorf("failed to parse multipart form: %w", err)
	}

	file, fileHeader, err = r.FormFile("file")
	if err != nil {
		return nil, nil, nil, fmt.Errorf("file field is required: %w", err)
	}

	if desc := r.FormValue("description"); desc != "" {
		description = &desc
	}

	return file, fileHeader, description, nil
}

// ValidateFileSize checks if file size is within limit
func ValidateFileSize(size int64, maxSize int64) error {
	if size > maxSize {
		return fmt.Errorf("file size exceeds maximum allowed size of %d bytes", maxSize)
	}
	return nil
}

// ValidateMIMEType checks if MIME type is allowed
func ValidateMIMEType(mimeType string) error {
	allowedTypes := []string{
		"image/jpeg",
		"image/jpg",
		"image/png",
		"image/gif",
		"image/webp",
	}

	mimeType = strings.ToLower(strings.TrimSpace(mimeType))
	for _, allowed := range allowedTypes {
		if mimeType == allowed {
			return nil
		}
	}

	return fmt.Errorf("MIME type %s is not allowed. Allowed types: %v", mimeType, allowedTypes)
}

// ImageMetadata contains extracted image metadata
type ImageMetadata struct {
	Width  *int
	Height *int
	EXIF   interface{} // EXIF data (can be enhanced later)
}

// ExtractImageMetadata extracts width, height, and optionally EXIF from image data
func ExtractImageMetadata(data []byte, mimeType string) (*ImageMetadata, error) {
	var img image.Image
	var err error

	mimeType = strings.ToLower(strings.TrimSpace(mimeType))
	reader := bytes.NewReader(data)

	switch mimeType {
	case "image/jpeg", "image/jpg":
		img, _, err = image.Decode(reader)
	case "image/png":
		img, _, err = image.Decode(reader)
	case "image/gif":
		img, _, err = image.Decode(reader)
	case "image/webp":
		// WebP support requires golang.org/x/image/webp
		// For now, skip metadata extraction for WebP
		return &ImageMetadata{
			Width:  nil,
			Height: nil,
			EXIF:   nil,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported image type: %s", mimeType)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	return &ImageMetadata{
		Width:  &width,
		Height: &height,
		EXIF:   nil, // EXIF extraction can be added later with goexif
	}, nil
}

// ReadFileContent reads all content from a file
func ReadFileContent(file multipart.File) ([]byte, error) {
	defer file.Close()
	return io.ReadAll(file)
}

// GetMIMETypeFromFilename extracts MIME type from filename
func GetMIMETypeFromFilename(filename string) string {
	ext := ""
	if idx := strings.LastIndex(filename, "."); idx != -1 {
		ext = filename[idx:]
	}
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		return "application/octet-stream"
	}
	return mimeType
}

