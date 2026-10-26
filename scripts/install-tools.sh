#!/usr/bin/env sh
# Install the development tools this repo needs (golangci-lint, GoReleaser, ShellCheck) into the main checkout's bin/ (git-ignored),
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
# ShellCheck from its GitHub release. The release publishes no checksums, so the
# sha256 of each supported asset is pinned here; another version needs SHELLCHECK_SHA256.
SHELLCHECK_VERSION="${SHELLCHECK_VERSION:-v0.10.0}"
if "$BIN/shellcheck" --version 2>/dev/null | grep -qx "version: ${SHELLCHECK_VERSION#v}"; then
  echo "shellcheck $SHELLCHECK_VERSION already in $BIN"
else
  case "$(uname -s)" in
    Linux) sc_os=linux ;;
    Darwin) sc_os=darwin ;;
    *) echo "install-tools: no ShellCheck release for $(uname -s); install shellcheck $SHELLCHECK_VERSION into $BIN by hand" >&2; exit 1 ;;
  esac
  case "$(uname -m)" in
    x86_64 | amd64) sc_arch=x86_64 ;;
    aarch64 | arm64) sc_arch=aarch64 ;;
    *) echo "install-tools: no ShellCheck release for $(uname -m); install shellcheck $SHELLCHECK_VERSION into $BIN by hand" >&2; exit 1 ;;
  esac
  sc_asset="shellcheck-$SHELLCHECK_VERSION.$sc_os.$sc_arch.tar.xz"
  case "$sc_asset" in
    shellcheck-v0.10.0.linux.x86_64.tar.xz) sc_pinned=6c881ab0698e4e6ea235245f22832860544f17ba386442fe7e9d629f8cbedf87 ;;
    shellcheck-v0.10.0.linux.aarch64.tar.xz) sc_pinned=324a7e89de8fa2aed0d0c28f3dab59cf84c6d74264022c00c22af665ed1a09bb ;;
    shellcheck-v0.10.0.darwin.x86_64.tar.xz) sc_pinned=ef27684f23279d112d8ad84e0823642e43f838993bbb8c0963db9b58a90464c2 ;;
    shellcheck-v0.10.0.darwin.aarch64.tar.xz) sc_pinned=bbd2f14826328eee7679da7221f2bc3afb011f6a928b848c80c321f6046ddf81 ;;
    *) sc_pinned="" ;;
  esac
  sc_sha256="${SHELLCHECK_SHA256:-$sc_pinned}"
  if [ -z "$sc_sha256" ]; then
    echo "install-tools: no pinned sha256 for $sc_asset; set SHELLCHECK_SHA256 to its sha256" >&2
    exit 1
  fi
  echo "installing shellcheck $SHELLCHECK_VERSION into $BIN"
  sc_tmp="$(mktemp -d)"
  trap 'rm -rf "$sc_tmp"' EXIT
  curl -fsSL -o "$sc_tmp/$sc_asset" "https://github.com/koalaman/shellcheck/releases/download/$SHELLCHECK_VERSION/$sc_asset"
  if command -v sha256sum >/dev/null 2>&1; then
    sc_got="$(sha256sum "$sc_tmp/$sc_asset")"
  else
    sc_got="$(shasum -a 256 "$sc_tmp/$sc_asset")"
  fi
  sc_got="${sc_got%% *}"
  if [ "$sc_got" != "$sc_sha256" ]; then
    echo "install-tools: $sc_asset has sha256 $sc_got, want $sc_sha256; not installed" >&2
    exit 1
  fi
  tar -xJf "$sc_tmp/$sc_asset" -C "$sc_tmp"
  install -m 0755 "$sc_tmp/shellcheck-$SHELLCHECK_VERSION/shellcheck" "$BIN/shellcheck"
  rm -rf "$sc_tmp"
  trap - EXIT
fi
PNPM_VERSION="${PNPM_VERSION:-10.34.5}"
echo "installing pnpm $PNPM_VERSION into $CACHE_ROOT/.flai-cache/pnpm"
npm install --prefix "$CACHE_ROOT/.flai-cache/pnpm" --no-audit --no-fund --loglevel=error "pnpm@$PNPM_VERSION"
"$CACHE_ROOT/.flai-cache/pnpm/node_modules/.bin/pnpm" --version
"$BIN/golangci-lint" version
"$BIN/goreleaser" --version | head -1
"$BIN/shellcheck" --version | sed -n 's/^version: /shellcheck /p'
