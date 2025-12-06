package directory

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Client interface {
	GetUserByID(ctx context.Context, userID string) (*UserResponse, error)
	UserExists(ctx context.Context, userID string) (bool, error)
}

type HTTPClient struct {
	baseURL    string
	httpClient *http.Client
	timeout    time.Duration
}

func NewHTTPClient(cfg DirectoryServiceConfig) *HTTPClient {
	return &HTTPClient{
		baseURL: cfg.URL,
		httpClient: &http.Client{
			Timeout: cfg.Timeout,
			Transport: &http.Transport{
				MaxIdleConns:    cfg.MaxIdleConns,
				MaxConnsPerHost: cfg.MaxConnsPerHost,
				IdleConnTimeout: cfg.IdleConnTimeout,
			},
		},
		timeout: cfg.Timeout,
	}
}

func (c *HTTPClient) GetUserByID(ctx context.Context, userID string) (*UserResponse, error) {
	url := fmt.Sprintf("%s%s", strings.Split(c.baseURL, "*user_id*")[0], userID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request^ %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do reqest: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("firectory service returned: %w", err)
	}

	var userResp UserResponse
	if err := json.NewDecoder(resp.Body).Decode(&userResp); err != nil {
		return nil, fmt.Errorf("failed to decode data to UserResponse: %w", err)
	}

	return &userResp, nil
}

func (c *HTTPClient) UserExists(ctx context.Context, userID string) (bool, error) {
	_, err := c.GetUserByID(ctx, userID)
	if err != nil {
		if err == ErrUserNotFound {
			return false, nil
		}
		return false, fmt.Errorf("failed to get user by id: %w", err)
	}

	return true, nil
}
