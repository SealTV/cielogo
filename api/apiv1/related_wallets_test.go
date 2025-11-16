package apiv1_test

import (
	"testing"

	"github.com/sealtv/cielogo/api/apiv1"
	"github.com/stretchr/testify/assert"
)

func TestRelatedWalletsSorting_Constants(t *testing.T) {
	tests := []struct {
		name     string
		sorting  apiv1.RelatedWalletsSorting
		expected string
	}{
		{"inflow_asc", apiv1.RelatedWalletsSortingInflowAsc, "inflow_asc"},
		{"inflow_desc", apiv1.RelatedWalletsSortingInflowDesc, "inflow_desc"},
		{"outflow_asc", apiv1.RelatedWalletsSortingOutflowAsc, "outflow_asc"},
		{"outflow_desc", apiv1.RelatedWalletsSortingOutflowDesc, "outflow_desc"},
		{"transactions_asc", apiv1.RelatedWalletsSortingTransactionsAsc, "transactions_asc"},
		{"transactions_desc", apiv1.RelatedWalletsSortingTransactionsDesc, "transactions_desc"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, string(tt.sorting))
		})
	}
}

func TestRelatedWalletsRequest_GetQueryString_WithSorting(t *testing.T) {
	sortCriteria := apiv1.RelatedWalletsSortingInflowDesc
	req := &apiv1.RelatedWalletsRequest{
		Wallet:       "0x123",
		SortCriteria: &sortCriteria,
	}

	queryString := req.GetQueryString()
	assert.Equal(t, "sort_criteria=inflow_desc", queryString)
}

func TestRelatedWalletsRequest_GetQueryString_AllSortOptions(t *testing.T) {
	tests := []struct {
		name         string
		sortCriteria apiv1.RelatedWalletsSorting
		expected     string
	}{
		{"inflow_asc", apiv1.RelatedWalletsSortingInflowAsc, "sort_criteria=inflow_asc"},
		{"inflow_desc", apiv1.RelatedWalletsSortingInflowDesc, "sort_criteria=inflow_desc"},
		{"outflow_asc", apiv1.RelatedWalletsSortingOutflowAsc, "sort_criteria=outflow_asc"},
		{"outflow_desc", apiv1.RelatedWalletsSortingOutflowDesc, "sort_criteria=outflow_desc"},
		{"transactions_asc", apiv1.RelatedWalletsSortingTransactionsAsc, "sort_criteria=transactions_asc"},
		{"transactions_desc", apiv1.RelatedWalletsSortingTransactionsDesc, "sort_criteria=transactions_desc"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &apiv1.RelatedWalletsRequest{
				Wallet:       "0x123",
				SortCriteria: &tt.sortCriteria,
			}

			queryString := req.GetQueryString()
			assert.Equal(t, tt.expected, queryString)
		})
	}
}

func TestRelatedWalletsRequest_GetQueryString_NilSorting(t *testing.T) {
	req := &apiv1.RelatedWalletsRequest{
		Wallet:       "0x123",
		SortCriteria: nil,
	}

	queryString := req.GetQueryString()
	assert.Equal(t, "", queryString)
}
