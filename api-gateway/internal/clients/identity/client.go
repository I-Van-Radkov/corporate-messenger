package identity

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
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

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed to do request: %w", err)
	}
	defer resp.Body.Close()

	// Читаем тело ответа для отладки
	bodyBytes, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		// Пытаемся распарсить ошибку из ответа
		var errorResp struct {
			Error string `json:"error"`
		}
		if err := json.Unmarshal(bodyBytes, &errorResp); err == nil && errorResp.Error != "" {
			return false, fmt.Errorf("identity service returned %d: %s", resp.StatusCode, errorResp.Error)
		}
		// Если не удалось распарсить JSON, возвращаем сырой ответ
		return false, fmt.Errorf("identity service returned %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var userResp IntrospectResponse
	if err := json.Unmarshal(bodyBytes, &userResp); err != nil {
		return false, fmt.Errorf("failed to decode data to IntrospectResponse: %w", err)
	}

	return userResp.Active, nil
}
