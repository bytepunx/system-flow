#!/usr/bin/env sh
# Lint plus all three test tiers in order of cost. Needs golangci-lint v2 (scripts/install-tools.sh).
set -eu
. "$(dirname "$0")/env.sh"
cd "$ROOT/flai"
echo "gofmt"; test -z "$(gofmt -l .)" || { gofmt -l .; echo "gofmt: files need formatting"; exit 1; }
echo "go vet"; go vet ./...
echo "golangci-lint"; golangci-lint run ./...
"$ROOT/scripts/test.sh"
"$ROOT/scripts/integration.sh"
"$ROOT/scripts/smoke.sh"
