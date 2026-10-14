# Release-candidate consumer validation

Date: 2026-10-14. SDK: published `github.com/sealtv/cielogo@v1.0.0-rc.1`, commit `dd52688`. Consumer: `github.com/SealTV/trading-bot`, checkout `0ee9232`, including its existing uncommitted entry-point change copied into an isolated temporary directory. The source consumer checkout was not modified.

## Migration found by the trial

Changing only the SDK dependency reproduced a failure in `TestPublicInventorySDKWirePagination`: the consumer's private HTTP transport rejected the corrected tracked-wallet URL before sending a request. The v0.2.0 shim accepted only `/api/v1/tracked-wallets`, appended `/`, and injected the missing cursor itself. The SDK now sends both the trailing slash and cursor directly.

In the temporary consumer copy:

1. Set the SDK dependency to `v1.0.0-rc.1` and download its checksum.
2. Accept `/api/v1/tracked-wallets/` in the transport's fixed-origin GET guard.
3. Replace the URL rewrite with validation that `next_object` matches the expected cursor from the request context.
4. Update the consumer's evidence version string to `cielogo@v1.0.0-rc.1`.

The consumer's response validation, credential redaction, transport ownership, no-retry behavior, pagination bounds and credit caps remain in place. These are consumer changes; no additional SDK contract change was needed.

## Checks

Executed in the isolated consumer with `GOTOOLCHAIN=go1.27.2`:

```sh
go test ./...
go test -race ./...
go vet ./...
go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest run
```

All checks passed. Public-inventory tests exercise three-page URL/cursor forwarding, complete/empty inventories, unknown wallet types, malformed envelopes and pagination, duplicates, HTTP failures without retries, redaction, cancellation, deadlines, body/page/row limits and TLS/transport ownership. Existing application tests also pass.

This is a runtime consumer trial using deterministic fake HTTP responses, not a production deployment. No external API calls, real trades or payments were made by the consumer tests. Database integration tests were not run. The bot's feed uses a separate client, so this trial does not validate SDK WebSocket events in that application; SDK fixtures and local WebSocket lifecycle/race tests provide that coverage.
