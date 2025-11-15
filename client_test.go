package cielogo_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/sealtv/cielogo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	client := cielogo.NewClient("test-api-key")
	require.NotNil(t, client)
}

func TestNewClient_WithCustomHTTPClient(t *testing.T) {
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	client := cielogo.NewClient("test-api-key", cielogo.WithHTTPClient(httpClient))
	require.NotNil(t, client)
}

func TestNewClient_WithBaseURL(t *testing.T) {
	customURL := "https://custom-api.example.com"
	client := cielogo.NewClient("test-api-key", cielogo.WithBaseURL(customURL))
	require.NotNil(t, client)
}

func TestNewClient_WithMultipleOptions(t *testing.T) {
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}
	customURL := "https://custom-api.example.com"

	client := cielogo.NewClient("test-api-key",
		cielogo.WithHTTPClient(httpClient),
		cielogo.WithBaseURL(customURL),
	)
	require.NotNil(t, client)
}

func TestDefaultClient(t *testing.T) {
	// Test that NewClient creates a client with default settings
	client := cielogo.NewClient("test-api-key")

	// Client should be created successfully
	assert.NotNil(t, client)
}

func TestClient_EmptyAPIKey(t *testing.T) {
	// Even with empty API key, client should be created
	// (API key validation happens on actual requests)
	client := cielogo.NewClient("")
	assert.NotNil(t, client)
}
