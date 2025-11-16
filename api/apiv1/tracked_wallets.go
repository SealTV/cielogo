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

// UpdateTrackedWalletRequest is used to update a tracked wallet by its ID (V1).
type UpdateTrackedWalletRequest struct {
	Wallet string `json:"wallet"`
	Label  string `json:"label"`
	ListID *int64 `json:"list_id,omitempty"`
}

// UpdateTrackedWalletV2Request is used to update a tracked wallet by its address (V2).
// All fields are optional and only provided fields will be updated.
type UpdateTrackedWalletV2Request struct {
	Label          *string  `json:"label,omitempty"`
	ListID         *int     `json:"list_id,omitempty"`
	MinUSD         *float64 `json:"min_usd,omitempty"`
	TxTypes        []string `json:"tx_types,omitempty"`
	Chains         []string `json:"chains,omitempty"`
	NewTrades      *bool    `json:"new_trades,omitempty"`
	TelegramBot    *string  `json:"telegram_bot,omitempty"`
	DiscordChannel *string  `json:"discord_channel,omitempty"`
}

// TelegramBot represents a Telegram bot available for notifications.
type TelegramBot struct {
	ID       interface{} `json:"id"` // Can be string or int
	Label    string      `json:"label"`
	Name     string      `json:"name"`
	Disabled bool        `json:"disabled"`
	Plan     string      `json:"plan"`
	Link     string      `json:"link"`
	Default  bool        `json:"default"`
	IsCustom bool        `json:"is_custom"`
	Usage    string      `json:"usage"`
}

// GetTelegramBotsResponse is the response from the GetTelegramBots endpoint.
type GetTelegramBotsResponse struct {
	Bots []TelegramBot `json:"bots"`
}
