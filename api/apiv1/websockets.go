package apiv1

import (
	"encoding/json"
	"fmt"
)

type CommandType string

const (
	WalletSubscribeCommandType   CommandType = "subscribe_wallet"
	WalletUnsubscribeCommandType CommandType = "unsubscribe_wallet"
	FeedSubscribeCommandType     CommandType = "subscribe_feed"
	FeedUnsubscribeCommandType   CommandType = "unsubscribe_feed"
)

type WebSocketsCommand interface {
	json.Marshaler
	GetType() CommandType
}

type WalletSubscribeCmd struct {
	Wallet string  `json:"wallet"`
	Filter *Filter `json:"filter,omitempty"`
}

func (c *WalletSubscribeCmd) GetType() CommandType {
	return WalletSubscribeCommandType
}

func (c *WalletSubscribeCmd) MarshalJSON() ([]byte, error) {
	body := map[string]any{
		"type":   c.GetType(),
		"wallet": c.Wallet,
		"filter": c.Filter,
	}

	return json.Marshal(body)
}

type WalletUnsubscribeCmd struct {
	Wallet string `json:"wallet"`
}

func (c *WalletUnsubscribeCmd) GetType() CommandType {
	return WalletUnsubscribeCommandType
}

func (c *WalletUnsubscribeCmd) MarshalJSON() ([]byte, error) {
	body := map[string]any{
		"type":   c.GetType(),
		"wallet": c.Wallet,
	}

	return json.Marshal(body)
}

type FeedSubscribeCmd struct {
	ListID *int64  `json:"list_id,omitempty"`
	Filter *Filter `json:"filter,omitempty"`
}

func (c *FeedSubscribeCmd) GetType() CommandType {
	return FeedSubscribeCommandType
}

func (c *FeedSubscribeCmd) MarshalJSON() ([]byte, error) {
	body := map[string]any{
		"type":    c.GetType(),
		"list_id": c.ListID,
		"filter":  c.Filter,
	}

	return json.Marshal(body)
}

type FeedUnsubscribeCmd struct{}

func (c *FeedUnsubscribeCmd) GetType() CommandType {
	return FeedUnsubscribeCommandType
}

func (c *FeedUnsubscribeCmd) MarshalJSON() ([]byte, error) {
	body := map[string]any{
		"type": c.GetType(),
	}

	return json.Marshal(body)
}

type Filter struct {
	TxTypes     []TxType `json:"tx_types,omitempty"`
	Chains      []string `json:"chains,omitempty"`
	Tokens      []string `json:"tokens,omitempty"`
	MinUsdValue float64  `json:"min_usd_value,omitempty"`
	NewTrade    bool     `json:"new_trade,omitempty"`
}

type EventType string

const (
	ErrEventType                EventType = "error"
	TxEventType                 EventType = "tx"
	WalletSubscribedEventType   EventType = "wallet_subscribed"
	WalletUnsubscribedEventType EventType = "wallet_unsubscribed"
	FeedSubscribedEventType     EventType = "feed_subscribed"
	FeedUnsubscribedEventType   EventType = "feed_unsubscribed"
)

type WSEvent struct {
	Type EventType `json:"type"`
	Data any       `json:"data"`
}

func (e *WSEvent) UnmarshalJSON(b []byte) error {
	tmp := struct {
		Type EventType       `json:"type"`
		Data json.RawMessage `json:"data"`
	}{}
	if err := json.Unmarshal(b, &tmp); err != nil {
		return fmt.Errorf("cannot unmarshal event: %w", err)
	}

	e.Type = tmp.Type

	var err error
	switch tmp.Type {
	case ErrEventType:
		var data WSEventError
		if err = json.Unmarshal(tmp.Data, &data); err == nil {
			e.Data = data
		}
	case TxEventType:
		var data TxEvent
		if err = json.Unmarshal(tmp.Data, &data); err == nil {
			e.Data = data
		}
	case WalletSubscribedEventType:
		var data WalletSubscribeCmd
		if err = json.Unmarshal(tmp.Data, &data); err == nil {
			e.Data = data
		}
	case WalletUnsubscribedEventType:
		var data WalletUnsubscribeCmd
		if err = json.Unmarshal(tmp.Data, &data); err == nil {
			e.Data = data
		}
	case FeedSubscribedEventType:
		var data FeedSubscribeCmd
		if err = json.Unmarshal(tmp.Data, &data); err == nil {
			e.Data = data
		}
	case FeedUnsubscribedEventType:
		var data FeedUnsubscribeCmd
		if err = json.Unmarshal(tmp.Data, &data); err == nil {
			e.Data = data
		}
	default:
		err = fmt.Errorf("unknown event type: %s", tmp.Type)
	}

	if err != nil {
		return fmt.Errorf("cannot unmarshal event data: %w", err)
	}

	return nil
}

type WSEventError string

func (err WSEventError) Error() string {
	return string(err)
}
