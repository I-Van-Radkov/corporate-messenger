package identity

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client interface {
	IntrospectToken(ctx context.Context, token string) (bool, error)
}

type HTTPClient struct {
	baseURL    string
	httpClient *http.Client
	timeout    time.Duration
}

func NewHTTPClient(cfg IdentityServiceConfig) *HTTPClient {
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

func (c *HTTPClient) IntrospectToken(ctx context.Context, token string) (bool, error) {
	introspectRequest := IntrospectRequest{
		Token: token,
	}
	jsonBody, _ := json.Marshal(introspectRequest)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return false, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed to do reqest: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("firectory service returned: %w", err)
	}

	var userResp IntrospectResponse
	if err := json.NewDecoder(resp.Body).Decode(&userResp); err != nil {
		return false, fmt.Errorf("failed to decode data to UserResponse: %w", err)
	}

	return userResp.Active, nil
}
