package cielogo

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sealtv/cielogo/api"
)

const apiBaseUrl = "https://feed-api.cielo.finance/api"

// DefaultMaxResponseBytes bounds a single HTTP response (16 MiB).
const DefaultMaxResponseBytes int64 = 16 << 20

var (
	ErrListenerRunning  = errors.New("a websocket listener is already running")
	ErrNilContext       = errors.New("context must not be nil")
	ErrUnsafeRedirect   = errors.New("redirect to another origin is not allowed")
	ErrNilRequest       = errors.New("request must not be nil")
	ErrInvalidClient    = errors.New("client is not initialized")
	ErrResponseTooLarge = errors.New("response exceeds configured size limit")
	ErrInvalidResponse  = api.ErrInvalidResponse
)

type Client struct {
	websocketDialer  *websocket.Dialer
	maxResponseBytes int64
	portfolioBaseURL string
	apiKey           string
	baseURL          string
	cli              *http.Client
}

// ClientOption is a function that configures a Client
type ClientOption func(*Client)

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(client *http.Client) ClientOption {
	return func(c *Client) {
		if client != nil {
			c.cli = client
		}
	}
}

// WithBaseURL sets a custom base URL (useful for testing)
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		c.baseURL = strings.TrimRight(baseURL, "/")
	}
}

// WithPortfolioBaseURL configures the private API base URL for portfolio V2.
// Include /api in the URL. The public API does not register this endpoint.
func WithPortfolioBaseURL(baseURL string) ClientOption {
	return func(c *Client) { c.portfolioBaseURL = strings.TrimRight(baseURL, "/") }
}

// WithMaxResponseBytes changes the HTTP response limit. Nonpositive values keep
// the default; callers needing larger feed responses can choose a positive limit.
func WithMaxResponseBytes(limit int64) ClientOption {
	return func(c *Client) {
		if limit > 0 && limit < 1<<63-1 {
			c.maxResponseBytes = limit
		}
	}
}

func NewClient(apiKey string, opts ...ClientOption) *Client {
	transport := http.DefaultTransport
	if original, ok := transport.(*http.Transport); ok {
		cloned := original.Clone()
		cloned.MaxIdleConnsPerHost = 100
		transport = cloned
	}

	client := &Client{
		apiKey:           apiKey,
		maxResponseBytes: DefaultMaxResponseBytes,
		baseURL:          apiBaseUrl,
		cli: &http.Client{
			Transport: transport,
			Timeout:   90 * time.Second,
		},
	}

	for _, opt := range opts {
		if opt != nil {
			opt(client)
		}
	}

	return client
}

func (c *Client) makeRequest(ctx context.Context, method, path string, bodyObj, out any) error {
	if c == nil || c.cli == nil {
		return ErrInvalidClient
	}
	return c.makeRequestAt(ctx, c.baseURL, method, path, bodyObj, out)
}

func (c *Client) makeRequestAt(ctx context.Context, baseURL, method, path string, bodyObj, out any) error {
	if c == nil || c.cli == nil {
		return ErrInvalidClient
	}
	if ctx == nil {
		return ErrNilContext
	}
	url := baseURL + path

	var body io.Reader
	if bodyObj != nil {
		buf := new(bytes.Buffer)
		if err := json.NewEncoder(buf).Encode(bodyObj); err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}

		body = buf
	}

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("X-Api-Key", c.apiKey)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err := c.doRequest(req)
	if err != nil {
		return fmt.Errorf("failed to do request: %w", err)
	}
	defer resp.Body.Close()

	return c.decodeResponse(resp, out)
}

func (c *Client) decodeResponse(resp *http.Response, out any) error {
	content, readErr := io.ReadAll(io.LimitReader(resp.Body, c.maxResponseBytes+1))
	if int64(len(content)) > c.maxResponseBytes {
		readErr = ErrResponseTooLarge
	}
	var envelope struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}
	decodeErr := json.Unmarshal(content, &envelope)
	failedStatus := resp.StatusCode == http.StatusAccepted || resp.StatusCode < 200 || resp.StatusCode >= 300
	if failedStatus || (decodeErr == nil && envelope.Status != "" && envelope.Status != "ok") {
		cErr := &api.Error{StatusCode: resp.StatusCode, Status: envelope.Status, Message: envelope.Message, RetryAfter: resp.Header.Get("Retry-After")}
		cErr.RetryAfterDuration = parseRetryAfter(cErr.RetryAfter, time.Now())
		if cErr.Message == "" {
			cErr.Message = http.StatusText(resp.StatusCode)
		}
		if readErr != nil {
			return fmt.Errorf("%w: %w", cErr, readErr)
		}
		return cErr
	}
	if readErr != nil {
		return fmt.Errorf("cannot read response body: %w", readErr)
	}
	if resp.StatusCode == http.StatusNoContent {
		if out != nil {
			return fmt.Errorf("%w: missing response data", ErrInvalidResponse)
		}
		return nil
	}
	if decodeErr != nil {
		return fmt.Errorf("%w: %w", ErrInvalidResponse, decodeErr)
	}
	if envelope.Status != "ok" {
		return fmt.Errorf("%w: missing or unknown status", ErrInvalidResponse)
	}
	if out != nil {
		if err := json.Unmarshal(content, out); err != nil {
			return fmt.Errorf("%w: %w", ErrInvalidResponse, err)
		}
	}

	return nil
}

func parseRetryAfter(header string, now time.Time) time.Duration {
	if seconds, err := strconv.ParseInt(header, 10, 64); err == nil && seconds >= 0 && seconds <= int64((1<<63-1)/time.Second) {
		return time.Duration(seconds) * time.Second
	}
	if deadline, err := http.ParseTime(header); err == nil && deadline.After(now) {
		return deadline.Sub(now)
	}
	return 0
}

// CloseIdleConnections releases pooled HTTP connections without interrupting requests.
func (c *Client) CloseIdleConnections() {
	if c != nil && c.cli != nil {
		c.cli.CloseIdleConnections()
	}
}

func (c *Client) doRequest(req *http.Request) (*http.Response, error) {
	// Work on a copy so shared caller-owned HTTP clients are never mutated.
	requestClient := *c.cli
	requestClient.CheckRedirect = func(next *http.Request, via []*http.Request) error {
		if len(via) > 0 && (!strings.EqualFold(next.URL.Scheme, via[0].URL.Scheme) || !strings.EqualFold(next.URL.Host, via[0].URL.Host)) {
			return ErrUnsafeRedirect
		}
		if c.cli.CheckRedirect != nil {
			return c.cli.CheckRedirect(next, via)
		}
		if len(via) >= 10 {
			return errors.New("stopped after 10 redirects")
		}
		return nil
	}
	return requestClient.Do(req)
}
