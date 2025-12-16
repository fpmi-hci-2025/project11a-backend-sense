package testers

import (
	"fmt"
	"net/http"

	"sense-backend/tests/manual/client"
	"sense-backend/tests/manual/testdata"
)

// TestSearchEndpoints tests all search endpoints
func TestSearchEndpoints(c *client.Client) error {
	fmt.Println("\n=== Testing Search Endpoints ===")

	data := testdata.GetTestData()
	if data == nil {
		return fmt.Errorf("test data not loaded")
	}

	c.SetToken(data.Tokens.User1)

	// Test 1: Search publications with query
	fmt.Println("\n1. Testing GET /search?q=test")
	resp, err := c.DoRequest("GET", "/search?q=test&limit=10&offset=0", nil)
	if err != nil {
		return fmt.Errorf("search publications failed: %w", err)
	}

	if err := client.CheckStatus(resp, http.StatusOK); err != nil {
		return fmt.Errorf("search publications status check failed: %w", err)
	}

	var searchResp struct {
		Items  []interface{} `json:"items"`
		Total  int           `json:"total"`
		Limit  int           `json:"limit"`
		Offset int           `json:"offset"`
	}

	if err := client.ParseResponse(resp, &searchResp); err != nil {
		return fmt.Errorf("search publications parse failed: %w", err)
	}

	fmt.Printf("   ✓ Search returned %d results (total: %d)\n", len(searchResp.Items), searchResp.Total)

	// Test 2: Search publications with filters
	fmt.Println("\n2. Testing GET /search?q=test&type=post&visibility=public")
	resp, err = c.DoRequest("GET", "/search?q=test&type=post&visibility=public&limit=5", nil)
	if err != nil {
		return fmt.Errorf("search with filters failed: %w", err)
	}

	if err := client.CheckStatus(resp, http.StatusOK); err != nil {
		return fmt.Errorf("search with filters status check failed: %w", err)
	}

	if err := client.ParseResponse(resp, &searchResp); err != nil {
		return fmt.Errorf("search with filters parse failed: %w", err)
	}

	fmt.Printf("   ✓ Filtered search returned %d results\n", len(searchResp.Items))

	// Test 3: Search publications without query (should fail)
	fmt.Println("\n3. Testing GET /search (no query parameter)")
	resp, err = c.DoRequest("GET", "/search", nil)
	if err != nil {
		return fmt.Errorf("search without query request failed: %w", err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		return fmt.Errorf("expected 400 for missing query, got %d", resp.StatusCode)
	}

	fmt.Println("   ✓ Missing query correctly returns 400")

	// Test 4: Search users
	fmt.Println("\n4. Testing GET /search/users?q=user")
	resp, err = c.DoRequest("GET", "/search/users?q=user&limit=10", nil)
	if err != nil {
		return fmt.Errorf("search users failed: %w", err)
	}

	if err := client.CheckStatus(resp, http.StatusOK); err != nil {
		return fmt.Errorf("search users status check failed: %w", err)
	}

	var userSearchResp struct {
		Items  []interface{} `json:"items"`
		Total  int           `json:"total"`
		Limit  int           `json:"limit"`
		Offset int           `json:"offset"`
	}

	if err := client.ParseResponse(resp, &userSearchResp); err != nil {
		return fmt.Errorf("search users parse failed: %w", err)
	}

	fmt.Printf("   ✓ User search returned %d results (total: %d)\n", len(userSearchResp.Items), userSearchResp.Total)

	// Test 5: Search users with role filter
	fmt.Println("\n5. Testing GET /search/users?q=user&role=user")
	resp, err = c.DoRequest("GET", "/search/users?q=user&role=user", nil)
	if err != nil {
		return fmt.Errorf("search users with role failed: %w", err)
	}

	if err := client.CheckStatus(resp, http.StatusOK); err != nil {
		return fmt.Errorf("search users with role status check failed: %w", err)
	}

	if err := client.ParseResponse(resp, &userSearchResp); err != nil {
		return fmt.Errorf("search users with role parse failed: %w", err)
	}

	fmt.Printf("   ✓ User search with role filter returned %d results\n", len(userSearchResp.Items))

	// Test 6: Search users without query (should fail)
	fmt.Println("\n6. Testing GET /search/users (no query parameter)")
	resp, err = c.DoRequest("GET", "/search/users", nil)
	if err != nil {
		return fmt.Errorf("search users without query request failed: %w", err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		return fmt.Errorf("expected 400 for missing query, got %d", resp.StatusCode)
	}

	fmt.Println("   ✓ Missing query correctly returns 400")

	// Test 7: Search warmup
	fmt.Println("\n7. Testing POST /search/warmup")
	warmupReq := map[string]interface{}{
		"filters": map[string]string{
			"type":       "post",
			"visibility": "public",
		},
	}

	resp, err = c.DoRequest("POST", "/search/warmup", warmupReq)
	if err != nil {
		return fmt.Errorf("search warmup failed: %w", err)
	}

	if err := client.CheckStatus(resp, http.StatusOK); err != nil {
		return fmt.Errorf("search warmup status check failed: %w", err)
	}

	var warmupResp struct {
		Message string `json:"message"`
		TaskID  string `json:"task_id"`
	}

	if err := client.ParseResponse(resp, &warmupResp); err != nil {
		return fmt.Errorf("search warmup parse failed: %w", err)
	}

	if warmupResp.TaskID == "" {
		return fmt.Errorf("search warmup did not return task_id")
	}

	fmt.Printf("   ✓ Search warmup initiated: TaskID=%s\n", warmupResp.TaskID)

	// Test 8: Get popular tags
	fmt.Println("\n8. Testing GET /tags?limit=10")
	resp, err = c.DoRequest("GET", "/tags?limit=10", nil)
	if err != nil {
		return fmt.Errorf("get tags failed: %w", err)
	}

	if err := client.CheckStatus(resp, http.StatusOK); err != nil {
		return fmt.Errorf("get tags status check failed: %w", err)
	}

	var tagsResp struct {
		Items []struct {
			ID          string `json:"id"`
			Name        string `json:"name"`
			Description string `json:"description"`
			UsageCount  int    `json:"usage_count"`
		} `json:"items"`
		Total int `json:"total"`
	}

	if err := client.ParseResponse(resp, &tagsResp); err != nil {
		return fmt.Errorf("get tags parse failed: %w", err)
	}

	fmt.Printf("   ✓ Retrieved %d tags (total: %d)\n", len(tagsResp.Items), tagsResp.Total)

	// Test 9: Get tags with search filter
	fmt.Println("\n9. Testing GET /tags?search=photo&limit=5")
	resp, err = c.DoRequest("GET", "/tags?search=photo&limit=5", nil)
	if err != nil {
		return fmt.Errorf("get tags with search failed: %w", err)
	}

	if err := client.CheckStatus(resp, http.StatusOK); err != nil {
		return fmt.Errorf("get tags with search status check failed: %w", err)
	}

	if err := client.ParseResponse(resp, &tagsResp); err != nil {
		return fmt.Errorf("get tags with search parse failed: %w", err)
	}

	fmt.Printf("   ✓ Tag search returned %d results\n", len(tagsResp.Items))

	// Test 10: Search publications without authentication (optional auth)
	fmt.Println("\n10. Testing GET /search?q=test (without auth)")
	c.SetToken("") // Clear token
	resp, err = c.DoRequest("GET", "/search?q=test&limit=5", nil)
	if err != nil {
		return fmt.Errorf("search without auth failed: %w", err)
	}

	// Should still work, but might not include like status
	if err := client.CheckStatus(resp, http.StatusOK); err != nil {
		return fmt.Errorf("search without auth status check failed: %w", err)
	}

	if err := client.ParseResponse(resp, &searchResp); err != nil {
		return fmt.Errorf("search without auth parse failed: %w", err)
	}

	fmt.Printf("   ✓ Search without auth returned %d results (optional auth works)\n", len(searchResp.Items))

	// Restore token
	c.SetToken(data.Tokens.User1)

	fmt.Println("\n=== Search Endpoints Testing Complete ===")
	return nil
}
