package apiv1

// WalletType represents the blockchain type of a wallet.
type WalletType string

const (
	// UnknownWalletType represents an unrecognized wallet type.
	UnknownWalletType WalletType = "unknown"
	// EvmWalletType represents EVM-compatible wallets (Ethereum, Polygon, BSC, etc.).
	EvmWalletType WalletType = "evm"
	// SolanaWalletType represents Solana blockchain wallets.
	SolanaWalletType WalletType = "solana"
	// DydxWalletType represents dYdX protocol wallets.
	DydxWalletType WalletType = "dydx"
	// BitcoinWalletType represents Bitcoin blockchain wallets.
	BitcoinWalletType WalletType = "bitcoin"
	// TronWalletType represents Tron blockchain wallets.
	TronWalletType WalletType = "tron"
)

// GetTrackedWalletsRequest is used to retrieve tracked wallets with pagination.
type GetTrackedWalletsRequest struct {
	// NextObject is the pagination cursor from the previous response.
	NextObject *string
	// ListID filters tracked wallets by a specific list ID.
	ListID *int64
}

// AddTrackedWalletRequest is used to add a new wallet to tracking with optional notification settings.
type AddTrackedWalletRequest struct {
	// Wallet is the wallet address to track (required).
	Wallet string `json:"wallet"`
	// Label is a custom label for the wallet (required).
	Label string `json:"label"`
	// ListID assigns the wallet to a specific list (optional).
	ListID *int64 `json:"list_id,omitempty"`

	// Notification settings (all optional):

	// BundleID is an alternative to ListID for grouping wallets.
	BundleID *int64 `json:"bundle_id,omitempty"`
	// MinAmountUSD sets the minimum USD threshold for transaction notifications.
	// Only transactions with value >= this amount will trigger notifications.
	MinAmountUSD *float64 `json:"min_amount_usd,omitempty"`
	// Filters is a list of transaction type filter IDs to apply.
	Filters []int `json:"filters,omitempty"`
	// Chains is a list of blockchain chain IDs to monitor.
	// Only transactions on these chains will be tracked.
	Chains []int `json:"chains,omitempty"`
	// NewTrades enables notifications for new token trades.
	NewTrades *bool `json:"new_trades,omitempty"`
	// TelegramBotID specifies the Telegram bot ID for sending notifications.
	TelegramBotID *int `json:"telegram_bot_id,omitempty"`
	// DiscordChannelID specifies the Discord channel ID for sending notifications.
	DiscordChannelID *string `json:"discord_channel_id,omitempty"`
}

// RemoveTrackedWalletsRequest is used to remove one or more tracked wallets by their IDs.
type RemoveTrackedWalletsRequest struct {
	// WalletIDs is the list of wallet IDs to remove from tracking.
	WalletIDs []int64 `json:"wallet_ids"`
}

// GetTrackedWalletsResponse contains the list of tracked wallets with pagination information.
type GetTrackedWalletsResponse struct {
	// TrackedWallets is the list of tracked wallets in the current page.
	TrackedWallets []TrackedWallet `json:"tracked_wallets"`
	// Pagination contains information for navigating through pages.
	Pagination TrackedWalletPagination `json:"paging"`
}

// TrackedWalletPagination contains pagination metadata for tracked wallets lists.
type TrackedWalletPagination struct {
	// HasNextPage indicates if there are more pages available.
	HasNextPage bool `json:"has_next_page"`
	// TotalRowsInPage is the number of wallets in the current page.
	TotalRowsInPage int `json:"total_rows_in_page"`
	// NextObject is the cursor value to use for fetching the next page.
	NextObject int `json:"next_object"`
}

// TrackedWallet represents a wallet being tracked for activity monitoring.
type TrackedWallet struct {
	// ID is the unique identifier for this tracked wallet.
	ID int64 `json:"id"`
	// Wallet is the blockchain address being tracked.
	Wallet string `json:"wallet"`
	// Label is the custom label assigned to this wallet.
	Label string `json:"label"`
	// Type is the blockchain type of this wallet.
	Type WalletType `json:"type"`
	// ListID is the ID of the list this wallet belongs to, if any.
	ListID *int64 `json:"list_id,omitempty"`
	// List contains the wallet list details if the wallet is part of a list.
	List *WalletList `json:"list,omitempty"`
}

// UpdateTrackedWalletRequest is used to update a tracked wallet by its ID (V1).
// This method requires the wallet ID and updates all fields (full update).
type UpdateTrackedWalletRequest struct {
	// Wallet is the new wallet address (required).
	Wallet string `json:"wallet"`
	// Label is the new label for the wallet (required).
	Label string `json:"label"`
	// ListID assigns the wallet to a different list (optional).
	ListID *int64 `json:"list_id,omitempty"`
}

// UpdateTrackedWalletV2Request is used to update a tracked wallet by its address (V2).
// All fields are optional and only provided fields will be updated (partial update).
// This is the preferred method for updating specific fields without affecting others.
type UpdateTrackedWalletV2Request struct {
	// Label updates the wallet's custom label.
	Label *string `json:"label,omitempty"`
	// ListID assigns the wallet to a different list.
	ListID *int `json:"list_id,omitempty"`
	// MinUSD sets the minimum USD threshold for transaction notifications.
	MinUSD *float64 `json:"min_usd,omitempty"`
	// TxTypes filters notifications by transaction types.
	TxTypes []string `json:"tx_types,omitempty"`
	// Chains filters notifications by blockchain chains.
	Chains []string `json:"chains,omitempty"`
	// NewTrades enables/disables notifications for new token trades.
	NewTrades *bool `json:"new_trades,omitempty"`
	// TelegramBot specifies the Telegram bot for notifications.
	TelegramBot *string `json:"telegram_bot,omitempty"`
	// DiscordChannel specifies the Discord channel for notifications.
	DiscordChannel *string `json:"discord_channel,omitempty"`
}

// TelegramBot represents a Telegram bot available for sending wallet notifications.
// The bot can be used to receive alerts about tracked wallet activity.
type TelegramBot struct {
	// ID is the bot's unique identifier. Can be either string or int depending on the API response.
	ID interface{} `json:"id"`
	// Label is the display name for the bot.
	Label string `json:"label"`
	// Name is the bot's username (without @).
	Name string `json:"name"`
	// Disabled indicates if the bot is currently disabled.
	Disabled bool `json:"disabled"`
	// Plan indicates the subscription plan associated with this bot (e.g., "free", "pro", "premium").
	Plan string `json:"plan"`
	// Link is the direct URL to start a conversation with the bot.
	Link string `json:"link"`
	// Default indicates if this is the default bot for notifications.
	Default bool `json:"default"`
	// IsCustom indicates if this is a custom user-configured bot.
	IsCustom bool `json:"is_custom"`
	// Usage describes how this bot is being used (e.g., "notifications", "alerts").
	Usage string `json:"usage"`
}

// GetTelegramBotsResponse contains the list of available Telegram bots for notifications.
// Only bots where available=true are returned by the API.
type GetTelegramBotsResponse struct {
	// Bots is the list of available Telegram bots.
	Bots []TelegramBot `json:"bots"`
}
