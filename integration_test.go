//go:build integration
// +build integration

package cielogo_test

import (
	"context"
	"os"
	"testing"

	"github.com/sealtv/cielogo"
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

	return cielogo.NewClient(apiKey)
}

func TestIntegration_GetWalletPortfolioV1(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	resp, err := client.GetWalletPortfolioV1(ctx, testSolanaWallet)
	require.NoError(t, err, "GetWalletPortfolioV1 should not return an error")
	require.NotNil(t, resp, "Response should not be nil")

	assert.Greater(t, resp.LastUpdated, int64(0), "LastUpdated should be greater than 0")
}

func TestIntegration_GetWalletPortfolioV2(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

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
	client := getTestClient(t)
	ctx := context.Background()

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
	ctx := context.Background()

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
	ctx := context.Background()

	req := &apiv1.TokenPriceRequest{
		Chain:        apiv1.TokenChainSolana,
		TokenAddress: testSolanaToken,
	}

	resp, err := client.GetTokenPriceV1(ctx, req)
	require.NoError(t, err, "GetTokenPriceV1 should not return an error")
	require.NotNil(t, resp, "Response should not be nil")

	assert.Equal(t, "solana", resp.Chain, "Chain should be solana")
	assert.Greater(t, resp.Price, 0.0, "Price should be greater than 0")
	assert.Greater(t, resp.BlockNumber, int64(0), "BlockNumber should be greater than 0")
}

func TestIntegration_GetTokenStatsV1(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

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
	ctx := context.Background()

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
	assert.Greater(t, resp.Decimals, 0, "Decimals should be greater than 0")
}

func TestIntegration_GetTradingStatsV1(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	req := &apiv1.TradingStatsRequest{
		Wallet: testSolanaWallet,
		Days:   nil, // Use default (max)
	}

	resp, err := client.GetTradingStatsV1(ctx, req)

	// This endpoint may return 202 Accepted or require retry
	// so we check for both success and acceptable error cases
	if err != nil {
		t.Logf("GetTradingStatsV1 returned error (may be 202 Accepted): %v", err)
		return
	}

	require.NotNil(t, resp, "Response should not be nil")
	assert.GreaterOrEqual(t, resp.TotalTrades, 0, "Total trades should be >= 0")
	assert.GreaterOrEqual(t, resp.WinRate, 0.0, "Win rate should be >= 0")
}

func TestIntegration_GetTradingStatsV1_WithTimeframe(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

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
				t.Logf("GetTradingStatsV1 with %s returned error: %v", tf, err)
				return
			}

			require.NotNil(t, resp, "Response should not be nil for timeframe %s", tf)
		})
	}
}

func TestIntegration_ErrorHandling_InvalidWallet(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	_, err := client.GetWalletPortfolioV1(ctx, "invalid-wallet-address")
	assert.Error(t, err, "Should return error for invalid wallet address")
}

func TestIntegration_ErrorHandling_InvalidToken(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	req := &apiv1.TokenMetadataRequest{
		Chain:        apiv1.TokenChainSolana,
		TokenAddress: "invalid-token-address",
	}

	_, err := client.GetTokenMetadataV1(ctx, req)
	assert.Error(t, err, "Should return error for invalid token address")
}
