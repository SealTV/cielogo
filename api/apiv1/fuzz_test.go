package apiv1_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/sealtv/cielogo/api/apiv1"
)

func FuzzTxEvent(f *testing.F) {
	for _, seed := range []string{`null`, `{}`, `{"tx_type":"swap","token0_amount":2}`, `{"tx_type":"future","extra":{"x":123}}`, `{"tx_type":"bridge","amount":1,"platform":"bridge"}`, `{"tx_type":"perp_order","oid":9223372036854775807}`} {
		f.Add([]byte(seed))
	}
	f.Fuzz(func(t *testing.T, content []byte) {
		var event apiv1.TxEvent
		if err := json.Unmarshal(content, &event); err != nil {
			return
		}
		if event.TxType == "" || event.Data == nil {
			t.Fatal("accepted event without type/payload")
		}
		encoded, err := json.Marshal(event)
		if err != nil {
			t.Fatalf("cannot re-encode accepted event: %v", err)
		}
		var decoded apiv1.TxEvent
		if err := json.Unmarshal(encoded, &decoded); err != nil {
			t.Fatalf("cannot decode re-encoded event: %v", err)
		}
		if event.TxType != decoded.TxType {
			t.Fatal("event type changed")
		}
		if _, unknown := event.Data.(*apiv1.UnknownTransactionEvent); !unknown && !reflect.DeepEqual(event.Data, decoded.Data) {
			t.Fatal("known transaction payload changed during roundtrip")
		}
	})
}

func FuzzWSEvent(f *testing.F) {
	for _, seed := range []string{`null`, `{}`, `{"type":"tx","data":{"tx_type":"swap","token0_amount":2}}`, `{"type":"future","data":{"new":true}}`, `{"type":"pong","data":"ok"}`, `{"type":"wallet_subscribed","data":{"wallet":"wallet1","filter":{"chains":{"solana":true}}}}`} {
		f.Add([]byte(seed))
	}
	f.Fuzz(func(t *testing.T, content []byte) {
		var event apiv1.WSEvent
		if err := json.Unmarshal(content, &event); err != nil {
			return
		}
		if event.Type == "" {
			t.Fatal("accepted event without type")
		}
		encoded, err := json.Marshal(event)
		if err != nil {
			t.Fatalf("cannot marshal accepted event: %v", err)
		}
		var decoded apiv1.WSEvent
		if err := json.Unmarshal(encoded, &decoded); err != nil {
			t.Fatalf("cannot decode accepted event after roundtrip: %v", err)
		}
		if event.Type != decoded.Type {
			t.Fatal("event type changed")
		}
	})
}

func FuzzPortfolioResponse(f *testing.F) {
	for _, seed := range []string{`null`, `{}`, `{"portfolio":[],"total_usd":123}`, `{"token_address":"SOL","balance":2}`} {
		f.Add([]byte(seed))
	}
	f.Fuzz(func(t *testing.T, content []byte) {
		var reused apiv1.WalletPortfolioV2Response
		_ = json.Unmarshal([]byte(`{"token_address":"old","balance":1}`), &reused)
		if err := json.Unmarshal(content, &reused); err != nil {
			return
		}
		var fresh apiv1.WalletPortfolioV2Response
		if err := json.Unmarshal(content, &fresh); err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(reused, fresh) {
			t.Fatal("stale portfolio fields survived reuse")
		}
	})
}
