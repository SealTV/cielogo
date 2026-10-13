package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

// ErrInvalidResponse identifies malformed or incomplete API envelopes.
var ErrInvalidResponse = errors.New("invalid API response")

type CieloResponse[T any] struct {
	Status  string `json:"status"`
	Data    T      `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
}

// UnmarshalJSON rejects missing object data in a successful response. Null/absent
// slice data is accepted as an empty collection, as returned by some list endpoints.
func (r *CieloResponse[T]) UnmarshalJSON(content []byte) error {
	var envelope struct {
		Status  string          `json:"status"`
		Data    json.RawMessage `json:"data"`
		Message string          `json:"message"`
	}
	if err := json.Unmarshal(content, &envelope); err != nil {
		return fmt.Errorf("decode envelope: %w", err)
	}
	kind := reflect.TypeOf((*T)(nil)).Elem().Kind()
	if envelope.Status == "ok" && (kind == reflect.Struct || kind == reflect.Pointer) && (len(envelope.Data) == 0 || string(envelope.Data) == "null") {
		return fmt.Errorf("%w: missing response data", ErrInvalidResponse)
	}
	var data T
	if len(envelope.Data) > 0 {
		if err := json.Unmarshal(envelope.Data, &data); err != nil {
			return fmt.Errorf("decode response data: %w", err)
		}
	}
	*r = CieloResponse[T]{Status: envelope.Status, Message: envelope.Message, Data: data}
	return nil
}
