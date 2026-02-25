package scrappey

import (
	"fmt"
)

// APIError is the base error type for client-side failures.
type APIError struct {
	Message    string
	StatusCode int
	Body       string
	Cause      error
}

func (e *APIError) Error() string {
	if e == nil {
		return "scrappey error"
	}
	if e.StatusCode > 0 {
		return fmt.Sprintf("%s (HTTP %d)", e.Message, e.StatusCode)
	}
	return e.Message
}

func (e *APIError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Cause
}

// AuthenticationError is returned when the API key is missing or rejected.
type AuthenticationError struct {
	*APIError
}

// ConnectionError is returned when the API is unreachable.
type ConnectionError struct {
	*APIError
}

// TimeoutError is returned when the request times out.
type TimeoutError struct {
	*APIError
}
