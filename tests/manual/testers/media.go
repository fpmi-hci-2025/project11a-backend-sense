package testers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"sense-backend/tests/manual/client"
	"sense-backend/tests/manual/testdata"
)

// TestMediaEndpoints tests all media endpoints
func TestMediaEndpoints(c *client.Client) error {
	fmt.Println("\n=== Testing Media Endpoints ===")

	data := testdata.GetTestData()
	if data == nil {
		return fmt.Errorf("test data not loaded")
	}

	c.SetToken(data.Tokens.User1)

	// Create a test image file
	testImagePath := filepath.Join(os.TempDir(), "test_image.jpg")
	if err := createTestImage(testImagePath); err != nil {
		return fmt.Errorf("failed to create test image: %w", err)
	}
	defer os.Remove(testImagePath)

	// Test 1: Upload media file
	fmt.Println("\n1. Testing POST /media/upload")
	fields := map[string]string{
		"description": "Test image description",
	}

	resp, err := c.DoMultipartRequest("POST", "/media/upload", testImagePath, fields)
	if err != nil {
		return fmt.Errorf("upload media failed: %w", err)
	}

	if err := client.CheckStatus(resp, http.StatusCreated); err != nil {
		return fmt.Errorf("upload media status check failed: %w", err)
	}

	var mediaResp struct {
		ID        string `json:"id"`
		OwnerID   string `json:"owner_id"`
		Filename  string `json:"filename"`
		MIME      string `json:"mime"`
		Width     *int   `json:"width"`
		Height    *int   `json:"height"`
		CreatedAt string `json:"created_at"`
	}

	if err := client.ParseResponse(resp, &mediaResp); err != nil {
		return fmt.Errorf("upload media parse failed: %w", err)
	}

	testdata.SetMediaID(mediaResp.ID)
	fmt.Printf("   ✓ Media uploaded: ID=%s, MIME=%s", mediaResp.ID, mediaResp.MIME)
	if mediaResp.Width != nil && mediaResp.Height != nil {
		fmt.Printf(", Size=%dx%d", *mediaResp.Width, *mediaResp.Height)
	}
	fmt.Println()

	// Test 2: Get media info
	fmt.Printf("\n2. Testing GET /media/%s\n", data.IDs.MediaID)
	resp, err = c.DoRequest("GET", "/media/"+data.IDs.MediaID, nil)
	if err != nil {
		return fmt.Errorf("get media failed: %w", err)
	}

	if err := client.CheckStatus(resp, http.StatusOK); err != nil {
		return fmt.Errorf("get media status check failed: %w", err)
	}

	if err := client.ParseResponse(resp, &mediaResp); err != nil {
		return fmt.Errorf("get media parse failed: %w", err)
	}

	if mediaResp.ID != data.IDs.MediaID {
		return fmt.Errorf("media ID mismatch")
	}

	fmt.Printf("   ✓ Media info retrieved: Filename=%s, MIME=%s\n", mediaResp.Filename, mediaResp.MIME)

	// Test 3: Get media file (binary)
	fmt.Printf("\n3. Testing GET /media/%s/file\n", data.IDs.MediaID)
	resp, err = c.DoRequest("GET", "/media/"+data.IDs.MediaID+"/file", nil)
	if err != nil {
		return fmt.Errorf("get media file failed: %w", err)
	}

	if err := client.CheckStatus(resp, http.StatusOK); err != nil {
		return fmt.Errorf("get media file status check failed: %w", err)
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		return fmt.Errorf("Content-Type header missing")
	}

	contentDisposition := resp.Header.Get("Content-Disposition")
	if contentDisposition == "" {
		return fmt.Errorf("Content-Disposition header missing")
	}

	// Read binary data
	fileData, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return fmt.Errorf("failed to read file data: %w", err)
	}

	if len(fileData) == 0 {
		return fmt.Errorf("file data is empty")
	}

	fmt.Printf("   ✓ Media file retrieved: Size=%d bytes, Content-Type=%s\n", len(fileData), contentType)

	// Test 4: Get non-existent media
	fmt.Println("\n4. Testing GET /media/00000000-0000-0000-0000-000000000000")
	resp, err = c.DoRequest("GET", "/media/00000000-0000-0000-0000-000000000000", nil)
	if err != nil {
		return fmt.Errorf("get non-existent media request failed: %w", err)
	}

	if resp.StatusCode != http.StatusNotFound {
		return fmt.Errorf("expected 404 for non-existent media, got %d", resp.StatusCode)
	}

	fmt.Println("   ✓ Non-existent media correctly returns 404")

	// Test 5: Delete media (as owner)
	fmt.Printf("\n5. Testing DELETE /media/%s (as owner)\n", data.IDs.MediaID)
	resp, err = c.DoRequest("DELETE", "/media/"+data.IDs.MediaID, nil)
	if err != nil {
		return fmt.Errorf("delete media failed: %w", err)
	}

	if err := client.CheckStatus(resp, http.StatusNoContent); err != nil {
		return fmt.Errorf("delete media status check failed: %w", err)
	}

	fmt.Println("   ✓ Media deleted")

	// Test 6: Verify deletion
	fmt.Printf("\n6. Testing GET /media/%s (after deletion)\n", data.IDs.MediaID)
	resp, err = c.DoRequest("GET", "/media/"+data.IDs.MediaID, nil)
	if err != nil {
		return fmt.Errorf("get deleted media request failed: %w", err)
	}

	if resp.StatusCode != http.StatusNotFound {
		return fmt.Errorf("expected 404 for deleted media, got %d", resp.StatusCode)
	}

	fmt.Println("   ✓ Deleted media correctly returns 404")

	// Test 7: Upload media and test deletion by another user
	fmt.Println("\n7. Testing media ownership (upload with user1, delete with user2)")
	resp, err = c.DoMultipartRequest("POST", "/media/upload", testImagePath, nil)
	if err != nil {
		return fmt.Errorf("upload media for ownership test failed: %w", err)
	}

	if err := client.ParseResponse(resp, &mediaResp); err != nil {
		return fmt.Errorf("upload media for ownership test parse failed: %w", err)
	}

	// Try to delete with user2
	c.SetToken(data.Tokens.User2)
	resp, err = c.DoRequest("DELETE", "/media/"+mediaResp.ID, nil)
	if err != nil {
		return fmt.Errorf("delete media as non-owner request failed: %w", err)
	}

	if resp.StatusCode != http.StatusForbidden {
		return fmt.Errorf("expected 403 for deleting media as non-owner, got %d", resp.StatusCode)
	}

	fmt.Println("   ✓ Non-owner correctly cannot delete media (403)")

	// Clean up: delete with owner
	c.SetToken(data.Tokens.User1)
	resp, err = c.DoRequest("DELETE", "/media/"+mediaResp.ID, nil)
	if err == nil {
		resp.Body.Close()
	}

	fmt.Println("\n=== Media Endpoints Testing Complete ===")
	return nil
}

// createTestImage creates a minimal valid JPEG image for testing
func createTestImage(path string) error {
	// Minimal valid JPEG (1x1 pixel, red)
	// JPEG header + minimal image data
	jpegData := []byte{
		0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 0x4A, 0x46, 0x49, 0x46, 0x00, 0x01,
		0x01, 0x01, 0x00, 0x48, 0x00, 0x48, 0x00, 0x00, 0xFF, 0xDB, 0x00, 0x43,
		0x00, 0x08, 0x06, 0x06, 0x07, 0x06, 0x05, 0x08, 0x07, 0x07, 0x07, 0x09,
		0x09, 0x08, 0x0A, 0x0C, 0x14, 0x0D, 0x0C, 0x0B, 0x0B, 0x0C, 0x19, 0x12,
		0x13, 0x0F, 0x14, 0x1D, 0x1A, 0x1F, 0x1E, 0x1D, 0x1A, 0x1C, 0x1C, 0x20,
		0x24, 0x2E, 0x27, 0x20, 0x22, 0x2C, 0x23, 0x1C, 0x1C, 0x28, 0x37, 0x29,
		0x2C, 0x30, 0x31, 0x34, 0x34, 0x34, 0x1F, 0x27, 0x39, 0x3D, 0x38, 0x32,
		0x3C, 0x2E, 0x33, 0x34, 0x32, 0xFF, 0xC0, 0x00, 0x0B, 0x08, 0x00, 0x01,
		0x00, 0x01, 0x01, 0x01, 0x11, 0x00, 0xFF, 0xC4, 0x00, 0x14, 0x00, 0x01,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x08, 0xFF, 0xC4, 0x00, 0x14, 0x10, 0x01, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0xFF, 0xDA, 0x00, 0x08, 0x01, 0x01, 0x00, 0x00, 0x3F, 0x00,
		0xFF, 0xD9,
	}

	return os.WriteFile(path, jpegData, 0644)
}

