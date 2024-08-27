package cielogo

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/sealtv/cielogo/api"
	"github.com/sealtv/cielogo/api/apiv1"
)

func (c *Client) GetFeedV1(ctx context.Context, req *apiv1.FeedRequest) (*apiv1.FeedResponse, error) {
	resp := api.CieloResponse[apiv1.FeedResponse]{}

	path := fmt.Sprintf("/v1/feed/?%s", req.GetQueryString())
	if err := c.makeRequest(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, fmt.Errorf("failed to get feed: %w", err)
	}

	return &resp.Data, nil
}

// GetNftsPnlV1 returns a list of nft pnl for a given wallet.
// https://developer.cielo.finance/reference/getnftspnl
func (c *Client) GetNftsPnlV1(ctx context.Context, req *apiv1.NftsPnLRequest) (*apiv1.NftsPnLResponse, error) {
	resp := api.CieloResponse[apiv1.NftsPnLResponse]{}

	path := fmt.Sprintf("/v1/%s/pnl/nfts?%s", req.Wallet, req.GetQueryString())
	if err := c.makeRequest(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, fmt.Errorf("failed to get  pnl: %w", err)
	}

	return &resp.Data, nil
}

// GetTokensPnlV1 returns a list of tokens pnl for a given wallet.
// https://developer.cielo.finance/reference/gettokenspnl
func (c *Client) GetTokensPnlV1(ctx context.Context, req *apiv1.TokensPnLRequest) (*apiv1.TokensPnLResponse, error) {
	resp := api.CieloResponse[apiv1.TokensPnLResponse]{}

	path := fmt.Sprintf("/v1/%s/pnl/tokens?%s", req.Wallet, req.GetQueryString())
	if err := c.makeRequest(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, fmt.Errorf("failed to get tokens pnl: %w", err)
	}

	return &resp.Data, nil
}

// GetAggregatedTokenPnLV1 returns a list of aggregated token pnl for a given wallet.
// https://developer.cielo.finance/reference/gettotalstats
func (c *Client) GetAggregatedTokenPnLV1(ctx context.Context, req *apiv1.AggregatedTokenPnLRequest) (*apiv1.AggregatedTokenPnLResponse, error) {
	resp := api.CieloResponse[apiv1.AggregatedTokenPnLResponse]{}

	path := fmt.Sprintf("/v1/%s/pnl/total-stats?%s", req.Wallet, req.GetQueryString())
	if err := c.makeRequest(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, fmt.Errorf("failed to get tokens pnl: %w", err)
	}

	return &resp.Data, nil
}

// GetRelatedWalletsV1 returns a list of related wallets for a given wallet.
// https://developer.cielo.finance/reference/getrelatedwalletsl
func (c *Client) GetRelatedWalletsV1(ctx context.Context, req *apiv1.RelatedWalletsRequest) (*apiv1.RelatedWalletsResponse, error) {
	resp := api.CieloResponse[apiv1.RelatedWalletsResponse]{}

	path := fmt.Sprintf("/v1/%s/related-wallets/?%s", req.Wallet, req.GetQueryString())
	if err := c.makeRequest(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, fmt.Errorf("failed to get tokens pnl: %w", err)
	}

	return &resp.Data, nil
}

// GetWalletTagsV1 returns a list of wallet tags for a given wallet.
// https://developer.cielo.finance/reference/getwallettags
func (c *Client) GetWalletTagsV1(ctx context.Context, req *apiv1.GetWalletTagsRequest) (*apiv1.GetWalletTagsResponse, error) {
	resp := api.CieloResponse[apiv1.GetWalletTagsResponse]{}

	path := fmt.Sprintf("/v1/%s/tags", req.Wallet)
	if err := c.makeRequest(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, fmt.Errorf("failed to get wallet tags: %w", err)
	}

	return &resp.Data, nil
}

// GetWalletsTagsV1 returns a list of wallet tags for a given list of wallets.
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

// GetWalletsByTagV1 returns a list of wallets for a given tag.
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
		values.Add("limit", fmt.Sprintf("%d", *req.Limit))
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

// GetAllWalletsListsV1 returns a list of all wallets.
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

// GetUserWalletsListsV1 returns a list of users wallets.
// https://developer.cielo.finance/reference/getuserlists
func (c *Client) GetUserWalletsListsV1(ctx context.Context) ([]apiv1.WalletList, error) {
	resp := api.CieloResponse[[]apiv1.WalletList]{}

	const path = "/v1/lists"

	if err := c.makeRequest(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, fmt.Errorf("failed to get users lists: %w", err)
	}

	return resp.Data, nil
}

// AddWalletsListV1 adds a wallet to a list.
// https://developer.cielo.finance/reference/adduserlist
func (c *Client) AddWalletsListV1(ctx context.Context, req *apiv1.AddWalletsListRequest) (*apiv1.WalletList, error) {
	const path = "/v1/lists"

	resp := api.CieloResponse[apiv1.WalletList]{}

	if err := c.makeRequest(ctx, http.MethodPost, path, req, &resp); err != nil {
		return nil, fmt.Errorf("failed to add wallet to list: %w", err)
	}

	return &resp.Data, nil
}

// UpdateWalletsListV1 updates a wallet list.
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
// https://developer.cielo.finance/reference/deleteuserlist
func (c *Client) DeleteWalletsListV1(ctx context.Context, listID int64) error {
	path := fmt.Sprintf("/v1/lists/%d", listID)

	if err := c.makeRequest(ctx, http.MethodDelete, path, nil, nil); err != nil {
		return fmt.Errorf("failed to delete wallet list: %w", err)
	}

	return nil
}

// ToggleFollowWalletsListV1 toggles the follow status of a wallet list.
// https://developer.cielo.finance/reference/togglefollowlist
func (c *Client) ToggleFollowWalletsListV1(ctx context.Context, listID int64) (*apiv1.ToggleFollowWalletsListResponce, error) {
	resp := apiv1.ToggleFollowWalletsListResponce{}

	path := fmt.Sprintf("/v1/lists/%d/toggle-follow", listID)

	if err := c.makeRequest(ctx, http.MethodPost, path, nil, &resp); err != nil {
		return nil, fmt.Errorf("failed to toggle follow wallet list: %w", err)
	}

	return &resp, nil
}

// GetTrackedWalletsV1 returns a list of tracked wallets.
// https://developer.cielo.finance/reference/gettrackedwallets
func (c *Client) GetTrackedWalletsV1(ctx context.Context, req *apiv1.GetTrackedWalletsRequest) (*apiv1.GetTrackedWalletsResponse, error) {
	const path = "/v1/tracked-wallets"

	values := url.Values{}

	if req.ListID != nil && *req.ListID != 0 {
		values.Add("list_id", fmt.Sprintf("%d", *req.ListID))
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

// AddTrackedWalletsV1 adds a wallet to tracked wallets.
// https://developer.cielo.finance/reference/addtrackedwallet
func (c *Client) AddTrackedWalletsV1(ctx context.Context, req *apiv1.AddTrackedWalletRequest) (*apiv1.TrackedWallet, error) {
	const path = "/v1/tracked-wallets"

	resp := api.CieloResponse[apiv1.TrackedWallet]{}
	if err := c.makeRequest(ctx, http.MethodPost, path, req, &resp); err != nil {
		return nil, fmt.Errorf("failed to add tracked wallet: %w", err)
	}

	return &resp.Data, nil
}

// RemoveTrackedWalletsV1 deletes a tracked wallet.
// https://developer.cielo.finance/reference/removetrackedwallets
func (c *Client) RemoveTrackedWalletsV1(ctx context.Context, req *apiv1.RemoveTrackedWalletsRequest) error {
	const path = "/v1/tracked-wallets"

	if err := c.makeRequest(ctx, http.MethodDelete, path, req, nil); err != nil {
		return fmt.Errorf("failed to delete tracked wallet: %w", err)
	}

	return nil
}
