#!/usr/bin/env sh
# Regenerate flai/internal/mdlint/testdata/cases/expected.txt: the rule and
# line markdownlint-cli2, at the version lint-md.sh runs, reports for each
# finding in mdlint's fixtures. flai's own lint is tested against it.
set -eu
. "$(dirname "$0")/env.sh"
cd "$ROOT/flai/internal/mdlint/testdata/cases"
command -v npx >/dev/null 2>&1 || { echo "mdlint-fixtures: npx not found; install Node" >&2; exit 1; }
MDLINT_VERSION="${MDLINT_VERSION:-0.20.0}"
npx --yes "markdownlint-cli2@$MDLINT_VERSION" "*.md" 2>&1 \
  | sed -n 's/^\([^ :]*\.md\):\([0-9]*\)[: ].* \(MD[0-9]*\)\/.*/\1:\2 \3/p' \
  | sort -t: -k1,1 -k2,2n > expected.txt
echo "mdlint-fixtures: $(wc -l < expected.txt | tr -d ' ') findings in expected.txt"
