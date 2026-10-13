package cielogo

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sealtv/cielogo/api/apiv1"
)

const (
	writeWait      = 1 * time.Second
	maxMessageSize = 512 * 1024
)

// WithWebSocketDialer configures TLS, proxy and handshake settings for WebSocket
// connections. A copy is retained; referenced TLS configurations must be immutable.
func WithWebSocketDialer(dialer *websocket.Dialer) ClientOption {
	return func(c *Client) {
		if dialer != nil {
			copied := *dialer
			c.websocketDialer = &copied
		}
	}
}

type WebsocketClient struct {
	conn       *websocket.Conn
	writeMu    sync.Mutex
	listenerMu sync.Mutex
}

func (c *Client) NewWebsocketConnection(ctx context.Context, opts ...WebsocketOption) (*WebsocketClient, error) {
	if c == nil || c.cli == nil {
		return nil, ErrInvalidClient
	}
	if ctx == nil {
		return nil, ErrNilContext
	}
	wsURL := strings.TrimRight(c.baseURL, "/") + "/v1/ws"
	wsURL = strings.Replace(wsURL, "https://", "wss://", 1)
	wsURL = strings.Replace(wsURL, "http://", "ws://", 1)
	dialer := c.websocketDialer
	if dialer == nil {
		dialer = websocket.DefaultDialer
	}
	conn, resp, err := dialer.DialContext(ctx, wsURL, http.Header{
		"X-Api-Key": []string{c.apiKey},
	})
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		if resp != nil && resp.Body != nil {
			if apiErr := c.decodeResponse(resp, nil); apiErr != nil {
				return nil, fmt.Errorf("failed to open websocket connection: %w", apiErr)
			}
		}
		return nil, fmt.Errorf("failed to open websocket connection: %w", err)
	}

	conn.SetReadLimit(maxMessageSize)

	ws := &WebsocketClient{
		conn: conn,
	}

	for _, opt := range opts {
		if opt != nil {
			opt(ws)
		}
	}

	return ws, nil
}

func (ws *WebsocketClient) Close() {
	if ws != nil && ws.conn != nil {
		_ = ws.conn.Close()
	}
}

func (ws *WebsocketClient) SendCommand(cmd apiv1.WebSocketsCommand) error {
	if ws == nil || ws.conn == nil {
		return ErrInvalidClient
	}
	if cmd == nil || (reflect.ValueOf(cmd).Kind() == reflect.Pointer && reflect.ValueOf(cmd).IsNil()) {
		return ErrNilRequest
	}
	ws.writeMu.Lock()
	defer ws.writeMu.Unlock()
	if err := ws.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
		return fmt.Errorf("cannot set write deadline: %w", err)
	}

	if err := ws.conn.WriteJSON(cmd); err != nil {
		return fmt.Errorf("cannot write json event: %w", err)
	}

	return nil
}

// RunListener reads events until cancellation, normal closure or a read/decode error.
// Cancellation closes the connection to unblock a pending read. Only one listener
// may run per connection; SendCommand is safe for concurrent callers.
func (ws *WebsocketClient) RunListener(ctx context.Context, out chan<- apiv1.WSEvent) error {
	if ws == nil || ws.conn == nil {
		return ErrInvalidClient
	}
	if ctx == nil {
		return ErrNilContext
	}
	if !ws.listenerMu.TryLock() {
		return ErrListenerRunning
	}
	defer ws.listenerMu.Unlock()
	if err := ctx.Err(); err != nil {
		ws.Close()
		return err
	}
	done := make(chan struct{})
	defer close(done)
	go func() {
		select {
		case <-ctx.Done():
			ws.Close()
		case <-done:
		}
	}()
	for {
		var event apiv1.WSEvent
		if err := ws.conn.ReadJSON(&event); err != nil {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) || errors.Is(err, websocket.ErrCloseSent) {
				return nil
			}
			return fmt.Errorf("cannot read websocket event: %w", err)
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case out <- event:
		}
	}
}

type WebsocketOption func(*WebsocketClient)

func WithCloseHandler(h func(code int, text string) error) WebsocketOption {
	return func(ws *WebsocketClient) {
		ws.conn.SetCloseHandler(h)
	}
}

func WithPingHandler(h func(appData string) error) WebsocketOption {
	return func(ws *WebsocketClient) {
		ws.conn.SetPingHandler(h)
	}
}

func WithPongHandler(h func(appData string) error) WebsocketOption {
	return func(ws *WebsocketClient) {
		ws.conn.SetPongHandler(h)
	}
}

func WithDeadline(t time.Time) WebsocketOption {
	return func(ws *WebsocketClient) {
		_ = ws.conn.SetReadDeadline(t) // Ignore error as it's used in option pattern
	}
}
