#!/usr/bin/env sh
# Smoke tests: end-to-end. Render the template and check it; check this repository.
# Under flai verify, in a close-out or a merge queue batch, the repository check
# and the markdown lint run here: TestMonorepoIsClean, in the integration tier,
# skips there, and so does TestRepositoryLintsClean unless the stories change
# mdlint or the markdownlint configuration (S-0548). The gateway checks are not
# here: scripts/gateway-smoke.sh is the gateway tier, which runs only by name,
# flai test gateway, or on the host's schedule, as it spends gateway money.
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
