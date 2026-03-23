package api

import (
	"testing"
	"time"
)

// TestNewClient verifies that NewClient creates a client with default settings
func TestNewClient(t *testing.T) {
	client := NewClient()

	if client == nil {
		t.Fatal("NewClient returned nil")
	}

	if client.baseURL != DefaultBaseURL {
		t.Errorf("expected baseURL %s, got %s", DefaultBaseURL, client.baseURL)
	}

	if client.timeout != DefaultTimeout {
		t.Errorf("expected timeout %v, got %v", DefaultTimeout, client.timeout)
	}

	if client.maxRetries != DefaultMaxRetries {
		t.Errorf("expected maxRetries %d, got %d", DefaultMaxRetries, client.maxRetries)
	}
}

// TestNewClientWithOptions verifies that options are applied correctly
func TestNewClientWithOptions(t *testing.T) {
	customTimeout := 15 * time.Second
	customRetries := 5
	customBaseURL := "https://custom.api.com"

	client := NewClient(
		WithTimeout(customTimeout),
		WithMaxRetries(customRetries),
		WithBaseURL(customBaseURL),
	)

	if client.timeout != customTimeout {
		t.Errorf("expected timeout %v, got %v", customTimeout, client.timeout)
	}

	if client.maxRetries != customRetries {
		t.Errorf("expected maxRetries %d, got %d", customRetries, client.maxRetries)
	}

	if client.baseURL != customBaseURL {
		t.Errorf("expected baseURL %s, got %s", customBaseURL, client.baseURL)
	}
}

// TestAPIErrorError verifies the APIError.Error() method
func TestAPIErrorError(t *testing.T) {
	err := &APIError{
		StatusCode: 404,
		Message:    "Pokemon not found",
	}

	expected := "PokeAPI error (status 404): Pokemon not found"
	if err.Error() != expected {
		t.Errorf("expected %q, got %q", expected, err.Error())
	}
}

// TestCalculateBackoff verifies exponential backoff calculation
func TestCalculateBackoff(t *testing.T) {
	client := NewClient()

	// Test that backoff increases exponentially
	backoff0 := client.calculateBackoff(0)
	backoff1 := client.calculateBackoff(1)
	backoff2 := client.calculateBackoff(2)

	// Each should be roughly 2x the previous (accounting for jitter)
	if backoff1 < backoff0 {
		t.Errorf("backoff should increase: %v < %v", backoff1, backoff0)
	}

	if backoff2 < backoff1 {
		t.Errorf("backoff should increase: %v < %v", backoff2, backoff1)
	}

	// Test that backoff is capped at MaxBackoff
	backoffLarge := client.calculateBackoff(20)
	if backoffLarge > MaxBackoff*2 {
		t.Errorf("backoff should not exceed MaxBackoff + jitter: %v > %v", backoffLarge, MaxBackoff*2)
	}
}
