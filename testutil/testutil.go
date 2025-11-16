package testutil

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

// MockServer creates a test HTTP server that returns predefined responses
type MockServer struct {
	*httptest.Server
	Responses    map[string]interface{}
	RequestCount map[string]int
	mu           sync.RWMutex
}

// NewMockServer creates a new mock Cielo API server
func NewMockServer(t *testing.T) *MockServer {
	ms := &MockServer{
		Responses:    make(map[string]interface{}),
		RequestCount: make(map[string]int),
	}

	ms.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if r.URL.RawQuery != "" {
			path = path + "?" + r.URL.RawQuery
		}

		ms.mu.Lock()
		ms.RequestCount[path]++
		ms.mu.Unlock()

		ms.mu.RLock()
		response, ok := ms.Responses[path]
		ms.mu.RUnlock()

		if !ok {
			w.WriteHeader(http.StatusNotFound)
			if err := json.NewEncoder(w).Encode(map[string]string{
				"error": "not found",
			}); err != nil {
				t.Logf("Failed to encode 404 response: %v", err)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	}))

	t.Cleanup(func() {
		ms.Close()
	})
	return ms
}

// SetResponse configures a mock response for a specific path
func (ms *MockServer) SetResponse(path string, response interface{}) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.Responses[path] = response
}

// GetRequestCount returns the number of times a path was called
func (ms *MockServer) GetRequestCount(path string) int {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	return ms.RequestCount[path]
}

// AssertRequestCount checks that a path was called a specific number of times
func (ms *MockServer) AssertRequestCount(t *testing.T, path string, expected int) {
	t.Helper()
	actual := ms.GetRequestCount(path)
	if actual != expected {
		t.Errorf("Expected %d requests to %s, got %d", expected, path, actual)
	}
}

// Ptr returns a pointer to the given value. Useful for setting optional fields in tests.
func Ptr[T any](v T) *T {
	return &v
}
