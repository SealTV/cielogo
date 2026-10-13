# Changelog

## v0.3.0

Correct existing API requests and response decoding against the server contract. This release has breaking Go type changes; see [MIGRATION.md](MIGRATION.md).

- Repair tracked-wallet pagination, feed transaction filtering/cursors, token PnL chain filtering and followed-list filtering.
- Align tracked-wallet V2 notification settings and list wallet IDs with the server.
- Correct trading stats, swap, NFT PnL, portfolio totals, wallet-by-tag and toggle-follow responses.
- Add new perpetual/Hyperliquid event payloads and preserve future unknown events, market caps, verification and migration.
- Extend portfolio, metadata, statistics and PnL models; support unsigned metadata values and dynamic chain discovery.
- Require explicit private configuration for portfolio V2 and expose token-filtered responses separately.
- Preserve HTTP status and Retry-After, without automatic paid retries.
- Correct WebSocket acknowledgements; respect cancellation and serialize concurrent writes.
- Add server-shaped fixtures, URL/body/pagination tests and WebSocket lifecycle/race tests.
- Correct examples, API credit documentation and migration instructions.

New analytics/trading endpoints outside the previously implemented SDK surface remain future work. This release fixes compatibility of existing methods; it does not claim complete coverage of every public or private endpoint.
