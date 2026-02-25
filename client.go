package scrappey

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

// Client is the main Scrappey API client.
type Client struct {
	apiKey     string
	baseURL    string
	timeout    time.Duration
	httpClient *http.Client
}

// NewClient creates a new Scrappey API client.
func NewClient(apiKey string, cfg *Config) (*Client, error) {
	if strings.TrimSpace(apiKey) == "" {
		return nil, &AuthenticationError{
			APIError: &APIError{
				Message: "API key is required",
			},
		}
	}

	baseURL := DefaultBaseURL
	timeout := DefaultTimeout
	var httpClient *http.Client

	if cfg != nil {
		if strings.TrimSpace(cfg.BaseURL) != "" {
			baseURL = strings.TrimRight(cfg.BaseURL, "/")
		}
		if cfg.Timeout > 0 {
			timeout = cfg.Timeout
		}
		httpClient = cfg.HTTPClient
	}

	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: timeout,
		}
	} else if httpClient.Timeout == 0 {
		httpClient.Timeout = timeout
	}

	return &Client{
		apiKey:     apiKey,
		baseURL:    baseURL,
		timeout:    timeout,
		httpClient: httpClient,
	}, nil
}

// APIKey returns the configured API key.
func (c *Client) APIKey() string {
	return c.apiKey
}

// BaseURL returns the configured API base URL.
func (c *Client) BaseURL() string {
	return c.baseURL
}

// Timeout returns the default request timeout.
func (c *Client) Timeout() time.Duration {
	return c.timeout
}

// Request sends a raw Scrappey command payload.
// The payload must include a non-empty "cmd" string key.
func (c *Client) Request(ctx context.Context, payload map[string]any) (*APIResponse, error) {
	if payload == nil {
		return nil, &APIError{Message: "request payload is required"}
	}

	cmd, ok := payload["cmd"].(string)
	if !ok || strings.TrimSpace(cmd) == "" {
		return nil, &APIError{Message: "command (cmd) is required"}
	}

	if ctx == nil {
		ctx = context.Background()
	}

	endpoint, err := c.requestURL()
	if err != nil {
		return nil, &APIError{
			Message: "failed to build request URL",
			Cause:   err,
		}
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, &APIError{
			Message: "failed to encode request payload",
			Cause:   err,
		}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, &APIError{
			Message: "failed to create HTTP request",
			Cause:   err,
		}
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, classifyTransportError(err, c.apiKey)
	}
	defer res.Body.Close()

	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, &APIError{
			Message:    "failed to read API response",
			StatusCode: res.StatusCode,
			Cause:      err,
		}
	}

	if res.StatusCode == http.StatusUnauthorized {
		return nil, &AuthenticationError{
			APIError: &APIError{
				Message:    "invalid API key",
				StatusCode: res.StatusCode,
				Body:       string(responseBody),
			},
		}
	}

	var parsed APIResponse
	if err := json.Unmarshal(responseBody, &parsed); err != nil {
		return nil, &APIError{
			Message:    "failed to parse API response",
			StatusCode: res.StatusCode,
			Body:       string(responseBody),
			Cause:      err,
		}
	}
	parsed.HTTPStatus = res.StatusCode

	var raw map[string]any
	if err := json.Unmarshal(responseBody, &raw); err == nil {
		parsed.Raw = raw
	}

	return &parsed, nil
}

// Get sends "request.get".
func (c *Client) Get(ctx context.Context, options RequestOptions) (*APIResponse, error) {
	return c.requestWithCommand(ctx, "request.get", options)
}

// Post sends "request.post".
func (c *Client) Post(ctx context.Context, options RequestOptions) (*APIResponse, error) {
	return c.requestWithCommand(ctx, "request.post", options)
}

// Put sends "request.put".
func (c *Client) Put(ctx context.Context, options RequestOptions) (*APIResponse, error) {
	return c.requestWithCommand(ctx, "request.put", options)
}

// Delete sends "request.delete".
func (c *Client) Delete(ctx context.Context, options RequestOptions) (*APIResponse, error) {
	return c.requestWithCommand(ctx, "request.delete", options)
}

// Patch sends "request.patch".
func (c *Client) Patch(ctx context.Context, options RequestOptions) (*APIResponse, error) {
	return c.requestWithCommand(ctx, "request.patch", options)
}

// CreateSession sends "sessions.create".
func (c *Client) CreateSession(ctx context.Context, options SessionOptions) (*APIResponse, error) {
	return c.requestWithCommand(ctx, "sessions.create", options)
}

// DestroySession sends "sessions.destroy".
func (c *Client) DestroySession(ctx context.Context, session string) (*APIResponse, error) {
	if strings.TrimSpace(session) == "" {
		return nil, &APIError{Message: "session is required"}
	}
	return c.requestWithCommand(ctx, "sessions.destroy", map[string]any{
		"session": session,
	})
}

// CreateWebSocket sends "websocket.create".
func (c *Client) CreateWebSocket(ctx context.Context, options WebSocketOptions) (*APIResponse, error) {
	return c.requestWithCommand(ctx, "websocket.create", options)
}

// CloseIdleConnections closes pooled idle HTTP connections.
func (c *Client) CloseIdleConnections() {
	if c == nil || c.httpClient == nil {
		return
	}
	c.httpClient.CloseIdleConnections()
}

func (c *Client) requestWithCommand(ctx context.Context, cmd string, options map[string]any) (*APIResponse, error) {
	payload := cloneMap(options)
	payload["cmd"] = cmd
	return c.Request(ctx, payload)
}

func (c *Client) requestURL() (string, error) {
	parsed, err := url.Parse(c.baseURL)
	if err != nil {
		return "", err
	}

	query := parsed.Query()
	query.Set("key", c.apiKey)
	parsed.RawQuery = query.Encode()
	return parsed.String(), nil
}

func cloneMap(input map[string]any) map[string]any {
	if input == nil {
		return map[string]any{}
	}
	output := make(map[string]any, len(input))
	for key, value := range input {
		output[key] = value
	}
	return output
}

var apiKeyQueryPattern = regexp.MustCompile(`([?&]key=)[^&\s]+`)

func classifyTransportError(err error, apiKey string) error {
	sanitized := sanitizeSensitive(err.Error(), apiKey)

	if errors.Is(err, context.DeadlineExceeded) {
		return &TimeoutError{
			APIError: &APIError{
				Message: "request timed out",
				Cause:   errors.New(sanitized),
			},
		}
	}

	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return &TimeoutError{
			APIError: &APIError{
				Message: "request timed out",
				Cause:   errors.New(sanitized),
			},
		}
	}

	return &ConnectionError{
		APIError: &APIError{
			Message: fmt.Sprintf("failed to connect to Scrappey API: %s", sanitized),
			Cause:   errors.New(sanitized),
		},
	}
}

func sanitizeSensitive(value string, apiKey string) string {
	sanitized := value
	if apiKey != "" {
		sanitized = strings.ReplaceAll(sanitized, apiKey, "REDACTED")
	}
	return apiKeyQueryPattern.ReplaceAllString(sanitized, "${1}REDACTED")
}
