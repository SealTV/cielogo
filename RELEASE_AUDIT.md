# 1.0 release audit

Date: 2026-10-13. Baseline: v0.3.0 (dc4781b). Supplied server source: 1b99824e. Scope: existing SDK operations, transport/error behavior, serialization, WebSocket lifecycle and release gates. No production API calls were made.

## Findings fixed in the candidate

1. Integration-tag tests did not compile after the trading-stats model change. CI now compiles them without secrets.
2. Integration tests accepted arbitrary trading errors and expected an absent balance decimals field. Unexpected errors now fail; only documented typed pending responses are accepted.
3. Missing/malformed success envelopes, null object data and concatenated JSON could hide failures. They return ErrInvalidResponse; API errors inside HTTP 200 are surfaced.
4. Response bodies were unbounded and the default HTTP client had no timeout. Defaults are 16 MiB / 90 seconds, with overrides.
5. Nil requests/options/contexts and uninitialized clients could panic. They are guarded or ignored as documented.
6. Wallet path parameters were interpolated without escaping. Query/fragment injection is prevented.
7. Redirects could forward the API-key header to another origin. Cross-origin redirects are rejected without mutating caller-owned clients.
8. TxEvent marshaling discarded its transaction payload. Flat serialization now preserves typed and unknown data.
9. Failed WebSocket decode could leave mismatched old Data and new Type. Decoding is atomic.
10. Concurrent listeners and nil commands/handlers/options were unsafe. Lifecycle guards and nil behavior are defined.
11. WebSocket TLS/proxy configuration and structured handshake errors were unavailable. A configurable dialer and status/Retry-After errors are supported.
12. Remaining Bridge spelling mistakes were unsuitable to freeze. Canonical names and migration instructions are provided; old type/constant aliases remain where possible.
13. CI lacked tagged compilation/fuzz/security gates; live tests masked failures and ran scheduled paid requests. CI now covers these gaps, and live smoke is manual and budgeted.

## Evidence

| Check | Result |
|---|---|
| Current HTTP wrappers | Literal verb/path/body matrix passes for 28 non-deprecated operations |
| Actual server handlers | 12 request/response scenarios pass with real binding, validation and serialization, using fake downstream services |
| Tag-batch boundary | 49-wallet success and 50-wallet HTTP 400 reproduced separately |
| Race repetitions | 20 runs of SDK packages pass |
| Local fuzz campaigns | 4,908,354 executions across five targets; no crash/invariant failure |
| Coverage | Root 85.7%; api 76.9%; apiv1 62.3%; SDK aggregate excluding examples/helpers 76.0% |
| Vet/build/lint | Pass for published packages |
| Vulnerability scan | No called vulnerabilities reported on Go 1.27.2; local Go 1.27.1 requires standard-library security fixes |
| Integration tag | Compiles without API calls |
| Production smoke | Not run; needs an explicit budget and configured key |
| Consumer runtime trial | Not performed |

The fuzz count describes a bounded campaign, not proof of absence of defects. Security results depend on the database and compiler. Local handlers do not exercise databases/business logic, authentication middleware or production ingress. Live mutations/private production endpoints were not tested.

Public CI must pass on the actual candidate commit before tagging. It checks Go 1.22/1.23 compatibility and patched 1.26.9/1.27.2, and enforces 70% SDK aggregate coverage.

## Known upstream limit

The supplied wallet-tag handler checks len(wallets) >= 50 while appending and writes an error without returning. The advertised 50-wallet boundary therefore produces HTTP 400 and may append a success body. The SDK surfaces the HTTP error. Use at most 49 wallets until the backend is fixed; the SDK does not silently split paid calls.

The server's total_supply_string is display-formatted, not an exact integer. The SDK preserves it as a string; use documented raw values within their supported range for precision-sensitive code.

## Final stable prerequisites

1. Merge only after every candidate CI job passes.
2. Publish v1.0.0-rc.1 with stability and migration documentation.
3. Explicitly run TestIntegrationReleaseSmoke with a 72-100 credit ceiling. At most 72 credits are reserved; no retries. Review pending/not-found outcomes before claiming live data paths are validated.
4. Try the candidate in a real consumer, including pagination, events and error branches.
5. If no required contract changes remain, tag the validated commit as v1.0.0. Retain the production/private/mutation limitations in release notes.
