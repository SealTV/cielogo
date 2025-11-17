package apiv1_test

import (
	"testing"

	"github.com/sealtv/cielogo/api/apiv1"
	"github.com/sealtv/cielogo/api/chains"
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

func TestFeedRequest_GetQueryString_Wallet(t *testing.T) {
	req := &apiv1.FeedRequest{
		Wallet: "0x123",
	}

	queryString := req.GetQueryString()
	assert.Contains(t, queryString, "wallet=0x123")
}

func TestFeedRequest_GetQueryString_Limit(t *testing.T) {
	limit := 50
	req := &apiv1.FeedRequest{
		Limit: &limit,
	}

	queryString := req.GetQueryString()
	assert.Contains(t, queryString, "limit=50")
}

func TestFeedRequest_GetQueryString_List(t *testing.T) {
	list := 123
	req := &apiv1.FeedRequest{
		List: &list,
	}

	queryString := req.GetQueryString()
	assert.Contains(t, queryString, "list=123")
}

func TestFeedRequest_GetQueryString_Chains(t *testing.T) {
	req := &apiv1.FeedRequest{
		Chains: []chains.ChainType{chains.Ethereum, chains.Solana},
	}

	queryString := req.GetQueryString()
	assert.Contains(t, queryString, "chains=ethereum%2Csolana")
}

func TestFeedRequest_GetQueryString_TxTypes(t *testing.T) {
	req := &apiv1.FeedRequest{
		TxTypes: []apiv1.TxType{apiv1.TxTypeSwap, apiv1.TxTypeTransfer},
	}

	queryString := req.GetQueryString()
	assert.Contains(t, queryString, "tx_types=swap%2Ctransfer")
}

func TestFeedRequest_GetQueryString_Tokens(t *testing.T) {
	req := &apiv1.FeedRequest{
		Tokens: []string{"USDC", "SOL", "ETH"},
	}

	queryString := req.GetQueryString()
	assert.Contains(t, queryString, "tokens=USDC%2CSOL%2CETH")
}

func TestFeedRequest_GetQueryString_NewTrades(t *testing.T) {
	newTrades := true
	req := &apiv1.FeedRequest{
		NewTrades: &newTrades,
	}

	queryString := req.GetQueryString()
	assert.Contains(t, queryString, "newTrades=true")
}

func TestFeedRequest_GetQueryString_StartFrom(t *testing.T) {
	startFrom := "cursor123"
	req := &apiv1.FeedRequest{
		StartFrom: &startFrom,
	}

	queryString := req.GetQueryString()
	assert.Contains(t, queryString, "startFrom=cursor123")
}

func TestFeedRequest_GetQueryString_FromTimestamp(t *testing.T) {
	timestamp := int64(1699564800)
	req := &apiv1.FeedRequest{
		FromTimestamp: &timestamp,
	}

	queryString := req.GetQueryString()
	assert.Contains(t, queryString, "fromTimestamp=1699564800")
}

func TestFeedRequest_GetQueryString_ToTimestamp(t *testing.T) {
	timestamp := int64(1699651200)
	req := &apiv1.FeedRequest{
		ToTimestamp: &timestamp,
	}

	queryString := req.GetQueryString()
	assert.Contains(t, queryString, "toTimestamp=1699651200")
}
