package apiv1

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

// PortfolioAsset represents a single token asset in a wallet portfolio.
type PortfolioAsset struct {
	Chain            string  `json:"chain"`
	TokenAddress     string  `json:"token_address"`
	TokenName        string  `json:"token_name"`
	TokenSymbol      string  `json:"token_symbol"`
	TokenPriceUSD    float64 `json:"token_price_usd"`
	Balance          float64 `json:"balance"`
	TotalUSDValue    float64 `json:"total_usd_value"`
	Supply           float64 `json:"supply"`
	SupplyOwned      float64 `json:"supply_owned"`
	Platform         string  `json:"platform"`
	PlatformUI       string  `json:"platform_ui"`
	BalanceRaw       string  `json:"balance_raw"`
	TotalNativeValue float64 `json:"total_native_value"`
	TokenPriceNative float64 `json:"token_price_native"`
	WalletAddress    string  `json:"wallet_address"`
}

// WalletPortfolioResponse represents the response for a single wallet portfolio.
// Used by the V1 portfolio endpoint.
type WalletPortfolioResponse struct {
	LastUpdated   int64               `json:"last_updated"`
	TotalUSDValue float64             `json:"total_usd"`
	Portfolio     []PortfolioAsset    `json:"portfolio,omitempty"`
	PendingChains []string            `json:"pending_chains,omitempty"`
	Chains        []ChainDistribution `json:"chains"`
}

// PortfolioAssetV2 extends PortfolioAsset with wallet address information
// for multi-wallet portfolio queries.
type PortfolioAssetV2 struct {
	Chain            string  `json:"chain"`
	TokenAddress     string  `json:"token_address"`
	TokenName        string  `json:"token_name"`
	TokenSymbol      string  `json:"token_symbol"`
	TokenPriceUSD    float64 `json:"token_price_usd"`
	Balance          float64 `json:"balance"`
	TotalUSDValue    float64 `json:"total_usd_value"`
	Supply           float64 `json:"supply"`
	SupplyOwned      float64 `json:"supply_owned"`
	Platform         string  `json:"platform"`
	PlatformUI       string  `json:"platform_ui"`
	BalanceRaw       string  `json:"balance_raw"`
	TotalNativeValue float64 `json:"total_native_value"`
	TokenPriceNative float64 `json:"token_price_native"`
	WalletAddress    string  `json:"wallet_address"`
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
	TotalUSDValue float64             `json:"total_usd"`
	Portfolio     []PortfolioAssetV2  `json:"portfolio"`
	Chains        []ChainDistribution `json:"chains,omitempty"`
	PendingChains []string            `json:"pending_chains,omitempty"`
	SingleToken   *PortfolioAssetV2   `json:"-"`
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

// UnmarshalJSON distinguishes a token-filtered private response from a full portfolio.
func (r *WalletPortfolioV2Response) UnmarshalJSON(data []byte) error {
	type plain WalletPortfolioV2Response
	var shape struct {
		TokenAddress *string `json:"token_address"`
	}
	if err := json.Unmarshal(data, &shape); err != nil {
		return fmt.Errorf("decode portfolio shape: %w", err)
	}
	if shape.TokenAddress != nil {
		var token PortfolioAssetV2
		if err := json.Unmarshal(data, &token); err != nil {
			return fmt.Errorf("decode portfolio token: %w", err)
		}
		*r = WalletPortfolioV2Response{SingleToken: &token}
		return nil
	}
	var decoded plain
	if err := json.Unmarshal(data, &decoded); err != nil {
		return fmt.Errorf("decode portfolio: %w", err)
	}
	*r = WalletPortfolioV2Response(decoded)
	return nil
}
