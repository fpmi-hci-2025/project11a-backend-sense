package comment

import (
	"context"
	"errors"
	"testing"
	"time"

	"sense-backend/internal/domain"
	"sense-backend/internal/usecase/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func createTestComment() *domain.Comment {
	return &domain.Comment{
		ID:           "comment-123",
		PublicationID: "pub-123",
		AuthorID:     "user-123",
		Text:         "Test comment",
		CreatedAt:    time.Now(),
		LikesCount:   0,
	}
}

func TestCreate_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	commentRepo := mocks.NewMockCommentRepository(ctrl)
	uc := NewUseCase(commentRepo)

	req := &CreateRequest{
		Text: "Test comment",
	}

	commentRepo.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, comment *domain.Comment) error {
			assert.Equal(t, "pub-123", comment.PublicationID)
			assert.Equal(t, "user-123", comment.AuthorID)
			assert.Equal(t, "Test comment", comment.Text)
			return nil
		})

	comment, err := uc.Create(context.Background(), "pub-123", "user-123", req)

	require.NoError(t, err)
	assert.NotNil(t, comment)
	assert.Equal(t, "Test comment", comment.Text)
}

func TestCreate_WithParent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	commentRepo := mocks.NewMockCommentRepository(ctrl)
	uc := NewUseCase(commentRepo)

	parentID := "parent-123"
	req := &CreateRequest{
		Text:     "Reply comment",
		ParentID: &parentID,
	}

	commentRepo.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, comment *domain.Comment) error {
			assert.Equal(t, "parent-123", *comment.ParentID)
			return nil
		})

	comment, err := uc.Create(context.Background(), "pub-123", "user-123", req)

	require.NoError(t, err)
	assert.NotNil(t, comment)
	assert.Equal(t, "parent-123", *comment.ParentID)
}

func TestGet_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	commentRepo := mocks.NewMockCommentRepository(ctrl)
	uc := NewUseCase(commentRepo)

	comment := createTestComment()

	commentRepo.EXPECT().
		GetByID(gomock.Any(), "comment-123").
		Return(comment, nil)

	result, err := uc.Get(context.Background(), "comment-123")

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "comment-123", result.ID)
}

func TestGet_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	commentRepo := mocks.NewMockCommentRepository(ctrl)
	uc := NewUseCase(commentRepo)

	commentRepo.EXPECT().
		GetByID(gomock.Any(), "nonexistent").
		Return(nil, errors.New("not found"))

	result, err := uc.Get(context.Background(), "nonexistent")

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestUpdate_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	commentRepo := mocks.NewMockCommentRepository(ctrl)
	uc := NewUseCase(commentRepo)

	comment := createTestComment()

	commentRepo.EXPECT().
		GetByID(gomock.Any(), "comment-123").
		Return(comment, nil)

	commentRepo.EXPECT().
		Update(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, comment *domain.Comment) error {
			assert.Equal(t, "Updated text", comment.Text)
			return nil
		})

	result, err := uc.Update(context.Background(), "comment-123", "user-123", "Updated text")

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Updated text", result.Text)
}

func TestUpdate_NotAuthor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	commentRepo := mocks.NewMockCommentRepository(ctrl)
	uc := NewUseCase(commentRepo)

	comment := createTestComment()

	commentRepo.EXPECT().
		GetByID(gomock.Any(), "comment-123").
		Return(comment, nil)

	result, err := uc.Update(context.Background(), "comment-123", "other-user", "Updated text")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "forbidden: not the author", err.Error())
}

func TestDelete_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	commentRepo := mocks.NewMockCommentRepository(ctrl)
	uc := NewUseCase(commentRepo)

	comment := createTestComment()

	commentRepo.EXPECT().
		GetByID(gomock.Any(), "comment-123").
		Return(comment, nil)

	commentRepo.EXPECT().
		Delete(gomock.Any(), "comment-123").
		Return(nil)

	err := uc.Delete(context.Background(), "comment-123", "user-123")

	require.NoError(t, err)
}

func TestDelete_NotAuthor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	commentRepo := mocks.NewMockCommentRepository(ctrl)
	uc := NewUseCase(commentRepo)

	comment := createTestComment()

	commentRepo.EXPECT().
		GetByID(gomock.Any(), "comment-123").
		Return(comment, nil)

	err := uc.Delete(context.Background(), "comment-123", "other-user")

	assert.Error(t, err)
	assert.Equal(t, "forbidden: not the author", err.Error())
}

func TestLike_Toggle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	commentRepo := mocks.NewMockCommentRepository(ctrl)
	uc := NewUseCase(commentRepo)

	commentRepo.EXPECT().
		Like(gomock.Any(), "user-123", "comment-123").
		Return(true, nil)

	commentRepo.EXPECT().
		GetLikesCount(gomock.Any(), "comment-123").
		Return(3, nil)

	liked, count, err := uc.Like(context.Background(), "comment-123", "user-123")

	require.NoError(t, err)
	assert.True(t, liked)
	assert.Equal(t, 3, count)
}

func TestGetByPublication_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	commentRepo := mocks.NewMockCommentRepository(ctrl)
	uc := NewUseCase(commentRepo)

	comments := []*domain.Comment{
		createTestComment(),
		createTestComment(),
	}

	commentRepo.EXPECT().
		GetByPublication(gomock.Any(), "pub-123", 10, 0).
		Return(comments, 2, nil)

	result, total, err := uc.GetByPublication(context.Background(), "pub-123", 10, 0)

	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, 2, total)
}

