package cielogo

import (
	"context"
	"fmt"
	"net/http"

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
func (c *Client) GetWalletTagsV1(ctx context.Context, req apiv1.GetWalletTagsRequest) (*apiv1.GetWalletTagsResponse, error) {
	resp := api.CieloResponse[apiv1.GetWalletTagsResponse]{}

	path := fmt.Sprintf("/v1/%s/tags", req.Wallet)
	if err := c.makeRequest(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, fmt.Errorf("failed to get wallet tags: %w", err)
	}

	return &resp.Data, nil
}
