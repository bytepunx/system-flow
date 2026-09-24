#!/usr/bin/env sh
# Install the latest published flai with install.sh into the repo cache and
# check it runs, then exercise the same path through `flai self-upgrade` from
# the source build. Needs GitHub API access: GITHUB_TOKEN, GH_TOKEN, or gh.
set -eu
. "$(dirname "$0")/env.sh"
DIR="$ROOT/.flai-cache/install-test"
rm -rf "$DIR"
if [ -z "${GITHUB_TOKEN:-}${GH_TOKEN:-}" ] && ! command -v gh >/dev/null 2>&1; then
  echo "install-test: no GITHUB_TOKEN, GH_TOKEN, or gh; cannot reach the private repository" >&2
  exit 1
fi

echo "install-test: install.sh never invokes sudo"
! grep -vE '^\s*#' "$ROOT/install.sh" | grep -qE '(^|[^A-Za-z0-9_-])sudo([^A-Za-z0-9_-]|$)' \
  || { echo "install-test: install.sh still invokes sudo outside a comment" >&2; exit 1; }

echo "install-test: install.sh with an explicit FLAI_INSTALL_DIR"
FLAI_INSTALL_DIR="$DIR" sh "$ROOT/install.sh"
"$DIR/flai" version | grep -q '^flai [0-9]' || { echo "install-test: installed binary does not report a version" >&2; exit 1; }

echo "install-test: install.sh with no FLAI_INSTALL_DIR installs under a fresh HOME/.flai/bin, without sudo"
# gh reads its own auth from $HOME, so the token is resolved now, before HOME
# is redirected below, and passed through explicitly.
TOKEN="${GITHUB_TOKEN:-${GH_TOKEN:-}}"
[ -n "$TOKEN" ] || TOKEN=$(gh auth token 2>/dev/null || true)
SCRATCH_HOME="$ROOT/.flai-cache/install-test-home"
rm -rf "$SCRATCH_HOME"
OUT="$ROOT/.flai-cache/install-test-home.log"
env -i HOME="$SCRATCH_HOME" PATH="$PATH" GITHUB_TOKEN="$TOKEN" \
  sh "$ROOT/install.sh" >"$OUT" 2>&1
[ -x "$SCRATCH_HOME/.flai/bin/flai" ] || { echo "install-test: default install did not land at \$HOME/.flai/bin" >&2; cat "$OUT" >&2; exit 1; }
"$SCRATCH_HOME/.flai/bin/flai" version | grep -q '^flai [0-9]' || { echo "install-test: default-path binary does not report a version" >&2; exit 1; }
grep -q 'export PATH=' "$OUT" || { echo "install-test: install.sh did not print a PATH line" >&2; cat "$OUT" >&2; exit 1; }
rm -rf "$SCRATCH_HOME" "$OUT"

echo "install-test: flai self-upgrade --check"
"$ROOT/scripts/flai-build.sh" >/dev/null
"$ROOT/bin/flai" self-upgrade --check
echo "install-test: flai self-upgrade --dir"
"$ROOT/bin/flai" self-upgrade --dir "$DIR"
"$DIR/flai" version | head -1

echo "install-test: flai self-upgrade with no --dir needs no sudo, once installed under a directory the user owns"
# $DIR is inside this project, where self-upgrade never installs (S-0111), so
# the binary's own path is checked on a copy outside any project.
OUTSIDE=$(mktemp -d)
trap 'rm -rf "$OUTSIDE"' EXIT
cp "$DIR/flai" "$OUTSIDE/flai"
"$OUTSIDE/flai" self-upgrade --check | grep -q 'installed at '"$OUTSIDE"'/flai' \
  || { echo "install-test: self-upgrade with no --dir did not resolve to the installed binary's own path" >&2; exit 1; }

echo "install-test: flai self-upgrade with no --dir, run from inside a project, resolves to HOME/.flai/bin"
env -i HOME="$SCRATCH_HOME" PATH="$PATH" GITHUB_TOKEN="$TOKEN" \
  "$DIR/flai" self-upgrade --check | grep -q 'installed at '"$SCRATCH_HOME"'/.flai/bin/flai' \
  || { echo "install-test: self-upgrade inside a project did not resolve to \$HOME/.flai/bin" >&2; exit 1; }
