# CieloGo

[![CI](https://github.com/SealTV/cielogo/actions/workflows/ci.yml/badge.svg)](https://github.com/SealTV/cielogo/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/SealTV/cielogo/branch/main/graph/badge.svg)](https://codecov.io/gh/SealTV/cielogo)
[![Go Report Card](https://goreportcard.com/badge/github.com/sealtv/cielogo)](https://goreportcard.com/report/github.com/sealtv/cielogo)
[![Go Reference](https://pkg.go.dev/badge/github.com/sealtv/cielogo.svg)](https://pkg.go.dev/github.com/sealtv/cielogo)

CieloGo is a Go client library for interacting with the [Cielo Finance API](https://api-info.cielo.finance/), allowing developers to easily access wallet analytics, pnl stats, related wallets, and wallet tags through a simple and intuitive interface. Built for efficiency and ease of use, CieloGo abstracts the complexity of direct API calls into straightforward Go functions.

Please note, CieloGo is an unofficial client library and is not endorsed or maintained by Cielo Finance.

## Features

- **Transaction Feed** - Get wallet or list transaction feeds with advanced filtering (chain, token, value range, market cap)
- **Portfolio Management** - Retrieve public single-wallet portfolios; multi-wallet/token-filtered portfolios require a configured private API
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
client := cielogo.NewClient("YOUR_API_KEY")
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

## Upgrading

v0.3.0 corrects the API contract and changes several public types. Read [MIGRATION.md](MIGRATION.md) before upgrading from v0.2.0. Release changes are listed in [CHANGELOG.md](CHANGELOG.md).

`GetWalletPortfolioV2` is a private API operation. Configure it explicitly with `WithPortfolioBaseURL`; token-filtered results are available as `SingleToken`. The default public host does not expose this route in the current API router.

```go
client := cielogo.NewClient("YOUR_API_KEY",
    cielogo.WithPortfolioBaseURL("https://your-private-api.example/api"),
)
```

HTTP 202 pending responses return `*api.Error` with `StatusCode`, `Status`, `RetryAfter` and `RetryAfterDuration`. Callers control retries and credit spending.

## API Credits

Public API operations consume credits according to the server router. GetChainsV1 is free. Below are the public-router costs for supported operations:

| Endpoint | Cost (Credits) | Notes |
|----------|----------------|-------|
| **Feed & PnL** |
| GetFeedV1 | 3 or 5 depending on filter combination | **10/6 with IncludeMarketCap=true** ⚠️ |
| GetNftsPnlV1 | 5 | |
| GetTokensPnlV1 | 5 | |
| GetAggregatedTokenPnLV1 | **20** | Expensive operation |
| **Portfolio & Token Info** |
| GetWalletPortfolioV1 | **20** | |
| GetWalletPortfolioV2 | Private endpoint | Deployment-specific |
| GetTokenMetadataV1 | 2 | |
| GetTokenPriceV1 | 2 | |
| GetTokenStatsV1 | 10 | |
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
| GetAllWalletsListV1 | 5 | |
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
2. **Use selective feed filters** - Chain/type filters can reduce cost to 3 credits; a wallet alone does not always reduce the price
3. **Avoid IncludeMarketCap** unless necessary - Doubles feed costs (5→10 or 3→6)
4. **Cache expensive results** - Portfolio (20 credits) and Trading Stats (30 credits) change infrequently
5. **Use V2 partial updates** - `UpdateTrackedWalletV2` allows updating specific fields without full replacement

## Contributing

Contributions to the CieloGo project are welcome. Please feel free to report any bugs, suggest features, or open pull requests.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
