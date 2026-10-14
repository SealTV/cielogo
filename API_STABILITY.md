# CieloGo 1.x contract

The 1.x line covers the HTTP operations and WebSocket protocol exposed by this SDK. It does not promise wrappers for every future server endpoint. New operations, response fields, chain/tag identifiers and future transaction variants may be added without changing existing callers' Go source.

## Compatibility

Existing exported types, methods, field types and documented behavior are preserved across 1.x minor releases. Breaking Go API changes require a new major version. Server-side changes, permissions, billing and downtime remain controlled by Cielo. Fixture/handler checks verify the source contract, not production deployment equivalence.

String enums are open sets: always handle unrecognized values. New transaction types are preserved as `UnknownTransactionEvent`; unknown WebSocket events carry `json.RawMessage`. Applications should include a default case in type switches. Verification and migration remain raw JSON to preserve their evolving shapes.

Monetary values follow the API's float64 wire representation; they are not arbitrary-precision accounting values. `BalanceRaw` preserves raw portfolio balances as strings. `TotalSupplyString` is the server's display-formatted supply (it may contain suffixes), not an exact integer representation.

## Requests and clients

`NewClient` accepts an API key; the server validates it. Nil options and `WithHTTPClient(nil)` are ignored. A nil request, nil context or uninitialized client returns a recognizable error instead of panicking. Client options must be set before sharing the client. The client supports concurrent HTTP requests. Caller-supplied HTTP clients, transports, dialers and referenced TLS configurations must not be mutated concurrently.

The default HTTP timeout is 90 seconds. A custom HTTP client retains its own timeout. Callers should set operation-specific context deadlines. Responses default to a 16 MiB maximum; use `WithMaxResponseBytes` to select another positive limit. `CloseIdleConnections` releases pooled connections without canceling requests.

Path parameters are escaped. Redirects to a different origin are rejected with `ErrUnsafeRedirect` so the API key is not forwarded elsewhere; same-origin redirects honor the caller's redirect policy. The SDK does not mutate the caller-owned HTTP client.

The SDK does not retry API errors or pending responses. The underlying HTTP transport may retry connection failures, including standard-library retries on reused GET connections; callers with strict per-attempt credit limits must configure their transport accordingly. The budgeted integration smoke uses fresh HTTP/1 connections to prevent those transport retries. A complete API error is available through `errors.As` as `*api.Error`, including numeric HTTP status, string status, message and Retry-After. Malformed/truncated success envelopes and missing object data return `ErrInvalidResponse`. Null collection data represents an empty list. HTTP 202 remains pending, even if the transport considers it successful.

Exported sentinel errors support `errors.Is`. Error strings are diagnostic text, not a stable branching API. Response models are intended for use with the provided API wrappers; custom decoders reject missing transaction/event types. A failed decode does not replace an existing event model. `TxEvent` marshals its payload in the API's flat shape, preserving typed transaction data and unknown raw payloads.

## WebSocket lifecycle

`NewWebsocketConnection` uses the configured public base URL, replacing HTTP(S) with WS(S), and respects its context. `WithWebSocketDialer` configures handshake/TLS/proxy behavior; HTTP-client settings do not implicitly configure the WebSocket dialer. Handshake failures preserve structured HTTP error information when a response is available.

One `RunListener` runs per connection. A concurrent listener returns `ErrListenerRunning`. `SendCommand` supports concurrent calls; writes have a one-second deadline. Incoming messages are limited to 512 KiB. Canceling the listener closes the connection and returns the context error, including when output is blocked. Read/decode errors terminate the listener. Normal close/going-away returns nil. An explicit abrupt close may return the underlying read error.

The caller owns the output channel: keep it open while the listener runs. The SDK does not close it. Handlers are configured during connection creation. Nil ping/pong/close handlers restore gorilla's defaults. Custom handlers execute in the reader and must preserve protocol responses as needed. Connection and event models do not provide automatic reconnect or subscriptions replay.

## Public/private boundary

The default base URL is the public API. Portfolio V2 requires an explicit `WithPortfolioBaseURL` for the private endpoint. This operation is supported for deployments with authorized access to that route; it is not a promise that the public host exposes it. Other private, trading and new analytics routes are outside the frozen initial surface.

## Verification and release gate

- Race tests and every current HTTP verb/path/body contract.
- Server-shaped fixtures and local checks against actual handlers, using fake downstream services.
- Fuzz smoke for transaction, WebSocket, portfolio, generic envelopes and Retry-After.
- Minimum Go 1.22/1.23 compatibility and patched Go 1.26.9/1.27.2 runs. Production should use a maintained, patched Go toolchain; old compatibility jobs are not a security recommendation.
- Vulnerability scanning on a patched toolchain, plus vet, build and lint.
- Integration tests compile without secrets on every change.
- Before the final v1.0.0 tag, run the budgeted read-only production smoke (upper bound 72 credits, explicit ceiling 72-100), review any pending/not-found results, and try the candidate in a consumer. Writes and private production endpoints remain untested unless separately configured and authorized.

A release candidate can pass all local gates while these final deployment/consumer checks remain open. Such a candidate must be labeled prerelease.

## Known upstream limits

The supplied wallet-tag handler rejects the 50th wallet despite the advertised maximum of 50. Use at most 49 until the server is fixed. The SDK surfaces the error and does not silently split paid requests. See [RELEASE_AUDIT.md](RELEASE_AUDIT.md) for evidence and remaining release gates.
