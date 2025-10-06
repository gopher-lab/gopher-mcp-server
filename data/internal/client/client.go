package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// GopherAIClient represents the client for interacting with Gopher AI API
type GopherAIClient struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
}

// NewGopherAIClient creates a new Gopher AI client instance
func NewGopherAIClient(apiKey string) *GopherAIClient {
	return &GopherAIClient{
		BaseURL: "https://data.gopher-ai.com/api/v1",
		APIKey:  apiKey,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// NewGopherAIClientWithURL creates a new Gopher AI client with custom base URL
func NewGopherAIClientWithURL(baseURL, apiKey string) *GopherAIClient {
	return &GopherAIClient{
		BaseURL: baseURL,
		APIKey:  apiKey,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SearchTwitter performs a Twitter search using the Gopher AI API
func (c *GopherAIClient) SearchTwitter(query string, maxResults int) (*TwitterSearchResponse, error) {
	request := TwitterSearchRequest{
		Type: "twitter",
		Arguments: TwitterSearchArguments{
			Type:       "searchbyquery",
			Query:      query,
			MaxResults: maxResults,
		},
	}

	// First, initiate the search and get UUID
	initResp, err := c.initiateSearch("/search/live/twitter", request)
	if err != nil {
		return nil, err
	}

	if initResp.Error != "" {
		return nil, fmt.Errorf("API error: %s", initResp.Error)
	}

	// Poll for results
	return c.pollForResults(initResp.UUID)
}

// initiateSearch starts a search and returns the UUID
func (c *GopherAIClient) initiateSearch(endpoint string, requestBody interface{}) (*TwitterSearchInitResponse, error) {
	// Marshal request body to JSON
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	// Create HTTP request
	url := c.BaseURL + endpoint
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.APIKey)

	// Make the request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for HTTP errors
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var apiErr APIError
		if err := json.Unmarshal(body, &apiErr); err != nil {
			return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
		}
		if apiErr.Message == "" {
			return nil, fmt.Errorf("API error (HTTP %d): %s", resp.StatusCode, string(body))
		}
		return nil, fmt.Errorf("API error: %s", apiErr.Message)
	}

	// Parse response
	var response TwitterSearchInitResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

// pollForResults polls for search results using the UUID
func (c *GopherAIClient) pollForResults(uuid string) (*TwitterSearchResponse, error) {
	maxAttempts := 30               // Maximum number of polling attempts
	pollInterval := 2 * time.Second // Poll every 2 seconds

	for attempt := 0; attempt < maxAttempts; attempt++ {
		// Create HTTP request for results
		url := c.BaseURL + "/search/live/twitter/result/" + uuid
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		// Set headers
		req.Header.Set("Authorization", "Bearer "+c.APIKey)

		// Make the request
		resp, err := c.HTTPClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to make request: %w", err)
		}

		// Read response body
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}

		// Check for HTTP errors
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			if resp.StatusCode == 404 {
				// Still processing, continue polling
				time.Sleep(pollInterval)
				continue
			}
			var apiErr APIError
			if err := json.Unmarshal(body, &apiErr); err != nil {
				return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
			}
			if apiErr.Message == "" {
				return nil, fmt.Errorf("API error (HTTP %d): %s", resp.StatusCode, string(body))
			}
			return nil, fmt.Errorf("API error: %s", apiErr.Message)
		}

		// Try to parse as array of tweets first
		var tweets []TwitterPost
		if err := json.Unmarshal(body, &tweets); err == nil {
			// Successfully parsed as tweets array
			return &TwitterSearchResponse{
				Status: "success",
				Data:   tweets,
			}, nil
		}

		// If not an array, check if it's still processing
		var processingResp map[string]interface{}
		if err := json.Unmarshal(body, &processingResp); err == nil {
			if status, ok := processingResp["status"].(string); ok && status == "processing" {
				// Still processing, continue polling
				time.Sleep(pollInterval)
				continue
			}
		}

		// If we get here, something unexpected happened
		time.Sleep(pollInterval)
	}

	return nil, fmt.Errorf("timeout waiting for search results after %d attempts", maxAttempts)
}

// SetTimeout sets a custom timeout for HTTP requests
func (c *GopherAIClient) SetTimeout(timeout time.Duration) {
	c.HTTPClient.Timeout = timeout
}
