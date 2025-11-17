package apiv1

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/sealtv/cielogo/api/chains"
)

// TokensPnLRequest is used to retrieve token profit and loss data for a wallet.
type TokensPnLRequest struct {
	// Wallet is the wallet address to get PnL data for (required).
	Wallet string
	// Chains filters PnL by specific blockchain chains.
	Chains []chains.ChainType
	// Timeframe specifies the time period for PnL calculation (e.g., "7d", "30d", "all").
	Timeframe *string
	// NextObject is the pagination cursor from the previous response.
	NextObject *string
	// CexTransfers includes centralized exchange transfers in PnL calculations when true.
	CexTransfers *bool
	// Tokens filters PnL by specific token addresses or symbols.
	Tokens []string
	// ActivePositionsOnly filters for positions with balance > 0 when true.
	// Set to true to see only currently held tokens, false or nil to see all historical positions.
	ActivePositionsOnly *bool
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
		values.Add("cex_transfers", strconv.FormatBool(*r.CexTransfers))
	}

	// Add NextObject if specified
	if r.NextObject != nil {
		values.Add("next_object", *r.NextObject)
	}

	// Add Tokens if specified
	if len(r.Tokens) > 0 {
		values.Add("token", strings.Join(r.Tokens, ","))
	}

	// Add ActivePositionsOnly if specified
	if r.ActivePositionsOnly != nil {
		values.Add("active_positions_only", strconv.FormatBool(*r.ActivePositionsOnly))
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

	TotalBuyUSD    float64 `json:"total_buy_usd"`
	TotalBuyAmount float64 `json:"total_buy_amount"`

	TotalSellUSD    float64 `json:"total_sell_usd"`
	TotalSellAmount float64 `json:"total_sell_amount"`

	AverageBuyPrice  float64 `json:"average_buy_price"`
	AverageSellPrice float64 `json:"average_sell_price"`

	TotalPnlUSD   float64 `json:"total_pnl_usd"`
	RoiPercentage float64 `json:"roi_percentage"`

	UnrealizedPnlUSD        float64 `json:"unrealized_pnl_usd"`
	UnrealizedRoiPercentage float64 `json:"unrealized_roi_percentage"`

	TokenPrice float64 `json:"token_price"`

	NumSwaps   int  `json:"num_swaps"`
	IsHoneypot bool `json:"is_honeypot"`
}
