# Migrating to the 1.0 release candidate

The candidate retains the corrected v0.3.0 request/response shapes and adds stricter failure handling. Read [API_STABILITY.md](API_STABILITY.md) for the planned 1.x compatibility contract.

Changes from v0.3.0:

- `BridgeEvent.Amoun` becomes `Amount`; `Platfrom` becomes `Platform`. These are the remaining pre-1.0 field corrections.
- Use canonical `ToggleFollowWalletsListResponse`, `TxTypeNftMint` and `LpTypeAdd`. Old type/constant spellings remain aliases for compatibility.
- Malformed/concatenated JSON, missing success status and missing object data now produce `ErrInvalidResponse`, rather than returning empty data. Null lists still decode as empty collections. An application error with HTTP 200 is surfaced as `*api.Error`.
- Default HTTP timeout is 90 seconds and maximum response size is 16 MiB. Use `WithHTTPClient` / `WithMaxResponseBytes` to change them.
- Cross-origin redirects return `ErrUnsafeRedirect`. A custom redirect policy is still honored within the same origin.
- Nil requests and contexts return sentinel errors. Nil client options and `WithHTTPClient(nil)` are ignored.
- `TxEvent` now preserves its flattened payload when marshaled; events without a type are rejected. Reused WebSocket models remain unchanged after a failed decode.
- A second concurrent listener returns `ErrListenerRunning`. Use `WithWebSocketDialer` for TLS/proxy/handshake settings; WebSocket handshake errors retain status and Retry-After.
- `CloseIdleConnections` releases pooled HTTP connections.
- Live integration tests now require a manual budgeted run. `CIELO_LIVE_CREDIT_BUDGET` must be 72-100 for `TestIntegrationReleaseSmoke`; no automatic paid retry occurs.

# Migrating from v0.2.0 to v0.3.0

v0.3.0 corrects the existing SDK contract against the API implementation. It changes some public Go types and decoded values. Update consumers before upgrading.

## Requests

- `AddWalletsListRequest.Wallets` and `UpdateWalletsListRequest.Wallets` are `[]int64` tracked-wallet IDs, not addresses. Track wallets first and pass their IDs.
- `UpdateTrackedWalletV2Request` uses `MinAmountUSD`, `Filters []int`, `Chains []int`, `TelegramBotID *int`, `DiscordChannelID *string` and optional `BundleID *int`. Remove `MinUSD`, `TxTypes`, `TelegramBot` and `DiscordChannel`. Use nil slices to leave filters/chains unchanged and non-nil empty slices to clear them. Notification filters and chain IDs come from the wallet-settings API; do not assume they are transaction names or blockchain chain IDs.
- The SDK now sends `txTypes` for HTTP feed, `chain` for token PnL, and `following_only` for the list directory. The existing request field names remain the same.
- `GetTrackedWalletsV1` sends `ListID` and `NextObject` correctly. Convert the numeric response cursor with `strconv.Itoa(resp.Pagination.NextObject)` when building the next request. Stop when `HasNextPage` is false.
- Related-wallet requests support `Timeframe`, `FromTimestamp`, `ToTimestamp` and `NextObject`.
- Token-stat requests support `IncludePools` and `IncludeLiquidity`.
- Token PnL supports `max` for all history, not `all`. Trading stats default to `30d` on the server.

## Responses

- Feed has a dedicated `FeedPagination`. Use `resp.Paging.NextCursor` as the next request's `StartFrom`. `NextObject` no longer exists on feed pagination; other endpoints retain their existing pagination models.
- Trading statistics now expose `PnL`, `WinRate`, `SwapsCount`, `BuyCount`, `SellCount`, buy/sell amounts, ROI and holding distributions. `WinRate` is a fraction (0..1). Remove uses of `TotalTrades`, `SuccessfulTrades`, `FailedTrades`, `AverageProfit`, `TotalProfit`, `TotalVolume`, `MostTradedToken`, `MostProfitableToken`, `TradingFrequency` and `LastTradeTimestamp`: the server does not return those fields. Do not interpret buy/sell counts as successful/failed trades.
- `SwapEvent` uses `Token0*` and `Token1*`, `Dex`, `IsSell`, `FirstInteraction` and platform metadata. The old single-token/bridge-shaped fields are removed.
- Feed events include market-cap enrichment, verification and migration; verification/migration are `json.RawMessage` so their full server payload is preserved. New perpetual and Hyperliquid events are typed. Future unknown event types use `*UnknownTransactionEvent`, preserving their raw JSON rather than failing the page.
- `TxEvent.Block` and event payload blocks are `uint`; unsigned metadata supply/block values and price block numbers are `uint64`.
- Portfolio `TotalUSDValue` now reads `total_usd` correctly; each asset still reads its own `total_usd_value`. `PendingChains` exposes partial results. Chain distribution, raw balances, wallet addresses and native values are available.
- NFT PnL keeps existing Go field names but reads the actual `token_name`, `token_symbol`, `average_buy_price`, `average_sell_price` and `total_pnl_usd` fields.
- Wallet-by-tag results now read `wallets`. Toggle-follow reads the nested list's `followed` flag correctly.
- Token metadata/stats, token PnL and aggregated PnL include additional server fields. Token balance `Decimals` is deprecated: the server does not provide it.

## Private portfolio V2

The public API router does not register `/api/v2/portfolio`. The SDK requires an explicit private API URL instead of sending this request to the public host:

```go
client := cielogo.NewClient(apiKey,
    cielogo.WithPortfolioBaseURL("https://your-private-api.example/api"),
)
```

Keep `/api` in the configured URL. Do not use the public host unless your deployment explicitly routes this private endpoint there. Other calls continue to use the public `baseURL`. `WithBaseURL` does not configure private portfolio V2.

`GetWalletPortfolioV2` returns either a full portfolio or, when a token filter is used, `resp.SingleToken`. Check `SingleToken` before accessing `Portfolio`. Integration tests and the example accept `CIELO_PORTFOLIO_BASE_URL`.

## HTTP pending and retries

Use `errors.As(err, &apiErr)` with `var apiErr *api.Error`. Errors now expose `StatusCode`, the raw `RetryAfter` header and parsed `RetryAfterDuration`. A 202 response retains `Status == "pending"`. No automatic retry occurs: repeated calls can consume credits. Context cancellation and transport errors remain wrapped for `errors.Is`.

## WebSocket

- The WebSocket URL derives from the public `WithBaseURL` setting; HTTP(S) changes to WS(S) automatically. Dialing respects the context.
- `WalletSubscribedEventType` data is now `WalletSubscriptionAck`; `FeedSubscribedEventType` data is `FeedSubscriptionAck`. Their filters contain maps, whereas outgoing `Filter` continues to use arrays.
- `FeedUnsubscribedEventType` and `PongEventType` carry the string `"ok"`. Send `PingCmd` for a JSON pong; protocol-level pings remain supported.
- Unknown future event types preserve `json.RawMessage` in `Data`.
- `FeedSubscribeCmd` supports `MultiBundleID`; do not set both `ListID` and `MultiBundleID`.
- Only one `RunListener` may run per connection. Canceling its context closes the connection and returns `ctx.Err()`. Normal close/going-away returns nil. Read/decode errors return immediately without sending a blank event. `SendCommand` is safe for concurrent callers.
