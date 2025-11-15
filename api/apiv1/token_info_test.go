package apiv1_test

import (
	"encoding/json"
	"testing"

	"github.com/sealtv/cielogo/api/apiv1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTokenChain_Constants(t *testing.T) {
	tests := []struct {
		name     string
		chain    apiv1.TokenChain
		expected string
	}{
		{"solana", apiv1.TokenChainSolana, "solana"},
		{"ethereum", apiv1.TokenChainEthereum, "ethereum"},
		{"base", apiv1.TokenChainBase, "base"},
		{"hyperevm", apiv1.TokenChainHyperevm, "hyperevm"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, string(tt.chain))
		})
	}
}

func TestTokenMetadataRequest_GetQueryString(t *testing.T) {
	tests := []struct {
		name     string
		req      *apiv1.TokenMetadataRequest
		expected string
	}{
		{
			name: "solana token",
			req: &apiv1.TokenMetadataRequest{
				Chain:        apiv1.TokenChainSolana,
				TokenAddress: "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
			},
			expected: "chain=solana&token_address=EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
		},
		{
			name: "ethereum token",
			req: &apiv1.TokenMetadataRequest{
				Chain:        apiv1.TokenChainEthereum,
				TokenAddress: "0xdAC17F958D2ee523a2206206994597C13D831ec7",
			},
			expected: "chain=ethereum&token_address=0xdAC17F958D2ee523a2206206994597C13D831ec7",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.req.GetQueryString()
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestTokenPriceRequest_GetQueryString(t *testing.T) {
	req := &apiv1.TokenPriceRequest{
		Chain:        apiv1.TokenChainBase,
		TokenAddress: "0x123",
	}

	expected := "chain=base&token_address=0x123"
	assert.Equal(t, expected, req.GetQueryString())
}

func TestTokenStatsRequest_GetQueryString(t *testing.T) {
	req := &apiv1.TokenStatsRequest{
		Chain:        apiv1.TokenChainHyperevm,
		TokenAddress: "0x456",
	}

	expected := "chain=hyperevm&token_address=0x456"
	assert.Equal(t, expected, req.GetQueryString())
}

func TestTokenBalanceRequest_GetQueryString(t *testing.T) {
	req := &apiv1.TokenBalanceRequest{
		Wallet:       "wallet123",
		TokenAddress: "token456",
		Chain:        apiv1.TokenChainSolana,
	}

	got := req.GetQueryString()
	assert.Contains(t, got, "chain=solana")
	assert.Contains(t, got, "token_address=token456")
}

func TestTokenMetadataResponse_JSON(t *testing.T) {
	response := apiv1.TokenMetadataResponse{
		Chain:          "solana",
		Address:        "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
		Name:           "USD Coin",
		Symbol:         "USDC",
		Decimals:       6,
		CreatedAtTs:    1599564800,
		CreatedAtBlock: 12345678,
		Twitter:        "https://twitter.com/circle",
		Telegram:       "https://t.me/circle",
		Website:        "https://www.circle.com/en/usdc",
		Supply:         1000000000000,
		LogoURI:        "https://example.com/usdc.png",
	}

	data, err := json.Marshal(response)
	require.NoError(t, err)

	var decoded apiv1.TokenMetadataResponse
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Equal(t, response.Chain, decoded.Chain)
	assert.Equal(t, response.Address, decoded.Address)
	assert.Equal(t, response.Name, decoded.Name)
	assert.Equal(t, response.Symbol, decoded.Symbol)
	assert.Equal(t, response.Decimals, decoded.Decimals)
	assert.Equal(t, response.Supply, decoded.Supply)
}

func TestTokenPriceResponse_JSON(t *testing.T) {
	response := apiv1.TokenPriceResponse{
		Chain:       "ethereum",
		Address:     "0xdAC17F958D2ee523a2206206994597C13D831ec7",
		BlockNumber: 18500000,
		Price:       0.9998,
	}

	data, err := json.Marshal(response)
	require.NoError(t, err)

	var decoded apiv1.TokenPriceResponse
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Equal(t, response.Chain, decoded.Chain)
	assert.Equal(t, response.Address, decoded.Address)
	assert.Equal(t, response.BlockNumber, decoded.BlockNumber)
	assert.InDelta(t, response.Price, decoded.Price, 0.0001)
}

func TestVolumeStats_JSON(t *testing.T) {
	stats := apiv1.VolumeStats{
		VolumeUSD:     1000000.50,
		BuyVolumeUSD:  600000.25,
		SellVolumeUSD: 400000.25,
		UniqueBuyers:  1500,
		UniqueSellers: 1200,
		Buys:          5000,
		Sells:         4000,
	}

	data, err := json.Marshal(stats)
	require.NoError(t, err)

	var decoded apiv1.VolumeStats
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Equal(t, stats.VolumeUSD, decoded.VolumeUSD)
	assert.Equal(t, stats.UniqueBuyers, decoded.UniqueBuyers)
	assert.Equal(t, stats.Buys, decoded.Buys)
}

func TestPriceChange_JSON(t *testing.T) {
	change := apiv1.PriceChange{
		FiveMin:     2.5,
		OneHour:     5.3,
		SixHours:    -3.2,
		TwentyFourH: 15.7,
	}

	data, err := json.Marshal(change)
	require.NoError(t, err)

	var decoded apiv1.PriceChange
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Equal(t, change.FiveMin, decoded.FiveMin)
	assert.Equal(t, change.OneHour, decoded.OneHour)
	assert.Equal(t, change.SixHours, decoded.SixHours)
	assert.Equal(t, change.TwentyFourH, decoded.TwentyFourH)
}

func TestTokenStatsResponse_JSON(t *testing.T) {
	response := apiv1.TokenStatsResponse{
		TokenAddress: "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
		PriceUSD:     1.0,
		MarketCapUSD: 25000000000,
		Change: apiv1.PriceChange{
			FiveMin:     0.01,
			OneHour:     0.02,
			SixHours:    -0.01,
			TwentyFourH: 0.05,
		},
		Volume: apiv1.VolumeBreakdown{
			FiveMin: apiv1.VolumeStats{
				VolumeUSD:     100000,
				BuyVolumeUSD:  60000,
				SellVolumeUSD: 40000,
				UniqueBuyers:  50,
				UniqueSellers: 40,
				Buys:          100,
				Sells:         80,
			},
			OneHour: apiv1.VolumeStats{
				VolumeUSD:     500000,
				BuyVolumeUSD:  300000,
				SellVolumeUSD: 200000,
				UniqueBuyers:  200,
				UniqueSellers: 150,
				Buys:          400,
				Sells:         300,
			},
			SixHours: apiv1.VolumeStats{
				VolumeUSD:     2000000,
				BuyVolumeUSD:  1200000,
				SellVolumeUSD: 800000,
				UniqueBuyers:  800,
				UniqueSellers: 600,
				Buys:          1500,
				Sells:         1200,
			},
			TwentyFourH: apiv1.VolumeStats{
				VolumeUSD:     10000000,
				BuyVolumeUSD:  6000000,
				SellVolumeUSD: 4000000,
				UniqueBuyers:  3000,
				UniqueSellers: 2500,
				Buys:          6000,
				Sells:         5000,
			},
		},
	}

	data, err := json.Marshal(response)
	require.NoError(t, err)

	var decoded apiv1.TokenStatsResponse
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Equal(t, response.TokenAddress, decoded.TokenAddress)
	assert.Equal(t, response.PriceUSD, decoded.PriceUSD)
	assert.Equal(t, response.MarketCapUSD, decoded.MarketCapUSD)
	assert.Equal(t, response.Volume.TwentyFourH.VolumeUSD, decoded.Volume.TwentyFourH.VolumeUSD)
}

func TestTokenBalanceResponse_JSON(t *testing.T) {
	response := apiv1.TokenBalanceResponse{
		Chain:         "solana",
		TokenAddress:  "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
		TokenName:     "USD Coin",
		TokenSymbol:   "USDC",
		TokenPriceUSD: 1.0,
		Balance:       1000.50,
		TotalUSDValue: 1000.50,
		Decimals:      6,
	}

	data, err := json.Marshal(response)
	require.NoError(t, err)

	var decoded apiv1.TokenBalanceResponse
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Equal(t, response.Chain, decoded.Chain)
	assert.Equal(t, response.TokenAddress, decoded.TokenAddress)
	assert.Equal(t, response.TokenName, decoded.TokenName)
	assert.Equal(t, response.Balance, decoded.Balance)
	assert.Equal(t, response.Decimals, decoded.Decimals)
}
