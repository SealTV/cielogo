package apiv1_test

import (
	"testing"

	"github.com/sealtv/cielogo/api/apiv1"
	"github.com/sealtv/cielogo/api/chains"
	"github.com/stretchr/testify/assert"
)

func TestTokensPnLRequest_GetQueryString_ActivePositionsOnly(t *testing.T) {
	activeOnly := true
	req := &apiv1.TokensPnLRequest{
		Wallet:              "0x123",
		ActivePositionsOnly: &activeOnly,
	}

	queryString := req.GetQueryString()
	assert.Contains(t, queryString, "active_positions_only=true")
}

func TestTokensPnLRequest_GetQueryString_ActivePositionsOnlyFalse(t *testing.T) {
	activeOnly := false
	req := &apiv1.TokensPnLRequest{
		Wallet:              "0x123",
		ActivePositionsOnly: &activeOnly,
	}

	queryString := req.GetQueryString()
	assert.Contains(t, queryString, "active_positions_only=false")
}

func TestTokensPnLRequest_GetQueryString_WithAllFilters(t *testing.T) {
	activeOnly := true
	cexTransfers := false
	timeframe := "7d"
	req := &apiv1.TokensPnLRequest{
		Wallet:              "0x123",
		Chains:              []chains.ChainType{chains.Ethereum, chains.Solana},
		Timeframe:           &timeframe,
		CexTransfers:        &cexTransfers,
		Tokens:              []string{"USDC", "SOL"},
		ActivePositionsOnly: &activeOnly,
	}

	queryString := req.GetQueryString()
	assert.Contains(t, queryString, "chains=ethereum%2Csolana")
	assert.Contains(t, queryString, "timeframe=7d")
	assert.Contains(t, queryString, "cex_transfers=false")
	assert.Contains(t, queryString, "token=USDC%2CSOL")
	assert.Contains(t, queryString, "active_positions_only=true")
}

func TestTokensPnLRequest_GetQueryString_NilActivePositionsOnly(t *testing.T) {
	req := &apiv1.TokensPnLRequest{
		Wallet:              "0x123",
		ActivePositionsOnly: nil,
	}

	queryString := req.GetQueryString()
	assert.NotContains(t, queryString, "active_positions_only")
}
