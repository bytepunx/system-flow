#!/usr/bin/env sh
# Lint plus all three test tiers in order of cost. Needs golangci-lint v2 (scripts/install-tools.sh).
set -eu
. "$(dirname "$0")/env.sh"
cd "$ROOT/flai"
echo "gofmt"; test -z "$(gofmt -l .)" || { gofmt -l .; echo "gofmt: files need formatting"; exit 1; }
echo "go vet"; go vet ./...
lint="$(command -v golangci-lint || true)"
[ -n "$lint" ] || { echo "golangci-lint: not found; run scripts/install-tools.sh (make install-tools) in the main checkout"; exit 1; }
version="$("$lint" version 2>&1 | sed -n 's/.*has version v\{0,1\}\([0-9][0-9.]*\).*/\1/p' | head -1)"
echo "golangci-lint $version ($lint)"
case "$version" in
  2.*) ;;
  *) echo "golangci-lint: $lint is version ${version:-unknown}; flai/.golangci.yaml needs v2: run scripts/install-tools.sh (make install-tools) in the main checkout"; exit 1 ;;
esac
"$lint" run ./...
"$ROOT/scripts/test.sh"
"$ROOT/scripts/integration.sh"
"$ROOT/scripts/smoke.sh"
