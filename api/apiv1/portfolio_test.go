package apiv1_test

import (
	"encoding/json"
	"testing"

	"github.com/sealtv/cielogo/api/apiv1"
	"github.com/sealtv/cielogo/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWalletPortfolioV2Request_GetQueryString(t *testing.T) {
	tests := []struct {
		name     string
		req      *apiv1.WalletPortfolioV2Request
		expected string
	}{
		{
			name: "single wallet",
			req: &apiv1.WalletPortfolioV2Request{
				Wallets: []string{"wallet1"},
			},
			expected: "wallet=wallet1",
		},
		{
			name: "multiple wallets",
			req: &apiv1.WalletPortfolioV2Request{
				Wallets: []string{"wallet1", "wallet2", "wallet3"},
			},
			expected: "wallet=wallet1%2Cwallet2%2Cwallet3",
		},
		{
			name: "with token filter",
			req: &apiv1.WalletPortfolioV2Request{
				Wallets: []string{"wallet1"},
				Token:   testutil.Ptr("token123"),
			},
			expected: "token=token123&wallet=wallet1",
		},
		{
			name: "empty wallets",
			req: &apiv1.WalletPortfolioV2Request{
				Wallets: []string{},
			},
			expected: "",
		},
		{
			name: "multiple wallets with token",
			req: &apiv1.WalletPortfolioV2Request{
				Wallets: []string{"wallet1", "wallet2"},
				Token:   testutil.Ptr("EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v"),
			},
			expected: "token=EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v&wallet=wallet1%2Cwallet2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.req.GetQueryString()
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestPortfolioAsset_JSON(t *testing.T) {
	asset := apiv1.PortfolioAsset{
		Chain:         "solana",
		TokenAddress:  "So11111111111111111111111111111111111111112",
		TokenName:     "Wrapped SOL",
		TokenSymbol:   "SOL",
		TokenPriceUSD: 23.45,
		Balance:       10.5,
		TotalUSDValue: 246.225,
		Supply:        1000000,
		SupplyOwned:   0.00105,
	}

	// Test marshaling
	data, err := json.Marshal(asset)
	require.NoError(t, err)

	// Test unmarshaling
	var decoded apiv1.PortfolioAsset
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Equal(t, asset.Chain, decoded.Chain)
	assert.Equal(t, asset.TokenAddress, decoded.TokenAddress)
	assert.Equal(t, asset.TokenName, decoded.TokenName)
	assert.Equal(t, asset.TokenSymbol, decoded.TokenSymbol)
	assert.Equal(t, asset.TokenPriceUSD, decoded.TokenPriceUSD)
	assert.Equal(t, asset.Balance, decoded.Balance)
	assert.Equal(t, asset.TotalUSDValue, decoded.TotalUSDValue)
	assert.Equal(t, asset.Supply, decoded.Supply)
	assert.Equal(t, asset.SupplyOwned, decoded.SupplyOwned)
}

func TestWalletPortfolioResponse_JSON(t *testing.T) {
	response := apiv1.WalletPortfolioResponse{
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
				Supply:        1000000,
				SupplyOwned:   0.01,
			},
		},
	}

	data, err := json.Marshal(response)
	require.NoError(t, err)

	var decoded apiv1.WalletPortfolioResponse
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Equal(t, response.LastUpdated, decoded.LastUpdated)
	assert.Equal(t, response.TotalUSDValue, decoded.TotalUSDValue)
	assert.Len(t, decoded.Portfolio, 1)
	assert.Equal(t, "solana", decoded.Portfolio[0].Chain)
}

func TestChainDistribution_JSON(t *testing.T) {
	dist := apiv1.ChainDistribution{
		Chain:      "ethereum",
		TotalUSD:   5000.0,
		Percentage: 45.5,
	}

	data, err := json.Marshal(dist)
	require.NoError(t, err)

	var decoded apiv1.ChainDistribution
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Equal(t, dist.Chain, decoded.Chain)
	assert.Equal(t, dist.TotalUSD, decoded.TotalUSD)
	assert.Equal(t, dist.Percentage, decoded.Percentage)
}

func TestWalletPortfolioV2Response_JSON(t *testing.T) {
	response := apiv1.WalletPortfolioV2Response{
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
				Supply:        1000000,
				SupplyOwned:   0.01,
				WalletAddress: "wallet1",
			},
		},
		Chains: []apiv1.ChainDistribution{
			{
				Chain:      "solana",
				TotalUSD:   6000.0,
				Percentage: 60.0,
			},
			{
				Chain:      "ethereum",
				TotalUSD:   4000.0,
				Percentage: 40.0,
			},
		},
	}

	data, err := json.Marshal(response)
	require.NoError(t, err)

	var decoded apiv1.WalletPortfolioV2Response
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Equal(t, response.LastUpdated, decoded.LastUpdated)
	assert.Equal(t, response.TotalUSDValue, decoded.TotalUSDValue)
	assert.Len(t, decoded.Portfolio, 1)
	assert.Equal(t, "wallet1", decoded.Portfolio[0].WalletAddress)
	assert.Len(t, decoded.Chains, 2)
}
