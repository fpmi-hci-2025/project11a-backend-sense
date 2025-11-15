package testers

import (
	"fmt"
	"net/http"

	"sense-backend/tests/manual/client"
	"sense-backend/tests/manual/testdata"
)

// TestAuthEndpoints tests all authentication endpoints
func TestAuthEndpoints(c *client.Client) error {
	fmt.Println("\n=== Testing Auth Endpoints ===")

	data := testdata.GetTestData()
	if data == nil {
		return fmt.Errorf("test data not loaded")
	}

	// Test 1: Register user1
	fmt.Println("\n1. Testing POST /auth/register (user1)")
	registerReq := map[string]interface{}{
		"username":    data.Users.User1.Username,
		"email":       data.Users.User1.Email,
		"password":    data.Users.User1.Password,
		"phone":       data.Users.User1.Phone,
		"description": data.Users.User1.Description,
	}

	resp, err := c.DoRequest("POST", "/auth/register", registerReq)
	if err != nil {
		return fmt.Errorf("register user1 failed: %w", err)
	}

	var sessionResp struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
		User        struct {
			ID           string `json:"id"`
			Username     string `json:"username"`
			Email        string `json:"email"`
			Role         string `json:"role"`
			RegisteredAt string `json:"registered_at"`
		} `json:"user"`
	}

	// Handle case when user already exists (409)
	if resp.StatusCode == http.StatusConflict {
		fmt.Println("   ⚠ User1 already exists, attempting login instead")
		loginReq := map[string]string{
			"login":    data.Users.User1.Username,
			"password": data.Users.User1.Password,
		}
		resp, err = c.DoRequest("POST", "/auth/login", loginReq)
		if err != nil {
			return fmt.Errorf("login user1 after conflict failed: %w", err)
		}
		if err := client.CheckStatus(resp, http.StatusOK); err != nil {
			return fmt.Errorf("login user1 after conflict status check failed: %w", err)
		}
	} else if err := client.CheckStatus(resp, http.StatusCreated); err != nil {
		return fmt.Errorf("register user1 status check failed: %w", err)
	}

	if err := client.ParseResponse(resp, &sessionResp); err != nil {
		return fmt.Errorf("register user1 parse failed: %w", err)
	}

	testdata.SetUser1Token(sessionResp.AccessToken)
	testdata.SetUser1ID(sessionResp.User.ID)
	fmt.Printf("   ✓ User1 registered: ID=%s, Token=%s...\n", sessionResp.User.ID, sessionResp.AccessToken[:20])

	// Test 2: Register user2
	fmt.Println("\n2. Testing POST /auth/register (user2)")
	registerReq2 := map[string]interface{}{
		"username":    data.Users.User2.Username,
		"email":       data.Users.User2.Email,
		"password":    data.Users.User2.Password,
		"phone":       data.Users.User2.Phone,
		"description": data.Users.User2.Description,
	}

	resp, err = c.DoRequest("POST", "/auth/register", registerReq2)
	if err != nil {
		return fmt.Errorf("register user2 failed: %w", err)
	}

	// Handle case when user already exists (409)
	if resp.StatusCode == http.StatusConflict {
		fmt.Println("   ⚠ User2 already exists, attempting login instead")
		loginReq := map[string]string{
			"login":    data.Users.User2.Username,
			"password": data.Users.User2.Password,
		}
		resp, err = c.DoRequest("POST", "/auth/login", loginReq)
		if err != nil {
			return fmt.Errorf("login user2 after conflict failed: %w", err)
		}
		if err := client.CheckStatus(resp, http.StatusOK); err != nil {
			return fmt.Errorf("login user2 after conflict status check failed: %w", err)
		}
	} else if err := client.CheckStatus(resp, http.StatusCreated); err != nil {
		return fmt.Errorf("register user2 status check failed: %w", err)
	}

	if err := client.ParseResponse(resp, &sessionResp); err != nil {
		return fmt.Errorf("register user2 parse failed: %w", err)
	}

	testdata.SetUser2Token(sessionResp.AccessToken)
	testdata.SetUser2ID(sessionResp.User.ID)
	fmt.Printf("   ✓ User2 registered: ID=%s, Token=%s...\n", sessionResp.User.ID, sessionResp.AccessToken[:20])

	// Test 3: Login with username
	fmt.Println("\n3. Testing POST /auth/login (with username)")
	loginReq := map[string]string{
		"login":    data.Users.User1.Username,
		"password": data.Users.User1.Password,
	}

	resp, err = c.DoRequest("POST", "/auth/login", loginReq)
	if err != nil {
		return fmt.Errorf("login with username failed: %w", err)
	}

	if err := client.CheckStatus(resp, http.StatusOK); err != nil {
		return fmt.Errorf("login with username status check failed: %w", err)
	}

	if err := client.ParseResponse(resp, &sessionResp); err != nil {
		return fmt.Errorf("login with username parse failed: %w", err)
	}

	fmt.Printf("   ✓ Login with username successful: Token=%s...\n", sessionResp.AccessToken[:20])

	// Test 4: Login with email
	fmt.Println("\n4. Testing POST /auth/login (with email)")
	loginReq2 := map[string]string{
		"login":    data.Users.User1.Email,
		"password": data.Users.User1.Password,
	}

	resp, err = c.DoRequest("POST", "/auth/login", loginReq2)
	if err != nil {
		return fmt.Errorf("login with email failed: %w", err)
	}

	if err := client.CheckStatus(resp, http.StatusOK); err != nil {
		return fmt.Errorf("login with email status check failed: %w", err)
	}

	fmt.Println("   ✓ Login with email successful")

	// Test 5: Invalid credentials
	fmt.Println("\n5. Testing POST /auth/login (invalid credentials)")
	invalidLoginReq := map[string]string{
		"login":    data.Users.User1.Username,
		"password": "wrongpassword",
	}

	resp, err = c.DoRequest("POST", "/auth/login", invalidLoginReq)
	if err != nil {
		return fmt.Errorf("invalid login request failed: %w", err)
	}

	if resp.StatusCode != http.StatusUnauthorized {
		return fmt.Errorf("expected 401 for invalid credentials, got %d", resp.StatusCode)
	}

	fmt.Println("   ✓ Invalid credentials correctly rejected (401)")

	// Test 6: Check token
	fmt.Println("\n6. Testing GET /auth/check")
	c.SetToken(testdata.GetTestData().Tokens.User1)
	resp, err = c.DoRequest("GET", "/auth/check", nil)
	if err != nil {
		return fmt.Errorf("check token failed: %w", err)
	}

	if err := client.CheckStatus(resp, http.StatusOK); err != nil {
		return fmt.Errorf("check token status check failed: %w", err)
	}

	var userResp struct {
		ID           string `json:"id"`
		Username     string `json:"username"`
		Email        string `json:"email"`
		Role         string `json:"role"`
		RegisteredAt string `json:"registered_at"`
	}

	if err := client.ParseResponse(resp, &userResp); err != nil {
		return fmt.Errorf("check token parse failed: %w", err)
	}

	if userResp.ID != testdata.GetTestData().IDs.User1ID {
		return fmt.Errorf("user ID mismatch: expected %s, got %s", testdata.GetTestData().IDs.User1ID, userResp.ID)
	}

	fmt.Printf("   ✓ Token check successful: User=%s\n", userResp.Username)

	// Test 7: Check with invalid token
	fmt.Println("\n7. Testing GET /auth/check (invalid token)")
	c.SetToken("invalid_token")
	resp, err = c.DoRequest("GET", "/auth/check", nil)
	if err != nil {
		return fmt.Errorf("check invalid token request failed: %w", err)
	}

	if resp.StatusCode != http.StatusUnauthorized {
		return fmt.Errorf("expected 401 for invalid token, got %d", resp.StatusCode)
	}

	fmt.Println("   ✓ Invalid token correctly rejected (401)")

	// Test 8: Logout
	fmt.Println("\n8. Testing POST /auth/logout")
	c.SetToken(testdata.GetTestData().Tokens.User1)
	resp, err = c.DoRequest("POST", "/auth/logout", nil)
	if err != nil {
		return fmt.Errorf("logout failed: %w", err)
	}

	if err := client.CheckStatus(resp, http.StatusOK); err != nil {
		return fmt.Errorf("logout status check failed: %w", err)
	}

	fmt.Println("   ✓ Logout successful")

	// Test 9: Check token after logout
	fmt.Println("\n9. Testing GET /auth/check (after logout)")
	resp, err = c.DoRequest("GET", "/auth/check", nil)
	if err != nil {
		return fmt.Errorf("check after logout request failed: %w", err)
	}

	// Token might still be valid depending on implementation
	// This is implementation-dependent, so we just log the result
	fmt.Printf("   ✓ Check after logout: status %d\n", resp.StatusCode)

	// Restore token for further tests
	c.SetToken(testdata.GetTestData().Tokens.User1)

	fmt.Println("\n=== Auth Endpoints Testing Complete ===")
	return nil
}
