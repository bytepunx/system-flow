#!/usr/bin/env sh
# Validate this repository against the system-flow standard.
# In a close-out, which exports the story as CLOSE_OUT_STORY, the check is
# scoped to it: findings outside the story are notes the run passes over, and
# each rule's are recorded in an issue under design/issues, which the
# close-out commits (S-0249). In a merge queue batch, which exports its
# stories comma-separated as MERGE_QUEUE_STORIES, the check is scoped to them:
# a finding outside every one is a note, and none is recorded, since the
# batch's checkout is never committed (I-0174). Neither set, as in CI and
# make, every finding counts.
set -eu
. "$(dirname "$0")/env.sh"
cd "$ROOT"
if [ -n "${CLOSE_OUT_STORY:-}" ]; then
  set -- --story "$CLOSE_OUT_STORY" --record-issues "$@"
elif [ -n "${MERGE_QUEUE_STORIES:-}" ]; then
  set -- --story "$MERGE_QUEUE_STORIES" "$@"
fi
exec "$ROOT/scripts/flai.sh" check --strict "$@"
