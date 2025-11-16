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
	Text string `json:"text,omitempty"`
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
func (c *Client) HealthCheck(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/api/health", nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("health check failed with status: %d", resp.StatusCode)
	}

	return nil
}

// Compose generates content using AI service
func (c *Client) Compose(ctx context.Context, req *ComposeRequest) (*ComposeResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/api/compose", bytes.NewBuffer(body))
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

	var composeResp ComposeResponse
	if err := json.NewDecoder(resp.Body).Decode(&composeResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &composeResp, nil
}

// Recommend gets recommendations from AI service
func (c *Client) Recommend(ctx context.Context, userID string) (*RecommendResponse, error) {
	req := &RecommendRequest{ID: userID}
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/api/recommend", bytes.NewBuffer(body))
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
		return nil, fmt.Errorf("recommend failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var recommendResp RecommendResponse
	if err := json.NewDecoder(resp.Body).Decode(&recommendResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &recommendResp, nil
}

