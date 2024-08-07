package apiv1

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/sealtv/cielogo/api/chains"
)

type FeedRequest struct {
	// Filter the feed by a specific wallet address.
	Wallet string
	// Limit the number of transactions returned in the feed.
	// The maximum limit is 100.
	Limit *int
	// Filter transactions by a specific List ID.
	List *int
	// Filter transactions by specific blockchain chains (e.g. ethereum), comma-separated for multiple values (e.g, ethereum,polygon)
	Chains []chains.ChainType
	// Filter transactions by types (e.g. swap, nft_trade), comma-separated for multiple values (e.g, swap,transfer,nft_trade)
	TxTypes []TxType
	// Filter transactions by specific tokens, identified by either their address or symbol,
	// comma-separated for multiple values (e.g, LINK,BITCOIN)
	Tokens []string
	// Set a minimum USD value for transactions. Default - 0
	MinUSD *float64
	// Filter transactions by new trades.
	NewTrades *bool
	// Set value from response 'paging.next_object_id' to get next page.
	StartFrom *string
	// Filter transactions from a specific UNIX timestamp.
	FromTimestamp *int64
	// Filter transactions to a specific UNIX timestamp.
	ToTimestamp *int64
}

func (r *FeedRequest) GetQueryString() string {
	values := url.Values{}

	if r.Wallet != "" {
		values.Add("wallet", r.Wallet)
	}

	if r.Limit != nil {
		values.Add("limit", fmt.Sprintf("%d", *r.Limit))
	}

	if r.List != nil {
		values.Add("list", fmt.Sprintf("%d", *r.List))
	}

	if len(r.Chains) > 0 {
		chainStrs := make([]string, 0, len(r.Chains))
		for _, chain := range r.Chains {
			chainStrs = append(chainStrs, string(chain))
		}

		values.Add("chains", strings.Join(chainStrs, ","))
	}

	if len(r.TxTypes) > 0 {
		txTypeStrs := make([]string, 0, len(r.TxTypes))
		for _, txType := range r.TxTypes {
			txTypeStrs = append(txTypeStrs, string(txType))
		}

		values.Add("tx_types", strings.Join(txTypeStrs, ","))
	}

	if len(r.Tokens) > 0 {
		values.Add("tokens", strings.Join(r.Tokens, ","))
	}

	if r.MinUSD != nil {
		values.Add("minUSD", fmt.Sprintf("%f", *r.MinUSD))
	}

	if r.NewTrades != nil {
		values.Add("newTrades", fmt.Sprintf("%t", *r.NewTrades))
	}

	if r.StartFrom != nil {
		values.Add("startFrom", *r.StartFrom)
	}

	if r.FromTimestamp != nil {
		values.Add("fromTimestamp", fmt.Sprintf("%d", *r.FromTimestamp))
	}

	if r.ToTimestamp != nil {
		values.Add("toTimestamp", fmt.Sprintf("%d", *r.ToTimestamp))
	}

	return values.Encode()
}

type FeedResponse struct {
	Items  []TxEvent  `json:"items"`
	Paging Pagination `json:"paging"`
}
