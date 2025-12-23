package testers

import (
	"fmt"
	"net/http"

	"sense-backend/tests/manual/client"
	"sense-backend/tests/manual/testdata"
)

// TestNotificationEndpoints tests all notification endpoints
func TestNotificationEndpoints(c *client.Client) error {
	fmt.Println("\n=== Testing Notification Endpoints ===")

	data := testdata.GetTestData()
	if data == nil {
		return fmt.Errorf("test data not loaded")
	}

	c.SetToken(data.Tokens.User1)

	// Test 1: Get all notifications
	fmt.Println("\n1. Testing GET /notifications")
	resp, err := c.DoRequest("GET", "/notifications?limit=20&offset=0", nil)
	if err != nil {
		return fmt.Errorf("get notifications failed: %w", err)
	}

	if err := client.CheckStatus(resp, http.StatusOK); err != nil {
		return fmt.Errorf("get notifications status check failed: %w", err)
	}

	var notifResp struct {
		Items []struct {
			ID        string                 `json:"id"`
			Type      string                 `json:"type"`
			Title     string                 `json:"title"`
			Message   string                 `json:"message"`
			IsRead    bool                   `json:"is_read"`
			CreatedAt string                 `json:"created_at"`
			Data      map[string]interface{} `json:"data"`
		} `json:"items"`
		Total  int `json:"total"`
		Limit  int `json:"limit"`
		Offset int `json:"offset"`
	}

	if err := client.ParseResponse(resp, &notifResp); err != nil {
		return fmt.Errorf("get notifications parse failed: %w", err)
	}

	fmt.Printf("   ✓ Retrieved %d notifications (total: %d)\n", len(notifResp.Items), notifResp.Total)

	// Display sample notifications if any exist
	if len(notifResp.Items) > 0 {
		fmt.Printf("   Sample notification: Type=%s, Title=%s, IsRead=%v\n",
			notifResp.Items[0].Type, notifResp.Items[0].Title, notifResp.Items[0].IsRead)
	}

	// Test 2: Get only unread notifications
	fmt.Println("\n2. Testing GET /notifications?unread_only=true")
	resp, err = c.DoRequest("GET", "/notifications?unread_only=true&limit=20", nil)
	if err != nil {
		return fmt.Errorf("get unread notifications failed: %w", err)
	}

	if err := client.CheckStatus(resp, http.StatusOK); err != nil {
		return fmt.Errorf("get unread notifications status check failed: %w", err)
	}

	var unreadResp struct {
		Items []struct {
			ID     string `json:"id"`
			Type   string `json:"type"`
			IsRead bool   `json:"is_read"`
		} `json:"items"`
		Total  int `json:"total"`
		Limit  int `json:"limit"`
		Offset int `json:"offset"`
	}

	if err := client.ParseResponse(resp, &unreadResp); err != nil {
		return fmt.Errorf("get unread notifications parse failed: %w", err)
	}

	fmt.Printf("   ✓ Retrieved %d unread notifications (total: %d)\n", len(unreadResp.Items), unreadResp.Total)

	// Verify all returned notifications are unread
	for _, notif := range unreadResp.Items {
		if notif.IsRead {
			return fmt.Errorf("unread_only filter returned read notification: %s", notif.ID)
		}
	}

	if len(unreadResp.Items) > 0 {
		fmt.Println("   ✓ All returned notifications are unread")
	}

	// Test 3: Pagination - get second page
	fmt.Println("\n3. Testing GET /notifications?limit=5&offset=5 (pagination)")
	resp, err = c.DoRequest("GET", "/notifications?limit=5&offset=5", nil)
	if err != nil {
		return fmt.Errorf("get notifications page 2 failed: %w", err)
	}

	if err := client.CheckStatus(resp, http.StatusOK); err != nil {
		return fmt.Errorf("get notifications page 2 status check failed: %w", err)
	}

	var pageResp struct {
		Items  []interface{} `json:"items"`
		Total  int           `json:"total"`
		Limit  int           `json:"limit"`
		Offset int           `json:"offset"`
	}

	if err := client.ParseResponse(resp, &pageResp); err != nil {
		return fmt.Errorf("get notifications page 2 parse failed: %w", err)
	}

	if pageResp.Limit != 5 {
		return fmt.Errorf("expected limit=5, got %d", pageResp.Limit)
	}

	if pageResp.Offset != 5 {
		return fmt.Errorf("expected offset=5, got %d", pageResp.Offset)
	}

	fmt.Printf("   ✓ Pagination works correctly: Retrieved %d items with offset=%d\n",
		len(pageResp.Items), pageResp.Offset)

	// Test 4: Get notifications with limit=1
	fmt.Println("\n4. Testing GET /notifications?limit=1 (single notification)")
	resp, err = c.DoRequest("GET", "/notifications?limit=1", nil)
	if err != nil {
		return fmt.Errorf("get single notification failed: %w", err)
	}

	if err := client.CheckStatus(resp, http.StatusOK); err != nil {
		return fmt.Errorf("get single notification status check failed: %w", err)
	}

	var singleResp struct {
		Items []interface{} `json:"items"`
		Total int           `json:"total"`
		Limit int           `json:"limit"`
	}

	if err := client.ParseResponse(resp, &singleResp); err != nil {
		return fmt.Errorf("get single notification parse failed: %w", err)
	}

	if len(singleResp.Items) > 1 {
		return fmt.Errorf("expected at most 1 item, got %d", len(singleResp.Items))
	}

	fmt.Printf("   ✓ Single notification limit works: Retrieved %d items\n", len(singleResp.Items))

	// Test 5: Get notifications with unread_only=false (explicit all)
	fmt.Println("\n5. Testing GET /notifications?unread_only=false (all notifications)")
	resp, err = c.DoRequest("GET", "/notifications?unread_only=false&limit=10", nil)
	if err != nil {
		return fmt.Errorf("get all notifications explicitly failed: %w", err)
	}

	if err := client.CheckStatus(resp, http.StatusOK); err != nil {
		return fmt.Errorf("get all notifications explicitly status check failed: %w", err)
	}

	if err := client.ParseResponse(resp, &notifResp); err != nil {
		return fmt.Errorf("get all notifications explicitly parse failed: %w", err)
	}

	fmt.Printf("   ✓ All notifications (unread_only=false): Retrieved %d items (total: %d)\n",
		len(notifResp.Items), notifResp.Total)

	// Test 6: Get notifications without authentication (should fail)
	fmt.Println("\n6. Testing GET /notifications (without auth)")
	c.SetToken("") // Clear token

	resp, err = c.DoRequest("GET", "/notifications", nil)
	if err != nil {
		return fmt.Errorf("get notifications without auth request failed: %w", err)
	}

	if resp.StatusCode != http.StatusUnauthorized {
		return fmt.Errorf("expected 401 for unauthenticated request, got %d", resp.StatusCode)
	}

	fmt.Println("   ✓ Notifications correctly require authentication (401)")

	// Test 7: Get notifications as different user
	fmt.Println("\n7. Testing GET /notifications (as User2)")
	c.SetToken(data.Tokens.User2)

	resp, err = c.DoRequest("GET", "/notifications?limit=10", nil)
	if err != nil {
		return fmt.Errorf("get notifications as user2 failed: %w", err)
	}

	if err := client.CheckStatus(resp, http.StatusOK); err != nil {
		return fmt.Errorf("get notifications as user2 status check failed: %w", err)
	}

	var user2NotifResp struct {
		Items  []interface{} `json:"items"`
		Total  int           `json:"total"`
		Limit  int           `json:"limit"`
		Offset int           `json:"offset"`
	}

	if err := client.ParseResponse(resp, &user2NotifResp); err != nil {
		return fmt.Errorf("get notifications as user2 parse failed: %w", err)
	}

	fmt.Printf("   ✓ User2 has %d notifications (total: %d)\n",
		len(user2NotifResp.Items), user2NotifResp.Total)

	// Test 8: Large limit test (should cap at max)
	fmt.Println("\n8. Testing GET /notifications?limit=1000 (large limit)")
	resp, err = c.DoRequest("GET", "/notifications?limit=1000", nil)
	if err != nil {
		return fmt.Errorf("get notifications with large limit failed: %w", err)
	}

	if err := client.CheckStatus(resp, http.StatusOK); err != nil {
		return fmt.Errorf("get notifications with large limit status check failed: %w", err)
	}

	var largeResp struct {
		Items []interface{} `json:"items"`
		Total int           `json:"total"`
		Limit int           `json:"limit"`
	}

	if err := client.ParseResponse(resp, &largeResp); err != nil {
		return fmt.Errorf("get notifications with large limit parse failed: %w", err)
	}

	// Limit should be capped (typically at 100)
	fmt.Printf("   ✓ Large limit handled: Returned limit=%d, items=%d\n",
		largeResp.Limit, len(largeResp.Items))

	// Test 9: Test with both unread_only and pagination
	fmt.Println("\n9. Testing GET /notifications?unread_only=true&limit=3&offset=0")
	resp, err = c.DoRequest("GET", "/notifications?unread_only=true&limit=3&offset=0", nil)
	if err != nil {
		return fmt.Errorf("get filtered paginated notifications failed: %w", err)
	}

	if err := client.CheckStatus(resp, http.StatusOK); err != nil {
		return fmt.Errorf("get filtered paginated notifications status check failed: %w", err)
	}

	if err := client.ParseResponse(resp, &unreadResp); err != nil {
		return fmt.Errorf("get filtered paginated notifications parse failed: %w", err)
	}

	fmt.Printf("   ✓ Filtered pagination works: Retrieved %d unread notifications\n",
		len(unreadResp.Items))

	// Restore User1 token
	c.SetToken(data.Tokens.User1)

	fmt.Println("\n=== Notification Endpoints Testing Complete ===")
	fmt.Println("\nNote: Actual notification creation is triggered by other actions")
	fmt.Println("      (likes, comments, follows, etc.) and not tested here directly.")
	return nil
}
