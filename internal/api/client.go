package api

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"time"
)

const (
	// DefaultBaseURL is the base URL for the PokeAPI v2
	DefaultBaseURL = "https://pokeapi.co/api/v2"

	// DefaultTimeout is the default HTTP request timeout
	DefaultTimeout = 10 * time.Second

	// DefaultMaxRetries is the default number of retry attempts
	DefaultMaxRetries = 3

	// InitialBackoff is the initial backoff duration for retries
	InitialBackoff = 100 * time.Millisecond

	// MaxBackoff is the maximum backoff duration for retries
	MaxBackoff = 10 * time.Second

	// UserAgent is the User-Agent header sent with requests
	UserAgent = "Pokedex CLI (+https://github.com/pfjndev/pokedex)"
)

// Client represents an HTTP client for the PokeAPI.
// It manages request/response handling, timeouts, and retry logic.
type Client struct {
	httpClient *http.Client
	baseURL    string
	timeout    time.Duration
	maxRetries int
	initialBO  time.Duration
}

// ClientOption is a functional option for configuring the Client.
type ClientOption func(*Client)

// WithTimeout sets a custom timeout for requests.
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.timeout = timeout
		c.httpClient.Timeout = timeout
	}
}

// WithMaxRetries sets a custom number of retry attempts.
func WithMaxRetries(maxRetries int) ClientOption {
	return func(c *Client) {
		c.maxRetries = maxRetries
	}
}

// WithBaseURL sets a custom base URL for the API.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

// NewClient creates a new PokeAPI client with default settings.
// It accepts optional ClientOption arguments to customize the client.
//
// Example:
//
//	client := NewClient()
//	client := NewClient(WithTimeout(15 * time.Second), WithMaxRetries(5))
func NewClient(opts ...ClientOption) *Client {
	c := &Client{
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
		baseURL:    DefaultBaseURL,
		timeout:    DefaultTimeout,
		maxRetries: DefaultMaxRetries,
		initialBO:  InitialBackoff,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// APIError represents an error from the PokeAPI.
type APIError struct {
	StatusCode int
	Message    string
	Err        error
}

// Error implements the error interface for APIError.
func (e *APIError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("PokeAPI error (status %d): %s (%v)", e.StatusCode, e.Message, e.Err)
	}
	return fmt.Sprintf("PokeAPI error (status %d): %s", e.StatusCode, e.Message)
}

// GetPokemon fetches a Pokemon by name or ID from the PokeAPI.
// The nameOrID parameter can be either a Pokemon name (e.g., "bulbasaur") or an ID (e.g., "1").
// It returns a PokemonResponse containing all Pokemon data or an error if the request fails.
//
// Error cases:
// - 404 Not Found: Pokemon does not exist
// - 429 Too Many Requests: Rate limited by PokeAPI
// - Network errors: Connection issues
// - Invalid JSON: Malformed response
//
// Example:
//
//	pokemon, err := client.GetPokemon("bulbasaur")
//	pokemon, err := client.GetPokemon("1")
func (c *Client) GetPokemon(nameOrID string) (*PokemonResponse, error) {
	url := fmt.Sprintf("%s/pokemon/%s", c.baseURL, nameOrID)
	resp := &PokemonResponse{}

	if err := c.doRequestWithRetry(url, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// ListPokemon fetches a paginated list of Pokemon from the PokeAPI.
// The limit parameter controls how many Pokemon are returned (default: 20, max: 10000).
// The offset parameter is used for pagination (default: 0).
// It returns a PokemonListResponse containing the list of Pokemon or an error if the request fails.
//
// Example:
//
//	list, err := client.ListPokemon(20, 0)   // First 20 Pokemon
//	list, err := client.ListPokemon(20, 20)  // Next 20 Pokemon
func (c *Client) ListPokemon(limit, offset int) (*PokemonListResponse, error) {
	url := fmt.Sprintf("%s/pokemon?limit=%d&offset=%d", c.baseURL, limit, offset)
	resp := &PokemonListResponse{}

	if err := c.doRequestWithRetry(url, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// doRequestWithRetry performs an HTTP GET request with exponential backoff retry logic.
// It automatically retries on transient failures (5xx, rate limiting 429, timeouts, etc).
// The maxAttempts is determined by c.maxRetries.
func (c *Client) doRequestWithRetry(url string, v any) error {
	var lastErr error

	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		err := c.doRequest(url, v)
		if err == nil {
			return nil
		}

		lastErr = err

		// Don't retry on 404 Not Found
		if apiErr, ok := err.(*APIError); ok && apiErr.StatusCode == http.StatusNotFound {
			return err
		}

		// Don't retry if we've exhausted our attempts
		if attempt == c.maxRetries {
			break
		}

		// Calculate exponential backoff with jitter
		backoff := c.calculateBackoff(attempt)
		time.Sleep(backoff)
	}

	return lastErr
}

// doRequest performs a single HTTP GET request and unmarshals the JSON response.
// It adds proper headers and handles various error conditions.
func (c *Client) doRequest(url string, v any) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return &APIError{
			Message: "failed to create request",
			Err:     err,
		}
	}

	// Add User-Agent header
	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return &APIError{
			Message: "failed to execute request",
			Err:     err,
		}
	}
	defer func() { _ = resp.Body.Close() }()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &APIError{
			StatusCode: resp.StatusCode,
			Message:    "failed to read response body",
			Err:        err,
		}
	}

	// Check HTTP status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return c.handleErrorStatus(resp.StatusCode, body)
	}

	// Unmarshal JSON response
	if err := json.Unmarshal(body, v); err != nil {
		return &APIError{
			StatusCode: resp.StatusCode,
			Message:    "failed to unmarshal response JSON",
			Err:        err,
		}
	}

	return nil
}

// handleErrorStatus processes HTTP error status codes and returns appropriate error messages.
// It provides specific messages for common error codes and handles rate limiting.
func (c *Client) handleErrorStatus(statusCode int, body []byte) error {
	var message string

	switch statusCode {
	case http.StatusNotFound:
		message = "Pokemon not found"
	case http.StatusTooManyRequests:
		message = "rate limited by PokeAPI; please try again later"
	case http.StatusBadRequest:
		message = "bad request"
	case http.StatusUnauthorized:
		message = "unauthorized"
	case http.StatusForbidden:
		message = "forbidden"
	case http.StatusInternalServerError:
		message = "internal server error"
	case http.StatusBadGateway:
		message = "bad gateway"
	case http.StatusServiceUnavailable:
		message = "service unavailable"
	case http.StatusGatewayTimeout:
		message = "gateway timeout"
	default:
		message = fmt.Sprintf("HTTP error: %d", statusCode)
	}

	// Try to include response body if available
	if len(body) > 0 {
		// Attempt to parse error details if the API returned JSON
		var errData map[string]any
		if err := json.Unmarshal(body, &errData); err == nil {
			if detail, ok := errData["detail"].(string); ok {
				message = detail
			}
		}
	}

	return &APIError{
		StatusCode: statusCode,
		Message:    message,
	}
}

// calculateBackoff calculates the exponential backoff duration with jitter for retry attempts.
// Formula: min(MaxBackoff, InitialBackoff * 2^attempt + jitter)
// Jitter helps prevent thundering herd problem when multiple clients retry simultaneously.
func (c *Client) calculateBackoff(attempt int) time.Duration {
	// Calculate exponential backoff: 2^attempt
	backoff := min(
		// Cap at MaxBackoff
		c.initialBO*time.Duration(math.Pow(2, float64(attempt))), MaxBackoff)

	// Add random jitter (±10% of backoff duration)
	jitter := time.Duration(rand.Int63n(int64(backoff / 5)))
	return backoff + jitter
}
