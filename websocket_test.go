package cielogo_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sealtv/cielogo"
	"github.com/sealtv/cielogo/api"
	"github.com/sealtv/cielogo/api/apiv1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWebsocketAcknowledgementContract(t *testing.T) {
	cases := []struct {
		name, payload string
		inspect       func(*testing.T, any)
	}{
		{"wallet", `{"type":"wallet_subscribed","data":{"wallet":"wallet1","filter":{"chains":{"solana":true},"tx_types":{"swap":true}}}}`, func(t *testing.T, v any) {
			ack, ok := v.(apiv1.WalletSubscriptionAck)
			require.True(t, ok)
			require.NotNil(t, ack.Filter)
			assert.True(t, ack.Filter.Chains["solana"])
			assert.True(t, ack.Filter.TxTypes[apiv1.TxTypeSwap])
		}},
		{"feed", `{"type":"feed_subscribed","data":{"multi_bundle_id":17,"filter":{"tokens":{"SOL":true}}}}`, func(t *testing.T, v any) {
			ack, ok := v.(apiv1.FeedSubscriptionAck)
			require.True(t, ok)
			require.NotNil(t, ack.MultiBundleID)
			assert.Equal(t, 17, *ack.MultiBundleID)
			assert.True(t, ack.Filter.Tokens["SOL"])
		}},
		{"unsubscribe", `{"type":"feed_unsubscribed","data":"ok"}`, func(t *testing.T, v any) { assert.Equal(t, "ok", v) }},
		{"pong", `{"type":"pong","data":"ok"}`, func(t *testing.T, v any) { assert.Equal(t, "ok", v) }},
		{"future", `{"type":"future_event","data":{"new":true}}`, func(t *testing.T, v any) {
			raw, ok := v.(json.RawMessage)
			require.True(t, ok)
			assert.JSONEq(t, `{"new":true}`, string(raw))
		}},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			var event apiv1.WSEvent
			require.NoError(t, json.Unmarshal([]byte(tt.payload), &event))
			tt.inspect(t, event.Data)
		})
	}
}

func wsTestServer(t *testing.T, serve func(*websocket.Conn)) *httptest.Server {
	t.Helper()
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v1/ws", r.URL.Path)
		assert.Equal(t, "key", r.Header.Get("X-Api-Key"))
		conn, err := (&websocket.Upgrader{}).Upgrade(w, r, nil)
		if err != nil {
			t.Error(err)
			return
		}
		defer conn.Close()
		serve(conn)
	}))
	t.Cleanup(s.Close)
	return s
}

func TestWebsocketListenerCancellation(t *testing.T) {
	s := wsTestServer(t, func(conn *websocket.Conn) { _, _, _ = conn.ReadMessage() })
	c := cielogo.NewClient("key", cielogo.WithBaseURL(s.URL))
	ws, err := c.NewWebsocketConnection(context.Background())
	require.NoError(t, err)
	defer ws.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan error, 1)
	go func() { done <- ws.RunListener(ctx, make(chan apiv1.WSEvent)) }()
	cancel()
	select {
	case err := <-done:
		require.ErrorIs(t, err, context.Canceled)
	case <-time.After(2 * time.Second):
		t.Fatal("listener did not unblock on cancellation")
	}
}

func TestWebsocketDialCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	ws, err := cielogo.NewClient("key", cielogo.WithBaseURL("http://127.0.0.1:1")).NewWebsocketConnection(ctx)
	assert.Nil(t, ws)
	require.True(t, errors.Is(err, context.Canceled))
}

func TestWebsocketListenerStopsOnMalformedEvent(t *testing.T) {
	s := wsTestServer(t, func(conn *websocket.Conn) {
		_ = conn.WriteMessage(websocket.TextMessage, []byte(`{"type":"tx","data":{"tx_type":"swap","token0_amount":"invalid"}}`))
		_, _, _ = conn.ReadMessage()
	})
	ws, err := cielogo.NewClient("key", cielogo.WithBaseURL(s.URL)).NewWebsocketConnection(context.Background())
	require.NoError(t, err)
	defer ws.Close()
	out := make(chan apiv1.WSEvent, 1)
	err = ws.RunListener(context.Background(), out)
	require.Error(t, err)
	assert.Empty(t, out)
}

func TestWebsocketNormalClose(t *testing.T) {
	s := wsTestServer(t, func(conn *websocket.Conn) {
		_ = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "done"))
	})
	ws, err := cielogo.NewClient("key", cielogo.WithBaseURL(s.URL)).NewWebsocketConnection(context.Background())
	require.NoError(t, err)
	defer ws.Close()
	require.NoError(t, ws.RunListener(context.Background(), make(chan apiv1.WSEvent)))
}

func TestWebsocketConcurrentCommands(t *testing.T) {
	received := make(chan struct{})
	s := wsTestServer(t, func(conn *websocket.Conn) {
		for range 20 {
			var cmd map[string]any
			if err := conn.ReadJSON(&cmd); err != nil {
				t.Error(err)
				return
			}
			assert.Equal(t, "ping", cmd["type"])
		}
		close(received)
	})
	ws, err := cielogo.NewClient("key", cielogo.WithBaseURL(s.URL)).NewWebsocketConnection(context.Background())
	require.NoError(t, err)
	defer ws.Close()
	var wg sync.WaitGroup
	for range 20 {
		wg.Add(1)
		go func() { defer wg.Done(); assert.NoError(t, ws.SendCommand(&apiv1.PingCmd{})) }()
	}
	wg.Wait()
	select {
	case <-received:
	case <-time.After(2 * time.Second):
		t.Fatal("commands not received")
	}
}

func TestWebsocketConcurrentListenerGuard(t *testing.T) {
	sent := make(chan struct{})
	s := wsTestServer(t, func(conn *websocket.Conn) {
		_ = conn.WriteJSON(map[string]any{"type": "pong", "data": "ok"})
		close(sent)
		_, _, _ = conn.ReadMessage()
	})
	ws, err := cielogo.NewClient("key", cielogo.WithBaseURL(s.URL)).NewWebsocketConnection(context.Background(), nil, cielogo.WithPingHandler(nil), cielogo.WithPongHandler(nil), cielogo.WithCloseHandler(nil))
	require.NoError(t, err)
	defer ws.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	out := make(chan apiv1.WSEvent, 1)
	done := make(chan error, 1)
	go func() { done <- ws.RunListener(ctx, out) }()
	<-sent
	select {
	case <-out:
	case <-time.After(2 * time.Second):
		t.Fatal("first listener did not start")
	}
	require.ErrorIs(t, ws.RunListener(context.Background(), make(chan apiv1.WSEvent)), cielogo.ErrListenerRunning)
	var nilCmd *apiv1.WalletSubscribeCmd
	require.ErrorIs(t, ws.SendCommand(nilCmd), cielogo.ErrNilRequest)
	require.ErrorIs(t, ws.SendCommand(nil), cielogo.ErrNilRequest)
	cancel()
	select {
	case err := <-done:
		require.ErrorIs(t, err, context.Canceled)
	case <-time.After(2 * time.Second):
		t.Fatal("listener did not stop")
	}
}

func TestWebsocketCancellationWhileOutputBlocked(t *testing.T) {
	s := wsTestServer(t, func(conn *websocket.Conn) {
		_ = conn.WriteJSON(map[string]any{"type": "pong", "data": "ok"})
		_, _, _ = conn.ReadMessage()
	})
	ws, err := cielogo.NewClient("key", cielogo.WithBaseURL(s.URL)).NewWebsocketConnection(context.Background())
	require.NoError(t, err)
	defer ws.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan error, 1)
	go func() { done <- ws.RunListener(ctx, make(chan apiv1.WSEvent)) }()
	cancel()
	select {
	case err := <-done:
		require.ErrorIs(t, err, context.Canceled)
	case <-time.After(2 * time.Second):
		t.Fatal("blocked output ignored cancellation")
	}
}

func TestWebsocketTLSAndHandshakeErrors(t *testing.T) {
	tlsServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := (&websocket.Upgrader{}).Upgrade(w, r, nil)
		if err != nil {
			t.Error(err)
			return
		}
		defer conn.Close()
		_ = conn.WriteJSON(map[string]any{"type": "pong", "data": "ok"})
		_ = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "done"))
	}))
	defer tlsServer.Close()
	transport, ok := tlsServer.Client().Transport.(*http.Transport)
	require.True(t, ok)
	dialer := &websocket.Dialer{TLSClientConfig: transport.TLSClientConfig}
	c := cielogo.NewClient("key", cielogo.WithBaseURL(tlsServer.URL), cielogo.WithWebSocketDialer(dialer))
	dialer.HandshakeTimeout = time.Nanosecond // The SDK must have retained its own copy.
	ws, err := c.NewWebsocketConnection(context.Background())
	require.NoError(t, err)
	defer ws.Close()
	out := make(chan apiv1.WSEvent, 1)
	require.NoError(t, ws.RunListener(context.Background(), out))
	assert.Equal(t, apiv1.PongEventType, (<-out).Type)

	denied := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", "10")
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = w.Write([]byte(`{"status":"error","message":"rate limited"}`))
	}))
	defer denied.Close()
	_, err = cielogo.NewClient("key", cielogo.WithBaseURL(denied.URL)).NewWebsocketConnection(context.Background())
	var apiErr *api.Error
	require.ErrorAs(t, err, &apiErr)
	assert.Equal(t, http.StatusTooManyRequests, apiErr.StatusCode)
	assert.Equal(t, "10", apiErr.RetryAfter)
}
