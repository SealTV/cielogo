package apiv1

import "fmt"

type RelatedWalletsSorting string

const (
	RelatedWalletsSortingInflowAsc        RelatedWalletsSorting = "inflow_asc"
	RelatedWalletsSortingInflowDesc       RelatedWalletsSorting = "inflow_desc"
	RelatedWalletsSortingOutflowAsc       RelatedWalletsSorting = "outflow_asc"
	RelatedWalletsSortingOutflowDesc      RelatedWalletsSorting = "outflow_desc"
	RelatedWalletsSortingTransactionsAsc  RelatedWalletsSorting = "transactions_asc"
	RelatedWalletsSortingTransactionsDesc RelatedWalletsSorting = "transactions_desc"
)

type RelatedWalletsRequest struct {
	Wallet       string                 `json:"wallet"`
	SortCriteria *RelatedWalletsSorting `json:"sort_criteria,omitempty"`
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
