// Package evilmail provides a Go client for the EvilMail API.
//
// The client supports all EvilMail API endpoints including temporary email
// creation, account management, inbox access, verification code extraction,
// random email generation, and domain listing.
//
// Basic usage:
//
//	client := evilmail.New("your-api-key")
//
//	temp, err := client.TempEmail.Create(ctx, &evilmail.CreateTempEmailRequest{
//	    Domain:     "evilmail.pro",
//	    TTLMinutes: 60,
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(temp.Email)
//
// The client can be customized using functional options:
//
//	client := evilmail.New("your-api-key",
//	    evilmail.WithBaseURL("https://custom.url"),
//	    evilmail.WithHTTPClient(customHTTPClient),
//	    evilmail.WithTimeout(30 * time.Second),
//	)
package evilmail

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	// DefaultBaseURL is the default base URL for the EvilMail API.
	DefaultBaseURL = "https://evilmail.pro"

	// DefaultTimeout is the default timeout for HTTP requests.
	DefaultTimeout = 30 * time.Second

	// Version is the current version of this SDK.
	Version = "1.0.0"
)

// Client is the top-level EvilMail API client. It holds configuration
// and provides access to each API resource through dedicated service fields.
type Client struct {
	// TempEmail provides access to the temporary email API.
	TempEmail *TempEmailService

	// Accounts provides access to the email accounts API.
	Accounts *AccountsService

	// Inbox provides access to the inbox and message API.
	Inbox *InboxService

	// Verification provides access to the verification code extraction API.
	Verification *VerificationService

	// RandomEmail provides access to the random email generation API.
	RandomEmail *RandomEmailService

	// Domains provides access to the authenticated domains API (includes customer domains).
	Domains *DomainsService

	// PublicDomains provides access to the public (unauthenticated) domains API.
	PublicDomains *PublicDomainsService

	// Shortlinks provides access to the shortlink creation API.
	Shortlinks *ShortlinksService

	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// Option is a functional option for configuring the Client.
type Option func(*Client)

// WithBaseURL sets a custom base URL for the API.
// This is useful for testing or for connecting to a self-hosted instance.
func WithBaseURL(baseURL string) Option {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

// WithHTTPClient sets a custom HTTP client for making requests.
// This allows full control over transport settings, proxies, TLS configuration,
// and other low-level HTTP behavior.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithTimeout sets the timeout for HTTP requests made by the client.
// If you need more control over the HTTP client, use WithHTTPClient instead.
// Note: if combined with WithHTTPClient, apply WithHTTPClient first.
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		if c.httpClient == nil {
			c.httpClient = &http.Client{}
		}
		c.httpClient.Timeout = timeout
	}
}

// New creates a new EvilMail API client with the given API key and options.
//
// The API key is used for authentication via the X-API-Key header on all requests.
// At minimum, an API key is required:
//
//	client := evilmail.New("your-api-key")
//
// The client can be further configured using functional options:
//
//	client := evilmail.New("your-api-key",
//	    evilmail.WithBaseURL("https://custom.url"),
//	    evilmail.WithTimeout(60 * time.Second),
//	)
func New(apiKey string, opts ...Option) *Client {
	c := &Client{
		baseURL: DefaultBaseURL,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
	}

	for _, opt := range opts {
		opt(c)
	}

	c.TempEmail = &TempEmailService{client: c}
	c.Accounts = &AccountsService{client: c}
	c.Inbox = &InboxService{client: c}
	c.Verification = &VerificationService{client: c}
	c.RandomEmail = &RandomEmailService{client: c}
	c.Domains = &DomainsService{client: c}
	c.PublicDomains = &PublicDomainsService{client: c}
	c.Shortlinks = &ShortlinksService{client: c}

	return c
}

// newRequest creates a new HTTP request with the appropriate headers set.
// The path may include query parameters (e.g., "/api/foo?bar=baz").
// The body parameter, if non-nil, will be JSON-encoded.
func (c *Client) newRequest(ctx context.Context, method, path string, body any) (*http.Request, error) {
	base, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, fmt.Errorf("evilmail: invalid base URL %q: %w", c.baseURL, err)
	}

	ref, err := url.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("evilmail: invalid path %q: %w", path, err)
	}

	u := base.ResolveReference(ref).String()

	var bodyReader io.Reader
	if body != nil {
		buf, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("evilmail: failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(buf)
	}

	req, err := http.NewRequestWithContext(ctx, method, u, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("evilmail: failed to create request: %w", err)
	}

	req.Header.Set("X-API-Key", c.apiKey)
	req.Header.Set("User-Agent", "evilmail-go/"+Version)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}

// do executes the given HTTP request and decodes the JSON response into v.
// If v is nil, the response body is discarded.
// Returns an *APIError for non-2xx status codes.
func (c *Client) do(req *http.Request, v any) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("evilmail: request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("evilmail: failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		apiErr := &APIError{
			StatusCode: resp.StatusCode,
			Status:     resp.Status,
			Body:       body,
		}

		// Attempt to extract a message from the JSON response.
		var errResp struct {
			Message string `json:"message"`
			Error   string `json:"error"`
		}
		if json.Unmarshal(body, &errResp) == nil {
			if errResp.Message != "" {
				apiErr.Message = errResp.Message
			} else if errResp.Error != "" {
				apiErr.Message = errResp.Error
			}
		}

		if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
			return &AuthError{Err: apiErr}
		}

		return apiErr
	}

	if v != nil && len(body) > 0 {
		if err := json.Unmarshal(body, v); err != nil {
			return fmt.Errorf("evilmail: failed to decode response: %w", err)
		}
	}

	return nil
}

// get is a convenience method for making GET requests.
func (c *Client) get(ctx context.Context, path string, v any) error {
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return err
	}
	return c.do(req, v)
}

// post is a convenience method for making POST requests with a JSON body.
func (c *Client) post(ctx context.Context, path string, body, v any) error {
	req, err := c.newRequest(ctx, http.MethodPost, path, body)
	if err != nil {
		return err
	}
	return c.do(req, v)
}

// put is a convenience method for making PUT requests with a JSON body.
func (c *Client) put(ctx context.Context, path string, body, v any) error {
	req, err := c.newRequest(ctx, http.MethodPut, path, body)
	if err != nil {
		return err
	}
	return c.do(req, v)
}

// del is a convenience method for making DELETE requests with a JSON body.
func (c *Client) del(ctx context.Context, path string, body, v any) error {
	req, err := c.newRequest(ctx, http.MethodDelete, path, body)
	if err != nil {
		return err
	}
	return c.do(req, v)
}

// getJSON makes a GET request and unwraps the standard { "status", "data" }
// API response envelope, decoding the "data" field into the value pointed
// to by result.
func getJSON[T any](c *Client, ctx context.Context, path string) (T, error) {
	var resp apiResponse[T]
	if err := c.get(ctx, path, &resp); err != nil {
		var zero T
		return zero, err
	}
	return resp.Data, nil
}

// postJSON makes a POST request and unwraps the standard API response envelope.
func postJSON[T any](c *Client, ctx context.Context, path string, body any) (T, error) {
	var resp apiResponse[T]
	if err := c.post(ctx, path, body, &resp); err != nil {
		var zero T
		return zero, err
	}
	return resp.Data, nil
}

// putJSON makes a PUT request and unwraps the standard API response envelope.
func putJSON[T any](c *Client, ctx context.Context, path string, body any) (T, error) {
	var resp apiResponse[T]
	if err := c.put(ctx, path, body, &resp); err != nil {
		var zero T
		return zero, err
	}
	return resp.Data, nil
}

// delJSON makes a DELETE request and unwraps the standard API response envelope.
func delJSON[T any](c *Client, ctx context.Context, path string, body any) (T, error) {
	var resp apiResponse[T]
	if err := c.del(ctx, path, body, &resp); err != nil {
		var zero T
		return zero, err
	}
	return resp.Data, nil
}
