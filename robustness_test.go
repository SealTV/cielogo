package cielogo_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/sealtv/cielogo"
	"github.com/sealtv/cielogo/api"
	"github.com/sealtv/cielogo/api/apiv1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInvalidSuccessResponses(t *testing.T) {
	for _, body := range []string{`null`, `{}`, `{"status":"ok"}`, `{"status":"ok","data":null}`, `{"status":"ok","data":[]}`, `{"status":"ok","data":{}} {}`, `{"status":"ok","data":{"pnl":"invalid"}}`} {
		t.Run(body, func(t *testing.T) {
			s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { _, _ = io.WriteString(w, body) }))
			defer s.Close()
			resp, err := cielogo.NewClient("key", cielogo.WithBaseURL(s.URL)).GetTradingStatsV1(context.Background(), &apiv1.TradingStatsRequest{Wallet: "wallet1"})
			require.ErrorIs(t, err, cielogo.ErrInvalidResponse)
			assert.Nil(t, resp)
		})
	}
}

func TestApplicationErrorWithHTTP200(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, `{"status":"error","message":"upstream error"}`)
	}))
	defer s.Close()
	resp, err := cielogo.NewClient("key", cielogo.WithBaseURL(s.URL)).GetWalletPortfolioV1(context.Background(), "wallet1")
	var apiErr *api.Error
	require.ErrorAs(t, err, &apiErr)
	assert.Equal(t, http.StatusOK, apiErr.StatusCode)
	assert.Equal(t, "upstream error", apiErr.Message)
	assert.Nil(t, resp)
}

func TestResponseLimit(t *testing.T) {
	for _, status := range []int{http.StatusOK, http.StatusTooManyRequests} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(status)
				_, _ = io.WriteString(w, strings.Repeat("x", 1000))
			}))
			defer s.Close()
			_, err := cielogo.NewClient("key", cielogo.WithBaseURL(s.URL), cielogo.WithMaxResponseBytes(100)).GetWalletPortfolioV1(context.Background(), "wallet1")
			require.ErrorIs(t, err, cielogo.ErrResponseTooLarge)
			if status == http.StatusTooManyRequests {
				var apiErr *api.Error
				require.ErrorAs(t, err, &apiErr)
				assert.Equal(t, status, apiErr.StatusCode)
			}
		})
	}
}

func TestCrossOriginRedirectDoesNotSendAPIKey(t *testing.T) {
	var hits atomic.Int64
	sink := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { hits.Add(1) }))
	defer sink.Close()
	source := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, sink.URL, http.StatusTemporaryRedirect)
	}))
	defer source.Close()
	custom := &http.Client{Timeout: time.Second}
	_, err := cielogo.NewClient("sensitive-test-key", cielogo.WithBaseURL(source.URL), cielogo.WithHTTPClient(custom)).GetWalletPortfolioV1(context.Background(), "wallet1")
	require.ErrorIs(t, err, cielogo.ErrUnsafeRedirect)
	assert.Zero(t, hits.Load())
	assert.Nil(t, custom.CheckRedirect)
}

func TestSameOriginRedirect(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/redirected" {
			http.Redirect(w, r, "/redirected", http.StatusTemporaryRedirect)
			return
		}
		assert.Equal(t, "key", r.Header.Get("X-Api-Key"))
		_, _ = io.WriteString(w, `{"status":"ok","data":{"id":17}}`)
	}))
	defer s.Close()
	resp, err := cielogo.NewClient("key", cielogo.WithBaseURL(s.URL)).GetWalletByAddressV1(context.Background(), "wallet1")
	require.NoError(t, err)
	assert.Equal(t, int64(17), resp.ID)
}

func TestPathParametersCannotInjectQuery(t *testing.T) {
	wallet := "wallet?injected=true#fragment%+"
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v1/"+url.PathEscape(wallet)+"/portfolio", r.URL.EscapedPath())
		assert.Empty(t, r.URL.RawQuery)
		assert.Empty(t, r.URL.Fragment)
		_, _ = io.WriteString(w, `{"status":"ok","data":{"portfolio":[]}}`)
	}))
	defer s.Close()
	_, err := cielogo.NewClient("key", cielogo.WithBaseURL(s.URL)).GetWalletPortfolioV1(context.Background(), wallet)
	require.NoError(t, err)
}

func TestNilRequestsAndClients(t *testing.T) {
	c := cielogo.NewClient("key", nil, cielogo.WithHTTPClient(nil))
	calls := []func() error{
		func() error { _, e := c.GetFeedV1(context.Background(), nil); return e },
		func() error { _, e := c.GetNftsPnlV1(context.Background(), nil); return e },
		func() error { _, e := c.GetTokensPnlV1(context.Background(), nil); return e },
		func() error { _, e := c.GetAggregatedTokenPnLV1(context.Background(), nil); return e },
		func() error { _, e := c.GetRelatedWalletsV1(context.Background(), nil); return e },
		func() error { _, e := c.GetWalletsTagsV1(context.Background(), nil); return e },
		func() error { _, e := c.GetWalletsByTagV1(context.Background(), nil); return e },
		func() error { _, e := c.GetAllWalletsListV1(context.Background(), nil); return e },
		func() error { _, e := c.AddWalletsListV1(context.Background(), nil); return e },
		func() error { _, e := c.UpdateWalletsListV1(context.Background(), nil); return e },
		func() error { _, e := c.GetTrackedWalletsV1(context.Background(), nil); return e },
		func() error { _, e := c.AddTrackedWalletsV1(context.Background(), nil); return e },
		func() error { return c.RemoveTrackedWalletsV1(context.Background(), nil) },
		func() error { _, e := c.UpdateTrackedWalletV1(context.Background(), 17, nil); return e },
		func() error { _, e := c.UpdateTrackedWalletV2(context.Background(), "wallet1", nil); return e },
		func() error { _, e := c.GetWalletPortfolioV2(context.Background(), nil); return e },
		func() error { _, e := c.GetTokenMetadataV1(context.Background(), nil); return e },
		func() error { _, e := c.GetTokenPriceV1(context.Background(), nil); return e },
		func() error { _, e := c.GetTokenStatsV1(context.Background(), nil); return e },
		func() error { _, e := c.GetTokenBalanceV1(context.Background(), nil); return e },
		func() error { _, e := c.GetTradingStatsV1(context.Background(), nil); return e },
	}
	for i, call := range calls {
		assert.ErrorIs(t, call(), cielogo.ErrNilRequest, "method %d", i)
	}
	_, err := (&cielogo.Client{}).GetChainsV1(context.Background())
	require.ErrorIs(t, err, cielogo.ErrInvalidClient)
	var nilClient *cielogo.Client
	_, err = nilClient.GetChainsV1(context.Background())
	require.ErrorIs(t, err, cielogo.ErrInvalidClient)
	_, err = c.GetChainsV1(nil) //nolint:staticcheck // Deliberately exercise the nil-context guard.
	require.ErrorIs(t, err, cielogo.ErrNilContext)
	nilClient.CloseIdleConnections()
}

func TestModelSerializationAndReuse(t *testing.T) {
	payload := `{"wallet":"wallet1","tx_type":"swap","token0_address":"SOL","token0_amount":2,"verification":{"name": "verified"}}`
	var event apiv1.TxEvent
	require.NoError(t, json.Unmarshal([]byte(payload), &event))
	encoded, err := json.Marshal(event)
	require.NoError(t, err)
	var after apiv1.TxEvent
	require.NoError(t, json.Unmarshal(encoded, &after))
	swap, ok := after.Data.(*apiv1.SwapEvent)
	require.True(t, ok)
	assert.Equal(t, "SOL", swap.Token0Address)
	assert.Equal(t, 2.0, swap.Token0Amount)
	var reused apiv1.WSEvent
	require.NoError(t, json.Unmarshal([]byte(`{"type":"pong","data":"ok"}`), &reused))
	require.Error(t, json.Unmarshal([]byte(`{"type":"tx","data":{"tx_type":"swap","token0_amount":"bad"}}`), &reused))
	assert.Equal(t, apiv1.PongEventType, reused.Type)
	assert.Equal(t, "ok", reused.Data)
}

type brokenReader struct{}

func (brokenReader) Read([]byte) (int, error) { return 0, io.ErrUnexpectedEOF }
func (brokenReader) Close() error             { return nil }

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }
func TestBodyReadFailure(t *testing.T) {
	c := cielogo.NewClient("key", cielogo.WithHTTPClient(&http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusOK, Body: brokenReader{}, Header: make(http.Header)}, nil
	})}))
	_, err := c.GetChainsV1(context.Background())
	require.ErrorIs(t, err, io.ErrUnexpectedEOF)
}

func TestNullCollections(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.Copy(w, bytes.NewBufferString(`{"status":"ok","data":null}`))
	}))
	defer s.Close()
	lists, err := cielogo.NewClient("key", cielogo.WithBaseURL(s.URL)).GetUserWalletsListsV1(context.Background())
	require.NoError(t, err)
	assert.Nil(t, lists)
	assert.False(t, errors.Is(err, cielogo.ErrInvalidResponse))
}

func TestCustomDefaultTransportIsRetained(t *testing.T) {
	original := http.DefaultTransport
	http.DefaultTransport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(`{"status":"ok","data":["solana"]}`)), Header: make(http.Header)}, nil
	})
	defer func() { http.DefaultTransport = original }()
	result, err := cielogo.NewClient("key").GetChainsV1(context.Background())
	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, "solana", string(result[0]))
}
