#!/usr/bin/env sh
# Install the development tools this repo needs into bin/ (git-ignored). Does not touch ~/go/bin.
set -eu
. "$(dirname "$0")/env.sh"
GOLANGCI_LINT_VERSION="${GOLANGCI_LINT_VERSION:-v2.5.0}"
GORELEASER_VERSION="${GORELEASER_VERSION:-v2.18.1}"  # pinned; needs Go 1.27, which the go toolchain fetches automatically
echo "installing golangci-lint $GOLANGCI_LINT_VERSION and goreleaser $GORELEASER_VERSION into $ROOT/bin"
GOBIN="$ROOT/bin" go install "github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$GOLANGCI_LINT_VERSION"
GOBIN="$ROOT/bin" go install "github.com/goreleaser/goreleaser/v2@$GORELEASER_VERSION"
"$ROOT/bin/golangci-lint" version
"$ROOT/bin/goreleaser" --version | head -1
