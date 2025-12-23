package testers

import (
	"fmt"
	"net/http"

	"sense-backend/tests/manual/client"
	"sense-backend/tests/manual/testdata"
)

// TestSocialEndpoints tests all social/follow endpoints
func TestSocialEndpoints(c *client.Client) error {
	fmt.Println("\n=== Testing Social/Follow Endpoints ===")

	data := testdata.GetTestData()
	if data == nil {
		return fmt.Errorf("test data not loaded")
	}

	// Ensure we have two different users
	if data.IDs.User1ID == "" || data.IDs.User2ID == "" {
		return fmt.Errorf("need two users for follow tests")
	}

	if data.IDs.User1ID == data.IDs.User2ID {
		return fmt.Errorf("user1 and user2 must be different")
	}

	// Test 1: User1 follows User2
	fmt.Printf("\n1. Testing POST /follow/%s (User1 follows User2)\n", data.IDs.User2ID)
	c.SetToken(data.Tokens.User1)

	resp, err := c.DoRequest("POST", "/follow/"+data.IDs.User2ID, nil)
	if err != nil {
		return fmt.Errorf("follow user failed: %w", err)
	}

	if err := client.CheckStatus(resp, http.StatusCreated); err != nil {
		return fmt.Errorf("follow user status check failed: %w", err)
	}

	var followResp struct {
		Message string `json:"message"`
	}

	if err := client.ParseResponse(resp, &followResp); err != nil {
		return fmt.Errorf("follow user parse failed: %w", err)
	}

	fmt.Printf("   ✓ User1 followed User2: %s\n", followResp.Message)

	// Test 2: Verify follow relationship in profile
	fmt.Printf("\n2. Testing GET /profile/%s (checking is_following)\n", data.IDs.User2ID)
	resp, err = c.DoRequest("GET", "/profile/"+data.IDs.User2ID, nil)
	if err != nil {
		return fmt.Errorf("get profile after follow failed: %w", err)
	}

	if err := client.CheckStatus(resp, http.StatusOK); err != nil {
		return fmt.Errorf("get profile after follow status check failed: %w", err)
	}

	var profileResp struct {
		ID          string `json:"id"`
		Username    string `json:"username"`
		IsFollowing bool   `json:"is_following"`
	}

	if err := client.ParseResponse(resp, &profileResp); err != nil {
		return fmt.Errorf("get profile after follow parse failed: %w", err)
	}

	if !profileResp.IsFollowing {
		return fmt.Errorf("is_following should be true after following")
	}

	fmt.Printf("   ✓ Profile shows is_following=true for User2\n")

	// Test 3: Verify follower count in stats
	fmt.Printf("\n3. Testing GET /profile/%s/stats (checking followers_count)\n", data.IDs.User2ID)
	resp, err = c.DoRequest("GET", "/profile/"+data.IDs.User2ID+"/stats", nil)
	if err != nil {
		return fmt.Errorf("get stats after follow failed: %w", err)
	}

	if err := client.CheckStatus(resp, http.StatusOK); err != nil {
		return fmt.Errorf("get stats after follow status check failed: %w", err)
	}

	var statsResp struct {
		PublicationsCount int `json:"publications_count"`
		FollowersCount    int `json:"followers_count"`
		FollowingCount    int `json:"following_count"`
		LikesReceived     int `json:"likes_received"`
		CommentsReceived  int `json:"comments_received"`
		SavedCount        int `json:"saved_count"`
	}

	if err := client.ParseResponse(resp, &statsResp); err != nil {
		return fmt.Errorf("get stats after follow parse failed: %w", err)
	}

	if statsResp.FollowersCount < 1 {
		return fmt.Errorf("followers_count should be at least 1 after follow")
	}

	fmt.Printf("   ✓ User2 has %d followers\n", statsResp.FollowersCount)

	// Test 4: Try to follow the same user again (should succeed - idempotent)
	fmt.Printf("\n4. Testing POST /follow/%s (duplicate follow)\n", data.IDs.User2ID)
	resp, err = c.DoRequest("POST", "/follow/"+data.IDs.User2ID, nil)
	if err != nil {
		return fmt.Errorf("duplicate follow request failed: %w", err)
	}

	// Should succeed (idempotent operation)
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("expected 201 or 200 for duplicate follow, got %d", resp.StatusCode)
	}

	fmt.Println("   ✓ Duplicate follow handled correctly")

	// Test 5: Try to follow yourself (should fail)
	fmt.Printf("\n5. Testing POST /follow/%s (self-follow)\n", data.IDs.User1ID)
	resp, err = c.DoRequest("POST", "/follow/"+data.IDs.User1ID, nil)
	if err != nil {
		return fmt.Errorf("self-follow request failed: %w", err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		return fmt.Errorf("expected 400 for self-follow, got %d", resp.StatusCode)
	}

	fmt.Println("   ✓ Self-follow correctly rejected with 400")

	// Test 6: Try to follow non-existent user (should fail)
	fmt.Println("\n6. Testing POST /follow/00000000-0000-0000-0000-000000000000")
	resp, err = c.DoRequest("POST", "/follow/00000000-0000-0000-0000-000000000000", nil)
	if err != nil {
		return fmt.Errorf("follow non-existent user request failed: %w", err)
	}

	if resp.StatusCode != http.StatusNotFound {
		return fmt.Errorf("expected 404 for non-existent user, got %d", resp.StatusCode)
	}

	fmt.Println("   ✓ Following non-existent user correctly returns 404")

	// Test 7: Unfollow user
	fmt.Printf("\n7. Testing DELETE /follow/%s (User1 unfollows User2)\n", data.IDs.User2ID)
	resp, err = c.DoRequest("DELETE", "/follow/"+data.IDs.User2ID, nil)
	if err != nil {
		return fmt.Errorf("unfollow user failed: %w", err)
	}

	if err := client.CheckStatus(resp, http.StatusNoContent); err != nil {
		return fmt.Errorf("unfollow user status check failed: %w", err)
	}

	fmt.Println("   ✓ User1 unfollowed User2 (204 No Content)")

	// Test 8: Verify unfollow in profile
	fmt.Printf("\n8. Testing GET /profile/%s (checking is_following after unfollow)\n", data.IDs.User2ID)
	resp, err = c.DoRequest("GET", "/profile/"+data.IDs.User2ID, nil)
	if err != nil {
		return fmt.Errorf("get profile after unfollow failed: %w", err)
	}

	if err := client.CheckStatus(resp, http.StatusOK); err != nil {
		return fmt.Errorf("get profile after unfollow status check failed: %w", err)
	}

	if err := client.ParseResponse(resp, &profileResp); err != nil {
		return fmt.Errorf("get profile after unfollow parse failed: %w", err)
	}

	if profileResp.IsFollowing {
		return fmt.Errorf("is_following should be false after unfollowing")
	}

	fmt.Printf("   ✓ Profile shows is_following=false after unfollow\n")

	// Test 9: Unfollow user that is not followed (should succeed - idempotent)
	fmt.Printf("\n9. Testing DELETE /follow/%s (duplicate unfollow)\n", data.IDs.User2ID)
	resp, err = c.DoRequest("DELETE", "/follow/"+data.IDs.User2ID, nil)
	if err != nil {
		return fmt.Errorf("duplicate unfollow request failed: %w", err)
	}

	// Should succeed (idempotent operation) - 204 or 404 acceptable
	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusNotFound {
		return fmt.Errorf("expected 204 or 404 for duplicate unfollow, got %d", resp.StatusCode)
	}

	fmt.Println("   ✓ Duplicate unfollow handled correctly")

	// Test 10: Unfollow non-existent user (should fail)
	fmt.Println("\n10. Testing DELETE /follow/00000000-0000-0000-0000-000000000000")
	resp, err = c.DoRequest("DELETE", "/follow/00000000-0000-0000-0000-000000000000", nil)
	if err != nil {
		return fmt.Errorf("unfollow non-existent user request failed: %w", err)
	}

	if resp.StatusCode != http.StatusNotFound {
		return fmt.Errorf("expected 404 for unfollowing non-existent user, got %d", resp.StatusCode)
	}

	fmt.Println("   ✓ Unfollowing non-existent user correctly returns 404")

	// Test 11: Cross-test - User2 follows User1
	fmt.Printf("\n11. Testing POST /follow/%s (User2 follows User1)\n", data.IDs.User1ID)
	c.SetToken(data.Tokens.User2)

	resp, err = c.DoRequest("POST", "/follow/"+data.IDs.User1ID, nil)
	if err != nil {
		return fmt.Errorf("user2 follow user1 failed: %w", err)
	}

	if err := client.CheckStatus(resp, http.StatusCreated); err != nil {
		return fmt.Errorf("user2 follow user1 status check failed: %w", err)
	}

	fmt.Println("   ✓ User2 followed User1")

	// Test 12: Verify mutual follow stats
	fmt.Printf("\n12. Testing GET /profile/%s/stats (User1 stats)\n", data.IDs.User1ID)
	resp, err = c.DoRequest("GET", "/profile/"+data.IDs.User1ID+"/stats", nil)
	if err != nil {
		return fmt.Errorf("get user1 stats failed: %w", err)
	}

	if err := client.CheckStatus(resp, http.StatusOK); err != nil {
		return fmt.Errorf("get user1 stats status check failed: %w", err)
	}

	if err := client.ParseResponse(resp, &statsResp); err != nil {
		return fmt.Errorf("get user1 stats parse failed: %w", err)
	}

	fmt.Printf("   ✓ User1 stats: Followers=%d, Following=%d\n",
		statsResp.FollowersCount, statsResp.FollowingCount)

	// Clean up - User2 unfollows User1
	fmt.Printf("\n13. Cleanup: DELETE /follow/%s (User2 unfollows User1)\n", data.IDs.User1ID)
	resp, err = c.DoRequest("DELETE", "/follow/"+data.IDs.User1ID, nil)
	if err != nil {
		return fmt.Errorf("cleanup unfollow failed: %w", err)
	}

	if err := client.CheckStatus(resp, http.StatusNoContent); err != nil {
		return fmt.Errorf("cleanup unfollow status check failed: %w", err)
	}

	fmt.Println("   ✓ Cleanup complete")

	// Restore User1 token
	c.SetToken(data.Tokens.User1)

	fmt.Println("\n=== Social/Follow Endpoints Testing Complete ===")
	return nil
}
