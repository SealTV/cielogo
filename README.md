# CieloGo

[![CI](https://github.com/SealTV/cielogo/actions/workflows/ci.yml/badge.svg)](https://github.com/SealTV/cielogo/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/SealTV/cielogo/branch/main/graph/badge.svg)](https://codecov.io/gh/SealTV/cielogo)
[![Go Report Card](https://goreportcard.com/badge/github.com/sealtv/cielogo)](https://goreportcard.com/report/github.com/sealtv/cielogo)
[![Go Reference](https://pkg.go.dev/badge/github.com/sealtv/cielogo.svg)](https://pkg.go.dev/github.com/sealtv/cielogo)

CieloGo is a Go client library for interacting with the [Cielo Finance API](https://api-info.cielo.finance/), allowing developers to easily access wallet analytics, pnl stats, related wallets, and wallet tags through a simple and intuitive interface. Built for efficiency and ease of use, CieloGo abstracts the complexity of direct API calls into straightforward Go functions.

Please note, CieloGo is an unofficial client library and is not endorsed or maintained by Cielo Finance.

## Features

- **Transaction Feed** - Get wallet or list transaction feeds with advanced filtering (chain, token, value range, market cap)
- **Portfolio Management** - Retrieve wallet portfolios with token balances and USD values (single or multi-wallet)
- **Token Information** - Access token metadata, real-time prices, statistics, and balances
- **PnL Analytics** - Analyze NFT and token profit/loss with aggregated statistics and performance metrics
- **Trading Statistics** - Get detailed trading performance stats including PnL, ROI, win rates, and more
- **Related Wallets** - Discover wallets with transaction relationships and customizable sorting
- **Wallet Tracking** - Full CRUD operations for tracked wallets with notification support (Telegram/Discord)
- **Wallet Tags** - Tag wallets for organization and filter by custom or system tags
- **Wallet Lists** - Create and manage wallet lists with follow/unfollow functionality
- **WebSocket Support** - Real-time wallet activity monitoring via WebSocket connections

## Installation

To use CieloGo in your Go project, install it as a module:

```bash
go get github.com/sealtv/cielogo
```

## Usage

To start using CieloGo, import it into your Go project:

```go
import "github.com/sealtv/cielogo"
```

### Initializing the Client

```go
client := cielogo.NewClient()
```

### Making Requests

Here are some examples of how you might call various methods of the CieloGo client.

**Get Feed:**

```go
ctx := context.Background()
req := &apiv1.FeedRequest{...} // setup your request
feed, err := client.GetFeedV1(ctx, req)
```

**Get NFTs PnL:**

```go
nftsPnLReq := apiv1.NftsPnLRequest{Wallet: "0xWALLET_ADDRESS"}
nftsPnl, err := client.GetNftsPnlV1(ctx, &nftsPnLReq)
```

**Get Tokens PnL:**

```go
tokensPnLReq := apiv1.TokensPnLRequest{Wallet: "0xWALLET_ADDRESS"}
tokensPnl, err := client.GetTokensPnlV1(ctx, &tokensPnLReq)
```

For more details on each request and response structure, refer to the [Cielo Finance API documentation](https://developer.cielo.finance).

## Breaking Changes

### v0.x.x → v1.0.0

#### 1. RelatedWallets Sorting Type Rename (Typo Fix)

**Breaking Change:** Type and constant names corrected from `RelatedWalletsSoring` to `RelatedWalletsSorting`.

**Migration:**
```go
// Before:
sortCriteria := apiv1.RelatedWalletsSoringInflowDesc

// After:
sortCriteria := apiv1.RelatedWalletsSortingInflowDesc
```

**All Constants Renamed:**
- `RelatedWalletsSoringInflowAsc` → `RelatedWalletsSortingInflowAsc`
- `RelatedWalletsSoringInflowDesc` → `RelatedWalletsSortingInflowDesc`
- `RelatedWalletsSoringOutflowAsc` → `RelatedWalletsSortingOutflowAsc`
- `RelatedWalletsSoringOutflowDesc` → `RelatedWalletsSortingOutflowDesc`
- `RelatedWalletsSoringTransactionsAsc` → `RelatedWalletsSortingTransactionsAsc`
- `RelatedWalletsSoringTransactionsDesc` → `RelatedWalletsSortingTransactionsDesc`

**Automated Fix:**
```bash
find . -type f -name "*.go" -exec sed -i 's/RelatedWalletsSoring/RelatedWalletsSorting/g' {} +
```

#### 2. DeleteWalletsListV1 Signature Change

**Breaking Change:** New required `deleteWallets` parameter controls cascade deletion.

**Migration:**
```go
// Before:
err := client.DeleteWalletsListV1(ctx, listID)

// After (preserve wallets):
err := client.DeleteWalletsListV1(ctx, listID, false)

// Or (delete wallets too):
err := client.DeleteWalletsListV1(ctx, listID, true)
```

**Default Behavior:** Pass `false` to maintain backward compatibility (delete list only, keep wallets).

## API Credits

All API endpoints consume credits from your Cielo Finance account. Below is a comprehensive cost table:

| Endpoint | Cost (Credits) | Notes |
|----------|----------------|-------|
| **Feed & PnL** |
| GetFeedV1 | 5 (3 with wallet filter) | **10/6 with IncludeMarketCap=true** ⚠️ |
| GetNftsPnlV1 | 5 | |
| GetTokensPnlV1 | 5 | |
| GetAggregatedTokenPnLV1 | **20** | Expensive operation |
| **Portfolio & Token Info** |
| GetWalletPortfolioV1 | **20** | |
| GetWalletPortfolioV2 | **20 per wallet** | |
| GetTokenMetadataV1 | 1 | |
| GetTokenPriceV1 | 1 | |
| GetTokenStatsV1 | 3 | |
| GetTokenBalanceV1 | 3 | |
| **Trading Stats** |
| GetTradingStatsV1 | **30** | Most expensive, may return 202 Accepted |
| **Related Wallets** |
| GetRelatedWalletsV1 | 10 | |
| **Tags** |
| GetWalletTagsV1 | 5 | Deprecated, use GetWalletsTagsV1 |
| GetWalletsTagsV1 | 5 | Batch operation (up to 50 wallets) |
| GetWalletsByTagV1 | 10 | |
| **Lists** |
| GetAllWalletsListsV1 | 5 | |
| GetUserWalletsListsV1 | 5 | |
| AddWalletsListV1 | 5 | |
| UpdateWalletsListV1 | 5 | |
| DeleteWalletsListV1 | 5 | |
| ToggleFollowWalletsListV1 | 5 | |
| **Tracked Wallets** |
| GetTrackedWalletsV1 | 5 | |
| AddTrackedWalletsV1 | 5 | |
| RemoveTrackedWalletsV1 | 5 | |
| GetWalletByAddressV1 | 5 | |
| UpdateTrackedWalletV1 | 5 | By ID (full update) |
| UpdateTrackedWalletV2 | 5 | By address (partial update) |
| GetTelegramBotsV1 | 5 | |

### Cost Optimization Tips

1. **Use batch operations** - `GetWalletsTagsV1` supports up to 50 wallets for 5 credits (vs 5 credits per wallet)
2. **Filter feeds by wallet** - Reduces cost from 5 to 3 credits
3. **Avoid IncludeMarketCap** unless necessary - Doubles feed costs (5→10 or 3→6)
4. **Cache expensive results** - Portfolio (20 credits) and Trading Stats (30 credits) change infrequently
5. **Use V2 partial updates** - `UpdateTrackedWalletV2` allows updating specific fields without full replacement

## Contributing

Contributions to the CieloGo project are welcome. Please feel free to report any bugs, suggest features, or open pull requests.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
