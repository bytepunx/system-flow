#!/usr/bin/env sh
# Install the latest flai release (or FLAI_VERSION) for this platform.
#
#   curl -fsSL https://raw.githubusercontent.com/bytepunx/system-flow/main/install.sh | sh
#
# Detects OS and architecture, resolves the newest flai/v* release, downloads
# the archive, checksums.txt, and checksums.txt.sig, verifies the signature
# over checksums.txt with openssl and the release key below (ADR-0070), then
# the archive's SHA-256 against checksums.txt, and installs flai into
# FLAI_INSTALL_DIR (default $HOME/.flai/bin, alongside flai's own config and
# cache under ~/.flai; ADR-0008). A release it cannot verify, such as one from
# before releases were signed, is not installed. The install directory is
# always one the user owns, so no sudo is ever used or needed. While the
# repository is private the GitHub API needs a token: GITHUB_TOKEN or
# GH_TOKEN, or a `gh auth login` session. `flai self-upgrade` does the same
# from an installed binary.
set -eu

REPO="${FLAI_REPO:-bytepunx/system-flow}"
BINARY="flai"
INSTALL_DIR="${FLAI_INSTALL_DIR:-$HOME/.flai/bin}"
API="${FLAI_API:-https://api.github.com}"

# The release key's public half: flai/internal/buildinfo/cosign.pub, which
# flai embeds and every release publishes as cosign.pub. A test in
# flai/internal/selfupgrade fails when the two differ.
RELEASE_KEY='-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEXGMlBGpRPBSMMOcQrk81afzRbx/M
izI2Hxp0X5OBfTx1U0q1kqKv2oMeL9d0FRsHo7kNz/pNi/gEs82LtNPbJA==
-----END PUBLIC KEY-----'

if [ -t 1 ]; then GREEN="\033[32m"; YELLOW="\033[33m"; RED="\033[31m"; RESET="\033[0m"; else GREEN=""; YELLOW=""; RED=""; RESET=""; fi
info()  { printf "  ${GREEN}+${RESET} %s\n" "$*"; }
warn()  { printf "  ${YELLOW}!${RESET} %s\n" "$*"; }
fatal() { printf "  ${RED}x${RESET} %s\n" "$*" >&2; exit 1; }
step()  { printf "  %s\n" "$*"; }

command -v curl >/dev/null 2>&1 || fatal "curl is required"
command -v tar >/dev/null 2>&1 || fatal "tar is required"
command -v openssl >/dev/null 2>&1 \
  || fatal "openssl is required to verify the release's signature; install it with your package manager (apt install openssl, dnf install openssl, apk add openssl, or brew install openssl) and run this again"

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
# api fetches $2 into $3 accepting $1. A call the network drops (GitHub cuts
# off a large response now and then, I-0086) is tried three times in all, two
# seconds apart; an HTTP error, curl's exit 22 under --fail, fails at once.
api() {
  API_TRY=1
  while :; do
    API_RC=0
    if [ -n "$TOKEN" ]; then
      curl -fsSL -H "Accept: $1" -H "Authorization: Bearer $TOKEN" "$2" -o "$3" || API_RC=$?
    else
      curl -fsSL -H "Accept: $1" "$2" -o "$3" || API_RC=$?
    fi
    if [ "$API_RC" -eq 0 ]; then return 0; fi
    if [ "$API_RC" -eq 22 ] || [ "$API_TRY" -ge 3 ]; then return "$API_RC"; fi
    warn "fetching $2 failed (curl exit $API_RC); trying again in 2s (attempt $((API_TRY + 1)) of 3)" >&2
    sleep 2
    API_TRY=$((API_TRY + 1))
  done
}

TMP=$(mktemp -d)
trap 'rm -rf "$TMP"' EXIT

# Release: the newest tag flai/v*, or the pinned one.
if [ -n "${FLAI_VERSION:-}" ]; then
  TAG="flai/v${FLAI_VERSION#v}"
  TAG_URL="$API/repos/$REPO/releases/tags/$(printf '%s' "$TAG" | sed 's|/|%2F|g')"
else
  # Small pages, since a large listing is what the network drops (I-0086):
  # the next is asked for only while no flai/v* tag has turned up, up to the
  # 50 newest releases, and not past a page that holds no release.
  step "Resolving latest release..."
  TAG=""
  PAGE=1
  while [ -z "$TAG" ] && [ "$PAGE" -le 5 ]; do
    api "application/vnd.github+json" "$API/repos/$REPO/releases?per_page=10&page=$PAGE" "$TMP/releases.json" \
      || fatal "could not list releases of $REPO (private repository? set GITHUB_TOKEN or run gh auth login)"
    grep -q '"tag_name"' "$TMP/releases.json" || break
    TAG=$(tr -d '\n' < "$TMP/releases.json" | grep -o '"tag_name": *"flai/v[^"]*"' | head -1 | sed -E 's/.*"(flai\/v[^"]+)".*/\1/')
    PAGE=$((PAGE + 1))
  done
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
SIG_URL=$(asset_url "checksums.txt.sig")
[ -n "$ARCHIVE_URL" ] || fatal "release $TAG has no asset $ARCHIVE"
[ -n "$SUMS_URL" ] || fatal "release $TAG has no checksums.txt"
[ -n "$SIG_URL" ] || fatal "release $TAG has no checksums.txt.sig, so it cannot be verified and is not installed; releases from before signing are unsigned: install a newer one (leave FLAI_VERSION unset for the newest)"

step "Downloading $ARCHIVE..."
api "application/octet-stream" "$ARCHIVE_URL" "$TMP/$ARCHIVE" || fatal "download failed"
api "application/octet-stream" "$SUMS_URL" "$TMP/checksums.txt" || fatal "could not download checksums.txt"
api "application/octet-stream" "$SIG_URL" "$TMP/checksums.txt.sig" || fatal "could not download checksums.txt.sig"

# Verification is mandatory: first that the release key signed checksums.txt,
# then that the archive is the one checksums.txt names. The signature is
# base64 of an ASN.1 DER ECDSA P-256 signature over the file's SHA-256, as
# cosign sign-blob writes it; openssl base64 decodes it the same everywhere,
# where base64's own flags differ between GNU and BSD.
printf '%s\n' "$RELEASE_KEY" >"$TMP/cosign.pub"
openssl base64 -d -A -in "$TMP/checksums.txt.sig" -out "$TMP/checksums.txt.sig.der" 2>/dev/null \
  || fatal "the signature of release $TAG is not base64; the release may have been altered, so it is not installed"
openssl dgst -sha256 -verify "$TMP/cosign.pub" -signature "$TMP/checksums.txt.sig.der" "$TMP/checksums.txt" >/dev/null 2>&1 \
  || fatal "the signature of release $TAG does not verify with the release key this install.sh carries, so it is not installed; the release may have been altered, or signed by a newer key: get the newest install.sh and run it again"
info "Signature verified"

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
