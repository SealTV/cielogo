package cielogo

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/sealtv/cielogo/api"
	"github.com/sealtv/cielogo/api/apiv1"
)

// GetFeedV1 retrieves the transaction feed for a wallet or list.
//
// Cost: 5 credits per request (3 credits when filtered by wallet).
// WARNING: Setting IncludeMarketCap to true doubles the cost (10 or 6 credits).
//
// https://developer.cielo.finance/reference/getfeed
func (c *Client) GetFeedV1(ctx context.Context, req *apiv1.FeedRequest) (*apiv1.FeedResponse, error) {
	resp := api.CieloResponse[apiv1.FeedResponse]{}

	path := fmt.Sprintf("/v1/feed/?%s", req.GetQueryString())
	if err := c.makeRequest(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, fmt.Errorf("failed to get feed: %w", err)
	}

	return &resp.Data, nil
}

// GetNftsPnlV1 retrieves NFT profit and loss data for a wallet.
//
// Cost: 5 credits per request
//
// https://developer.cielo.finance/reference/getnftspnl
func (c *Client) GetNftsPnlV1(ctx context.Context, req *apiv1.NftsPnLRequest) (*apiv1.NftsPnLResponse, error) {
	resp := api.CieloResponse[apiv1.NftsPnLResponse]{}

	path := fmt.Sprintf("/v1/%s/pnl/nfts?%s", req.Wallet, req.GetQueryString())
	if err := c.makeRequest(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, fmt.Errorf("failed to get  pnl: %w", err)
	}

	return &resp.Data, nil
}

// GetTokensPnlV1 retrieves token profit and loss data for a wallet.
//
// Cost: 5 credits per request
//
// https://developer.cielo.finance/reference/gettokenspnl
func (c *Client) GetTokensPnlV1(ctx context.Context, req *apiv1.TokensPnLRequest) (*apiv1.TokensPnLResponse, error) {
	resp := api.CieloResponse[apiv1.TokensPnLResponse]{}

	path := fmt.Sprintf("/v1/%s/pnl/tokens?%s", req.Wallet, req.GetQueryString())
	if err := c.makeRequest(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, fmt.Errorf("failed to get tokens pnl: %w", err)
	}

	return &resp.Data, nil
}

// GetAggregatedTokenPnLV1 retrieves aggregated token statistics and performance metrics for a wallet.
//
// Cost: 20 credits per request
//
// https://developer.cielo.finance/reference/gettotalstats
func (c *Client) GetAggregatedTokenPnLV1(ctx context.Context, req *apiv1.AggregatedTokenPnLRequest) (*apiv1.AggregatedTokenPnLResponse, error) {
	resp := api.CieloResponse[apiv1.AggregatedTokenPnLResponse]{}

	path := fmt.Sprintf("/v1/%s/pnl/total-stats?%s", req.Wallet, req.GetQueryString())
	if err := c.makeRequest(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, fmt.Errorf("failed to get tokens pnl: %w", err)
	}

	return &resp.Data, nil
}

// GetRelatedWalletsV1 finds wallets that have transacted with the specified wallet.
//
// Cost: 10 credits per request
//
// https://developer.cielo.finance/reference/getrelatedwalletsl
func (c *Client) GetRelatedWalletsV1(ctx context.Context, req *apiv1.RelatedWalletsRequest) (*apiv1.RelatedWalletsResponse, error) {
	resp := api.CieloResponse[apiv1.RelatedWalletsResponse]{}

	path := fmt.Sprintf("/v1/%s/related-wallets/?%s", req.Wallet, req.GetQueryString())
	if err := c.makeRequest(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, fmt.Errorf("failed to get tokens pnl: %w", err)
	}

	return &resp.Data, nil
}

// GetWalletTagsV1 retrieves tags for a single wallet.
//
// Deprecated: This endpoint is deprecated by the Cielo Finance API.
// Use GetWalletsTagsV1 instead, which supports batch operations for up to 50 wallets.
//
// Cost: 5 credits per request
//
// https://developer.cielo.finance/reference/getwallettags
func (c *Client) GetWalletTagsV1(ctx context.Context, req *apiv1.GetWalletTagsRequest) (*apiv1.GetWalletTagsResponse, error) {
	resp := api.CieloResponse[apiv1.GetWalletTagsResponse]{}

	path := fmt.Sprintf("/v1/%s/tags", req.Wallet)
	if err := c.makeRequest(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, fmt.Errorf("failed to get wallet tags: %w", err)
	}

	return &resp.Data, nil
}

// GetWalletsTagsV1 retrieves tags for multiple wallets in a single request (up to 50 wallets).
//
// Cost: 5 credits per request
//
// https://developer.cielo.finance/reference/getwalletstags
func (c *Client) GetWalletsTagsV1(ctx context.Context, req *apiv1.GetWalletsTagsRequest) ([]apiv1.WalletTags, error) {
	resp := api.CieloResponse[[]apiv1.WalletTags]{}

	values := url.Values{}
	for _, wallet := range req.Wallets {
		values.Add("wallet", wallet)
	}

	path := fmt.Sprintf("/v1/tags?%s", values.Encode())

	if err := c.makeRequest(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, fmt.Errorf("failed to get wallets tags: %w", err)
	}

	return resp.Data, nil
}

// GetWalletsByTagV1 retrieves all wallets that have a specific tag.
//
// Cost: 10 credits per request
//
// https://developer.cielo.finance/reference/getwalletsbytag
func (c *Client) GetWalletsByTagV1(ctx context.Context, req *apiv1.GetWalletsByTagRequest) (*apiv1.GetWalletsByTagResponse, error) {
	resp := api.CieloResponse[apiv1.GetWalletsByTagResponse]{}

	values := url.Values{}
	for _, tag := range req.Tags {
		values.Add("tags", string(tag))
	}

	if req.WalletType != nil && *req.WalletType != "" {
		values.Add("wallet_type", string(*req.WalletType))
	}

	if req.Limit != nil {
		values.Add("limit", strconv.Itoa(*req.Limit))
	}

	if req.NextObject != nil && *req.NextObject != "" {
		values.Add("next_object", *req.NextObject)
	}

	path := fmt.Sprintf("/v1/tags/wallets?%s", values.Encode())

	if err := c.makeRequest(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, fmt.Errorf("failed to get wallets by tag: %w", err)
	}

	return &resp.Data, nil
}

// Wallet Lists

// GetAllWalletsListV1 retrieves all public wallet lists with optional filtering and sorting.
//
// Cost: 5 credits per request
//
// https://developer.cielo.finance/reference/getalllists
func (c *Client) GetAllWalletsListV1(ctx context.Context, req *apiv1.GetAllWalletsListsRequest) (*apiv1.GetAllWalletsListsResponse, error) {
	resp := api.CieloResponse[apiv1.GetAllWalletsListsResponse]{}

	values := url.Values{}

	if req.Order != nil && *req.Order != "" {
		values.Add("order", string(*req.Order))
	}

	if req.FollowOnly {
		values.Add("follow_only", "true")
	}

	if req.NextObject != nil && *req.NextObject != "" {
		values.Add("next_object", *req.NextObject)
	}

	path := fmt.Sprintf("/v1/lists/all?%s", values.Encode())

	if err := c.makeRequest(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, fmt.Errorf("failed to get all wallets list: %w", err)
	}

	return &resp.Data, nil
}

// GetUserWalletsListsV1 retrieves all wallet lists owned by the authenticated user.
//
// Cost: 5 credits per request
//
// https://developer.cielo.finance/reference/getuserlists
func (c *Client) GetUserWalletsListsV1(ctx context.Context) ([]apiv1.WalletList, error) {
	resp := api.CieloResponse[[]apiv1.WalletList]{}

	const path = "/v1/lists"

	if err := c.makeRequest(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, fmt.Errorf("failed to get users lists: %w", err)
	}

	return resp.Data, nil
}

// AddWalletsListV1 creates a new wallet list.
//
// Cost: 5 credits per request
//
// https://developer.cielo.finance/reference/adduserlist
func (c *Client) AddWalletsListV1(ctx context.Context, req *apiv1.AddWalletsListRequest) (*apiv1.WalletList, error) {
	const path = "/v1/lists"

	resp := api.CieloResponse[apiv1.WalletList]{}

	if err := c.makeRequest(ctx, http.MethodPost, path, req, &resp); err != nil {
		return nil, fmt.Errorf("failed to add wallet to list: %w", err)
	}

	return &resp.Data, nil
}

// UpdateWalletsListV1 updates an existing wallet list's properties.
//
// Cost: 5 credits per request
//
// https://developer.cielo.finance/reference/updateuserlist
func (c *Client) UpdateWalletsListV1(ctx context.Context, req *apiv1.UpdateWalletsListRequest) (*apiv1.WalletList, error) {
	resp := api.CieloResponse[apiv1.WalletList]{}
	path := fmt.Sprintf("/v1/lists/%d", req.ListID)

	if err := c.makeRequest(ctx, http.MethodPut, path, req, &resp); err != nil {
		return nil, fmt.Errorf("failed to update wallet list: %w", err)
	}

	return &resp.Data, nil
}

// DeleteWalletsListV1 deletes a wallet list.
// If deleteWallets is true, all wallets in the list will also be deleted.
// If deleteWallets is false, only the list itself is deleted, wallets are preserved.
//
// Cost: 5 credits per request
//
// https://developer.cielo.finance/reference/deleteuserlist
func (c *Client) DeleteWalletsListV1(ctx context.Context, listID int64, deleteWallets bool) error {
	path := fmt.Sprintf("/v1/lists/%d", listID)

	// Add query parameter if deleteWallets is true
	if deleteWallets {
		path += "?delete_wallets=true"
	}

	if err := c.makeRequest(ctx, http.MethodDelete, path, nil, nil); err != nil {
		return fmt.Errorf("failed to delete wallet list: %w", err)
	}

	return nil
}

// ToggleFollowWalletsListV1 toggles the follow status of a public wallet list.
//
// Cost: 5 credits per request
//
// https://developer.cielo.finance/reference/togglefollowlist
func (c *Client) ToggleFollowWalletsListV1(ctx context.Context, listID int64) (*apiv1.ToggleFollowWalletsListResponce, error) {
	resp := apiv1.ToggleFollowWalletsListResponce{}

	path := fmt.Sprintf("/v1/lists/%d/toggle-follow", listID)

	if err := c.makeRequest(ctx, http.MethodPut, path, nil, &resp); err != nil {
		return nil, fmt.Errorf("failed to toggle follow wallet list: %w", err)
	}

	return &resp, nil
}

// GetTrackedWalletsV1 retrieves all wallets being tracked by the user, with optional filtering by list.
//
// Cost: 5 credits per request
//
// https://developer.cielo.finance/reference/gettrackedwallets
func (c *Client) GetTrackedWalletsV1(ctx context.Context, req *apiv1.GetTrackedWalletsRequest) (*apiv1.GetTrackedWalletsResponse, error) {
	const path = "/v1/tracked-wallets"

	values := url.Values{}

	if req.ListID != nil && *req.ListID != 0 {
		values.Add("list_id", strconv.FormatInt(*req.ListID, 10))
	}

	if req.NextObject != nil && *req.NextObject != "" {
		values.Add("next_object", *req.NextObject)
	}

	resp := api.CieloResponse[apiv1.GetTrackedWalletsResponse]{}
	if err := c.makeRequest(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, fmt.Errorf("failed to get tracked wallets: %w", err)
	}

	return &resp.Data, nil
}

// AddTrackedWalletsV1 adds a new wallet to tracking with optional notification settings.
//
// Cost: 5 credits per request
//
// https://developer.cielo.finance/reference/addtrackedwallet
func (c *Client) AddTrackedWalletsV1(ctx context.Context, req *apiv1.AddTrackedWalletRequest) (*apiv1.TrackedWallet, error) {
	const path = "/v1/tracked-wallets"

	resp := api.CieloResponse[apiv1.TrackedWallet]{}
	if err := c.makeRequest(ctx, http.MethodPost, path, req, &resp); err != nil {
		return nil, fmt.Errorf("failed to add tracked wallet: %w", err)
	}

	return &resp.Data, nil
}

// RemoveTrackedWalletsV1 removes one or more wallets from tracking by their IDs.
//
// Cost: 5 credits per request
//
// https://developer.cielo.finance/reference/removetrackedwallets
func (c *Client) RemoveTrackedWalletsV1(ctx context.Context, req *apiv1.RemoveTrackedWalletsRequest) error {
	const path = "/v1/tracked-wallets"

	if err := c.makeRequest(ctx, http.MethodDelete, path, req, nil); err != nil {
		return fmt.Errorf("failed to delete tracked wallet: %w", err)
	}

	return nil
}

// GetWalletByAddressV1 retrieves a tracked wallet by its wallet address.
// Returns 404 if the wallet is not being tracked.
//
// Cost: 5 credits per request
//
// https://developer.cielo.finance/reference/getWalletByAddress
func (c *Client) GetWalletByAddressV1(ctx context.Context, wallet string) (*apiv1.TrackedWallet, error) {
	resp := api.CieloResponse[apiv1.TrackedWallet]{}

	path := fmt.Sprintf("/v1/tracked-wallets/address/%s", wallet)
	if err := c.makeRequest(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, fmt.Errorf("failed to get wallet by address: %w", err)
	}

	return &resp.Data, nil
}

// UpdateTrackedWalletV1 updates a tracked wallet by its wallet ID.
//
// Cost: 5 credits per request
//
// https://developer.cielo.finance/reference/updateTrackedWalletV1
func (c *Client) UpdateTrackedWalletV1(ctx context.Context, walletID int64, req *apiv1.UpdateTrackedWalletRequest) (*apiv1.TrackedWallet, error) {
	resp := api.CieloResponse[apiv1.TrackedWallet]{}

	path := fmt.Sprintf("/v1/tracked-wallets/%d", walletID)
	if err := c.makeRequest(ctx, http.MethodPut, path, req, &resp); err != nil {
		return nil, fmt.Errorf("failed to update tracked wallet: %w", err)
	}

	return &resp.Data, nil
}

// UpdateTrackedWalletV2 updates a tracked wallet by its wallet address with support for partial updates.
// All fields are optional and only provided fields will be updated.
//
// Cost: 5 credits per request
//
// https://developer.cielo.finance/reference/updateTrackedWalletV2
func (c *Client) UpdateTrackedWalletV2(ctx context.Context, wallet string, req *apiv1.UpdateTrackedWalletV2Request) (*apiv1.TrackedWallet, error) {
	resp := api.CieloResponse[apiv1.TrackedWallet]{}

	path := fmt.Sprintf("/v2/tracked-wallets/%s", wallet)
	if err := c.makeRequest(ctx, http.MethodPut, path, req, &resp); err != nil {
		return nil, fmt.Errorf("failed to update tracked wallet v2: %w", err)
	}

	return &resp.Data, nil
}

// GetTelegramBotsV1 retrieves the list of available Telegram bots for notifications.
// Only returns bots where available=true.
//
// Cost: 5 credits per request
//
// https://developer.cielo.finance/reference/getTelegramBots
func (c *Client) GetTelegramBotsV1(ctx context.Context) (*apiv1.GetTelegramBotsResponse, error) {
	resp := api.CieloResponse[apiv1.GetTelegramBotsResponse]{}

	const path = "/v1/tracked-wallets/telegram-bots"
	if err := c.makeRequest(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, fmt.Errorf("failed to get telegram bots: %w", err)
	}

	return &resp.Data, nil
}

// Portfolio

// GetWalletPortfolioV1 retrieves the portfolio of a wallet including token balances and USD values.
// For Solana wallets, the response also includes the native SOL balance.
// Portfolio assets with a total_usd_value of zero are excluded from the response.
//
// Supported chains: Solana, EVM (Ethereum, Base, HyperEVM)
// Cost: 20 credits per request
//
// https://developer.cielo.finance/reference/getWalletPortfolio
func (c *Client) GetWalletPortfolioV1(ctx context.Context, wallet string) (*apiv1.WalletPortfolioResponse, error) {
	resp := api.CieloResponse[apiv1.WalletPortfolioResponse]{}

	path := fmt.Sprintf("/v1/%s/portfolio", wallet)
	if err := c.makeRequest(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, fmt.Errorf("failed to get wallet portfolio: %w", err)
	}

	return &resp.Data, nil
}

// GetWalletPortfolioV2 retrieves the portfolio of one or multiple wallets with optional token filtering.
// Supports comma-separated wallet addresses for aggregated portfolio view.
// Each token in the response includes a wallet_address field indicating which wallet holds it.
//
// When a token parameter is provided (Solana wallets only):
//   - Only supported for Solana wallets - returns error for EVM/Sui wallets
//   - Returns only the specified token's balance information
//   - Response format changes to a single token object instead of portfolio object
//   - Returns 404 if the token is not found in the wallet
//
// When multiple wallets are provided:
//   - Token filtering is not supported and will return an error
//   - Portfolios are fetched in parallel for better performance
//   - Tokens are sorted by total_usd_value in descending order
//   - Total USD and chain distributions are aggregated across all wallets
//
// Supported chains: Solana, EVM (Ethereum, Base, HyperEVM), Sui
// Cost: 20 credits per wallet
//
// https://developer.cielo.finance/reference/getWalletPortfolioV2
func (c *Client) GetWalletPortfolioV2(ctx context.Context, req *apiv1.WalletPortfolioV2Request) (*apiv1.WalletPortfolioV2Response, error) {
	resp := api.CieloResponse[apiv1.WalletPortfolioV2Response]{}

	path := fmt.Sprintf("/v2/portfolio?%s", req.GetQueryString())
	if err := c.makeRequest(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, fmt.Errorf("failed to get wallet portfolio v2: %w", err)
	}

	return &resp.Data, nil
}

// Token Information

// GetTokenMetadataV1 retrieves detailed metadata for a specified token including
// name, symbol, decimals, creation details, social links, and supply information.
//
// Supported chains: solana, ethereum, base, hyperevm
// Cost: 1 credit per request
//
// https://developer.cielo.finance/reference/getTokenMetadata
func (c *Client) GetTokenMetadataV1(ctx context.Context, req *apiv1.TokenMetadataRequest) (*apiv1.TokenMetadataResponse, error) {
	resp := api.CieloResponse[apiv1.TokenMetadataResponse]{}

	path := fmt.Sprintf("/v1/token/metadata?%s", req.GetQueryString())
	if err := c.makeRequest(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, fmt.Errorf("failed to get token metadata: %w", err)
	}

	return &resp.Data, nil
}

// GetTokenPriceV1 retrieves the current price in USD for a specified token.
//
// Supported chains: solana, ethereum, base, hyperevm
// Cost: 1 credit per request
//
// https://developer.cielo.finance/reference/getTokenPrice
func (c *Client) GetTokenPriceV1(ctx context.Context, req *apiv1.TokenPriceRequest) (*apiv1.TokenPriceResponse, error) {
	resp := api.CieloResponse[apiv1.TokenPriceResponse]{}

	path := fmt.Sprintf("/v1/token/price?%s", req.GetQueryString())
	if err := c.makeRequest(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, fmt.Errorf("failed to get token price: %w", err)
	}

	return &resp.Data, nil
}

// GetTokenStatsV1 retrieves comprehensive statistics for a specified token including
// price changes, volume metrics, market cap, and unique trader counts across
// multiple time periods (5m, 1h, 6h, 24h).
//
// Supported chains: solana, ethereum, base, hyperevm
// Cost: 3 credits per request
//
// https://developer.cielo.finance/reference/getTokenStats
func (c *Client) GetTokenStatsV1(ctx context.Context, req *apiv1.TokenStatsRequest) (*apiv1.TokenStatsResponse, error) {
	resp := api.CieloResponse[apiv1.TokenStatsResponse]{}

	path := fmt.Sprintf("/v1/token/stats?%s", req.GetQueryString())
	if err := c.makeRequest(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, fmt.Errorf("failed to get token stats: %w", err)
	}

	return &resp.Data, nil
}

// GetTokenBalanceV1 retrieves the balance of a specific token for a given wallet,
// including the token's current price and total USD value.
//
// Supported chains: solana, ethereum, base, hyperevm
// Cost: 3 credits per request
//
// https://developer.cielo.finance/reference/getTokenBalance
func (c *Client) GetTokenBalanceV1(ctx context.Context, req *apiv1.TokenBalanceRequest) (*apiv1.TokenBalanceResponse, error) {
	resp := api.CieloResponse[apiv1.TokenBalanceResponse]{}

	path := fmt.Sprintf("/v1/%s/token-balance?%s", req.Wallet, req.GetQueryString())
	if err := c.makeRequest(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, fmt.Errorf("failed to get token balance: %w", err)
	}

	return &resp.Data, nil
}

// Trading Analytics

// GetTradingStatsV1 retrieves detailed performance statistics for a wallet's trading activity,
// including PnL, ROI, win rate, and trading behavior insights.
//
// Note: This endpoint may return 202 Accepted if data is not ready yet.
// In that case, retry the request after 10 seconds.
//
// Cost: 30 credits per request
//
// https://developer.cielo.finance/reference/getTradingStats
func (c *Client) GetTradingStatsV1(ctx context.Context, req *apiv1.TradingStatsRequest) (*apiv1.TradingStatsResponse, error) {
	resp := api.CieloResponse[apiv1.TradingStatsResponse]{}

	path := fmt.Sprintf("/v1/%s/trading-stats", req.Wallet)
	queryString := req.GetQueryString()
	if queryString != "" {
		path = fmt.Sprintf("%s?%s", path, queryString)
	}

	if err := c.makeRequest(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, fmt.Errorf("failed to get trading stats: %w", err)
	}

	return &resp.Data, nil
}
