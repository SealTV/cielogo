package apiv1_test

import (
	"testing"

	"github.com/sealtv/cielogo/api/apiv1"
	"github.com/stretchr/testify/assert"
)

func TestFeedRequest_GetQueryString_MaxUSD(t *testing.T) {
	maxUSD := 1000.50
	req := &apiv1.FeedRequest{
		MaxUSD: &maxUSD,
	}

	queryString := req.GetQueryString()
	assert.Contains(t, queryString, "maxUSD=1000.5")
}

func TestFeedRequest_GetQueryString_IncludeMarketCap(t *testing.T) {
	includeMarketCap := true
	req := &apiv1.FeedRequest{
		IncludeMarketCap: &includeMarketCap,
	}

	queryString := req.GetQueryString()
	assert.Contains(t, queryString, "includeMarketCap=true")
}

func TestFeedRequest_GetQueryString_IncludeMarketCapFalse(t *testing.T) {
	includeMarketCap := false
	req := &apiv1.FeedRequest{
		IncludeMarketCap: &includeMarketCap,
	}

	queryString := req.GetQueryString()
	assert.Contains(t, queryString, "includeMarketCap=false")
}

func TestFeedRequest_GetQueryString_BothNewParams(t *testing.T) {
	maxUSD := 5000.0
	includeMarketCap := true
	req := &apiv1.FeedRequest{
		MaxUSD:           &maxUSD,
		IncludeMarketCap: &includeMarketCap,
	}

	queryString := req.GetQueryString()
	assert.Contains(t, queryString, "maxUSD=5000")
	assert.Contains(t, queryString, "includeMarketCap=true")
}

func TestFeedRequest_GetQueryString_MinAndMaxUSD(t *testing.T) {
	minUSD := 100.0
	maxUSD := 10000.0
	req := &apiv1.FeedRequest{
		MinUSD: &minUSD,
		MaxUSD: &maxUSD,
	}

	queryString := req.GetQueryString()
	assert.Contains(t, queryString, "minUSD=100")
	assert.Contains(t, queryString, "maxUSD=10000")
}

func TestFeedRequest_GetQueryString_NilValues(t *testing.T) {
	req := &apiv1.FeedRequest{
		MaxUSD:           nil,
		IncludeMarketCap: nil,
	}

	queryString := req.GetQueryString()
	assert.NotContains(t, queryString, "maxUSD")
	assert.NotContains(t, queryString, "includeMarketCap")
}
