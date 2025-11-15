package apiv1

import "net/url"

// TimeframeDays represents the time period for trading statistics queries.
type TimeframeDays string

const (
	Timeframe1Day   TimeframeDays = "1d"
	Timeframe7Days  TimeframeDays = "7d"
	Timeframe30Days TimeframeDays = "30d"
	TimeframeMax    TimeframeDays = "max"
)

// TradingStatsRequest represents a request for trading performance statistics.
type TradingStatsRequest struct {
	Wallet string

	// Days specifies the timeframe for the statistics.
	// Supported values: 1d, 7d, 30d, max
	// Defaults to max if not specified.
	Days *TimeframeDays
}

// GetQueryString builds the query string for trading stats requests.
func (r *TradingStatsRequest) GetQueryString() string {
	values := url.Values{}

	if r.Days != nil && *r.Days != "" {
		values.Add("days", string(*r.Days))
	}

	return values.Encode()
}

// TradingStatsResponse represents detailed performance statistics for a wallet's trading activity.
// Includes PnL, ROI, win rate, and trading behavior insights.
type TradingStatsResponse struct {
	TotalTrades         int     `json:"total_trades"`
	SuccessfulTrades    int     `json:"successful_trades"`
	FailedTrades        int     `json:"failed_trades"`
	WinRate             float64 `json:"win_rate"`
	AverageProfit       float64 `json:"average_profit"`
	TotalProfit         float64 `json:"total_profit"`
	TotalVolume         float64 `json:"total_volume"`
	MostTradedToken     string  `json:"most_traded_token"`
	MostProfitableToken string  `json:"most_profitable_token"`
	TradingFrequency    string  `json:"trading_frequency"`
	LastTradeTimestamp  int64   `json:"last_trade_timestamp"`
}
