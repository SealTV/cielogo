package apiv1

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/sealtv/cielogo/api/chains"
)

type AggregatedTokenPnLTimeframe string

const (
	AggregatedTokenPnLTimeframe1Day  AggregatedTokenPnLTimeframe = "1d"
	AggregatedTokenPnLTimeframe7Day  AggregatedTokenPnLTimeframe = "7d"
	AggregatedTokenPnLTimeframe30Day AggregatedTokenPnLTimeframe = "30d"
)

type AggregatedTokenPnLRequest struct {
	Wallet       string
	Chains       []chains.ChainType
	Timeframe    *AggregatedTokenPnLTimeframe
	CexTransfers *bool
}

func (r *AggregatedTokenPnLRequest) GetQueryString() string {
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
		values.Add("timeframe", string(*r.Timeframe))
	}

	// Add CexTransfers if specified
	if r.CexTransfers != nil {
		values.Add("cex_transfers", fmt.Sprintf("%t", *r.CexTransfers))
	}

	return values.Encode()
}

type AggregatedTokenPnLResponse struct {
	Wallet       string `json:"wallet"`
	TokensTraded int    `json:"tokens_traded"`

	RealizedPnlUsd        float64 `json:"realized_pnl_usd"`
	RealizedRoiPercentage float64 `json:"realized_roi_percentage"`

	UnrealizedPnlUsd        float64 `json:"unrealized_pnl_usd"`
	UnrealizedRoiPercentage float64 `json:"unrealized_roi_percentage"`

	CombinedPnlUsd        float64 `json:"combined_pnl_usd"`
	CombinedRoiPercentage float64 `json:"combined_roi_percentage"`

	Winrate            float64 `json:"winrate"`
	AverageHoldingTime float64 `json:"average_holding_time"`
}
