#!/usr/bin/env sh
# Smoke tests: end-to-end. Render the template and check it; check this repository.
set -eu
. "$(dirname "$0")/env.sh"
echo "smoke: template render and check"
"$ROOT/scripts/template-test.sh"
echo "smoke: repository check"
"$ROOT/scripts/check.sh"
echo "smoke: markdown lint"
"$ROOT/scripts/lint-md.sh"
echo "smoke: installer and self-upgrade, from a release built from the tree and served on 127.0.0.1"
"$ROOT/scripts/install-test.sh"
