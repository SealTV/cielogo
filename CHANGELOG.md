# Changelog

## v1.0.0 — 2026-10-14

Promote the audited SDK surface to the first stable 1.x contract. The SDK implementation is unchanged from v1.0.0-rc.1; this release finalizes the stability, migration and validation documentation.

- Automated checks pass on four Go versions, including race, build, lint, vulnerability and fuzz gates.
- An isolated trading-bot trial passes tests, race checks, vet and lint after removing its obsolete tracked-wallet URL shim; document that migration.
- Bound paid integration request attempts with fresh HTTP/1 connections and report pending/not-found outcomes and reserved credits explicitly.
- Production smoke stopped at invalid-key authentication before any paid request. Successful live data paths, production writes and private endpoints remain unverified; see [RELEASE_AUDIT.md](RELEASE_AUDIT.md).

## v1.0.0-rc.1

Stabilize the existing SDK surface before the first stable release. See [API_STABILITY.md](API_STABILITY.md) and [MIGRATION.md](MIGRATION.md).

- Strict, bounded HTTP decoding, default timeout, recognizable nil/configuration errors and cross-origin redirect protection.
- Complete HTTP method matrix, malformed-response and serialization regression tests.
- Preserve flattened transaction payloads on marshal; atomic WebSocket decode; concurrent-listener guard and configurable TLS dialer.
- Correct remaining Bridge field spellings and expose canonical aliases.
- Fix integration-tag compilation and unexpected-error masking; add a manual production smoke with an explicit 72-100 credit ceiling.
- Add actual-server-handler checks, five fuzz targets, multi-Go CI and a vulnerability gate.
- Candidate only: production smoke and consumer validation are required before the final stable tag.

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
