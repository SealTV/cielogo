package apiv1

import (
	"net/url"
	"strings"
)

// PortfolioAsset represents a single token asset in a wallet portfolio.
type PortfolioAsset struct {
	Chain         string  `json:"chain"`
	TokenAddress  string  `json:"token_address"`
	TokenName     string  `json:"token_name"`
	TokenSymbol   string  `json:"token_symbol"`
	TokenPriceUSD float64 `json:"token_price_usd"`
	Balance       float64 `json:"balance"`
	TotalUSDValue float64 `json:"total_usd_value"`
	Supply        float64 `json:"supply"`
	SupplyOwned   float64 `json:"supply_owned"`
}

// WalletPortfolioResponse represents the response for a single wallet portfolio.
// Used by the V1 portfolio endpoint.
type WalletPortfolioResponse struct {
	LastUpdated   int64            `json:"last_updated"`
	TotalUSDValue float64          `json:"total_usd_value,omitempty"`
	Portfolio     []PortfolioAsset `json:"portfolio,omitempty"`
}

// PortfolioAssetV2 extends PortfolioAsset with wallet address information
// for multi-wallet portfolio queries.
type PortfolioAssetV2 struct {
	Chain         string  `json:"chain"`
	TokenAddress  string  `json:"token_address"`
	TokenName     string  `json:"token_name"`
	TokenSymbol   string  `json:"token_symbol"`
	TokenPriceUSD float64 `json:"token_price_usd"`
	Balance       float64 `json:"balance"`
	TotalUSDValue float64 `json:"total_usd_value"`
	Supply        float64 `json:"supply"`
	SupplyOwned   float64 `json:"supply_owned"`
	WalletAddress string  `json:"wallet_address"`
}

// ChainDistribution represents the token distribution across different blockchains
// in an aggregated portfolio.
type ChainDistribution struct {
	Chain      string  `json:"chain"`
	TotalUSD   float64 `json:"total_usd"`
	Percentage float64 `json:"percentage"`
}

// WalletPortfolioV2Response represents the response for multi-wallet portfolio queries.
// Supports aggregation across multiple wallets and includes chain distribution.
type WalletPortfolioV2Response struct {
	LastUpdated   int64               `json:"last_updated"`
	TotalUSDValue float64             `json:"total_usd_value"`
	Portfolio     []PortfolioAssetV2  `json:"portfolio"`
	Chains        []ChainDistribution `json:"chains,omitempty"`
}

// WalletPortfolioV2Request represents a request for multi-wallet portfolio data.
// Supports comma-separated wallet addresses and optional token filtering (Solana only).
type WalletPortfolioV2Request struct {
	// Wallets is a list of wallet addresses to retrieve portfolios for.
	// Multiple wallets will be aggregated in the response.
	Wallets []string

	// Token is an optional token mint address to filter for a specific token.
	// Only supported for single Solana wallet queries.
	// Returns 400 error if used with multiple wallets or non-Solana wallets.
	Token *string
}

// GetQueryString builds the query string for V2 portfolio requests.
// Wallets are joined with commas, and token filter is added if present.
func (r *WalletPortfolioV2Request) GetQueryString() string {
	values := url.Values{}

	if len(r.Wallets) > 0 {
		values.Add("wallet", strings.Join(r.Wallets, ","))
	}

	if r.Token != nil && *r.Token != "" {
		values.Add("token", *r.Token)
	}

	return values.Encode()
}
