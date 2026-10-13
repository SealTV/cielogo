package cielogo_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sealtv/cielogo"
	"github.com/sealtv/cielogo/api/apiv1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// A literal method/route/body matrix complements field-specific fixture tests.
// Every existing HTTP operation must reach the server with its registered verb.
func TestAllHTTPMethodContracts(t *testing.T) {
	ctx := context.Background()
	cases := []struct {
		name, method, path, body string
		collection               bool
		call                     func(*cielogo.Client) error
	}{
		{"feed", http.MethodGet, "/v1/feed", "", false, func(c *cielogo.Client) error { _, e := c.GetFeedV1(ctx, &apiv1.FeedRequest{}); return e }},
		{"NFT pnl", http.MethodGet, "/v1/wallet1/pnl/nfts", "", false, func(c *cielogo.Client) error {
			_, e := c.GetNftsPnlV1(ctx, &apiv1.NftsPnLRequest{Wallet: "wallet1"})
			return e
		}},
		{"token pnl", http.MethodGet, "/v1/wallet1/pnl/tokens", "", false, func(c *cielogo.Client) error {
			_, e := c.GetTokensPnlV1(ctx, &apiv1.TokensPnLRequest{Wallet: "wallet1"})
			return e
		}},
		{"aggregate pnl", http.MethodGet, "/v1/wallet1/pnl/total-stats", "", false, func(c *cielogo.Client) error {
			_, e := c.GetAggregatedTokenPnLV1(ctx, &apiv1.AggregatedTokenPnLRequest{Wallet: "wallet1"})
			return e
		}},
		{"related", http.MethodGet, "/v1/wallet1/related-wallets", "", false, func(c *cielogo.Client) error {
			_, e := c.GetRelatedWalletsV1(ctx, &apiv1.RelatedWalletsRequest{Wallet: "wallet1"})
			return e
		}},
		{"tags batch", http.MethodGet, "/v1/tags", "", true, func(c *cielogo.Client) error {
			_, e := c.GetWalletsTagsV1(ctx, &apiv1.GetWalletsTagsRequest{Wallets: []string{"wallet1"}})
			return e
		}},
		{"tags search", http.MethodGet, "/v1/tags/wallets", "", false, func(c *cielogo.Client) error {
			_, e := c.GetWalletsByTagV1(ctx, &apiv1.GetWalletsByTagRequest{Tags: []apiv1.TagType{apiv1.TagTypeGemFinder}})
			return e
		}},
		{"all lists", http.MethodGet, "/v1/lists/all", "", false, func(c *cielogo.Client) error {
			_, e := c.GetAllWalletsListV1(ctx, &apiv1.GetAllWalletsListsRequest{})
			return e
		}},
		{"user lists", http.MethodGet, "/v1/lists/", "", true, func(c *cielogo.Client) error { _, e := c.GetUserWalletsListsV1(ctx); return e }},
		{"add list", http.MethodPost, "/v1/lists/", `{"name":"list","wallets":[17]}`, false, func(c *cielogo.Client) error {
			_, e := c.AddWalletsListV1(ctx, &apiv1.AddWalletsListRequest{Name: "list", Wallets: []int64{17}})
			return e
		}},
		{"update list", http.MethodPut, "/v1/lists/17", `{"name":"list","wallets":[17]}`, false, func(c *cielogo.Client) error {
			_, e := c.UpdateWalletsListV1(ctx, &apiv1.UpdateWalletsListRequest{ListID: 17, Name: "list", Wallets: []int64{17}})
			return e
		}},
		{"delete list", http.MethodDelete, "/v1/lists/17", "", false, func(c *cielogo.Client) error { return c.DeleteWalletsListV1(ctx, 17, true) }},
		{"toggle follow", http.MethodPut, "/v1/lists/17/toggle-follow", "", false, func(c *cielogo.Client) error { _, e := c.ToggleFollowWalletsListV1(ctx, 17); return e }},
		{"tracked", http.MethodGet, "/v1/tracked-wallets/", "", false, func(c *cielogo.Client) error {
			_, e := c.GetTrackedWalletsV1(ctx, &apiv1.GetTrackedWalletsRequest{})
			return e
		}},
		{"add tracked", http.MethodPost, "/v1/tracked-wallets/", `{"wallet":"wallet1","label":"Wallet"}`, false, func(c *cielogo.Client) error {
			_, e := c.AddTrackedWalletsV1(ctx, &apiv1.AddTrackedWalletRequest{Wallet: "wallet1", Label: "Wallet"})
			return e
		}},
		{"remove tracked", http.MethodDelete, "/v1/tracked-wallets/", `{"wallet_ids":[17]}`, false, func(c *cielogo.Client) error {
			return c.RemoveTrackedWalletsV1(ctx, &apiv1.RemoveTrackedWalletsRequest{WalletIDs: []int64{17}})
		}},
		{"tracked address", http.MethodGet, "/v1/tracked-wallets/address/wallet1", "", false, func(c *cielogo.Client) error { _, e := c.GetWalletByAddressV1(ctx, "wallet1"); return e }},
		{"update tracked", http.MethodPut, "/v1/tracked-wallets/17", `{"wallet":"wallet1","label":"Wallet"}`, false, func(c *cielogo.Client) error {
			_, e := c.UpdateTrackedWalletV1(ctx, 17, &apiv1.UpdateTrackedWalletRequest{Wallet: "wallet1", Label: "Wallet"})
			return e
		}},
		{"update tracked V2", http.MethodPut, "/v2/tracked-wallets/wallet1", `{"label":"Wallet"}`, false, func(c *cielogo.Client) error {
			_, e := c.UpdateTrackedWalletV2(ctx, "wallet1", &apiv1.UpdateTrackedWalletV2Request{Label: apiv1.ToRef("Wallet")})
			return e
		}},
		{"telegram bots", http.MethodGet, "/v1/tracked-wallets/telegram-bots", "", false, func(c *cielogo.Client) error { _, e := c.GetTelegramBotsV1(ctx); return e }},
		{"portfolio", http.MethodGet, "/v1/wallet1/portfolio", "", false, func(c *cielogo.Client) error { _, e := c.GetWalletPortfolioV1(ctx, "wallet1"); return e }},
		{"private portfolio", http.MethodGet, "/v2/portfolio", "", false, func(c *cielogo.Client) error {
			_, e := c.GetWalletPortfolioV2(ctx, &apiv1.WalletPortfolioV2Request{Wallets: []string{"wallet1"}})
			return e
		}},
		{"metadata", http.MethodGet, "/v1/token/metadata", "", false, func(c *cielogo.Client) error {
			_, e := c.GetTokenMetadataV1(ctx, &apiv1.TokenMetadataRequest{Chain: apiv1.TokenChainEthereum, TokenAddress: "token1"})
			return e
		}},
		{"price", http.MethodGet, "/v1/token/price", "", false, func(c *cielogo.Client) error {
			_, e := c.GetTokenPriceV1(ctx, &apiv1.TokenPriceRequest{Chain: apiv1.TokenChainEthereum, TokenAddress: "token1"})
			return e
		}},
		{"stats", http.MethodGet, "/v1/token/stats", "", false, func(c *cielogo.Client) error {
			_, e := c.GetTokenStatsV1(ctx, &apiv1.TokenStatsRequest{Chain: apiv1.TokenChainEthereum, TokenAddress: "token1"})
			return e
		}},
		{"balance", http.MethodGet, "/v1/wallet1/token-balance", "", false, func(c *cielogo.Client) error {
			_, e := c.GetTokenBalanceV1(ctx, &apiv1.TokenBalanceRequest{Wallet: "wallet1", Chain: apiv1.TokenChainEthereum, TokenAddress: "token1"})
			return e
		}},
		{"trading", http.MethodGet, "/v1/wallet1/trading-stats", "", false, func(c *cielogo.Client) error {
			_, e := c.GetTradingStatsV1(ctx, &apiv1.TradingStatsRequest{Wallet: "wallet1"})
			return e
		}},
		{"chains", http.MethodGet, "/v1/chains", "", true, func(c *cielogo.Client) error { _, e := c.GetChainsV1(ctx); return e }},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tt.method, r.Method)
				assert.Equal(t, tt.path, r.URL.Path)
				assert.Equal(t, "key", r.Header.Get("X-Api-Key"))
				assert.Equal(t, "application/json", r.Header.Get("Accept"))
				body, err := io.ReadAll(r.Body)
				assert.NoError(t, err)
				if tt.body != "" {
					assert.JSONEq(t, tt.body, string(body))
				} else {
					assert.Empty(t, body)
				}
				if tt.name == "delete list" {
					assert.Equal(t, "true", r.URL.Query().Get("delete_wallets"))
				}
				data := any(map[string]any{})
				if tt.collection {
					data = []any{}
				}
				_ = json.NewEncoder(w).Encode(map[string]any{"status": "ok", "data": data})
			}))
			defer s.Close()
			c := cielogo.NewClient("key", cielogo.WithBaseURL(s.URL), cielogo.WithPortfolioBaseURL(s.URL))
			require.NoError(t, tt.call(c))
		})
	}
}
