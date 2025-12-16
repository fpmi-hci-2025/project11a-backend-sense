package testers

import (
	"fmt"
	"net/http"

	"sense-backend/tests/manual/client"
	"sense-backend/tests/manual/testdata"
)

// TestPublicationEndpoints tests all publication endpoints
func TestPublicationEndpoints(c *client.Client) error {
	fmt.Println("\n=== Testing Publication Endpoints ===")

	data := testdata.GetTestData()
	if data == nil {
		return fmt.Errorf("test data not loaded")
	}

	c.SetToken(data.Tokens.User1)

	// Test 1: Create post publication
	fmt.Println("\n1. Testing POST /publication/create (post)")
	createReq := map[string]interface{}{
		"type":       "post",
		"content":    "This is a test post publication",
		"visibility": "public",
	}

	resp, err := c.DoRequest("POST", "/publication/create", createReq)
	if err != nil {
		return fmt.Errorf("create post failed: %w", err)
	}

	if err := client.CheckStatus(resp, http.StatusCreated); err != nil {
		return fmt.Errorf("create post status check failed: %w", err)
	}

	var pubResp struct {
		ID              string `json:"id"`
		AuthorID        string `json:"author_id"`
		Type            string `json:"type"`
		Content         string `json:"content"`
		Visibility      string `json:"visibility"`
		PublicationDate string `json:"publication_date"`
		LikesCount      int    `json:"likes_count"`
		CommentsCount   int    `json:"comments_count"`
		SavedCount      int    `json:"saved_count"`
		Author          struct {
			ID       string `json:"id"`
			Username string `json:"username"`
		} `json:"author"`
		IsLiked bool `json:"is_liked"`
		IsSaved bool `json:"is_saved"`
	}

	if err := client.ParseResponse(resp, &pubResp); err != nil {
		return fmt.Errorf("create post parse failed: %w", err)
	}

	testdata.SetPublicationID(pubResp.ID)
	fmt.Printf("   ✓ Post created: ID=%s, Type=%s\n", pubResp.ID, pubResp.Type)

	// Test 2: Create article publication
	fmt.Println("\n2. Testing POST /publication/create (article)")
	createReq2 := map[string]interface{}{
		"type":       "article",
		"content":    "This is a test article with longer content. It contains multiple sentences and paragraphs.",
		"visibility": "public",
	}

	resp, err = c.DoRequest("POST", "/publication/create", createReq2)
	if err != nil {
		return fmt.Errorf("create article failed: %w", err)
	}

	if err := client.CheckStatus(resp, http.StatusCreated); err != nil {
		return fmt.Errorf("create article status check failed: %w", err)
	}

	if err := client.ParseResponse(resp, &pubResp); err != nil {
		return fmt.Errorf("create article parse failed: %w", err)
	}

	fmt.Printf("   ✓ Article created: ID=%s\n", pubResp.ID)

	// Test 3: Create quote publication
	fmt.Println("\n3. Testing POST /publication/create (quote)")
	createReq3 := map[string]interface{}{
		"type":       "quote",
		"content":    "The only way to do great work is to love what you do.",
		"source":     "Steve Jobs",
		"visibility": "public",
	}

	resp, err = c.DoRequest("POST", "/publication/create", createReq3)
	if err != nil {
		return fmt.Errorf("create quote failed: %w", err)
	}

	if err := client.CheckStatus(resp, http.StatusCreated); err != nil {
		return fmt.Errorf("create quote status check failed: %w", err)
	}

	fmt.Println("   ✓ Quote created")

	// Test 4: Get publication
	fmt.Printf("\n4. Testing GET /publication/%s\n", data.IDs.PublicationID)
	resp, err = c.DoRequest("GET", "/publication/"+data.IDs.PublicationID, nil)
	if err != nil {
		return fmt.Errorf("get publication failed: %w", err)
	}

	if err := client.CheckStatus(resp, http.StatusOK); err != nil {
		return fmt.Errorf("get publication status check failed: %w", err)
	}

	if err := client.ParseResponse(resp, &pubResp); err != nil {
		return fmt.Errorf("get publication parse failed: %w", err)
	}

	if pubResp.ID != data.IDs.PublicationID {
		return fmt.Errorf("publication ID mismatch")
	}

	fmt.Printf("   ✓ Publication retrieved: Content=%s\n", pubResp.Content[:30])

	// Test 5: Get non-existent publication
	fmt.Println("\n5. Testing GET /publication/non-existent-id")
	resp, err = c.DoRequest("GET", "/publication/00000000-0000-0000-0000-000000000000", nil)
	if err != nil {
		return fmt.Errorf("get non-existent publication request failed: %w", err)
	}

	if resp.StatusCode != http.StatusNotFound {
		return fmt.Errorf("expected 404 for non-existent publication, got %d", resp.StatusCode)
	}

	fmt.Println("   ✓ Non-existent publication correctly returns 404")

	// Test 6: Update publication
	fmt.Printf("\n6. Testing PUT /publication/%s\n", data.IDs.PublicationID)
	updateReq := map[string]interface{}{
		"content":    "Updated content for the publication",
		"visibility": "public",
	}

	resp, err = c.DoRequest("PUT", "/publication/"+data.IDs.PublicationID, updateReq)
	if err != nil {
		return fmt.Errorf("update publication failed: %w", err)
	}

	if err := client.CheckStatus(resp, http.StatusOK); err != nil {
		return fmt.Errorf("update publication status check failed: %w", err)
	}

	if err := client.ParseResponse(resp, &pubResp); err != nil {
		return fmt.Errorf("update publication parse failed: %w", err)
	}

	if pubResp.Content != "Updated content for the publication" {
		return fmt.Errorf("content not updated correctly")
	}

	fmt.Println("   ✓ Publication updated")

	// Test 7: Like publication
	fmt.Printf("\n7. Testing POST /publication/%s/like\n", data.IDs.PublicationID)
	resp, err = c.DoRequest("POST", "/publication/"+data.IDs.PublicationID+"/like", nil)
	if err != nil {
		return fmt.Errorf("like publication failed: %w", err)
	}

	if err := client.CheckStatus(resp, http.StatusOK); err != nil {
		return fmt.Errorf("like publication status check failed: %w", err)
	}

	var likeResp struct {
		Liked      bool `json:"liked"`
		LikesCount int  `json:"likes_count"`
	}

	if err := client.ParseResponse(resp, &likeResp); err != nil {
		return fmt.Errorf("like publication parse failed: %w", err)
	}

	if !likeResp.Liked {
		return fmt.Errorf("publication should be liked")
	}

	fmt.Printf("   ✓ Publication liked: Count=%d\n", likeResp.LikesCount)

	// Test 8: Unlike publication (toggle)
	fmt.Printf("\n8. Testing POST /publication/%s/like (toggle)\n", data.IDs.PublicationID)
	resp, err = c.DoRequest("POST", "/publication/"+data.IDs.PublicationID+"/like", nil)
	if err != nil {
		return fmt.Errorf("unlike publication failed: %w", err)
	}

	if err := client.ParseResponse(resp, &likeResp); err != nil {
		return fmt.Errorf("unlike publication parse failed: %w", err)
	}

	if likeResp.Liked {
		return fmt.Errorf("publication should be unliked")
	}

	fmt.Printf("   ✓ Publication unliked: Count=%d\n", likeResp.LikesCount)

	// Test 9: Get likes list
	fmt.Printf("\n9. Testing GET /publication/%s/likes\n", data.IDs.PublicationID)
	resp, err = c.DoRequest("GET", "/publication/"+data.IDs.PublicationID+"/likes?limit=10&offset=0", nil)
	if err != nil {
		return fmt.Errorf("get likes failed: %w", err)
	}

	if err := client.CheckStatus(resp, http.StatusOK); err != nil {
		return fmt.Errorf("get likes status check failed: %w", err)
	}

	var likesResp struct {
		Items []struct {
			ID       string `json:"id"`
			Username string `json:"username"`
		} `json:"items"`
		Total int `json:"total"`
		Limit int `json:"limit"`
		Offset int `json:"offset"`
	}

	if err := client.ParseResponse(resp, &likesResp); err != nil {
		return fmt.Errorf("get likes parse failed: %w", err)
	}

	fmt.Printf("   ✓ Likes list retrieved: Total=%d\n", likesResp.Total)

	// Test 10: Save publication
	fmt.Printf("\n10. Testing POST /publication/%s/save\n", data.IDs.PublicationID)
	saveReq := map[string]string{
		"note": "This is a test note for saved publication",
	}

	resp, err = c.DoRequest("POST", "/publication/"+data.IDs.PublicationID+"/save", saveReq)
	if err != nil {
		return fmt.Errorf("save publication failed: %w", err)
	}

	if err := client.CheckStatus(resp, http.StatusCreated); err != nil {
		return fmt.Errorf("save publication status check failed: %w", err)
	}

	fmt.Println("   ✓ Publication saved")

	// Test 11: Delete publication (will be recreated if needed)
	fmt.Printf("\n11. Testing DELETE /publication/%s\n", data.IDs.PublicationID)
	resp, err = c.DoRequest("DELETE", "/publication/"+data.IDs.PublicationID, nil)
	if err != nil {
		return fmt.Errorf("delete publication failed: %w", err)
	}

	if err := client.CheckStatus(resp, http.StatusNoContent); err != nil {
		return fmt.Errorf("delete publication status check failed: %w", err)
	}

	fmt.Println("   ✓ Publication deleted")

	// Test 12: Verify deletion
	fmt.Printf("\n12. Testing GET /publication/%s (after deletion)\n", data.IDs.PublicationID)
	resp, err = c.DoRequest("GET", "/publication/"+data.IDs.PublicationID, nil)
	if err != nil {
		return fmt.Errorf("get deleted publication request failed: %w", err)
	}

	if resp.StatusCode != http.StatusNotFound {
		return fmt.Errorf("expected 404 for deleted publication, got %d", resp.StatusCode)
	}

	fmt.Println("   ✓ Deleted publication correctly returns 404")

	// Create a new publication for further tests
	fmt.Println("\n13. Creating new publication for further tests")
	resp, err = c.DoRequest("POST", "/publication/create", createReq)
	if err != nil {
		return fmt.Errorf("create new publication failed: %w", err)
	}

	if err := client.ParseResponse(resp, &pubResp); err != nil {
		return fmt.Errorf("create new publication parse failed: %w", err)
	}

	testdata.SetPublicationID(pubResp.ID)
	fmt.Printf("   ✓ New publication created: ID=%s\n", pubResp.ID)

	fmt.Println("\n=== Publication Endpoints Testing Complete ===")
	return nil
}

