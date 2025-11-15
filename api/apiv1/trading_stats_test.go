package apiv1_test

import (
	"encoding/json"
	"testing"

	"github.com/sealtv/cielogo/api/apiv1"
	"github.com/sealtv/cielogo/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTimeframeDays_Constants(t *testing.T) {
	tests := []struct {
		name      string
		timeframe apiv1.TimeframeDays
		expected  string
	}{
		{"1 day", apiv1.Timeframe1Day, "1d"},
		{"7 days", apiv1.Timeframe7Days, "7d"},
		{"30 days", apiv1.Timeframe30Days, "30d"},
		{"max", apiv1.TimeframeMax, "max"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, string(tt.timeframe))
		})
	}
}

func TestTradingStatsRequest_GetQueryString(t *testing.T) {
	tests := []struct {
		name     string
		req      *apiv1.TradingStatsRequest
		expected string
	}{
		{
			name: "no days specified",
			req: &apiv1.TradingStatsRequest{
				Wallet: "wallet123",
				Days:   nil,
			},
			expected: "",
		},
		{
			name: "1 day",
			req: &apiv1.TradingStatsRequest{
				Wallet: "wallet123",
				Days:   testutil.Ptr(apiv1.Timeframe1Day),
			},
			expected: "days=1d",
		},
		{
			name: "7 days",
			req: &apiv1.TradingStatsRequest{
				Wallet: "wallet123",
				Days:   testutil.Ptr(apiv1.Timeframe7Days),
			},
			expected: "days=7d",
		},
		{
			name: "30 days",
			req: &apiv1.TradingStatsRequest{
				Wallet: "wallet123",
				Days:   testutil.Ptr(apiv1.Timeframe30Days),
			},
			expected: "days=30d",
		},
		{
			name: "max",
			req: &apiv1.TradingStatsRequest{
				Wallet: "wallet123",
				Days:   testutil.Ptr(apiv1.TimeframeMax),
			},
			expected: "days=max",
		},
		{
			name: "empty days string",
			req: &apiv1.TradingStatsRequest{
				Wallet: "wallet123",
				Days:   testutil.Ptr(apiv1.TimeframeDays("")),
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.req.GetQueryString()
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestTradingStatsResponse_JSON(t *testing.T) {
	response := apiv1.TradingStatsResponse{
		TotalTrades:         150,
		SuccessfulTrades:    120,
		FailedTrades:        30,
		WinRate:             80.0,
		AverageProfit:       125.50,
		TotalProfit:         15060.0,
		TotalVolume:         250000.0,
		MostTradedToken:     "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
		MostProfitableToken: "So11111111111111111111111111111111111111112",
		TradingFrequency:    "high",
		LastTradeTimestamp:  1699564800,
	}

	data, err := json.Marshal(response)
	require.NoError(t, err)

	var decoded apiv1.TradingStatsResponse
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Equal(t, response.TotalTrades, decoded.TotalTrades)
	assert.Equal(t, response.SuccessfulTrades, decoded.SuccessfulTrades)
	assert.Equal(t, response.FailedTrades, decoded.FailedTrades)
	assert.Equal(t, response.WinRate, decoded.WinRate)
	assert.Equal(t, response.AverageProfit, decoded.AverageProfit)
	assert.Equal(t, response.TotalProfit, decoded.TotalProfit)
	assert.Equal(t, response.TotalVolume, decoded.TotalVolume)
	assert.Equal(t, response.MostTradedToken, decoded.MostTradedToken)
	assert.Equal(t, response.MostProfitableToken, decoded.MostProfitableToken)
	assert.Equal(t, response.TradingFrequency, decoded.TradingFrequency)
	assert.Equal(t, response.LastTradeTimestamp, decoded.LastTradeTimestamp)
}

func TestTradingStatsResponse_EmptyStats(t *testing.T) {
	// Test with a wallet that has no trades
	response := apiv1.TradingStatsResponse{
		TotalTrades:         0,
		SuccessfulTrades:    0,
		FailedTrades:        0,
		WinRate:             0.0,
		AverageProfit:       0.0,
		TotalProfit:         0.0,
		TotalVolume:         0.0,
		MostTradedToken:     "",
		MostProfitableToken: "",
		TradingFrequency:    "none",
		LastTradeTimestamp:  0,
	}

	data, err := json.Marshal(response)
	require.NoError(t, err)

	var decoded apiv1.TradingStatsResponse
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Equal(t, 0, decoded.TotalTrades)
	assert.Equal(t, 0.0, decoded.WinRate)
	assert.Equal(t, "", decoded.MostTradedToken)
}
