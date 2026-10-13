# Contract checks

Run the SDK against the supplied server's actual Gin handlers without a database, external API key or paid requests:

```sh
scripts/check-server-contract.sh /path/to/cielo-api-microservice
```

The script creates and removes a temporary module, retaining the server's Go dependency manifest and replacements. It points that module at both local repositories and invokes a race-enabled test of 12 HTTP scenarios. Downstream service interfaces are stubbed; request binding, validation, server response serialization and SDK decoding use the real code.

No files in the server repository or SDK go.mod/go.sum are changed. The server's required Go toolchain and dependencies must be available; dependencies may be downloaded. This local gate intentionally is not part of public CI, since CI does not have the server checkout. The latest validated server revision is recorded in the release audit.
