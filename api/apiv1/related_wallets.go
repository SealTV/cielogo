package apiv1

import "fmt"

type RelatedWalletsSoring string

const (
	RelatedWalletsSoringInflowAsc        RelatedWalletsSoring = "inflow_asc"
	RelatedWalletsSoringInflowDesc       RelatedWalletsSoring = "inflow_desc"
	RelatedWalletsSoringOutflowAsc       RelatedWalletsSoring = "outflow_asc"
	RelatedWalletsSoringOutflowDesc      RelatedWalletsSoring = "outflow_desc"
	RelatedWalletsSoringTransactionsAsc  RelatedWalletsSoring = "transactions_asc"
	RelatedWalletsSoringTransactionsDesc RelatedWalletsSoring = "transactions_desc"
)

type RelatedWalletsRequest struct {
	Wallet       string                `json:"wallet"`
	SortCriteria *RelatedWalletsSoring `json:"sort_criteria,omitempty"`
}

func (r *RelatedWalletsRequest) GetQueryString() string {
	if r.SortCriteria != nil {
		return fmt.Sprintf("sort_criteria=%s", *r.SortCriteria)
	}

	return ""
}

type RelatedWalletsResponse struct {
	RelatedWallets []RelatedWallet `json:"related_wallets"`
}

type RelatedWallet struct {
	Wallet     string  `json:"wallet"`
	Label      string  `json:"label"`
	InflowUSD  float64 `json:"inflow_usd"`
	OutflowUSD float64 `json:"outflow_usd"`
	TotalUSD   float64 `json:"total_usd"`
	Count      int     `json:"count"`
}
