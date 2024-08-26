package cielogo

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sealtv/cielogo/api/apiv1"
)

const (
	wsURL = "wss://feed-api.cielo.finance/api/v1/ws"

	pingPeriod     = 30 * time.Second
	pongWait       = 60 * time.Second
	writeWait      = 1 * time.Second
	maxMessageSize = 512 * 1024
)

type WebsocketClient struct {
	conn *websocket.Conn
}

func (c *Client) NewWebsocketConnection(ctx context.Context, opts ...WebsocketOption) (*WebsocketClient, error) {
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, http.Header{
		"X-API-KEY": []string{c.apiKey},
	})
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
	if err := ws.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
		return fmt.Errorf("cannot set write deadline: %w", err)
	}

	if err := ws.conn.WriteJSON(cmd); err != nil {
		return fmt.Errorf("cannot write json event: %w", err)
	}

	return nil
}

func (ws *WebsocketClient) RunListener(ctx context.Context, out chan<- apiv1.WSEvent) error {
	for {
		var event apiv1.WSEvent
		err := ws.conn.ReadJSON(&event)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				return fmt.Errorf("unexpected close error: %w", err)
			}

			if err == websocket.ErrCloseSent {
				return nil
			}
		}

		select {
		case out <- event:
		case <-ctx.Done():
			return nil
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
		ws.conn.SetReadDeadline(t)
	}
}
