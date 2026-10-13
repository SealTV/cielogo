package cielogo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/sealtv/cielogo/api"
)

const apiBaseUrl = "https://feed-api.cielo.finance/api"

type Client struct {
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
		c.cli = client
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

func NewClient(apiKey string, opts ...ClientOption) *Client {
	transport, ok := http.DefaultTransport.(*http.Transport)
	if !ok {
		transport = &http.Transport{}
	}
	transport = transport.Clone()
	transport.MaxIdleConnsPerHost = 100

	client := &Client{
		apiKey:  apiKey,
		baseURL: apiBaseUrl,
		cli: &http.Client{
			Transport: transport,
		},
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

func (c *Client) makeRequest(ctx context.Context, method, path string, bodyObj, out any) error {
	return c.makeRequestAt(ctx, c.baseURL, method, path, bodyObj, out)
}

func (c *Client) makeRequestAt(ctx context.Context, baseURL, method, path string, bodyObj, out any) error {
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

	resp, err := c.cli.Do(req)
	if err != nil {
		return fmt.Errorf("failed to do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusAccepted || resp.StatusCode < 200 || resp.StatusCode >= 300 {
		cErr := &api.Error{StatusCode: resp.StatusCode, RetryAfter: resp.Header.Get("Retry-After")}
		if seconds, err := strconv.ParseInt(cErr.RetryAfter, 10, 64); err == nil && seconds >= 0 && seconds <= int64((1<<63-1)/time.Second) {
			cErr.RetryAfterDuration = time.Duration(seconds) * time.Second
		} else if deadline, err := http.ParseTime(cErr.RetryAfter); err == nil {
			cErr.RetryAfterDuration = time.Until(deadline)
			if cErr.RetryAfterDuration < 0 {
				cErr.RetryAfterDuration = 0
			}
		}
		if err := json.NewDecoder(resp.Body).Decode(cErr); err != nil {
			cErr.Message = http.StatusText(resp.StatusCode)
		}
		return cErr
	}

	if out != nil && resp.StatusCode != http.StatusNoContent {
		if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
			return fmt.Errorf("failed to unmarshal response body: %w", err)
		}
	}

	return nil
}
