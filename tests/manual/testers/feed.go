package testers

import (
	"fmt"
	"net/http"

	"sense-backend/tests/manual/client"
	"sense-backend/tests/manual/testdata"
)

// TestFeedEndpoints tests all feed endpoints
func TestFeedEndpoints(c *client.Client) error {
	fmt.Println("\n=== Testing Feed Endpoints ===")

	data := testdata.GetTestData()
	if data == nil {
		return fmt.Errorf("test data not loaded")
	}

	c.SetToken(data.Tokens.User1)

	// Test 1: Get main feed
	fmt.Println("\n1. Testing GET /feed")
	resp, err := c.DoRequest("GET", "/feed?limit=10&offset=0", nil)
	if err != nil {
		return fmt.Errorf("get feed failed: %w", err)
	}

	if err := client.CheckStatus(resp, http.StatusOK); err != nil {
		return fmt.Errorf("get feed status check failed: %w", err)
	}

	var feedResp struct {
		Items  []interface{} `json:"items"`
		Total  int           `json:"total"`
		Limit  int           `json:"limit"`
		Offset int           `json:"offset"`
	}

	if err := client.ParseResponse(resp, &feedResp); err != nil {
		return fmt.Errorf("get feed parse failed: %w", err)
	}

	fmt.Printf("   ✓ Feed retrieved: Total=%d, Items=%d\n", feedResp.Total, len(feedResp.Items))

	// Test 2: Get feed with filters
	fmt.Println("\n2. Testing GET /feed (with type filter)")
	resp, err = c.DoRequest("GET", "/feed?type=post&limit=5&offset=0", nil)
	if err != nil {
		return fmt.Errorf("get filtered feed failed: %w", err)
	}

	if err := client.CheckStatus(resp, http.StatusOK); err != nil {
		return fmt.Errorf("get filtered feed status check failed: %w", err)
	}

	if err := client.ParseResponse(resp, &feedResp); err != nil {
		return fmt.Errorf("get filtered feed parse failed: %w", err)
	}

	fmt.Printf("   ✓ Filtered feed retrieved: Total=%d\n", feedResp.Total)

	// Test 3: Get my publications
	fmt.Println("\n3. Testing GET /feed/me")
	resp, err = c.DoRequest("GET", "/feed/me?limit=10&offset=0", nil)
	if err != nil {
		return fmt.Errorf("get my feed failed: %w", err)
	}

	if err := client.CheckStatus(resp, http.StatusOK); err != nil {
		return fmt.Errorf("get my feed status check failed: %w", err)
	}

	if err := client.ParseResponse(resp, &feedResp); err != nil {
		return fmt.Errorf("get my feed parse failed: %w", err)
	}

	fmt.Printf("   ✓ My publications retrieved: Total=%d\n", feedResp.Total)

	// Test 4: Get saved publications
	fmt.Println("\n4. Testing GET /feed/me/saved")
	resp, err = c.DoRequest("GET", "/feed/me/saved?limit=10&offset=0", nil)
	if err != nil {
		return fmt.Errorf("get saved feed failed: %w", err)
	}

	if err := client.CheckStatus(resp, http.StatusOK); err != nil {
		return fmt.Errorf("get saved feed status check failed: %w", err)
	}

	var savedFeedResp struct {
		Items  []interface{} `json:"items"`
		Total  int           `json:"total"`
		Limit  int           `json:"limit"`
		Offset int           `json:"offset"`
	}

	if err := client.ParseResponse(resp, &savedFeedResp); err != nil {
		return fmt.Errorf("get saved feed parse failed: %w", err)
	}

	fmt.Printf("   ✓ Saved publications retrieved: Total=%d\n", savedFeedResp.Total)

	// Test 5: Get user publications
	if data.IDs.User1ID != "" {
		fmt.Printf("\n5. Testing GET /feed/user/%s\n", data.IDs.User1ID)
		resp, err = c.DoRequest("GET", "/feed/user/"+data.IDs.User1ID+"?limit=10&offset=0", nil)
		if err != nil {
			return fmt.Errorf("get user feed failed: %w", err)
		}

		if err := client.CheckStatus(resp, http.StatusOK); err != nil {
			return fmt.Errorf("get user feed status check failed: %w", err)
		}

		if err := client.ParseResponse(resp, &feedResp); err != nil {
			return fmt.Errorf("get user feed parse failed: %w", err)
		}

		fmt.Printf("   ✓ User publications retrieved: Total=%d\n", feedResp.Total)
	}

	// Test 6: Get non-existent user publications
	fmt.Println("\n6. Testing GET /feed/user/00000000-0000-0000-0000-000000000000")
	resp, err = c.DoRequest("GET", "/feed/user/00000000-0000-0000-0000-000000000000?limit=10&offset=0", nil)
	if err != nil {
		return fmt.Errorf("get non-existent user feed request failed: %w", err)
	}

	// This might return 200 with empty items or 404, depending on implementation
	fmt.Printf("   ✓ Non-existent user feed: status %d\n", resp.StatusCode)

	// Test 7: Get feed with visibility filter
	fmt.Println("\n7. Testing GET /feed (with visibility filter)")
	resp, err = c.DoRequest("GET", "/feed?visibility=public&limit=5&offset=0", nil)
	if err != nil {
		return fmt.Errorf("get feed with visibility filter failed: %w", err)
	}

	if err := client.CheckStatus(resp, http.StatusOK); err != nil {
		return fmt.Errorf("get feed with visibility filter status check failed: %w", err)
	}

	fmt.Println("   ✓ Feed with visibility filter retrieved")

	// Test 8: Get feed with pagination
	fmt.Println("\n8. Testing GET /feed (pagination)")
	resp, err = c.DoRequest("GET", "/feed?limit=5&offset=5", nil)
	if err != nil {
		return fmt.Errorf("get feed with pagination failed: %w", err)
	}

	if err := client.CheckStatus(resp, http.StatusOK); err != nil {
		return fmt.Errorf("get feed with pagination status check failed: %w", err)
	}

	if err := client.ParseResponse(resp, &feedResp); err != nil {
		return fmt.Errorf("get feed with pagination parse failed: %w", err)
	}

	if feedResp.Offset != 5 {
		return fmt.Errorf("pagination offset mismatch: expected 5, got %d", feedResp.Offset)
	}

	fmt.Printf("   ✓ Feed pagination works: Offset=%d, Limit=%d\n", feedResp.Offset, feedResp.Limit)

	fmt.Println("\n=== Feed Endpoints Testing Complete ===")
	return nil
}

