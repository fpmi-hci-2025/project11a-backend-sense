package testers

import (
	"fmt"
	"net/http"

	"sense-backend/tests/manual/client"
	"sense-backend/tests/manual/testdata"
)

// TestAIEndpoints tests all AI endpoints
func TestAIEndpoints(c *client.Client) error {
	fmt.Println("\n=== Testing AI Endpoints ===")

	data := testdata.GetTestData()
	if data == nil {
		return fmt.Errorf("test data not loaded")
	}

	c.SetToken(data.Tokens.User1)

	// Test 1: Get recommendations
	fmt.Println("\n1. Testing POST /recommendations")
	recommendReq := map[string]interface{}{
		"limit":     20,
		"algorithm": "collaborative_filtering",
	}

	resp, err := c.DoRequest("POST", "/recommendations", recommendReq)
	if err != nil {
		return fmt.Errorf("get recommendations failed: %w", err)
	}

	if err := client.CheckStatus(resp, http.StatusOK); err != nil {
		return fmt.Errorf("get recommendations status check failed: %w", err)
	}

	var recommendationsResp struct {
		Items []struct {
			ID              string `json:"id"`
			AuthorID        string `json:"author_id"`
			Type            string `json:"type"`
			Content         string `json:"content,omitempty"`
			PublicationDate string `json:"publication_date"`
			Visibility      string `json:"visibility"`
			LikesCount      int    `json:"likes_count"`
			CommentsCount   int    `json:"comments_count"`
			SavedCount      int    `json:"saved_count"`
			Recommendation  struct {
				Algorithm string `json:"algorithm"`
				Reason    string `json:"reason"`
			} `json:"recommendation"`
		} `json:"items"`
		Total int `json:"total"`
	}

	if err := client.ParseResponse(resp, &recommendationsResp); err != nil {
		return fmt.Errorf("get recommendations parse failed: %w", err)
	}

	fmt.Printf("   ✓ Recommendations retrieved: Count=%d\n", len(recommendationsResp.Items))

	// Test 2: Get recommendations with default limit
	fmt.Println("\n2. Testing POST /recommendations (with default limit)")
	resp, err = c.DoRequest("POST", "/recommendations", nil)
	if err != nil {
		return fmt.Errorf("get recommendations with defaults failed: %w", err)
	}

	if err := client.CheckStatus(resp, http.StatusOK); err != nil {
		return fmt.Errorf("get recommendations with defaults status check failed: %w", err)
	}

	if err := client.ParseResponse(resp, &recommendationsResp); err != nil {
		return fmt.Errorf("get recommendations with defaults parse failed: %w", err)
	}

	fmt.Printf("   ✓ Recommendations with defaults retrieved: Count=%d\n", len(recommendationsResp.Items))

	// Test 3: Get recommendations feed
	fmt.Println("\n3. Testing GET /recommendations/feed")
	resp, err = c.DoRequest("GET", "/recommendations/feed?limit=10&offset=0", nil)
	if err != nil {
		return fmt.Errorf("get recommendations feed failed: %w", err)
	}

	if err := client.CheckStatus(resp, http.StatusOK); err != nil {
		return fmt.Errorf("get recommendations feed status check failed: %w", err)
	}

	var feedResp struct {
		Items []struct {
			ID              string `json:"id"`
			AuthorID        string `json:"author_id"`
			Type            string `json:"type"`
			Content         string `json:"content,omitempty"`
			PublicationDate string `json:"publication_date"`
			Visibility      string `json:"visibility"`
			LikesCount      int    `json:"likes_count"`
			CommentsCount   int    `json:"comments_count"`
			SavedCount      int    `json:"saved_count"`
			Recommendation  struct {
				Algorithm string `json:"algorithm"`
				Reason    string `json:"reason"`
			} `json:"recommendation"`
		} `json:"items"`
		Total  int `json:"total"`
		Limit  int `json:"limit"`
		Offset int `json:"offset"`
	}

	if err := client.ParseResponse(resp, &feedResp); err != nil {
		return fmt.Errorf("get recommendations feed parse failed: %w", err)
	}

	fmt.Printf("   ✓ Recommendations feed retrieved: Count=%d, Total=%d\n", len(feedResp.Items), feedResp.Total)

	// Test 4: Get recommendations feed with pagination
	fmt.Println("\n4. Testing GET /recommendations/feed (with pagination)")
	resp, err = c.DoRequest("GET", "/recommendations/feed?limit=5&offset=0&algorithm=collaborative_filtering", nil)
	if err != nil {
		return fmt.Errorf("get recommendations feed with pagination failed: %w", err)
	}

	if err := client.CheckStatus(resp, http.StatusOK); err != nil {
		return fmt.Errorf("get recommendations feed with pagination status check failed: %w", err)
	}

	if err := client.ParseResponse(resp, &feedResp); err != nil {
		return fmt.Errorf("get recommendations feed with pagination parse failed: %w", err)
	}

	fmt.Printf("   ✓ Recommendations feed with pagination retrieved: Count=%d, Limit=%d, Offset=%d\n", len(feedResp.Items), feedResp.Limit, feedResp.Offset)

	// Test 5: Hide recommendation (if we have recommendations)
	if len(feedResp.Items) > 0 {
		// First, we need to get the recommendation ID from the database
		// For now, we'll test with a non-existent ID to verify error handling
		fmt.Println("\n5. Testing POST /recommendations/{id}/hide (with non-existent ID)")
		resp, err = c.DoRequest("POST", "/recommendations/00000000-0000-0000-0000-000000000000/hide", nil)
		if err != nil {
			return fmt.Errorf("hide recommendation request failed: %w", err)
		}

		// This should return 404 if recommendation doesn't exist
		if resp.StatusCode != http.StatusNotFound {
			fmt.Printf("   ⚠ Expected 404 for non-existent recommendation, got %d\n", resp.StatusCode)
		} else {
			fmt.Println("   ✓ Non-existent recommendation correctly returns 404")
		}
	} else {
		fmt.Println("\n5. Skipping hide recommendation test (no recommendations available)")
	}

	// Test 6: Purify text
	fmt.Println("\n6. Testing POST /purify")
	purifyReq := map[string]string{
		"text": "This is a test text that needs to be purified and cleaned.",
	}

	resp, err = c.DoRequest("POST", "/purify", purifyReq)
	if err != nil {
		return fmt.Errorf("purify text failed: %w", err)
	}

	if err := client.CheckStatus(resp, http.StatusOK); err != nil {
		return fmt.Errorf("purify text status check failed: %w", err)
	}

	var purifyResp struct {
		CleanedText string  `json:"cleaned_text"`
		IsClean     bool    `json:"is_clean"`
		Confidence  float64 `json:"confidence,omitempty"`
	}

	if err := client.ParseResponse(resp, &purifyResp); err != nil {
		return fmt.Errorf("purify text parse failed: %w", err)
	}

	fmt.Printf("   ✓ Text purified: IsClean=%v, Confidence=%.2f\n", purifyResp.IsClean, purifyResp.Confidence)

	// Test 7: Purify text with potentially problematic content
	fmt.Println("\n7. Testing POST /purify (with longer text)")
	purifyReq2 := map[string]string{
		"text": "This is a longer text that contains multiple sentences. It should be processed correctly by the AI service. The text purification should handle various types of content appropriately.",
	}

	resp, err = c.DoRequest("POST", "/purify", purifyReq2)
	if err != nil {
		return fmt.Errorf("purify longer text failed: %w", err)
	}

	if err := client.CheckStatus(resp, http.StatusOK); err != nil {
		return fmt.Errorf("purify longer text status check failed: %w", err)
	}

	if err := client.ParseResponse(resp, &purifyResp); err != nil {
		return fmt.Errorf("purify longer text parse failed: %w", err)
	}

	fmt.Printf("   ✓ Longer text purified: IsClean=%v\n", purifyResp.IsClean)

	// Test 8: Test without authentication
	fmt.Println("\n8. Testing POST /recommendations (without authentication)")
	c.SetToken("")
	resp, err = c.DoRequest("POST", "/recommendations", recommendReq)
	if err != nil {
		return fmt.Errorf("get recommendations without auth request failed: %w", err)
	}

	if resp.StatusCode != http.StatusUnauthorized {
		return fmt.Errorf("expected 401 for unauthenticated request, got %d", resp.StatusCode)
	}

	fmt.Println("   ✓ Unauthenticated request correctly returns 401")

	// Restore token for further tests
	c.SetToken(data.Tokens.User1)

	// Test 9: Test purify without authentication
	fmt.Println("\n9. Testing POST /purify (without authentication)")
	c.SetToken("")
	resp, err = c.DoRequest("POST", "/purify", purifyReq)
	if err != nil {
		return fmt.Errorf("purify without auth request failed: %w", err)
	}

	if resp.StatusCode != http.StatusUnauthorized {
		return fmt.Errorf("expected 401 for unauthenticated purify request, got %d", resp.StatusCode)
	}

	fmt.Println("   ✓ Unauthenticated purify request correctly returns 401")

	// Restore token
	c.SetToken(data.Tokens.User1)

	fmt.Println("\n=== AI Endpoints Testing Complete ===")
	return nil
}


