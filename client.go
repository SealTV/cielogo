package cielogo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sealtv/cielogo/api"
)

const apiBaseUrl = "https://feed-api.cielo.finance/api"

type Client struct {
	apiKey string
	cli    *http.Client
}

func NewClient(apiKey string) *Client {
	transport := http.DefaultTransport.(*http.Transport).Clone()

	transport.MaxIdleConnsPerHost = 100

	return &Client{
		apiKey: apiKey,
		cli: &http.Client{
			Transport: transport,
		},
	}
}

func (c *Client) makeRequest(ctx context.Context, method, path string, body, out any) error {
	url := apiBaseUrl + path

	var buf *bytes.Buffer
	if body != nil {
		buf = new(bytes.Buffer)
		if err := json.NewEncoder(buf).Encode(body); err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("X-API-KEY", c.apiKey)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err := c.cli.Do(req)
	if err != nil {
		return fmt.Errorf("failed to do request: %w", err)
	}
	defer resp.Body.Close()

	fmt.Println(resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		var cErr error = &api.Error{}
		if err := json.NewDecoder(resp.Body).Decode(cErr); err != nil {
			return fmt.Errorf("failed to unmarshal error response body: %w, status code: %d", err, resp.StatusCode)
		}

		return cErr
	}

	if resp != nil {
		if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
			return fmt.Errorf("failed to unmarshal response body: %w", err)
		}
	}

	return nil
}
