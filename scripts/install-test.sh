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
echo "install-test: install.sh"
FLAI_INSTALL_DIR="$DIR" sh "$ROOT/install.sh"
"$DIR/flai" version | grep -q '^flai [0-9]' || { echo "install-test: installed binary does not report a version" >&2; exit 1; }
echo "install-test: flai self-upgrade --check"
"$ROOT/scripts/flai-build.sh" >/dev/null
"$ROOT/bin/flai" self-upgrade --check
echo "install-test: flai self-upgrade --dir"
"$ROOT/bin/flai" self-upgrade --dir "$DIR"
"$DIR/flai" version | head -1
