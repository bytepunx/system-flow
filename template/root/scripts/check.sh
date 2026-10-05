#!/usr/bin/env sh
# Validate this repository against the system-flow standard.
# In a close-out, which exports the story as CLOSE_OUT_STORY, the check is
# scoped to it: findings outside the story are notes the run passes over, and
# each rule's are recorded in an issue under design/issues, which the
# close-out commits. Unset, as in CI and make, every finding counts.
set -eu
cd "$(dirname "$0")/.."
if [ -n "${CLOSE_OUT_STORY:-}" ]; then
  set -- --story "$CLOSE_OUT_STORY" --record-issues "$@"
fi
exec flai check --strict "$@"
