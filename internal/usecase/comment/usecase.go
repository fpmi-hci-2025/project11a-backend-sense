package comment

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"sense-backend/internal/domain"
)

// UseCase handles comment use cases
type UseCase struct {
	commentRepo domain.CommentRepository
}

// NewUseCase creates a new comment use case
func NewUseCase(commentRepo domain.CommentRepository) *UseCase {
	return &UseCase{commentRepo: commentRepo}
}

// CreateRequest represents create comment request
type CreateRequest struct {
	Text     string  `json:"text" validate:"required,max=2000"`
	ParentID *string `json:"parent_id,omitempty"`
}

// Create creates a new comment
func (uc *UseCase) Create(ctx context.Context, publicationID, authorID string, req *CreateRequest) (*domain.Comment, error) {
	comment := &domain.Comment{
		ID:           uuid.New().String(),
		PublicationID: publicationID,
		ParentID:     req.ParentID,
		AuthorID:     authorID,
		Text:         req.Text,
		CreatedAt:    time.Now(),
		LikesCount:   0,
	}

	if err := uc.commentRepo.Create(ctx, comment); err != nil {
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}

	return comment, nil
}

// Get retrieves comment by ID
func (uc *UseCase) Get(ctx context.Context, id string) (*domain.Comment, error) {
	return uc.commentRepo.GetByID(ctx, id)
}

// Update updates comment
func (uc *UseCase) Update(ctx context.Context, id, userID string, text string) (*domain.Comment, error) {
	comment, err := uc.commentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if comment.AuthorID != userID {
		return nil, errors.New("forbidden: not the author")
	}

	comment.Text = text
	if err := uc.commentRepo.Update(ctx, comment); err != nil {
		return nil, fmt.Errorf("failed to update comment: %w", err)
	}

	return comment, nil
}

// Delete deletes comment
func (uc *UseCase) Delete(ctx context.Context, id, userID string) error {
	comment, err := uc.commentRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if comment.AuthorID != userID {
		return errors.New("forbidden: not the author")
	}

	return uc.commentRepo.Delete(ctx, id)
}

// Like toggles like on comment
func (uc *UseCase) Like(ctx context.Context, commentID, userID string) (bool, int, error) {
	liked, err := uc.commentRepo.Like(ctx, userID, commentID)
	if err != nil {
		return false, 0, err
	}

	count, err := uc.commentRepo.GetLikesCount(ctx, commentID)
	if err != nil {
		return false, 0, err
	}

	return liked, count, nil
}

// GetByPublication retrieves comments for publication
func (uc *UseCase) GetByPublication(ctx context.Context, publicationID string, limit, offset int) ([]*domain.Comment, int, error) {
	return uc.commentRepo.GetByPublication(ctx, publicationID, limit, offset)
}

