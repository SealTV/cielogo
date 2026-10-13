package apiv1

import (
	"net/url"

	"github.com/sealtv/cielogo/api/chains"
)

type NftsPnLRequest struct {
	Wallet     string
	Timeframe  *string
	NextObject *string
}

func (r *NftsPnLRequest) GetQueryString() string {
	values := url.Values{}

	// Add timeframe if specified
	if r.Timeframe != nil {
		values.Add("timeframe", *r.Timeframe)
	}

	// Add NextObject if specified
	if r.NextObject != nil {
		values.Add("next_object", *r.NextObject)
	}

	return values.Encode()
}

// NftsPnLResponse
type NftsPnLResponse struct {
	Items  []NftPnl   `json:"items"`
	Paging Pagination `json:"paging"`
}

type NftPnl struct {
	Address string           `json:"token_address"`
	Symbol  string           `json:"token_symbol"`
	Name    string           `json:"token_name"`
	Chain   chains.ChainType `json:"chain"`

	TotalBuyAmount  float64 `json:"total_buy_amount"`
	TotalBuyAvg     float64 `json:"average_buy_price"`
	TotalSellAmount float64 `json:"total_sell_amount"`
	TotalBuyUsd     float64 `json:"total_buy_usd"`
	TotalSellAvg    float64 `json:"average_sell_price"`
	TotalSellUsd    float64 `json:"total_sell_usd"`
	ProfitLoss      float64 `json:"total_pnl_usd"`
	RoiPercentage   float64 `json:"roi_percentage"`
	NumSwaps        int     `json:"num_swaps"`
}
