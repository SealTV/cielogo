package cielogo_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/sealtv/cielogo"
	"github.com/sealtv/cielogo/api"
	"github.com/sealtv/cielogo/api/apiv1"
	"github.com/sealtv/cielogo/api/chains"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Fixtures use the active handler/model contract at cielo-api-microservice
// 1b99824e, rather than JSON marshaled from the SDK's own response structs.
func fixtureClient(t *testing.T, file string, inspect func(*http.Request)) *cielogo.Client {
	t.Helper()
	body, err := os.ReadFile("testdata/" + file)
	require.NoError(t, err)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		inspect(r)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(body)
	}))
	t.Cleanup(server.Close)
	return cielogo.NewClient("fixture-key", cielogo.WithBaseURL(server.URL), cielogo.WithPortfolioBaseURL(server.URL))
}

func TestFeedServerContract(t *testing.T) {
	cursor := "YWJjKysvPQ"
	c := fixtureClient(t, "feed.json", func(r *http.Request) {
		assert.Equal(t, "/v1/feed", r.URL.Path)
		assert.Equal(t, "swap", r.URL.Query().Get("txTypes"))
		assert.False(t, r.URL.Query().Has("tx_types"))
		assert.Equal(t, cursor, r.URL.Query().Get("startFrom"))
		assert.Equal(t, "true", r.URL.Query().Get("includeMarketCap"))
	})
	resp, err := c.GetFeedV1(context.Background(), &apiv1.FeedRequest{TxTypes: []apiv1.TxType{apiv1.TxTypeSwap}, StartFrom: &cursor, IncludeMarketCap: apiv1.ToRef(true)})
	require.NoError(t, err)
	require.Len(t, resp.Items, 6)
	assert.Equal(t, cursor, resp.Paging.NextCursor)
	swap, ok := resp.Items[0].Data.(*apiv1.SwapEvent)
	require.True(t, ok)
	assert.Equal(t, "SOL", swap.Token0Address)
	assert.Equal(t, 2.0, swap.Token0Amount)
	assert.Equal(t, "USDC", swap.Token1Address)
	assert.True(t, swap.IsSell)
	require.NotNil(t, resp.Items[0].TokenMarketCap)
	assert.Equal(t, 123.0, resp.Items[0].TokenMarketCap.MarketCap)
	assert.JSONEq(t, `{"name":"verified"}`, string(resp.Items[0].Verification))
	assert.IsType(t, &apiv1.PerpOrderEvent{}, resp.Items[1].Data)
	assert.IsType(t, &apiv1.PerpTwapEvent{}, resp.Items[2].Data)
	assert.IsType(t, &apiv1.PerpImminentLiqEvent{}, resp.Items[3].Data)
	assert.IsType(t, &apiv1.HyperliquidLedgerEvent{}, resp.Items[4].Data)
	unknown, ok := resp.Items[5].Data.(*apiv1.UnknownTransactionEvent)
	require.True(t, ok)
	assert.Equal(t, apiv1.TxType("future_event"), unknown.GetType())
	assert.JSONEq(t, `{"tx_type":"future_event","custom":123}`, string(unknown.Raw))
}

func TestAnalyticsServerContracts(t *testing.T) {
	t.Run("trading", func(t *testing.T) {
		c := fixtureClient(t, "trading_stats.json", func(r *http.Request) { assert.Equal(t, "/v1/wallet1/trading-stats", r.URL.Path) })
		resp, err := c.GetTradingStatsV1(context.Background(), &apiv1.TradingStatsRequest{Wallet: "wallet1"})
		require.NoError(t, err)
		assert.Equal(t, 123.5, resp.PnL)
		assert.Equal(t, 0.8, resp.WinRate)
		assert.Equal(t, int64(10), resp.SwapsCount)
		assert.Equal(t, int64(2), resp.HoldingDistribution.HoldGT24Hours)
	})
	t.Run("portfolio", func(t *testing.T) {
		c := fixtureClient(t, "portfolio.json", func(r *http.Request) { assert.Equal(t, "/v1/wallet1/portfolio", r.URL.Path) })
		resp, err := c.GetWalletPortfolioV1(context.Background(), "wallet1")
		require.NoError(t, err)
		assert.Equal(t, 1234.56, resp.TotalUSDValue)
		assert.Equal(t, []string{"ethereum"}, resp.PendingChains)
		require.Len(t, resp.Chains, 1)
		assert.Equal(t, "2000000000", resp.Portfolio[0].BalanceRaw)
	})
	t.Run("nft", func(t *testing.T) {
		c := fixtureClient(t, "nft_pnl.json", func(r *http.Request) { assert.Equal(t, "/v1/wallet1/pnl/nfts", r.URL.Path) })
		resp, err := c.GetNftsPnlV1(context.Background(), &apiv1.NftsPnLRequest{Wallet: "wallet1"})
		require.NoError(t, err)
		require.Len(t, resp.Items, 1)
		assert.Equal(t, "Collection", resp.Items[0].Name)
		assert.Equal(t, "NFT", resp.Items[0].Symbol)
		assert.Equal(t, 2.0, resp.Items[0].TotalBuyAvg)
		assert.Equal(t, 3.0, resp.Items[0].TotalSellAvg)
		assert.Equal(t, 5.0, resp.Items[0].ProfitLoss)
	})
	t.Run("tags", func(t *testing.T) {
		c := fixtureClient(t, "wallets_by_tag.json", func(r *http.Request) { assert.Equal(t, "/v1/tags/wallets", r.URL.Path) })
		resp, err := c.GetWalletsByTagV1(context.Background(), &apiv1.GetWalletsByTagRequest{Tags: []apiv1.TagType{apiv1.TagTypeGemFinder}})
		require.NoError(t, err)
		require.Len(t, resp.Wallets, 1)
		assert.Equal(t, "wallet1", resp.Wallets[0].Wallet)
		assert.Equal(t, 2, resp.Paging.NextPage)
	})
}

func TestTrackedWalletPaginationContract(t *testing.T) {
	var pages []int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v1/tracked-wallets/", r.URL.Path)
		assert.Equal(t, "17", r.URL.Query().Get("list_id"))
		page := 1
		if raw := r.URL.Query().Get("next_object"); raw != "" {
			parsed, err := strconv.Atoi(raw)
			assert.NoError(t, err)
			page = parsed
		}
		pages = append(pages, page)
		_ = json.NewEncoder(w).Encode(map[string]any{"status": "ok", "data": map[string]any{"tracked_wallets": []any{map[string]any{"id": page, "wallet": "wallet" + strconv.Itoa(page)}}, "paging": map[string]any{"next_object": page + 1, "has_next_page": page < 3}}})
	}))
	defer server.Close()
	c := cielogo.NewClient("key", cielogo.WithBaseURL(server.URL))
	req := &apiv1.GetTrackedWalletsRequest{ListID: apiv1.ToRef(int64(17))}
	for i := 1; i <= 3; i++ {
		resp, err := c.GetTrackedWalletsV1(context.Background(), req)
		require.NoError(t, err)
		require.Len(t, resp.TrackedWallets, 1)
		assert.Equal(t, int64(i), resp.TrackedWallets[0].ID)
		assert.Equal(t, i < 3, resp.Pagination.HasNextPage)
		cursor := strconv.Itoa(resp.Pagination.NextObject)
		req.NextObject = &cursor
	}
	assert.Equal(t, []int{1, 2, 3}, pages)
}

func TestManagementWireContract(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/lists/all":
			assert.Equal(t, "true", r.URL.Query().Get("following_only"))
			assert.False(t, r.URL.Query().Has("follow_only"))
			_, _ = w.Write([]byte(`{"status":"ok","data":{"list":[],"paging":{"next_object":"abc"}}}`))
		case "/v1/lists/":
			var body struct {
				Wallets []int64 `json:"wallets"`
			}
			assert.NoError(t, json.NewDecoder(r.Body).Decode(&body))
			assert.Equal(t, []int64{17, 18}, body.Wallets)
			_, _ = w.Write([]byte(`{"status":"ok","data":{"id":17}}`))
		case "/v1/lists/17/toggle-follow":
			_, _ = w.Write([]byte(`{"status":"ok","data":{"id":17,"followed":true}}`))
		case "/v2/tracked-wallets/wallet1":
			var body map[string]any
			assert.NoError(t, json.NewDecoder(r.Body).Decode(&body))
			assert.Equal(t, float64(100), body["min_amount_usd"])
			assert.Equal(t, []any{float64(1), float64(2)}, body["filters"])
			assert.Equal(t, []any{float64(1)}, body["chains"])
			assert.Equal(t, float64(17), body["telegram_bot_id"])
			assert.Equal(t, "channel", body["discord_channel_id"])
			assert.Equal(t, false, body["new_trades"])
			for _, key := range []string{"min_usd", "tx_types", "telegram_bot", "discord_channel"} {
				assert.NotContains(t, body, key)
			}
			_, _ = w.Write([]byte(`{"status":"ok","data":{"id":17,"wallet":"wallet1"}}`))
		case "/v1/wallet1/pnl/tokens":
			assert.Equal(t, "solana", r.URL.Query().Get("chain"))
			assert.False(t, r.URL.Query().Has("chains"))
			_, _ = w.Write([]byte(`{"status":"ok","data":{"items":[]}}`))
		default:
			t.Errorf("unexpected path: %s", r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()
	c := cielogo.NewClient("key", cielogo.WithBaseURL(server.URL))
	_, err := c.GetAllWalletsListV1(context.Background(), &apiv1.GetAllWalletsListsRequest{FollowOnly: true})
	require.NoError(t, err)
	_, err = c.AddWalletsListV1(context.Background(), &apiv1.AddWalletsListRequest{Name: "list", Wallets: []int64{17, 18}})
	require.NoError(t, err)
	ack, err := c.ToggleFollowWalletsListV1(context.Background(), 17)
	require.NoError(t, err)
	assert.True(t, ack.Followed)
	_, err = c.UpdateTrackedWalletV2(context.Background(), "wallet1", &apiv1.UpdateTrackedWalletV2Request{MinAmountUSD: apiv1.ToRef(100.0), Filters: []int{1, 2}, Chains: []int{1}, TelegramBotID: apiv1.ToRef(17), DiscordChannelID: apiv1.ToRef("channel"), NewTrades: apiv1.ToRef(false)})
	require.NoError(t, err)
	_, err = c.GetTokensPnlV1(context.Background(), &apiv1.TokensPnLRequest{Wallet: "wallet1", Chains: []chains.ChainType{chains.Solana}})
	require.NoError(t, err)
}

func TestPendingHTTPContract(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", "10")
		w.WriteHeader(http.StatusAccepted)
		_, _ = w.Write([]byte(`{"status":"pending","message":"Data is not ready yet"}`))
	}))
	defer server.Close()
	c := cielogo.NewClient("key", cielogo.WithBaseURL(server.URL))
	_, err := c.GetTradingStatsV1(context.Background(), &apiv1.TradingStatsRequest{Wallet: "wallet1"})
	var apiErr *api.Error
	require.True(t, errors.As(err, &apiErr))
	assert.Equal(t, 202, apiErr.StatusCode)
	assert.Equal(t, "pending", apiErr.Status)
	assert.Equal(t, "10", apiErr.RetryAfter)
	assert.Equal(t, 10*time.Second, apiErr.RetryAfterDuration)
}

func TestPrivatePortfolioContract(t *testing.T) {
	_, err := cielogo.NewClient("key").GetWalletPortfolioV2(context.Background(), &apiv1.WalletPortfolioV2Request{Wallets: []string{"wallet1"}})
	require.ErrorContains(t, err, "WithPortfolioBaseURL")
	t.Run("full", func(t *testing.T) {
		c := fixtureClient(t, "portfolio.json", func(r *http.Request) {
			assert.Equal(t, "/v2/portfolio", r.URL.Path)
			assert.Equal(t, "wallet1,wallet2", r.URL.Query().Get("wallet"))
		})
		resp, err := c.GetWalletPortfolioV2(context.Background(), &apiv1.WalletPortfolioV2Request{Wallets: []string{"wallet1", "wallet2"}})
		require.NoError(t, err)
		assert.Equal(t, 1234.56, resp.TotalUSDValue)
		assert.Nil(t, resp.SingleToken)
	})
	t.Run("single token", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "SOL", r.URL.Query().Get("token"))
			_, _ = w.Write([]byte(`{"status":"ok","data":{"token_address":"SOL","balance":2,"total_usd_value":300}}`))
		}))
		defer server.Close()
		c := cielogo.NewClient("key", cielogo.WithPortfolioBaseURL(server.URL))
		resp, err := c.GetWalletPortfolioV2(context.Background(), &apiv1.WalletPortfolioV2Request{Wallets: []string{"wallet1"}, Token: apiv1.ToRef("SOL")})
		require.NoError(t, err)
		require.NotNil(t, resp.SingleToken)
		assert.Equal(t, 300.0, resp.SingleToken.TotalUSDValue)
		assert.Empty(t, resp.Portfolio)
	})
}

func TestMetadataUint64Contract(t *testing.T) {
	var metadata apiv1.TokenMetadataResponse
	require.NoError(t, json.Unmarshal([]byte(`{"supply":18446744073709551615,"created_at_block":18446744073709551615,"total_supply_string":"18446744073709551615"}`), &metadata))
	assert.Equal(t, ^uint64(0), metadata.Supply)
	assert.Equal(t, ^uint64(0), metadata.CreatedAtBlock)
	assert.Equal(t, "18446744073709551615", metadata.TotalSupplyString)
}

func TestTrackedWalletEmptyCursor(t *testing.T) {
	for _, cursor := range []*string{nil, apiv1.ToRef("")} {
		t.Run("empty", func(t *testing.T) {
			s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Empty(t, r.URL.RawQuery)
				_, _ = w.Write([]byte(`{"status":"ok","data":{"tracked_wallets":[]}}`))
			}))
			defer s.Close()
			_, err := cielogo.NewClient("key", cielogo.WithBaseURL(s.URL)).GetTrackedWalletsV1(context.Background(), &apiv1.GetTrackedWalletsRequest{NextObject: cursor})
			require.NoError(t, err)
		})
	}
}

func TestV2ClearNotificationFilters(t *testing.T) {
	body, err := json.Marshal(apiv1.UpdateTrackedWalletV2Request{Filters: []int{}, Chains: []int{}})
	require.NoError(t, err)
	assert.JSONEq(t, `{"filters":[],"chains":[]}`, string(body))
	body, err = json.Marshal(apiv1.UpdateTrackedWalletV2Request{})
	require.NoError(t, err)
	assert.JSONEq(t, `{}`, string(body))
}

func TestHTTPStatusContracts(t *testing.T) {
	for _, tt := range []struct {
		name        string
		status      int
		body, retry string
	}{
		{"created", http.StatusCreated, `{"status":"ok","data":{"id":17}}`, ""},
		{"empty delete", http.StatusNoContent, "", ""},
		{"non-json rate limit", http.StatusTooManyRequests, "upstream busy", "10"},
		{"pending date", http.StatusAccepted, `{"status":"pending"}`, time.Now().Add(time.Hour).UTC().Format(http.TimeFormat)},
	} {
		t.Run(tt.name, func(t *testing.T) {
			s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Retry-After", tt.retry)
				w.WriteHeader(tt.status)
				_, _ = w.Write([]byte(tt.body))
			}))
			defer s.Close()
			c := cielogo.NewClient("key", cielogo.WithBaseURL(s.URL))
			if tt.status == http.StatusCreated {
				resp, err := c.AddWalletsListV1(context.Background(), &apiv1.AddWalletsListRequest{Name: "list"})
				require.NoError(t, err)
				assert.Equal(t, int64(17), resp.ID)
				return
			}
			err := c.DeleteWalletsListV1(context.Background(), 17, false)
			if tt.status == http.StatusNoContent {
				require.NoError(t, err)
				return
			}
			var apiErr *api.Error
			require.ErrorAs(t, err, &apiErr)
			assert.Equal(t, tt.status, apiErr.StatusCode)
			assert.Equal(t, tt.retry, apiErr.RetryAfter)
			assert.Positive(t, apiErr.RetryAfterDuration)
		})
	}
}

func TestChainDiscoveryAndTokenStatsOptions(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/chains" {
			_, _ = w.Write([]byte(`{"status":"ok","data":["solana","arc","future_chain"]}`))
			return
		}
		assert.Equal(t, "/v1/token/stats", r.URL.Path)
		assert.Equal(t, "true", r.URL.Query().Get("include_pools"))
		assert.Equal(t, "false", r.URL.Query().Get("include_liquidity"))
		_, _ = w.Write([]byte(`{"status":"ok","data":{"ath_price_usd":123,"pools":[{"pool_address":"pool1"}],"active_users":17}}`))
	}))
	defer s.Close()
	c := cielogo.NewClient("key", cielogo.WithBaseURL(s.URL))
	chainsResult, err := c.GetChainsV1(context.Background())
	require.NoError(t, err)
	assert.Equal(t, []chains.ChainType{"solana", "arc", "future_chain"}, chainsResult)
	stats, err := c.GetTokenStatsV1(context.Background(), &apiv1.TokenStatsRequest{Chain: apiv1.TokenChainArc, TokenAddress: "token1", IncludePools: apiv1.ToRef(true), IncludeLiquidity: apiv1.ToRef(false)})
	require.NoError(t, err)
	assert.Equal(t, 123.0, stats.AthPriceUSD)
	require.Len(t, stats.Pools, 1)
	assert.Equal(t, "pool1", stats.Pools[0].PoolAddress)
	assert.Equal(t, uint64(17), stats.ActiveUsers)
}
