package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"sense-backend/internal/delivery/http/middleware"
	mediaUsecase "sense-backend/internal/usecase/media"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

// MediaHandler handles media endpoints
type MediaHandler struct {
	mediaUC   *mediaUsecase.UseCase
	validator *validator.Validate
	maxSize   int64
}

// NewMediaHandler creates a new media handler
func NewMediaHandler(mediaUC *mediaUsecase.UseCase, validator *validator.Validate, maxFileSize int64) *MediaHandler {
	return &MediaHandler{
		mediaUC:   mediaUC,
		validator: validator,
		maxSize:   maxFileSize,
	}
}

// RegisterRoutes registers media routes
func (h *MediaHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/upload", h.Upload).Methods("POST")
	r.HandleFunc("/{id}", h.Get).Methods("GET")
	r.HandleFunc("/{id}", h.Delete).Methods("DELETE")
	r.HandleFunc("/{id}/file", h.GetFile).Methods("GET")
}

// Upload handles POST /media/upload
func (h *MediaHandler) Upload(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		WriteError(w, http.StatusUnauthorized, "unauthorized", "Требуется аутентификация", nil)
		return
	}

	// Parse multipart form (max 32MB in memory)
	const maxMemory = 32 << 20 // 32MB
	file, fileHeader, _, err := ParseMultipartForm(r, maxMemory)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "validation_error", "Неверные данные в запросе", nil)
		return
	}
	defer file.Close()

	// Validate file size
	if err := ValidateFileSize(fileHeader.Size, h.maxSize); err != nil {
		details := fmt.Sprintf("Максимальный размер файла: %d байт", h.maxSize)
		WriteError(w, http.StatusRequestEntityTooLarge, "file_too_large", "Размер файла превышает максимально допустимый", &details)
		return
	}

	// Read file content
	fileData, err := ReadFileContent(file)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "validation_error", "Не удалось прочитать файл", nil)
		return
	}

	// Determine MIME type
	mimeType := fileHeader.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = GetMIMETypeFromFilename(fileHeader.Filename)
	}
	mimeType = strings.ToLower(strings.TrimSpace(mimeType))

	// Validate MIME type
	if err := ValidateMIMEType(mimeType); err != nil {
		details := err.Error()
		WriteError(w, http.StatusBadRequest, "validation_error", "Недопустимый тип файла", &details)
		return
	}

	// Extract image metadata if it's an image
	var width, height *int
	if strings.HasPrefix(mimeType, "image/") {
		metadata, err := ExtractImageMetadata(fileData, mimeType)
		if err == nil {
			width = metadata.Width
			height = metadata.Height
			// EXIF can be processed later if needed
		}
		// Ignore metadata extraction errors - not critical
	}

	// Prepare filename
	var filename *string
	if fileHeader.Filename != "" {
		filename = &fileHeader.Filename
	}

	// Create upload request
	uploadReq := &mediaUsecase.UploadRequest{
		Data:     fileData,
		Filename: filename,
		MIME:     mimeType,
		Width:    width,
		Height:   height,
		EXIF:     nil, // EXIF extraction can be added later
	}

	// Upload media
	media, err := h.mediaUC.Upload(r.Context(), userID, uploadReq)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "validation_error", "Не удалось загрузить файл", nil)
		return
	}

	// Return media asset (Data field is excluded from JSON)
	WriteJSON(w, http.StatusCreated, media)
}

// Get handles GET /media/{id}
func (h *MediaHandler) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	media, err := h.mediaUC.Get(r.Context(), id)
	if err != nil {
		WriteError(w, http.StatusNotFound, "not_found", "Медиа-файл не найден", nil)
		return
	}

	// Data field is excluded from JSON via json:"-" tag
	WriteJSON(w, http.StatusOK, media)
}

// GetFile handles GET /media/{id}/file
func (h *MediaHandler) GetFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	media, err := h.mediaUC.Get(r.Context(), id)
	if err != nil {
		WriteError(w, http.StatusNotFound, "not_found", "Медиа-файл не найден", nil)
		return
	}

	// Set headers
	w.Header().Set("Content-Type", media.MIME)
	if media.Filename != nil {
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", *media.Filename))
	}

	// Write binary data
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(media.Data); err != nil {
		// Error writing response - connection may be closed
		return
	}
}

// Delete handles DELETE /media/{id}
func (h *MediaHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		WriteError(w, http.StatusUnauthorized, "unauthorized", "Требуется аутентификация", nil)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.mediaUC.Delete(r.Context(), id, userID); err != nil {
		if err.Error() == "forbidden: not the owner" {
			WriteError(w, http.StatusForbidden, "forbidden", "Недостаточно прав для выполнения операции", nil)
			return
		}
		WriteError(w, http.StatusNotFound, "not_found", "Медиа-файл не найден", nil)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
