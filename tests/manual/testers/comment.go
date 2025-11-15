package testers

import (
	"fmt"
	"net/http"

	"sense-backend/tests/manual/client"
	"sense-backend/tests/manual/testdata"
)

// TestCommentEndpoints tests all comment endpoints
func TestCommentEndpoints(c *client.Client) error {
	fmt.Println("\n=== Testing Comment Endpoints ===")

	data := testdata.GetTestData()
	if data == nil {
		return fmt.Errorf("test data not loaded")
	}

	c.SetToken(data.Tokens.User1)

	// Ensure we have a publication
	if data.IDs.PublicationID == "" {
		return fmt.Errorf("no publication ID available, create a publication first")
	}

	// Test 1: Create comment
	fmt.Printf("\n1. Testing POST /publication/%s/comments\n", data.IDs.PublicationID)
	createReq := map[string]string{
		"text": "This is a test comment",
	}

	resp, err := c.DoRequest("POST", "/publication/"+data.IDs.PublicationID+"/comments", createReq)
	if err != nil {
		return fmt.Errorf("create comment failed: %w", err)
	}

	if err := client.CheckStatus(resp, http.StatusCreated); err != nil {
		return fmt.Errorf("create comment status check failed: %w", err)
	}

	var commentResp struct {
		ID            string `json:"id"`
		PublicationID string `json:"publication_id"`
		AuthorID      string `json:"author_id"`
		Text          string `json:"text"`
		CreatedAt     string `json:"created_at"`
		LikesCount    int    `json:"likes_count"`
		Author        struct {
			ID       string `json:"id"`
			Username string `json:"username"`
		} `json:"author"`
		IsLiked bool `json:"is_liked"`
	}

	if err := client.ParseResponse(resp, &commentResp); err != nil {
		return fmt.Errorf("create comment parse failed: %w", err)
	}

	testdata.SetCommentID(commentResp.ID)
	fmt.Printf("   ✓ Comment created: ID=%s\n", commentResp.ID)

	// Test 2: Get comments for publication
	fmt.Printf("\n2. Testing GET /publication/%s/comments\n", data.IDs.PublicationID)
	resp, err = c.DoRequest("GET", "/publication/"+data.IDs.PublicationID+"/comments?limit=10&offset=0", nil)
	if err != nil {
		return fmt.Errorf("get comments failed: %w", err)
	}

	if err := client.CheckStatus(resp, http.StatusOK); err != nil {
		return fmt.Errorf("get comments status check failed: %w", err)
	}

	var commentsResp struct {
		Items  []interface{} `json:"items"`
		Total  int           `json:"total"`
		Limit  int           `json:"limit"`
		Offset int           `json:"offset"`
	}

	if err := client.ParseResponse(resp, &commentsResp); err != nil {
		return fmt.Errorf("get comments parse failed: %w", err)
	}

	fmt.Printf("   ✓ Comments retrieved: Total=%d\n", commentsResp.Total)

	// Test 3: Get comment by ID
	fmt.Printf("\n3. Testing GET /comment/%s\n", data.IDs.CommentID)
	resp, err = c.DoRequest("GET", "/comment/"+data.IDs.CommentID, nil)
	if err != nil {
		return fmt.Errorf("get comment failed: %w", err)
	}

	if err := client.CheckStatus(resp, http.StatusOK); err != nil {
		return fmt.Errorf("get comment status check failed: %w", err)
	}

	if err := client.ParseResponse(resp, &commentResp); err != nil {
		return fmt.Errorf("get comment parse failed: %w", err)
	}

	if commentResp.ID != data.IDs.CommentID {
		return fmt.Errorf("comment ID mismatch")
	}

	textPreview := commentResp.Text
	if len(textPreview) > 30 {
		textPreview = textPreview[:30]
	}
	fmt.Printf("   ✓ Comment retrieved: Text=%s\n", textPreview)

	// Test 4: Update comment
	fmt.Printf("\n4. Testing PUT /comment/%s\n", data.IDs.CommentID)
	updateReq := map[string]string{
		"text": "Updated comment text",
	}

	resp, err = c.DoRequest("PUT", "/comment/"+data.IDs.CommentID, updateReq)
	if err != nil {
		return fmt.Errorf("update comment failed: %w", err)
	}

	if err := client.CheckStatus(resp, http.StatusOK); err != nil {
		return fmt.Errorf("update comment status check failed: %w", err)
	}

	if err := client.ParseResponse(resp, &commentResp); err != nil {
		return fmt.Errorf("update comment parse failed: %w", err)
	}

	if commentResp.Text != "Updated comment text" {
		return fmt.Errorf("comment text not updated correctly")
	}

	fmt.Println("   ✓ Comment updated")

	// Test 5: Like comment
	fmt.Printf("\n5. Testing POST /comment/%s/like\n", data.IDs.CommentID)
	resp, err = c.DoRequest("POST", "/comment/"+data.IDs.CommentID+"/like", nil)
	if err != nil {
		return fmt.Errorf("like comment failed: %w", err)
	}

	if err := client.CheckStatus(resp, http.StatusOK); err != nil {
		return fmt.Errorf("like comment status check failed: %w", err)
	}

	var likeResp struct {
		Liked      bool `json:"liked"`
		LikesCount int  `json:"likes_count"`
	}

	if err := client.ParseResponse(resp, &likeResp); err != nil {
		return fmt.Errorf("like comment parse failed: %w", err)
	}

	if !likeResp.Liked {
		return fmt.Errorf("comment should be liked")
	}

	fmt.Printf("   ✓ Comment liked: Count=%d\n", likeResp.LikesCount)

	// Test 6: Reply to comment
	fmt.Printf("\n6. Testing POST /comment/%s/reply\n", data.IDs.CommentID)
	replyReq := map[string]string{
		"text": "This is a reply to the comment",
	}

	resp, err = c.DoRequest("POST", "/comment/"+data.IDs.CommentID+"/reply", replyReq)
	if err != nil {
		return fmt.Errorf("reply to comment failed: %w", err)
	}

	if err := client.CheckStatus(resp, http.StatusCreated); err != nil {
		return fmt.Errorf("reply to comment status check failed: %w", err)
	}

	var replyResp struct {
		ID        string  `json:"id"`
		ParentID  *string `json:"parent_id"`
		Text      string  `json:"text"`
		AuthorID  string  `json:"author_id"`
		CreatedAt string  `json:"created_at"`
	}

	if err := client.ParseResponse(resp, &replyResp); err != nil {
		return fmt.Errorf("reply to comment parse failed: %w", err)
	}

	if replyResp.ParentID == nil || *replyResp.ParentID != data.IDs.CommentID {
		return fmt.Errorf("reply parent_id not set correctly")
	}

	fmt.Printf("   ✓ Reply created: ID=%s, ParentID=%s\n", replyResp.ID, *replyResp.ParentID)

	// Test 7: Delete comment
	fmt.Printf("\n7. Testing DELETE /comment/%s\n", data.IDs.CommentID)
	resp, err = c.DoRequest("DELETE", "/comment/"+data.IDs.CommentID, nil)
	if err != nil {
		return fmt.Errorf("delete comment failed: %w", err)
	}

	if err := client.CheckStatus(resp, http.StatusNoContent); err != nil {
		return fmt.Errorf("delete comment status check failed: %w", err)
	}

	fmt.Println("   ✓ Comment deleted")

	// Test 8: Verify deletion
	fmt.Printf("\n8. Testing GET /comment/%s (after deletion)\n", data.IDs.CommentID)
	resp, err = c.DoRequest("GET", "/comment/"+data.IDs.CommentID, nil)
	if err != nil {
		return fmt.Errorf("get deleted comment request failed: %w", err)
	}

	if resp.StatusCode != http.StatusNotFound {
		return fmt.Errorf("expected 404 for deleted comment, got %d", resp.StatusCode)
	}

	fmt.Println("   ✓ Deleted comment correctly returns 404")

	fmt.Println("\n=== Comment Endpoints Testing Complete ===")
	return nil
}
