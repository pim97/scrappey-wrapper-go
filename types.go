package scrappey

import (
	"net/http"
	"time"
)

const (
	// DefaultBaseURL is the official Scrappey API endpoint.
	DefaultBaseURL = "https://publisher.scrappey.com/api/v1"
	// DefaultTimeout is applied when no timeout is provided.
	DefaultTimeout = 5 * time.Minute
)

// Config controls client behavior.
type Config struct {
	BaseURL    string
	Timeout    time.Duration
	HTTPClient *http.Client
}

// RequestOptions contains parameters sent to request.* commands.
// It intentionally stays dynamic to match Scrappey's evolving API surface.
type RequestOptions = map[string]any

// SessionOptions contains parameters for sessions.create.
type SessionOptions = map[string]any

// WebSocketOptions contains parameters for websocket.create.
type WebSocketOptions = map[string]any

// SessionInfo describes one active session.
type SessionInfo struct {
	Session      string `json:"session"`
	LastAccessed int64  `json:"lastAccessed,omitempty"`
}

// APIResponse is the common response envelope returned by Scrappey.
type APIResponse struct {
	Solution    map[string]any `json:"solution,omitempty"`
	TimeElapsed int            `json:"timeElapsed,omitempty"`
	Data        string         `json:"data,omitempty"`
	Session     string         `json:"session,omitempty"`
	Error       string         `json:"error,omitempty"`
	Info        string         `json:"info,omitempty"`
	Active      bool           `json:"active,omitempty"`
	Sessions    []SessionInfo  `json:"sessions,omitempty"`
	Open        int            `json:"open,omitempty"`
	Limit       int            `json:"limit,omitempty"`
	Fingerprint map[string]any `json:"fingerprint,omitempty"`
	Context     map[string]any `json:"context,omitempty"`

	// HTTPStatus is not part of API JSON. It captures the HTTP status code.
	HTTPStatus int `json:"-"`
	// Raw contains the full decoded JSON body.
	Raw map[string]any `json:"-"`
}

// SolutionString returns solution[key] as string when possible.
func (r *APIResponse) SolutionString(key string) string {
	if r == nil || r.Solution == nil {
		return ""
	}
	value, ok := r.Solution[key]
	if !ok {
		return ""
	}
	asString, ok := value.(string)
	if !ok {
		return ""
	}
	return asString
}

// SolutionInt returns solution[key] as int when possible.
func (r *APIResponse) SolutionInt(key string) int {
	if r == nil || r.Solution == nil {
		return 0
	}
	value, ok := r.Solution[key]
	if !ok {
		return 0
	}
	switch number := value.(type) {
	case int:
		return number
	case int8:
		return int(number)
	case int16:
		return int(number)
	case int32:
		return int(number)
	case int64:
		return int(number)
	case uint:
		return int(number)
	case uint8:
		return int(number)
	case uint16:
		return int(number)
	case uint32:
		return int(number)
	case uint64:
		return int(number)
	case float32:
		return int(number)
	case float64:
		return int(number)
	default:
		return 0
	}
}
