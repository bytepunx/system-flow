#!/usr/bin/env sh
# Install the development tools this repo needs into the main checkout's bin/ (git-ignored),
# which env.sh puts on PATH in every story worktree too. Does not touch ~/go/bin.
set -eu
. "$(dirname "$0")/env.sh"
GOLANGCI_LINT_VERSION="${GOLANGCI_LINT_VERSION:-v2.5.0}"
GORELEASER_VERSION="${GORELEASER_VERSION:-v2.18.1}"  # pinned; needs Go 1.27, which the go toolchain fetches automatically
BIN="$CACHE_ROOT/bin"
mkdir -p "$BIN"
echo "installing golangci-lint $GOLANGCI_LINT_VERSION and goreleaser $GORELEASER_VERSION into $BIN"
GOBIN="$BIN" go install "github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$GOLANGCI_LINT_VERSION"
GOBIN="$BIN" go install "github.com/goreleaser/goreleaser/v2@$GORELEASER_VERSION"
PNPM_VERSION="${PNPM_VERSION:-10.34.5}"
echo "installing pnpm $PNPM_VERSION into $CACHE_ROOT/.flai-cache/pnpm"
npm install --prefix "$CACHE_ROOT/.flai-cache/pnpm" --no-audit --no-fund --loglevel=error "pnpm@$PNPM_VERSION"
"$CACHE_ROOT/.flai-cache/pnpm/node_modules/.bin/pnpm" --version
"$BIN/golangci-lint" version
"$BIN/goreleaser" --version | head -1
