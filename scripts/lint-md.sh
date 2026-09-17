#!/usr/bin/env sh
# Lint every markdown file with the same globs and config as system-flow-check.yml.
# Runs markdownlint-cli2 through npx (Node only; the npm cache lives in .flai-cache
# via env.sh), so it works on a fresh clone and in CI without pnpm.
set -eu
. "$(dirname "$0")/env.sh"
cd "$ROOT"
command -v npx >/dev/null 2>&1 || { echo "lint-md: npx not found; install Node" >&2; exit 1; }
MDLINT_VERSION="${MDLINT_VERSION:-0.20.0}"
echo "lint-md: markdownlint-cli2 $MDLINT_VERSION"
npx --yes "markdownlint-cli2@$MDLINT_VERSION" \
  "**/*.md" "!**/node_modules/**" "!**/testdata/**" "!bin/**" "!flaiover/build/**" "!**/.svelte-kit/**" "!.flai-cache/**"
