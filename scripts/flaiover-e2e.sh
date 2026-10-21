#!/usr/bin/env sh
# flaiover's end-to-end tests: the production build served by server.js, driven
# in Playwright's Chromium (flaiover/playwright.config.ts, flaiover/e2e). It
# installs flaiover's dependencies when missing or stale, and Playwright's
# Chromium when it is not under ~/.cache/ms-playwright yet (playwright install
# skips a browser it has), then builds and runs playwright test with any
# arguments given, such as a test file or --headed.
# The server listens on 127.0.0.1:4173; FLAIOVER_E2E_PORT picks another port.
set -eu
. "$(dirname "$0")/env.sh"
"$ROOT/scripts/flaiover-install.sh" --if-needed
cd "$ROOT/flaiover"
pnpm exec playwright install chromium
echo "flaiover build"
pnpm build
echo "playwright test"
exec pnpm exec playwright test "$@"
