package scrappey

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClientRequiresAPIKey(t *testing.T) {
	t.Parallel()

	_, err := NewClient("", nil)
	if err == nil {
		t.Fatal("expected error when api key is missing")
	}

	var authErr *AuthenticationError
	if !errors.As(err, &authErr) {
		t.Fatalf("expected AuthenticationError, got %T", err)
	}
}

func TestGetBuildsExpectedPayload(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if got := r.URL.Query().Get("key"); got != "test-key" {
			t.Fatalf("expected key query test-key, got %s", got)
		}

		defer r.Body.Close()
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed decoding request body: %v", err)
		}

		if got := body["cmd"]; got != "request.get" {
			t.Fatalf("expected cmd request.get, got %v", got)
		}
		if got := body["url"]; got != "https://example.com" {
			t.Fatalf("expected url https://example.com, got %v", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":"success","solution":{"statusCode":200,"response":"ok"}}`))
	}))
	defer server.Close()

	client, err := NewClient("test-key", &Config{BaseURL: server.URL})
	if err != nil {
		t.Fatalf("failed creating client: %v", err)
	}

	res, err := client.Get(context.Background(), RequestOptions{
		"url": "https://example.com",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Data != "success" {
		t.Fatalf("expected data=success, got %s", res.Data)
	}
	if got := res.SolutionInt("statusCode"); got != 200 {
		t.Fatalf("expected statusCode 200, got %d", got)
	}
}

func TestRequestRequiresCmd(t *testing.T) {
	t.Parallel()

	client, err := NewClient("test-key", nil)
	if err != nil {
		t.Fatalf("failed creating client: %v", err)
	}

	_, err = client.Request(context.Background(), map[string]any{
		"url": "https://example.com",
	})
	if err == nil {
		t.Fatal("expected error when cmd is missing")
	}
}

func TestAuthenticationErrorOnHTTP401(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":"invalid api key"}`))
	}))
	defer server.Close()

	client, err := NewClient("bad-key", &Config{BaseURL: server.URL})
	if err != nil {
		t.Fatalf("failed creating client: %v", err)
	}

	_, err = client.Get(context.Background(), RequestOptions{
		"url": "https://example.com",
	})
	if err == nil {
		t.Fatal("expected authentication error")
	}

	var authErr *AuthenticationError
	if !errors.As(err, &authErr) {
		t.Fatalf("expected AuthenticationError, got %T", err)
	}
}
