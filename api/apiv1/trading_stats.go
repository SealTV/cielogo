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
	// Defaults to 30d if not specified.
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
	PnL float64 `json:"pnl"`
	// WinRate is a fraction (0..1), not a percentage.
	WinRate               float64             `json:"winrate"`
	SwapsCount            int64               `json:"swaps_count"`
	BuyCount              int64               `json:"buy_count"`
	SellCount             int64               `json:"sell_count"`
	ROIDistribution       ROIDistribution     `json:"roi_distribution"`
	TotalBuyAmountUSD     float64             `json:"total_buy_amount_usd"`
	TotalSellAmountUSD    float64             `json:"total_sell_amount_usd"`
	AverageBuyAmountUSD   float64             `json:"average_buy_amount_usd"`
	AverageSellAmountUSD  float64             `json:"average_sell_amount_usd"`
	AverageHoldingTimeSec float64             `json:"average_holding_time_sec"`
	ConsecutiveDays       int64               `json:"consecutive_trading_days"`
	HoldingDistribution   HoldingDistribution `json:"holding_distribution"`
}

type ROIDistribution struct {
	ROIAbove500   float64 `json:"roi_above_500"`
	ROI200To500   float64 `json:"roi_200_to_500"`
	ROI50To200    float64 `json:"roi_50_to_200"`
	ROI0To50      float64 `json:"roi_0_to_50"`
	ROINeg50To0   float64 `json:"roi_neg50_to_0"`
	ROIBelowNeg50 float64 `json:"roi_below_neg50"`
}

type HoldingDistribution struct {
	Hold0To3Min    int64 `json:"hold_0_3_min"`
	Hold3To20Min   int64 `json:"hold_3_20_min"`
	Hold20To60Min  int64 `json:"hold_20_60_min"`
	Hold1To6Hours  int64 `json:"hold_1_6_hours"`
	Hold6To24Hours int64 `json:"hold_6_24_hours"`
	HoldGT24Hours  int64 `json:"hold_gt_24_hours"`
	TotalTokens    int64 `json:"total_tokens"`
}
