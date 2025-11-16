package cielogo_test

import (
	"context"
	"testing"

	"github.com/sealtv/cielogo"
	"github.com/sealtv/cielogo/api"
	"github.com/sealtv/cielogo/api/apiv1"
	"github.com/sealtv/cielogo/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetWalletPortfolioV1(t *testing.T) {
	mockResp := api.CieloResponse[apiv1.WalletPortfolioResponse]{
		Data: apiv1.WalletPortfolioResponse{
			LastUpdated:   1699564800,
			TotalUSDValue: 1234.56,
			Portfolio: []apiv1.PortfolioAsset{
				{
					Chain:         "solana",
					TokenAddress:  "token1",
					TokenName:     "Token 1",
					TokenSymbol:   "TK1",
					TokenPriceUSD: 10.0,
					Balance:       100.0,
					TotalUSDValue: 1000.0,
				},
			},
		},
	}

	server := testutil.NewMockServer(t)
	server.SetResponse("/v1/wallet123/portfolio", mockResp)

	client := cielogo.NewClient("test-key", cielogo.WithBaseURL(server.URL))
	ctx := context.Background()

	resp, err := client.GetWalletPortfolioV1(ctx, "wallet123")
	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.Equal(t, int64(1699564800), resp.LastUpdated)
	assert.Equal(t, 1234.56, resp.TotalUSDValue)
	assert.Len(t, resp.Portfolio, 1)
}

func TestGetWalletPortfolioV2(t *testing.T) {
	mockResp := api.CieloResponse[apiv1.WalletPortfolioV2Response]{
		Data: apiv1.WalletPortfolioV2Response{
			LastUpdated:   1699564800,
			TotalUSDValue: 10000.0,
			Portfolio: []apiv1.PortfolioAssetV2{
				{
					Chain:         "solana",
					TokenAddress:  "token1",
					TokenName:     "Token 1",
					TokenSymbol:   "TK1",
					TokenPriceUSD: 10.0,
					Balance:       100.0,
					TotalUSDValue: 1000.0,
					WalletAddress: "wallet1",
				},
			},
		},
	}

	server := testutil.NewMockServer(t)
	server.SetResponse("/v2/portfolio?wallet=wallet1", mockResp)

	client := cielogo.NewClient("test-key", cielogo.WithBaseURL(server.URL))
	ctx := context.Background()

	req := &apiv1.WalletPortfolioV2Request{
		Wallets: []string{"wallet1"},
	}

	resp, err := client.GetWalletPortfolioV2(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.Equal(t, int64(1699564800), resp.LastUpdated)
	assert.Equal(t, 10000.0, resp.TotalUSDValue)
}

func TestGetTokenMetadataV1(t *testing.T) {
	mockResp := api.CieloResponse[apiv1.TokenMetadataResponse]{
		Data: apiv1.TokenMetadataResponse{
			Chain:    "solana",
			Address:  "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
			Name:     "USD Coin",
			Symbol:   "USDC",
			Decimals: 6,
			Supply:   1000000000000,
		},
	}

	server := testutil.NewMockServer(t)
	server.SetResponse("/v1/token/metadata?chain=solana&token_address=EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v", mockResp)

	client := cielogo.NewClient("test-key", cielogo.WithBaseURL(server.URL))
	ctx := context.Background()

	req := &apiv1.TokenMetadataRequest{
		Chain:        apiv1.TokenChainSolana,
		TokenAddress: "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
	}

	resp, err := client.GetTokenMetadataV1(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.Equal(t, "solana", resp.Chain)
	assert.Equal(t, "USD Coin", resp.Name)
	assert.Equal(t, "USDC", resp.Symbol)
}

func TestGetTokenPriceV1(t *testing.T) {
	mockResp := api.CieloResponse[apiv1.TokenPriceResponse]{
		Data: apiv1.TokenPriceResponse{
			Chain:       "solana",
			Address:     "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
			BlockNumber: 123456,
			Price:       0.9998,
		},
	}

	server := testutil.NewMockServer(t)
	server.SetResponse("/v1/token/price?chain=solana&token_address=EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v", mockResp)

	client := cielogo.NewClient("test-key", cielogo.WithBaseURL(server.URL))
	ctx := context.Background()

	req := &apiv1.TokenPriceRequest{
		Chain:        apiv1.TokenChainSolana,
		TokenAddress: "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
	}

	resp, err := client.GetTokenPriceV1(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.Equal(t, "solana", resp.Chain)
	assert.InDelta(t, 0.9998, resp.Price, 0.0001)
}

func TestGetTokenStatsV1(t *testing.T) {
	mockResp := api.CieloResponse[apiv1.TokenStatsResponse]{
		Data: apiv1.TokenStatsResponse{
			TokenAddress: "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
			PriceUSD:     1.0,
			MarketCapUSD: 25000000000,
			Change: apiv1.PriceChange{
				FiveMin:     0.01,
				OneHour:     0.02,
				SixHours:    -0.01,
				TwentyFourH: 0.05,
			},
		},
	}

	server := testutil.NewMockServer(t)
	server.SetResponse("/v1/token/stats?chain=solana&token_address=EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v", mockResp)

	client := cielogo.NewClient("test-key", cielogo.WithBaseURL(server.URL))
	ctx := context.Background()

	req := &apiv1.TokenStatsRequest{
		Chain:        apiv1.TokenChainSolana,
		TokenAddress: "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
	}

	resp, err := client.GetTokenStatsV1(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.Equal(t, "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v", resp.TokenAddress)
	assert.Equal(t, 1.0, resp.PriceUSD)
}

func TestGetTokenBalanceV1(t *testing.T) {
	mockResp := api.CieloResponse[apiv1.TokenBalanceResponse]{
		Data: apiv1.TokenBalanceResponse{
			Chain:         "solana",
			TokenAddress:  "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
			TokenName:     "USD Coin",
			TokenSymbol:   "USDC",
			TokenPriceUSD: 1.0,
			Balance:       1000.50,
			TotalUSDValue: 1000.50,
			Decimals:      6,
		},
	}

	server := testutil.NewMockServer(t)
	server.SetResponse("/v1/wallet123/token-balance?chain=solana&token_address=EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v", mockResp)

	client := cielogo.NewClient("test-key", cielogo.WithBaseURL(server.URL))
	ctx := context.Background()

	req := &apiv1.TokenBalanceRequest{
		Wallet:       "wallet123",
		TokenAddress: "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
		Chain:        apiv1.TokenChainSolana,
	}

	resp, err := client.GetTokenBalanceV1(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.Equal(t, "solana", resp.Chain)
	assert.Equal(t, 1000.50, resp.Balance)
}

func TestGetTradingStatsV1(t *testing.T) {
	mockResp := api.CieloResponse[apiv1.TradingStatsResponse]{
		Data: apiv1.TradingStatsResponse{
			TotalTrades:      150,
			SuccessfulTrades: 120,
			FailedTrades:     30,
			WinRate:          80.0,
			AverageProfit:    125.50,
			TotalProfit:      15060.0,
		},
	}

	server := testutil.NewMockServer(t)
	server.SetResponse("/v1/wallet123/trading-stats?days=7d", mockResp)

	client := cielogo.NewClient("test-key", cielogo.WithBaseURL(server.URL))
	ctx := context.Background()

	req := &apiv1.TradingStatsRequest{
		Wallet: "wallet123",
		Days:   testutil.Ptr(apiv1.Timeframe7Days),
	}

	resp, err := client.GetTradingStatsV1(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.Equal(t, 150, resp.TotalTrades)
	assert.Equal(t, 80.0, resp.WinRate)
}

func TestErrorResponse(t *testing.T) {
	server := testutil.NewMockServer(t)
	// Don't set any response - will return 404

	client := cielogo.NewClient("test-key", cielogo.WithBaseURL(server.URL))
	ctx := context.Background()

	_, err := client.GetWalletPortfolioV1(ctx, "wallet123")
	assert.Error(t, err)
}

func TestMakeRequest_WithAPIError(t *testing.T) {
	// MockServer returns 404 for unknown paths, which should be treated as an error
	server := testutil.NewMockServer(t)
	// Don't set any response for this path - will return 404

	client := cielogo.NewClient("test-key", cielogo.WithBaseURL(server.URL))
	ctx := context.Background()

	_, err := client.GetWalletPortfolioV1(ctx, "unknown-wallet")
	assert.Error(t, err, "Should return error for 404 response")
}

func TestClient_ContextCancellation(t *testing.T) {
	server := testutil.NewMockServer(t)
	mockResp := api.CieloResponse[apiv1.WalletPortfolioResponse]{
		Data: apiv1.WalletPortfolioResponse{},
	}
	server.SetResponse("/v1/wallet123/portfolio", mockResp)

	client := cielogo.NewClient("test-key", cielogo.WithBaseURL(server.URL))

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	_, err := client.GetWalletPortfolioV1(ctx, "wallet123")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context canceled")
}
