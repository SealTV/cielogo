package apiv1

type WalletType string

const (
	UnknownWalletType WalletType = "unknown"
	EvmWalletType     WalletType = "evm"
	SolanaWalletType  WalletType = "solana"
	DydxWalletType    WalletType = "dydx"
	BitcoinWalletType WalletType = "bitcoin"
	TronWalletType    WalletType = "tron"
)

type GetTrackedWalletsRequest struct {
	NextObject *string
	ListID     *int64
}

type AddTrackedWalletRequest struct {
	Wallet string `json:"wallet"`
	Label  string `json:"label"`
	ListID *int64 `json:"list_id,omitempty"`
}

type RemoveTrackedWalletsRequest struct {
	WalletIDs []int64 `json:"wallet_ids"`
}

type GetTrackedWalletsResponse struct {
	TrackedWallets []TrackedWallet         `json:"tracked_wallets"`
	Pagination     TrackedWalletPagination `json:"paging"`
}

type TrackedWalletPagination struct {
	HasNextPage     bool `json:"has_next_page"`
	TotalRowsInPage int  `json:"total_rows_in_page"`
	NextObject      int  `json:"next_object"`
}

type TrackedWallet struct {
	ID     int64       `json:"id"`
	Wallet string      `json:"wallet"`
	Label  string      `json:"label"`
	Type   WalletType  `json:"type"`
	ListID *int64      `json:"list_id,omitempty"`
	List   *WalletList `json:"list,omitempty"`
}
