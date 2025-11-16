package cielogo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/sealtv/cielogo/api"
)

const apiBaseUrl = "https://feed-api.cielo.finance/api"

type Client struct {
	apiKey  string
	baseURL string
	cli     *http.Client
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
		c.baseURL = baseURL
	}
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
	url := c.baseURL + path

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

	if resp.StatusCode != http.StatusOK {
		var cErr error = &api.Error{}
		if err := json.NewDecoder(resp.Body).Decode(cErr); err != nil {
			return fmt.Errorf("failed to unmarshal error response body: %w, status code: %d", err, resp.StatusCode)
		}

		return cErr
	}

	if resp != nil && out != nil {
		if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
			return fmt.Errorf("failed to unmarshal response body: %w", err)
		}
	}

	return nil
}
