#!/usr/bin/env sh
# Fetch the remote's branch and tags into this checkout: the first step of
# the publish formula (design/workflows/formulas/publish.md, ADR-0240), which
# runs it from the project's main checkout before flai release --pending,
# since flai never fetches.
set -eu
cd "$(dirname "$0")/.."
echo "fetch.sh: git fetch --tags"
git fetch --tags
