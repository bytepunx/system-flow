#!/usr/bin/env sh
# Lint every markdown file with .markdownlint.yaml (markdownlint-cli2 through npx).
set -eu
cd "$(dirname "$0")/.."
command -v npx >/dev/null 2>&1 || { echo "lint-md: npx not found; install Node" >&2; exit 1; }
npx --yes markdownlint-cli2 "**/*.md" "!**/node_modules/**" "!**/testdata/**" "!bin/**"
