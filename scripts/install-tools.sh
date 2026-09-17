#!/usr/bin/env sh
# Install the development tools this repo needs into bin/ (git-ignored). Does not touch ~/go/bin.
set -eu
. "$(dirname "$0")/env.sh"
GOLANGCI_LINT_VERSION="${GOLANGCI_LINT_VERSION:-v2.5.0}"
GORELEASER_VERSION="${GORELEASER_VERSION:-v2.18.1}"  # pinned; needs Go 1.27, which the go toolchain fetches automatically
echo "installing golangci-lint $GOLANGCI_LINT_VERSION and goreleaser $GORELEASER_VERSION into $ROOT/bin"
GOBIN="$ROOT/bin" go install "github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$GOLANGCI_LINT_VERSION"
GOBIN="$ROOT/bin" go install "github.com/goreleaser/goreleaser/v2@$GORELEASER_VERSION"
PNPM_VERSION="${PNPM_VERSION:-10.34.5}"
echo "installing pnpm $PNPM_VERSION into $ROOT/.flai-cache/pnpm"
npm install --prefix "$ROOT/.flai-cache/pnpm" --no-audit --no-fund --loglevel=error "pnpm@$PNPM_VERSION"
"$ROOT/.flai-cache/pnpm/node_modules/.bin/pnpm" --version
"$ROOT/bin/golangci-lint" version
"$ROOT/bin/goreleaser" --version | head -1
