package apiv1

import (
	"net/url"
)

// TokenChain represents the blockchain networks supported by token endpoints.
type TokenChain string

const (
	TokenChainSolana   TokenChain = "solana"
	TokenChainEthereum TokenChain = "ethereum"
	TokenChainBase     TokenChain = "base"
	TokenChainHyperevm TokenChain = "hyperevm"
)

// Token Metadata

// TokenMetadataRequest represents a request for token metadata.
type TokenMetadataRequest struct {
	Chain        TokenChain
	TokenAddress string
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
	Chain          string `json:"chain"`
	Address        string `json:"address"`
	Name           string `json:"name"`
	Symbol         string `json:"symbol"`
	Decimals       int    `json:"decimals"`
	CreatedAtTs    int64  `json:"created_at_ts"`
	CreatedAtBlock int64  `json:"created_at_block"`
	Twitter        string `json:"twitter,omitempty"`
	Telegram       string `json:"telegram,omitempty"`
	Website        string `json:"website,omitempty"`
	Supply         int64  `json:"supply"`
	LogoURI        string `json:"logo_uri"`
}

// Token Price

// TokenPriceRequest represents a request for current token price.
type TokenPriceRequest struct {
	Chain        TokenChain
	TokenAddress string
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
	BlockNumber int64   `json:"block_number"`
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
	Chain        TokenChain
	TokenAddress string
}

// GetQueryString builds the query string for token stats requests.
func (r *TokenStatsRequest) GetQueryString() string {
	values := url.Values{}
	values.Add("chain", string(r.Chain))
	values.Add("token_address", r.TokenAddress)
	return values.Encode()
}

// TokenStatsResponse represents comprehensive statistics for a token including
// price changes, market cap, and trading volume metrics across multiple time periods.
type TokenStatsResponse struct {
	TokenAddress string          `json:"token_address"`
	PriceUSD     float64         `json:"price_usd"`
	MarketCapUSD float64         `json:"market_cap_usd"`
	Change       PriceChange     `json:"change"`
	Volume       VolumeBreakdown `json:"volume"`
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
	Decimals      int     `json:"decimals"`
}
