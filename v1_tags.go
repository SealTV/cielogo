package cielogo

import (
	"context"
	"fmt"
	"net/http"

	"github.com/sealtv/cielogo/types"
	"github.com/sealtv/cielogo/types/apiv1"
)

func (c *Client) GetTagsV1(ctx context.Context, wallet string) ([]apiv1.Tag, error) {
	resp := types.CieloResponse[apiv1.TagsResponse]{}

	path := fmt.Sprintf("/v1/%s/tags", wallet)
	if err := c.makeRequest(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, fmt.Errorf("failed to get tags: %w", err)
	}

	return resp.Data.Tags, nil
}
