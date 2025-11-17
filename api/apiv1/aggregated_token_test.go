package apiv1_test

import (
	"testing"

	"github.com/sealtv/cielogo/api/apiv1"
	"github.com/sealtv/cielogo/api/chains"
	"github.com/stretchr/testify/assert"
)

func TestAggregatedTokenPnLTimeframe_Constants(t *testing.T) {
	tests := []struct {
		name      string
		timeframe apiv1.AggregatedTokenPnLTimeframe
		expected  string
	}{
		{"1 day", apiv1.AggregatedTokenPnLTimeframe1Day, "1d"},
		{"7 days", apiv1.AggregatedTokenPnLTimeframe7Day, "7d"},
		{"30 days", apiv1.AggregatedTokenPnLTimeframe30Day, "30d"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, string(tt.timeframe))
		})
	}
}

func TestAggregatedTokenPnLRequest_GetQueryString_Empty(t *testing.T) {
	req := &apiv1.AggregatedTokenPnLRequest{
		Wallet: "0x123",
	}

	queryString := req.GetQueryString()
	assert.Equal(t, "", queryString)
}

func TestAggregatedTokenPnLRequest_GetQueryString_WithChains(t *testing.T) {
	req := &apiv1.AggregatedTokenPnLRequest{
		Wallet: "0x123",
		Chains: []chains.ChainType{chains.Ethereum, chains.Solana},
	}

	queryString := req.GetQueryString()
	assert.Contains(t, queryString, "chains=ethereum%2Csolana")
}

func TestAggregatedTokenPnLRequest_GetQueryString_WithTimeframe(t *testing.T) {
	timeframe := apiv1.AggregatedTokenPnLTimeframe7Day
	req := &apiv1.AggregatedTokenPnLRequest{
		Wallet:    "0x123",
		Timeframe: &timeframe,
	}

	queryString := req.GetQueryString()
	assert.Contains(t, queryString, "timeframe=7d")
}

func TestAggregatedTokenPnLRequest_GetQueryString_WithCexTransfers(t *testing.T) {
	cexTransfers := true
	req := &apiv1.AggregatedTokenPnLRequest{
		Wallet:       "0x123",
		CexTransfers: &cexTransfers,
	}

	queryString := req.GetQueryString()
	assert.Contains(t, queryString, "cex_transfers=true")
}

func TestAggregatedTokenPnLRequest_GetQueryString_AllParameters(t *testing.T) {
	timeframe := apiv1.AggregatedTokenPnLTimeframe30Day
	cexTransfers := false
	req := &apiv1.AggregatedTokenPnLRequest{
		Wallet:       "0x123",
		Chains:       []chains.ChainType{chains.Ethereum, chains.Base, chains.Solana},
		Timeframe:    &timeframe,
		CexTransfers: &cexTransfers,
	}

	queryString := req.GetQueryString()
	assert.Contains(t, queryString, "chains=ethereum%2Cbase%2Csolana")
	assert.Contains(t, queryString, "timeframe=30d")
	assert.Contains(t, queryString, "cex_transfers=false")
}

func TestAggregatedTokenPnLRequest_GetQueryString_NilValues(t *testing.T) {
	req := &apiv1.AggregatedTokenPnLRequest{
		Wallet:       "0x123",
		Timeframe:    nil,
		CexTransfers: nil,
	}

	queryString := req.GetQueryString()
	assert.NotContains(t, queryString, "timeframe")
	assert.NotContains(t, queryString, "cex_transfers")
}
