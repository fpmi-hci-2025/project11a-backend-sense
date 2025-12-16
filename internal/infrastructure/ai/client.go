package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client represents AI service client
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a new AI service client
func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ComposeRequest represents request for compose endpoint
type ComposeRequest struct {
	Query    string                 `json:"query"`
	Metadata map[string]interface{} `json:"metadata"`
}

// ComposeResponse represents response from compose endpoint
type ComposeResponse struct {
	Text     string `json:"text,omitempty"`     // For backward compatibility
	Response string `json:"response,omitempty"` // Actual server response field
}

// RecommendRequest represents request for recommend endpoint
type RecommendRequest struct {
	ID string `json:"id"`
}

// RecommendResponse represents response from recommend endpoint
type RecommendResponse struct {
	Publications []string `json:"publications,omitempty"`
}

// HealthCheck checks if AI service is available
// Note: The AI server doesn't implement a health endpoint, so we try to connect to /compose
// with an empty request to verify the service is running
func (c *Client) HealthCheck(ctx context.Context) error {
	// Since /api/health doesn't exist, we'll try a simple request to verify connectivity
	// We could also just return nil if health check is not critical
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/compose", nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to connect to AI service: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	// Any response (even 400/405) means the service is reachable
	// We just need to verify it's not a connection error
	return nil
}

// Compose generates content using AI service
func (c *Client) Compose(ctx context.Context, req *ComposeRequest) (*ComposeResponse, error) {
	// Server expects {"text": "..."} format, not {"query": "...", "metadata": {...}}
	// Use query as the text field
	serverReq := map[string]interface{}{
		"text": req.Query,
	}
	
	body, err := json.Marshal(serverReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Server uses /compose, not /api/compose
	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/compose", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("compose failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var serverResp struct {
		Response string `json:"response"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&serverResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Convert server response format to our expected format
	return &ComposeResponse{
		Text:     serverResp.Response,
		Response: serverResp.Response,
	}, nil
}

// Recommend gets recommendations from AI service
// Note: The AI server doesn't implement /api/recommend endpoint
// This is a placeholder that returns an error indicating the endpoint is not available
func (c *Client) Recommend(ctx context.Context, userID string) (*RecommendResponse, error) {
	// The actual AI server doesn't have a recommend endpoint
	// We return an error indicating this feature is not implemented on the AI service
	return nil, fmt.Errorf("recommend endpoint is not implemented on the AI service. The server only provides /compose endpoint")
}

