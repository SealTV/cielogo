#!/usr/bin/env bash
set -euo pipefail
if [[ $# != 1 ]]; then
  echo "Usage: $0 /path/to/cielo-api-microservice" >&2
  exit 2
fi
sdk_dir=$(cd "$(dirname "$0")/.." && pwd)
api_dir=$(cd "$1" && pwd)
fixture_dir=$(mktemp -d "${TMPDIR:-/tmp}/cielogo-server-contract.XXXXXX")
trap 'rm -rf "$fixture_dir"' EXIT
cp "$api_dir/go.mod" "$api_dir/go.sum" "$fixture_dir/"
cp "$sdk_dir/scripts/server-contract/contract_test.go.txt" "$fixture_dir/contract_test.go"
cd "$fixture_dir"
go mod edit -module=cielogo-server-contract
# Copying the server's dependency manifest also preserves its required replaces.
go mod edit -require=github.com/sealtv/cielogo@v0.0.0 -replace="github.com/sealtv/cielogo=$sdk_dir"
go mod edit -require=github.com/uniwhales/cielo-api-microservice@v0.0.0 -replace="github.com/uniwhales/cielo-api-microservice=$api_dir"
go test -mod=mod -race -count=1 -timeout=3m -v .
