package testers

import (
	"fmt"
	"net/http"

	"sense-backend/tests/manual/client"
	"sense-backend/tests/manual/testdata"
)

// TestProfileEndpoints tests all profile endpoints
func TestProfileEndpoints(c *client.Client) error {
	fmt.Println("\n=== Testing Profile Endpoints ===")

	data := testdata.GetTestData()
	if data == nil {
		return fmt.Errorf("test data not loaded")
	}

	c.SetToken(data.Tokens.User1)

	// Test 1: Get my profile
	fmt.Println("\n1. Testing GET /profile/me")
	resp, err := c.DoRequest("GET", "/profile/me", nil)
	if err != nil {
		return fmt.Errorf("get my profile failed: %w", err)
	}

	if err := client.CheckStatus(resp, http.StatusOK); err != nil {
		return fmt.Errorf("get my profile status check failed: %w", err)
	}

	var userResp struct {
		ID           string `json:"id"`
		Username     string `json:"username"`
		Email        string `json:"email"`
		Phone        string `json:"phone"`
		Description  string `json:"description"`
		Role         string `json:"role"`
		RegisteredAt string `json:"registered_at"`
	}

	if err := client.ParseResponse(resp, &userResp); err != nil {
		return fmt.Errorf("get my profile parse failed: %w", err)
	}

	if userResp.ID != data.IDs.User1ID {
		return fmt.Errorf("user ID mismatch")
	}

	fmt.Printf("   ✓ My profile retrieved: Username=%s, Email=%s\n", userResp.Username, userResp.Email)

	// Test 2: Update my profile
	fmt.Println("\n2. Testing POST /profile/me")
	updateReq := map[string]string{
		"description": "Updated profile description",
		"icon_url":    "https://cdn.sense.social/avatars/user123.jpg",
	}

	resp, err = c.DoRequest("POST", "/profile/me", updateReq)
	if err != nil {
		return fmt.Errorf("update my profile failed: %w", err)
	}

	if err := client.CheckStatus(resp, http.StatusOK); err != nil {
		return fmt.Errorf("update my profile status check failed: %w", err)
	}

	var updatedUserResp struct {
		ID          string `json:"id"`
		Description string `json:"description"`
		IconURL     string `json:"icon_url"`
	}

	if err := client.ParseResponse(resp, &updatedUserResp); err != nil {
		return fmt.Errorf("update my profile parse failed: %w", err)
	}

	if updatedUserResp.Description != "Updated profile description" {
		return fmt.Errorf("description not updated correctly")
	}

	fmt.Printf("   ✓ Profile updated: Description=%s\n", updatedUserResp.Description)

	// Test 3: Get user profile
	if data.IDs.User1ID != "" {
		fmt.Printf("\n3. Testing GET /profile/%s\n", data.IDs.User1ID)
		resp, err = c.DoRequest("GET", "/profile/"+data.IDs.User1ID, nil)
		if err != nil {
			return fmt.Errorf("get user profile failed: %w", err)
		}

		if err := client.CheckStatus(resp, http.StatusOK); err != nil {
			return fmt.Errorf("get user profile status check failed: %w", err)
		}

		var profileResp struct {
			ID          string `json:"id"`
			Username    string `json:"username"`
			Description string `json:"description"`
			IsFollowing bool   `json:"is_following"`
		}

		if err := client.ParseResponse(resp, &profileResp); err != nil {
			return fmt.Errorf("get user profile parse failed: %w", err)
		}

		if profileResp.ID != data.IDs.User1ID {
			return fmt.Errorf("profile ID mismatch")
		}

		fmt.Printf("   ✓ User profile retrieved: Username=%s, IsFollowing=%v\n", profileResp.Username, profileResp.IsFollowing)
	}

	// Test 4: Get non-existent user profile
	fmt.Println("\n4. Testing GET /profile/00000000-0000-0000-0000-000000000000")
	resp, err = c.DoRequest("GET", "/profile/00000000-0000-0000-0000-000000000000", nil)
	if err != nil {
		return fmt.Errorf("get non-existent user profile request failed: %w", err)
	}

	if resp.StatusCode != http.StatusNotFound {
		return fmt.Errorf("expected 404 for non-existent user, got %d", resp.StatusCode)
	}

	fmt.Println("   ✓ Non-existent user correctly returns 404")

	// Test 5: Get user stats
	if data.IDs.User1ID != "" {
		fmt.Printf("\n5. Testing GET /profile/%s/stats\n", data.IDs.User1ID)
		resp, err = c.DoRequest("GET", "/profile/"+data.IDs.User1ID+"/stats", nil)
		if err != nil {
			return fmt.Errorf("get user stats failed: %w", err)
		}

		if err := client.CheckStatus(resp, http.StatusOK); err != nil {
			return fmt.Errorf("get user stats status check failed: %w", err)
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
			return fmt.Errorf("get user stats parse failed: %w", err)
		}

		fmt.Printf("   ✓ User stats retrieved: Publications=%d, Followers=%d, Following=%d\n",
			statsResp.PublicationsCount, statsResp.FollowersCount, statsResp.FollowingCount)
	}

	// Test 6: Get stats for non-existent user
	fmt.Println("\n6. Testing GET /profile/00000000-0000-0000-0000-000000000000/stats")
	resp, err = c.DoRequest("GET", "/profile/00000000-0000-0000-0000-000000000000/stats", nil)
	if err != nil {
		return fmt.Errorf("get non-existent user stats request failed: %w", err)
	}

	if resp.StatusCode != http.StatusNotFound {
		return fmt.Errorf("expected 404 for non-existent user stats, got %d", resp.StatusCode)
	}

	fmt.Println("   ✓ Non-existent user stats correctly returns 404")

	fmt.Println("\n=== Profile Endpoints Testing Complete ===")
	return nil
}

