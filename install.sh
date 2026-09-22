#!/usr/bin/env sh
# Install the latest flai release (or FLAI_VERSION) for this platform.
#
#   curl -fsSL https://raw.githubusercontent.com/bytepunx/system-flow/main/install.sh | sh
#
# Detects OS and architecture, resolves the newest flai/v* release, downloads
# the archive and checksums.txt, verifies the SHA-256, and installs flai into
# FLAI_INSTALL_DIR (default $HOME/.flai/bin, alongside flai's own config and
# cache under ~/.flai; ADR-0008). This is always a directory the user owns,
# so no sudo is ever used or needed. While the repository is private the
# GitHub API needs a token: GITHUB_TOKEN or GH_TOKEN, or a `gh auth login`
# session. `flai self-upgrade` does the same from an installed binary.
set -eu

REPO="${FLAI_REPO:-bytepunx/system-flow}"
BINARY="flai"
INSTALL_DIR="${FLAI_INSTALL_DIR:-$HOME/.flai/bin}"
API="${FLAI_API:-https://api.github.com}"

if [ -t 1 ]; then GREEN="\033[32m"; YELLOW="\033[33m"; RED="\033[31m"; RESET="\033[0m"; else GREEN=""; YELLOW=""; RED=""; RESET=""; fi
info()  { printf "  ${GREEN}+${RESET} %s\n" "$*"; }
warn()  { printf "  ${YELLOW}!${RESET} %s\n" "$*"; }
fatal() { printf "  ${RED}x${RESET} %s\n" "$*" >&2; exit 1; }
step()  { printf "  %s\n" "$*"; }

command -v curl >/dev/null 2>&1 || fatal "curl is required"
command -v tar >/dev/null 2>&1 || fatal "tar is required"

# Platform.
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
case "$OS" in
  linux|darwin) ;;
  *) fatal "unsupported OS: $OS (on Windows download the zip from https://github.com/$REPO/releases)" ;;
esac
case "$(uname -m)" in
  x86_64|amd64) ARCH="amd64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  *) fatal "unsupported architecture: $(uname -m)" ;;
esac
step "Platform: ${OS}/${ARCH}"

# Token: explicit, or borrowed from gh. Needed while the repository is private.
TOKEN="${GITHUB_TOKEN:-${GH_TOKEN:-}}"
if [ -z "$TOKEN" ] && command -v gh >/dev/null 2>&1; then
  TOKEN=$(gh auth token 2>/dev/null || true)
fi
api() {
  if [ -n "$TOKEN" ]; then
    curl -fsSL -H "Accept: $1" -H "Authorization: Bearer $TOKEN" "$2" -o "$3"
  else
    curl -fsSL -H "Accept: $1" "$2" -o "$3"
  fi
}

TMP=$(mktemp -d)
trap 'rm -rf "$TMP"' EXIT

# Release: the newest tag flai/v*, or the pinned one.
if [ -n "${FLAI_VERSION:-}" ]; then
  TAG="flai/v${FLAI_VERSION#v}"
  TAG_URL="$API/repos/$REPO/releases/tags/$(printf '%s' "$TAG" | sed 's|/|%2F|g')"
else
  step "Resolving latest release..."
  api "application/vnd.github+json" "$API/repos/$REPO/releases?per_page=50" "$TMP/releases.json" \
    || fatal "could not list releases of $REPO (private repository? set GITHUB_TOKEN or run gh auth login)"
  TAG=$(tr -d '\n' < "$TMP/releases.json" | grep -o '"tag_name": *"flai/v[^"]*"' | head -1 | sed -E 's/.*"(flai\/v[^"]+)".*/\1/')
  [ -n "$TAG" ] || fatal "no flai release found in $REPO"
  TAG_URL="$API/repos/$REPO/releases/tags/$(printf '%s' "$TAG" | sed 's|/|%2F|g')"
fi
VERSION="${TAG#flai/v}"
step "Release: flai $VERSION"
api "application/vnd.github+json" "$TAG_URL" "$TMP/release.json" || fatal "release $TAG not found"

# Asset API URLs, paired with names in document order. The asset API works
# for private repositories with a token and for public ones without.
ARCHIVE="${BINARY}_${VERSION}_${OS}_${ARCH}.tar.gz"
asset_url() {
  tr -d '\n' < "$TMP/release.json" | grep -oE '"(url|name)": *"[^"]*"' \
    | awk -v want="$1" '
      /"url": *".*\/releases\/assets\/[0-9]+"/ { sub(/^"url": *"/, ""); sub(/"$/, ""); url=$0; next }
      /"name":/ && url != "" { sub(/^"name": *"/, ""); sub(/"$/, ""); if ($0 == want) { print url; exit } url="" }'
}
ARCHIVE_URL=$(asset_url "$ARCHIVE")
SUMS_URL=$(asset_url "checksums.txt")
[ -n "$ARCHIVE_URL" ] || fatal "release $TAG has no asset $ARCHIVE"
[ -n "$SUMS_URL" ] || fatal "release $TAG has no checksums.txt"

step "Downloading $ARCHIVE..."
api "application/octet-stream" "$ARCHIVE_URL" "$TMP/$ARCHIVE" || fatal "download failed"
api "application/octet-stream" "$SUMS_URL" "$TMP/checksums.txt" || fatal "could not download checksums.txt"

# Verification is mandatory.
if command -v sha256sum >/dev/null 2>&1; then
  ACTUAL=$(sha256sum "$TMP/$ARCHIVE" | awk '{print $1}')
elif command -v shasum >/dev/null 2>&1; then
  ACTUAL=$(shasum -a 256 "$TMP/$ARCHIVE" | awk '{print $1}')
else
  fatal "no sha256sum or shasum available; cannot verify the download"
fi
EXPECTED=$(grep " $ARCHIVE\$" "$TMP/checksums.txt" | awk '{print $1}')
[ -n "$EXPECTED" ] || fatal "no checksum entry for $ARCHIVE"
[ "$ACTUAL" = "$EXPECTED" ] || fatal "checksum mismatch for $ARCHIVE"
info "Checksum verified"

tar -xzf "$TMP/$ARCHIVE" -C "$TMP" "$BINARY" || fatal "archive does not contain $BINARY"

DEST="$INSTALL_DIR/$BINARY"
mkdir -p "$INSTALL_DIR" 2>/dev/null || fatal "could not create $INSTALL_DIR (point FLAI_INSTALL_DIR at a directory you can write to)"
install -m 0755 "$TMP/$BINARY" "$DEST" \
  || fatal "could not install to $DEST (point FLAI_INSTALL_DIR at a directory you can write to)"
info "Installed $DEST ($("$DEST" version 2>/dev/null | head -1))"

PATH_LINE="export PATH=\"$INSTALL_DIR:\$PATH\""
FOUND=$(command -v "$BINARY" 2>/dev/null || true)
if [ "$FOUND" = "$DEST" ]; then
  info "flai $VERSION is ready; upgrade later with: flai self-upgrade"
elif [ -n "$FOUND" ]; then
  warn "'command -v flai' resolves to $FOUND, not $DEST; an earlier PATH entry takes precedence"
else
  warn "$INSTALL_DIR is not on your PATH yet"
fi
step "add $INSTALL_DIR to your PATH by adding this line to your shell profile:"
step "  $PATH_LINE"
