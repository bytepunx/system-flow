#!/usr/bin/env sh
# Lint every markdown file with the same globs and config as system-flow-check.yml.
# Runs markdownlint-cli2 through pnpm dlx with the project's pnpm cache, so nothing
# is written outside the repository.
set -eu
. "$(dirname "$0")/env.sh"
cd "$ROOT"
MDLINT_VERSION="${MDLINT_VERSION:-0.20.0}"
echo "lint-md: markdownlint-cli2 $MDLINT_VERSION"
pnpm --dir "$ROOT/flaiover" dlx "markdownlint-cli2@$MDLINT_VERSION" \
  "**/*.md" "!**/node_modules/**" "!**/testdata/**" "!bin/**" "!flaiover/build/**" "!**/.svelte-kit/**" "!.flai-cache/**"
