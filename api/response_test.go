package api_test

import (
	"encoding/json"
	"testing"

	"github.com/sealtv/cielogo/api"
)

func FuzzResponseEnvelope(f *testing.F) {
	for _, seed := range []string{`null`, `{}`, `{"status":"ok","data":{}}`, `{"status":"ok","data":null}`, `{"status":"pending","message":"retry"}`} {
		f.Add([]byte(seed))
	}
	f.Fuzz(func(t *testing.T, content []byte) {
		var response api.CieloResponse[struct {
			Value string `json:"value"`
		}]
		if err := json.Unmarshal(content, &response); err != nil {
			return
		}
		if response.Status == "ok" {
			var raw struct {
				Data json.RawMessage `json:"data"`
			}
			_ = json.Unmarshal(content, &raw)
			if len(raw.Data) == 0 || string(raw.Data) == "null" {
				t.Fatal("accepted successful response without object data")
			}
		}
	})
}
