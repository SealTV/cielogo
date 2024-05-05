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
	Wallet                  string  `json:"wallet"`
	PnlUsd                  float64 `json:"pnl_usd"`
	RoiPercentage           float64 `json:"roi_percentage"`
	TokensTraded            float64 `json:"tokens_traded"`
	UnrealizedPnlUsd        float64 `json:"unrealized_pnl_usd"`
	UnrealizedRoiPercentage float64 `json:"unrealized_roi_percentage"`
	Winrate                 float64 `json:"winrate"`
}
