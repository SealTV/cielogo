package apiv1

import "fmt"

// RelatedWalletsSorting defines the sorting criteria for related wallets.
// Use these constants to sort wallets by inflow, outflow, or transaction count.
type RelatedWalletsSorting string

const (
	// RelatedWalletsSortingInflowAsc sorts wallets by inflow USD in ascending order (lowest first).
	RelatedWalletsSortingInflowAsc RelatedWalletsSorting = "inflow_asc"
	// RelatedWalletsSortingInflowDesc sorts wallets by inflow USD in descending order (highest first).
	RelatedWalletsSortingInflowDesc RelatedWalletsSorting = "inflow_desc"
	// RelatedWalletsSortingOutflowAsc sorts wallets by outflow USD in ascending order (lowest first).
	RelatedWalletsSortingOutflowAsc RelatedWalletsSorting = "outflow_asc"
	// RelatedWalletsSortingOutflowDesc sorts wallets by outflow USD in descending order (highest first).
	RelatedWalletsSortingOutflowDesc RelatedWalletsSorting = "outflow_desc"
	// RelatedWalletsSortingTransactionsAsc sorts wallets by transaction count in ascending order (fewest first).
	RelatedWalletsSortingTransactionsAsc RelatedWalletsSorting = "transactions_asc"
	// RelatedWalletsSortingTransactionsDesc sorts wallets by transaction count in descending order (most first).
	RelatedWalletsSortingTransactionsDesc RelatedWalletsSorting = "transactions_desc"
)

// RelatedWalletsRequest is used to find wallets that have transacted with the specified wallet.
type RelatedWalletsRequest struct {
	// Wallet is the wallet address to find related wallets for (required).
	Wallet string `json:"wallet"`
	// SortCriteria determines how the results are sorted (optional).
	SortCriteria *RelatedWalletsSorting `json:"sort_criteria,omitempty"`
}

func (r *RelatedWalletsRequest) GetQueryString() string {
	if r.SortCriteria != nil {
		return fmt.Sprintf("sort_criteria=%s", *r.SortCriteria)
	}

	return ""
}

// RelatedWalletsResponse contains the list of wallets that have transacted with the queried wallet.
type RelatedWalletsResponse struct {
	// RelatedWallets is the list of wallets with transaction relationships.
	RelatedWallets []RelatedWallet `json:"related_wallets"`
}

// RelatedWallet represents a wallet that has had transactions with the queried wallet.
type RelatedWallet struct {
	// Wallet is the address of the related wallet.
	Wallet string `json:"wallet"`
	// Label is the custom label for this wallet, if any.
	Label string `json:"label"`
	// InflowUSD is the total USD value of inbound transactions from this wallet.
	InflowUSD float64 `json:"inflow_usd"`
	// OutflowUSD is the total USD value of outbound transactions to this wallet.
	OutflowUSD float64 `json:"outflow_usd"`
	// TotalUSD is the combined USD value of all transactions (inflow + outflow).
	TotalUSD float64 `json:"total_usd"`
	// Count is the total number of transactions with this wallet.
	Count int `json:"count"`
}
