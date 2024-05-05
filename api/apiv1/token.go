package apiv1

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/sealtv/cielogo/api/chains"
)

type TokensPnLRequest struct {
	Wallet       string
	Chains       []chains.ChainType
	Timeframe    *string
	NextObject   *string
	CexTransfers *bool
}

func (r *TokensPnLRequest) GetQueryString() string {
	values := url.Values{}

	// Add chains if specified
	if len(r.Chains) > 0 {
		var chainStrs []string
		for _, chain := range r.Chains {
			chainStrs = append(chainStrs, string(chain))
		}
		values.Add("chains", strings.Join(chainStrs, ","))
	}

	// Add timeframe if specified
	if r.Timeframe != nil {
		values.Add("timeframe", *r.Timeframe)
	}

	// Add CexTransfers if specified
	if r.CexTransfers != nil {
		values.Add("cex_transfers", fmt.Sprintf("%t", *r.CexTransfers))
	}

	// Add NextObject if specified
	if r.NextObject != nil {
		values.Add("next_object", *r.NextObject)
	}

	return values.Encode()
}

type TokensPnLResponse struct {
	Items  []TokenPnl `json:"items"`
	Paging Pagination `json:"paging"`
}

type TokenPnl struct {
	Address string           `json:"token_address"`
	Symbol  string           `json:"token_symbol"`
	Name    string           `json:"token_name"`
	Chain   chains.ChainType `json:"chain"`

	TotalBuyUSD             float64 `json:"total_buy_usd"`
	TotalBuyAmount          float64 `json:"total_buy_amount"`
	TotalSellUSD            float64 `json:"total_sell_usd"`
	TotalSellAmount         float64 `json:"total_sell_amount"`
	AverageBuyPrice         float64 `json:"average_buy_price"`
	AverageSellPrice        float64 `json:"average_sell_price"`
	TotalPnlUSD             float64 `json:"total_pnl_usd"`
	RoiPercentage           float64 `json:"roi_percentage"`
	UnrealizedPnlUSD        float64 `json:"unrealized_pnl_usd"`
	UnrealizedRoiPercentage float64 `json:"unrealized_roi_percentage"`
	NumSwaps                int     `json:"num_swaps"`
	IsHoneypot              bool    `json:"is_honeypot"`
}
