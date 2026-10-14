//go:build integration
// +build integration

package cielogo_test

import (
	"context"
	"crypto/tls"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/sealtv/cielogo"
	"github.com/sealtv/cielogo/api"
	"github.com/sealtv/cielogo/api/apiv1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	// Test wallet addresses - these should be real wallets for integration testing
	testSolanaWallet   = "GJRhJvmjKBg9DTz1YEFkgPTF1hphCMPZQsJdBKxSLqQZ"
	testEthereumWallet = "0x8a90cab2b38dba80c64b7734e58ee1db38b8992e"

	// Test token addresses
	testSolanaToken   = "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v" // USDC on Solana
	testEthereumToken = "0xdAC17F958D2ee523a2206206994597C13D831ec7"   // USDT on Ethereum
)

func getTestClient(t *testing.T) *cielogo.Client {
	t.Helper()

	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	apiKey := os.Getenv("CIELO_API_KEY")
	if apiKey == "" {
		t.Skip("CIELO_API_KEY not set, skipping integration test")
	}

	// net/http can retry GETs after failures on reused connections, and HTTP/2
	// has additional retry paths. Paid smoke calls use fresh HTTP/1 connections
	// so the explicit credit reservation bounds actual request attempts.
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.DisableKeepAlives = true
	transport.ForceAttemptHTTP2 = false
	transport.TLSNextProto = map[string]func(string, *tls.Conn) http.RoundTripper{}
	transport.TLSClientConfig = &tls.Config{MinVersion: tls.VersionTLS12, NextProtos: []string{"http/1.1"}}
	client := cielogo.NewClient(apiKey,
		cielogo.WithHTTPClient(&http.Client{Transport: transport, Timeout: 90 * time.Second}),
		cielogo.WithPortfolioBaseURL(os.Getenv("CIELO_PORTFOLIO_BASE_URL")),
	)
	t.Cleanup(client.CloseIdleConnections)
	return client
}

func TestIntegration_GetWalletPortfolioV1(t *testing.T) {
	client := getTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	resp, err := client.GetWalletPortfolioV1(ctx, testSolanaWallet)
	require.NoError(t, err, "GetWalletPortfolioV1 should not return an error")
	require.NotNil(t, resp, "Response should not be nil")

	assert.Greater(t, resp.LastUpdated, int64(0), "LastUpdated should be greater than 0")
}

func TestIntegration_GetWalletPortfolioV2(t *testing.T) {
	if os.Getenv("CIELO_PORTFOLIO_BASE_URL") == "" {
		t.Skip("private portfolio base URL not configured")
	}
	client := getTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	req := &apiv1.WalletPortfolioV2Request{
		Wallets: []string{testSolanaWallet},
	}

	resp, err := client.GetWalletPortfolioV2(ctx, req)
	require.NoError(t, err, "GetWalletPortfolioV2 should not return an error")
	require.NotNil(t, resp, "Response should not be nil")

	assert.Greater(t, resp.LastUpdated, int64(0), "LastUpdated should be greater than 0")
	assert.NotNil(t, resp.Portfolio, "Portfolio should not be nil")
}

func TestIntegration_GetWalletPortfolioV2_MultipleWallets(t *testing.T) {
	if os.Getenv("CIELO_PORTFOLIO_BASE_URL") == "" {
		t.Skip("private portfolio base URL not configured")
	}
	client := getTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	req := &apiv1.WalletPortfolioV2Request{
		Wallets: []string{testSolanaWallet, testEthereumWallet},
	}

	resp, err := client.GetWalletPortfolioV2(ctx, req)
	require.NoError(t, err, "GetWalletPortfolioV2 with multiple wallets should not return an error")
	require.NotNil(t, resp, "Response should not be nil")

	assert.Greater(t, resp.TotalUSDValue, 0.0, "TotalUSDValue should be greater than 0")
	if len(resp.Chains) > 0 {
		assert.NotEmpty(t, resp.Chains[0].Chain, "Chain name should not be empty")
	}
}

func TestIntegration_GetTokenMetadataV1(t *testing.T) {
	client := getTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	req := &apiv1.TokenMetadataRequest{
		Chain:        apiv1.TokenChainSolana,
		TokenAddress: testSolanaToken,
	}

	resp, err := client.GetTokenMetadataV1(ctx, req)
	require.NoError(t, err, "GetTokenMetadataV1 should not return an error")
	require.NotNil(t, resp, "Response should not be nil")

	assert.Equal(t, "solana", resp.Chain, "Chain should be solana")
	assert.NotEmpty(t, resp.Name, "Token name should not be empty")
	assert.NotEmpty(t, resp.Symbol, "Token symbol should not be empty")
	assert.Greater(t, resp.Decimals, 0, "Decimals should be greater than 0")
}

func TestIntegration_GetTokenPriceV1(t *testing.T) {
	client := getTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	req := &apiv1.TokenPriceRequest{
		Chain:        apiv1.TokenChainSolana,
		TokenAddress: testSolanaToken,
	}

	resp, err := client.GetTokenPriceV1(ctx, req)
	require.NoError(t, err, "GetTokenPriceV1 should not return an error")
	require.NotNil(t, resp, "Response should not be nil")

	assert.Equal(t, "solana", resp.Chain, "Chain should be solana")
	assert.Greater(t, resp.Price, 0.0, "Price should be greater than 0")
	assert.Greater(t, resp.BlockNumber, uint64(0), "BlockNumber should be greater than 0")
}

func TestIntegration_GetTokenStatsV1(t *testing.T) {
	client := getTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	req := &apiv1.TokenStatsRequest{
		Chain:        apiv1.TokenChainSolana,
		TokenAddress: testSolanaToken,
	}

	resp, err := client.GetTokenStatsV1(ctx, req)
	require.NoError(t, err, "GetTokenStatsV1 should not return an error")
	require.NotNil(t, resp, "Response should not be nil")

	assert.Equal(t, testSolanaToken, resp.TokenAddress, "Token address should match")
	assert.Greater(t, resp.PriceUSD, 0.0, "Price should be greater than 0")
	assert.Greater(t, resp.MarketCapUSD, 0.0, "Market cap should be greater than 0")

	// Check volume breakdown structure
	assert.NotNil(t, resp.Volume, "Volume should not be nil")
}

func TestIntegration_GetTokenBalanceV1(t *testing.T) {
	client := getTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	req := &apiv1.TokenBalanceRequest{
		Wallet:       testSolanaWallet,
		TokenAddress: testSolanaToken,
		Chain:        apiv1.TokenChainSolana,
	}

	resp, err := client.GetTokenBalanceV1(ctx, req)
	require.NoError(t, err, "GetTokenBalanceV1 should not return an error")
	require.NotNil(t, resp, "Response should not be nil")

	assert.Equal(t, "solana", resp.Chain, "Chain should be solana")
	assert.Equal(t, testSolanaToken, resp.TokenAddress, "Token address should match")
	assert.NotEmpty(t, resp.TokenName, "Token name should not be empty")
	assert.GreaterOrEqual(t, resp.Balance, 0.0, "Balance should be nonnegative")
}

func TestIntegration_GetTradingStatsV1(t *testing.T) {
	client := getTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	req := &apiv1.TradingStatsRequest{
		Wallet: testSolanaWallet,
		Days:   nil, // Server default is 30d
	}

	resp, err := client.GetTradingStatsV1(ctx, req)

	// This endpoint may return 202 Accepted or require retry
	// so we check for both success and acceptable error cases
	if err != nil {
		var apiErr *api.Error
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, 202, apiErr.StatusCode)
		require.Equal(t, "pending", apiErr.Status)
		require.Positive(t, apiErr.RetryAfterDuration)
		return
	}

	require.NotNil(t, resp, "Response should not be nil")
	assert.GreaterOrEqual(t, resp.SwapsCount, 0, "Total trades should be >= 0")
	assert.GreaterOrEqual(t, resp.WinRate, 0.0, "Win rate should be >= 0")
}

func TestIntegration_GetTradingStatsV1_WithTimeframe(t *testing.T) {
	client := getTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	timeframes := []apiv1.TimeframeDays{
		apiv1.Timeframe1Day,
		apiv1.Timeframe7Days,
		apiv1.Timeframe30Days,
		apiv1.TimeframeMax,
	}

	for _, tf := range timeframes {
		t.Run(string(tf), func(t *testing.T) {
			req := &apiv1.TradingStatsRequest{
				Wallet: testSolanaWallet,
				Days:   &tf,
			}

			resp, err := client.GetTradingStatsV1(ctx, req)

			if err != nil {
				var apiErr *api.Error
				require.ErrorAs(t, err, &apiErr)
				require.Equal(t, 202, apiErr.StatusCode)
				require.Equal(t, "pending", apiErr.Status)
				require.Positive(t, apiErr.RetryAfterDuration)
				return
			}

			require.NotNil(t, resp, "Response should not be nil for timeframe %s", tf)
		})
	}
}

func TestIntegration_ErrorHandling_InvalidWallet(t *testing.T) {
	client := getTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	_, err := client.GetWalletPortfolioV1(ctx, "invalid-wallet-address")
	assert.Error(t, err, "Should return error for invalid wallet address")
}

func TestIntegration_ErrorHandling_InvalidToken(t *testing.T) {
	client := getTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	req := &apiv1.TokenMetadataRequest{
		Chain:        apiv1.TokenChainSolana,
		TokenAddress: "invalid-token-address",
	}

	_, err := client.GetTokenMetadataV1(ctx, req)
	assert.Error(t, err, "Should return error for invalid token address")
}

// TestIntegrationReleaseSmoke is a read-only, sequential release gate. Upper-bound
// costs are reserved before each call, and SDK calls never automatically retry.
func TestIntegrationReleaseSmoke(t *testing.T) {
	client := getTestClient(t)
	raw := os.Getenv("CIELO_LIVE_CREDIT_BUDGET")
	budget, err := strconv.Atoi(raw)
	require.NoError(t, err, "set CIELO_LIVE_CREDIT_BUDGET explicitly")
	require.GreaterOrEqual(t, budget, 72, "smoke suite reserves at most 72 credits")
	require.LessOrEqual(t, budget, 100, "release smoke hard limit is 100 credits")
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	reserved := 0
	t.Cleanup(func() {
		t.Logf("Reserved credit upper bound: %d; limit: %d", reserved, budget)
	})
	call := func(name string, cost int, fn func(*testing.T)) {
		require.LessOrEqual(t, reserved+cost, budget, "credit budget exhausted before %s", name)
		reserved += cost
		if !t.Run(name, func(t *testing.T) {
			t.Logf("Reserved credit upper bound for this request: %d; cumulative: %d", cost, reserved)
			fn(t)
		}) {
			t.FailNow()
		}
	}
	pending := func(t *testing.T, err error) bool {
		if err == nil {
			return false
		}
		var apiErr *api.Error
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, http.StatusAccepted, apiErr.StatusCode)
		require.Equal(t, "pending", apiErr.Status)
		require.Positive(t, apiErr.RetryAfterDuration)
		t.Log("API returned pending; no paid retry was made")
		return true
	}
	call("chains", 0, func(t *testing.T) {
		result, err := client.GetChainsV1(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, result)
	})
	call("metadata", 2, func(t *testing.T) {
		result, err := client.GetTokenMetadataV1(ctx, &apiv1.TokenMetadataRequest{Chain: apiv1.TokenChainSolana, TokenAddress: testSolanaToken})
		require.NoError(t, err)
		require.Equal(t, "solana", result.Chain)
		require.NotEmpty(t, result.Symbol)
	})
	call("price", 2, func(t *testing.T) {
		result, err := client.GetTokenPriceV1(ctx, &apiv1.TokenPriceRequest{Chain: apiv1.TokenChainSolana, TokenAddress: testSolanaToken})
		require.NoError(t, err)
		require.Greater(t, result.Price, 0.0)
	})
	call("stats", 10, func(t *testing.T) {
		result, err := client.GetTokenStatsV1(ctx, &apiv1.TokenStatsRequest{Chain: apiv1.TokenChainSolana, TokenAddress: testSolanaToken})
		require.NoError(t, err)
		require.Equal(t, testSolanaToken, result.TokenAddress)
	})
	call("feed", 5, func(t *testing.T) {
		result, err := client.GetFeedV1(ctx, &apiv1.FeedRequest{Wallet: testSolanaWallet, Limit: apiv1.ToRef(5)})
		if pending(t, err) {
			return
		}
		require.NotNil(t, result)
		if result.Paging.HasNextPage {
			require.NotEmpty(t, result.Paging.NextCursor)
		}
	})
	call("portfolio", 20, func(t *testing.T) {
		result, err := client.GetWalletPortfolioV1(ctx, testSolanaWallet)
		require.NoError(t, err)
		require.GreaterOrEqual(t, result.TotalUSDValue, 0.0)
		for _, asset := range result.Portfolio {
			require.NotEmpty(t, asset.TokenAddress)
		}
	})
	call("balance", 3, func(t *testing.T) {
		result, err := client.GetTokenBalanceV1(ctx, &apiv1.TokenBalanceRequest{Wallet: testSolanaWallet, TokenAddress: testSolanaToken, Chain: apiv1.TokenChainSolana})
		if err != nil {
			var apiErr *api.Error
			require.ErrorAs(t, err, &apiErr)
			require.Equal(t, http.StatusNotFound, apiErr.StatusCode)
			t.Log("API returned balance not found; successful balance data was not validated")
			return
		}
		require.Equal(t, testSolanaToken, result.TokenAddress)
		require.GreaterOrEqual(t, result.Balance, 0.0)
	})
	call("trading", 30, func(t *testing.T) {
		result, err := client.GetTradingStatsV1(ctx, &apiv1.TradingStatsRequest{Wallet: testSolanaWallet, Days: apiv1.ToRef(apiv1.Timeframe7Days)})
		if pending(t, err) {
			return
		}
		require.GreaterOrEqual(t, result.WinRate, 0.0)
		require.LessOrEqual(t, result.WinRate, 1.0)
	})
}
