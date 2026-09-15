#!/usr/bin/env sh
# Smoke tests: end-to-end paths. Run after integration passes and in CI before merge.
# Replace the body with this project's command; keep the exit code honest.
set -eu
cd "$(dirname "$0")/.."
echo "smoke.sh: no tests configured yet; see design/conventions/code-quality.md"
