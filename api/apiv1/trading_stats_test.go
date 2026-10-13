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
	var decoded apiv1.TradingStatsResponse
	err := json.Unmarshal([]byte(`{"pnl":15060,"winrate":0.8,"swaps_count":150,"buy_count":120,"sell_count":30,"average_buy_amount_usd":125.5,"roi_distribution":{"roi_above_500":2},"holding_distribution":{"hold_gt_24_hours":4}}`), &decoded)
	require.NoError(t, err)
	assert.Equal(t, int64(150), decoded.SwapsCount)
	assert.Equal(t, 0.8, decoded.WinRate)
	assert.Equal(t, 15060.0, decoded.PnL)
	assert.Equal(t, 125.5, decoded.AverageBuyAmountUSD)
	assert.Equal(t, 2.0, decoded.ROIDistribution.ROIAbove500)
	assert.Equal(t, int64(4), decoded.HoldingDistribution.HoldGT24Hours)
}
