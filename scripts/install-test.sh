#!/usr/bin/env sh
# Install flai with install.sh and upgrade it with `flai self-upgrade` from a
# release built from this tree and served on 127.0.0.1, under the stand-in
# repository local/flai (S-0340). It builds two releases with
# scripts/flai-snapshot.sh --local, an older and a newer, serves them with
# scripts/release-server.sh, and points install.sh (FLAI_API, FLAI_REPO) and
# flai self-upgrade (FLAI_RELEASES_API, --repo) at that server. Nothing reaches
# GitHub and no token is used. The check that the latest published release
# installs from GitHub is scripts/install-published-test.sh.
#
# When a step fails it prints the step, its output, and the server's request
# log, and exits non-zero, leaving its scratch folders under .flai-cache.
set -eu
. "$(dirname "$0")/env.sh"
# Versions above any project's flai.minimum, so that the installed binaries
# run in this project, and above any published release, so that neither is
# mistaken for one. OLD is installed and upgraded to NEW; NEW is the latest.
OLD=999.0.0
NEW=999.0.1
REPO=local/flai
DIR="$ROOT/.flai-cache/install-test"
SCRATCH_HOME="$ROOT/.flai-cache/install-test-home"
DIST="$ROOT/.flai-cache/install-test-dist"
GH_EMPTY="$ROOT/.flai-cache/install-test-gh"
LOG="$ROOT/.flai-cache/install-test.log"
SERVER_LOG="$ROOT/.flai-cache/install-test-server.log"
OUTSIDE=""
trap 'if [ -n "$OUTSIDE" ]; then rm -rf "$OUTSIDE"; fi' EXIT
rm -rf "$DIR" "$SCRATCH_HOME" "$DIST" "$GH_EMPTY" "$LOG" "$LOG.install" "$SERVER_LOG"
# A story worktree has no .flai-cache of its own until something makes it.
mkdir -p "$ROOT/.flai-cache" "$GH_EMPTY"

# No token: install.sh and flai read GITHUB_TOKEN and GH_TOKEN, and borrow gh's
# session, which an empty gh configuration folder leaves them without.
unset GITHUB_TOKEN GH_TOKEN
GH_CONFIG_DIR="$GH_EMPTY"
export GH_CONFIG_DIR
# The server is plain HTTP on the loopback: no proxy is to carry its requests.
no_proxy="127.0.0.1,localhost${no_proxy:+,$no_proxy}"
NO_PROXY="$no_proxy"
export no_proxy NO_PROXY

say() { echo "install-test: $*"; }
warn() { echo "install-test: $*" >&2; }

# fail STEP STATUS prints the failing step, its output, and the server's
# request log, and exits non-zero.
fail() {
  warn "FAILED: $1 (exit status $2); its output:"
  cat "$LOG" >&2
  if [ -s "$SERVER_LOG" ]; then
    warn "the release server's log ($SERVER_LOG):"
    cat "$SERVER_LOG" >&2
  fi
  exit 1
}

# step NAME CMD [ARG...] runs the command with its output in LOG, prints the
# output when it passes, and fails with NAME when it does not. A step's
# function runs where set -e does not apply, so each returns its own failure.
step() {
  name=$1
  shift
  say "$name"
  status=0
  "$@" >"$LOG" 2>&1 || status=$?
  if [ "$status" -ne 0 ]; then
    fail "$name" "$status"
  fi
  cat "$LOG"
}

# quiet_step NAME CMD [ARG...] is step, printing the output only on failure.
quiet_step() {
  name=$1
  shift
  say "$name"
  status=0
  "$@" >"$LOG" 2>&1 || status=$?
  if [ "$status" -ne 0 ]; then
    fail "$name" "$status"
  fi
}

# expect_version BINARY VERSION checks the binary's first line of flai version.
expect_version() {
  first=$("$1" version | head -n 1)
  echo "$1: $first"
  [ "$first" = "flai $2" ] || { echo "$1 reports \"$first\", not \"flai $2\""; return 1; }
}

no_sudo() {
  if grep -vE '^\s*#' "$ROOT/install.sh" | grep -qE '(^|[^A-Za-z0-9_-])sudo([^A-Za-z0-9_-]|$)'; then
    echo "install.sh still invokes sudo outside a comment"
    return 1
  fi
}

build_releases() {
  "$ROOT/scripts/flai-build.sh" || return 1
  for v in "$OLD" "$NEW"; do
    "$ROOT/scripts/flai-snapshot.sh" --local "$v" --dist "$DIST/$v" || return 1
  done
}

install_explicit_dir() {
  FLAI_API="$URL" FLAI_REPO="$REPO" FLAI_VERSION="$OLD" FLAI_INSTALL_DIR="$DIR" \
    sh "$ROOT/install.sh" || return 1
  expect_version "$DIR/flai" "$OLD"
}

install_default_dir() {
  install_status=0
  env -i HOME="$SCRATCH_HOME" PATH="$PATH" FLAI_API="$URL" FLAI_REPO="$REPO" \
    sh "$ROOT/install.sh" >"$LOG.install" 2>&1 || install_status=$?
  cat "$LOG.install"
  [ "$install_status" -eq 0 ] || return "$install_status"
  [ -x "$SCRATCH_HOME/.flai/bin/flai" ] \
    || { echo "the default install did not land at \$HOME/.flai/bin"; return 1; }
  expect_version "$SCRATCH_HOME/.flai/bin/flai" "$NEW" || return 1
  grep -q 'export PATH=' "$LOG.install" \
    || { echo "install.sh did not print a PATH line"; return 1; }
}

upgrade_check() {
  out=$("$ROOT/bin/flai" self-upgrade --check --repo "$REPO" 2>&1)
  check_status=$?
  printf '%s\n' "$out"
  [ "$check_status" -eq 0 ] || return "$check_status"
  printf '%s\n' "$out" | grep -qF "latest is $NEW" \
    || { echo "self-upgrade --check did not report $NEW as the latest"; return 1; }
}

# $DIR holds OLD, so self-upgrade downloads NEW over it.
upgrade_dir() {
  "$ROOT/bin/flai" self-upgrade --dir "$DIR" --repo "$REPO" || return 1
  expect_version "$DIR/flai" "$NEW"
}

# $DIR is inside this project, where self-upgrade never installs (S-0111), so
# the binary's own path is checked on a copy outside any project. Resolved, as
# flai resolves its own path: macOS's temporary directory is under /var, a link
# to /private/var (I-0045).
upgrade_own_path() {
  cp "$DIR/flai" "$OUTSIDE/flai" || return 1
  out=$("$OUTSIDE/flai" self-upgrade --check --repo "$REPO" 2>&1)
  check_status=$?
  printf '%s\n' "$out"
  [ "$check_status" -eq 0 ] || return "$check_status"
  printf '%s\n' "$out" | grep -qF "flai $NEW installed at $OUTSIDE/flai;" \
    || { echo "self-upgrade with no --dir did not resolve to the installed binary's own path"; return 1; }
}

upgrade_in_project() {
  out=$(env -i HOME="$SCRATCH_HOME" PATH="$PATH" FLAI_RELEASES_API="$URL" \
    "$DIR/flai" self-upgrade --check --repo "$REPO" 2>&1)
  check_status=$?
  printf '%s\n' "$out"
  [ "$check_status" -eq 0 ] || return "$check_status"
  printf '%s\n' "$out" | grep -qF "flai $NEW installed at $SCRATCH_HOME/.flai/bin/flai;" \
    || { echo "self-upgrade inside a project did not resolve to \$HOME/.flai/bin"; return 1; }
}

step "install.sh never invokes sudo" no_sudo
quiet_step "build bin/flai, and releases $OLD and $NEW from the tree" build_releases

say "serve releases $OLD and $NEW as $REPO on 127.0.0.1"
# The server's stderr goes to a file, not into $(...), which would wait for the
# server to close it. It stops within a second of this script's exit.
server_status=0
URL=$("$ROOT/scripts/release-server.sh" --owner "$$" --repo "$REPO" "$DIST/$OLD" "$DIST/$NEW" 2>"$SERVER_LOG") \
  || server_status=$?
if [ "$server_status" -ne 0 ]; then
  : >"$LOG"
  fail "serve the releases" "$server_status"
fi
say "release server at $URL"
FLAI_RELEASES_API="$URL"
export FLAI_RELEASES_API

step "install.sh with an explicit FLAI_INSTALL_DIR installs release $OLD" install_explicit_dir
step "install.sh with no FLAI_INSTALL_DIR installs the latest under a fresh HOME/.flai/bin, without sudo" install_default_dir
step "flai self-upgrade --check" upgrade_check
step "flai self-upgrade --dir upgrades $OLD to $NEW" upgrade_dir
OUTSIDE=$(cd "$(mktemp -d)" && pwd -P)
step "flai self-upgrade with no --dir, outside a project, resolves to the binary's own path" upgrade_own_path
step "flai self-upgrade with no --dir, run from inside a project, resolves to HOME/.flai/bin" upgrade_in_project
rm -rf "$DIR" "$SCRATCH_HOME" "$DIST" "$GH_EMPTY" "$LOG" "$LOG.install" "$SERVER_LOG"
say "passed"
