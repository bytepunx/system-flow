#!/usr/bin/env sh
# Smoke tests: end-to-end. Render the template and check it; check this repository.
set -eu
. "$(dirname "$0")/env.sh"
echo "smoke: template render and check"
"$ROOT/scripts/template-test.sh"
echo "smoke: repository check"
"$ROOT/scripts/check.sh"
