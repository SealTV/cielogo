package cielogo

import (
	"context"
	"errors"
	"fmt"
	"net/http"
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

type WebsocketClient struct {
	conn    *websocket.Conn
	writeMu sync.Mutex
}

func (c *Client) NewWebsocketConnection(ctx context.Context, opts ...WebsocketOption) (*WebsocketClient, error) {
	wsURL := strings.TrimRight(c.baseURL, "/") + "/v1/ws"
	wsURL = strings.Replace(wsURL, "https://", "wss://", 1)
	wsURL = strings.Replace(wsURL, "http://", "ws://", 1)
	conn, resp, err := websocket.DefaultDialer.DialContext(ctx, wsURL, http.Header{
		"X-API-KEY": []string{c.apiKey},
	})
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, fmt.Errorf("failed to open websocket connection: %w", err)
	}

	conn.SetReadLimit(maxMessageSize)

	ws := &WebsocketClient{
		conn: conn,
	}

	for _, opt := range opts {
		opt(ws)
	}

	return ws, nil
}

func (ws *WebsocketClient) Close() {
	ws.conn.Close()
}

func (ws *WebsocketClient) SendCommand(cmd apiv1.WebSocketsCommand) error {
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
		ws.conn.SetCloseHandler(func(code int, text string) error {
			return h(code, text)
		})
	}
}

func WithPingHandler(h func(appData string) error) WebsocketOption {
	return func(ws *WebsocketClient) {
		ws.conn.SetPingHandler(func(appData string) error {
			return h(appData)
		})
	}
}

func WithPongHandler(h func(appData string) error) WebsocketOption {
	return func(ws *WebsocketClient) {
		ws.conn.SetPongHandler(func(appData string) error {
			return h(appData)
		})
	}
}

func WithDeadline(t time.Time) WebsocketOption {
	return func(ws *WebsocketClient) {
		_ = ws.conn.SetReadDeadline(t) // Ignore error as it's used in option pattern
	}
}
