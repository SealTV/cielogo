package apiv1

import (
	"net/url"
	"strconv"
)

// TokenChain represents the blockchain networks supported by token endpoints.
type TokenChain string

const (
	TokenChainSolana    TokenChain = "solana"
	TokenChainEthereum  TokenChain = "ethereum"
	TokenChainBase      TokenChain = "base"
	TokenChainHyperevm  TokenChain = "hyperevm"
	TokenChainRobinhood TokenChain = "robinhood"
	TokenChainArc       TokenChain = "arc"
)

// Token Metadata

// TokenMetadataRequest represents a request for token metadata.
type TokenMetadataRequest struct {
	Chain            TokenChain
	TokenAddress     string
	IncludePools     *bool
	IncludeLiquidity *bool
}

// GetQueryString builds the query string for token metadata requests.
func (r *TokenMetadataRequest) GetQueryString() string {
	values := url.Values{}
	values.Add("chain", string(r.Chain))
	values.Add("token_address", r.TokenAddress)
	return values.Encode()
}

// TokenMetadataResponse represents detailed metadata for a token including
// social links, supply information, and creation details.
type TokenMetadataResponse struct {
	Decimals              uint8    `json:"decimals"`
	Supply                uint64   `json:"supply"`
	TotalSupply           uint64   `json:"total_supply"`
	TotalSupplyString     string   `json:"total_supply_string"`
	Platform              string   `json:"platform"`
	PlatformUI            string   `json:"platform_ui"`
	PlatformURI           string   `json:"platform_uri"`
	LogoURI               string   `json:"logo_uri"`
	Chain                 string   `json:"chain"`
	Address               string   `json:"address"`
	Name                  string   `json:"name"`
	Symbol                string   `json:"symbol"`
	DeployerAddress       string   `json:"deployer_address"`
	CreatedAtBlock        uint64   `json:"created_at_block"`
	CreatedAtTs           int64    `json:"created_at_ts"`
	Website               string   `json:"website"`
	Twitter               string   `json:"twitter"`
	Telegram              string   `json:"telegram"`
	MigratedAt            int64    `json:"migrated_at"`
	MigratedFromPool      string   `json:"migrated_from_pool"`
	MigratedToPool        string   `json:"migrated_to_pool"`
	MigrationTx           string   `json:"migration_tx"`
	MigrationPriceUSD     *float64 `json:"migration_price_usd"`
	MigrationPriceNative  *float64 `json:"migration_price_native"`
	MigrationMarketCapUSD *float64 `json:"migration_market_cap_usd"`
	IsMayhemMode          bool     `json:"is_mayhem_mode"`
}

// Token Price

// TokenPriceRequest represents a request for current token price.
type TokenPriceRequest struct {
	Chain            TokenChain
	TokenAddress     string
	IncludePools     *bool
	IncludeLiquidity *bool
}

// GetQueryString builds the query string for token price requests.
func (r *TokenPriceRequest) GetQueryString() string {
	values := url.Values{}
	values.Add("chain", string(r.Chain))
	values.Add("token_address", r.TokenAddress)
	return values.Encode()
}

// TokenPriceResponse represents the current price of a token in USD.
type TokenPriceResponse struct {
	Chain       string  `json:"chain"`
	Address     string  `json:"address"`
	BlockNumber uint64  `json:"block_number"`
	Price       float64 `json:"price"`
}

// Token Statistics

// VolumeStats represents trading volume statistics for a specific time period.
type VolumeStats struct {
	VolumeUSD     float64 `json:"volume_usd"`
	BuyVolumeUSD  float64 `json:"buy_volume_usd"`
	SellVolumeUSD float64 `json:"sell_volume_usd"`
	UniqueBuyers  int     `json:"unique_buyers"`
	UniqueSellers int     `json:"unique_sellers"`
	Buys          int     `json:"buys"`
	Sells         int     `json:"sells"`
}

// PriceChange represents price percentage changes over different time periods.
type PriceChange struct {
	FiveMin     float64 `json:"5m"`
	OneHour     float64 `json:"1h"`
	SixHours    float64 `json:"6h"`
	TwentyFourH float64 `json:"24h"`
}

// VolumeBreakdown represents trading volume statistics broken down by time periods.
type VolumeBreakdown struct {
	FiveMin     VolumeStats `json:"5m"`
	OneHour     VolumeStats `json:"1h"`
	SixHours    VolumeStats `json:"6h"`
	TwentyFourH VolumeStats `json:"24h"`
}

// TokenStatsRequest represents a request for comprehensive token statistics.
type TokenStatsRequest struct {
	Chain            TokenChain
	TokenAddress     string
	IncludePools     *bool
	IncludeLiquidity *bool
}

// GetQueryString builds the query string for token stats requests.
func (r *TokenStatsRequest) GetQueryString() string {
	values := url.Values{}
	values.Add("chain", string(r.Chain))
	values.Add("token_address", r.TokenAddress)
	if r.IncludePools != nil {
		values.Add("include_pools", strconv.FormatBool(*r.IncludePools))
	}
	if r.IncludeLiquidity != nil {
		values.Add("include_liquidity", strconv.FormatBool(*r.IncludeLiquidity))
	}
	return values.Encode()
}

// TokenStatsResponse represents comprehensive statistics for a token including
// price changes, market cap, and trading volume metrics across multiple time periods.
type TokenStatsResponse struct {
	TokenAddress                string          `json:"token_address"`
	PriceUSD                    float64         `json:"price_usd"`
	MarketCapUSD                float64         `json:"market_cap_usd"`
	AthPriceUSD                 float64         `json:"ath_price_usd"`
	AthMarketCapUSD             float64         `json:"ath_market_cap_usd"`
	AthDate                     int64           `json:"ath_date,omitempty"`
	TotalFeesSOL                *float64        `json:"total_fees_sol,omitempty"`
	TotalFeesUSD                *float64        `json:"total_fees_usd,omitempty"`
	TaxRate                     *float64        `json:"tax_rate,omitempty"`
	Change                      PriceChange     `json:"change"`
	Volume                      VolumeBreakdown `json:"volume"`
	Pools                       []Pool          `json:"pools,omitempty"`
	MigrationPriceUSD           *float64        `json:"migration_price_usd"`
	MigrationPriceNative        *float64        `json:"migration_price_native"`
	MigrationMarketCapUSD       *float64        `json:"migration_market_cap_usd"`
	TargetMigrationPriceUSD     *float64        `json:"target_migration_price_usd"`
	TargetMigrationPriceNative  *float64        `json:"target_migration_price_native"`
	TargetMigrationMarketCapUSD *float64        `json:"target_migration_market_cap_usd"`
	Twitter                     string          `json:"twitter"`
	XLinkType                   string          `json:"x_link_type"`
	ActiveUsers                 uint64          `json:"active_users"`
}

// Token Balance

// TokenBalanceRequest represents a request for a specific token balance in a wallet.
type TokenBalanceRequest struct {
	Wallet       string
	TokenAddress string
	Chain        TokenChain
}

// GetQueryString builds the query string for token balance requests.
func (r *TokenBalanceRequest) GetQueryString() string {
	values := url.Values{}
	values.Add("token_address", r.TokenAddress)
	values.Add("chain", string(r.Chain))
	return values.Encode()
}

// TokenBalanceResponse represents the balance of a specific token in a wallet.
type TokenBalanceResponse struct {
	Chain         string  `json:"chain"`
	TokenAddress  string  `json:"token_address"`
	TokenName     string  `json:"token_name"`
	TokenSymbol   string  `json:"token_symbol"`
	TokenPriceUSD float64 `json:"token_price_usd"`
	Balance       float64 `json:"balance"`
	TotalUSDValue float64 `json:"total_usd_value"`
	// Deprecated: the server does not return decimals for token balances.
	Decimals    int     `json:"decimals,omitempty"`
	Supply      float64 `json:"supply"`
	SupplyOwned float64 `json:"supply_owned"`
}

// Pool is a liquidity pool returned when IncludePools is enabled.
type Pool struct {
	LikelySpam        bool     `json:"likely_spam"`
	LiquidityNative   float64  `json:"liquidity_native"`
	Liquidity         float64  `json:"liquidity"`
	LiquidityUSD      float64  `json:"liquidity_usd"`
	PairedToken       string   `json:"paired_token"`
	PoolAddress       string   `json:"pool_address"`
	Protocol          string   `json:"protocol"`
	ProtocolUI        string   `json:"protocol_ui"`
	PairedTokenSymbol string   `json:"paired_token_symbol"`
	PairedTokenName   string   `json:"paired_token_name"`
	TaxRate           *float64 `json:"tax_rate"`
}
