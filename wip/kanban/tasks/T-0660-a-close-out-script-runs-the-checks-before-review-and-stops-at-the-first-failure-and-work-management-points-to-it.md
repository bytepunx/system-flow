---
id: T-0660
type: task
nature: feature
title: A close-out script runs the checks before review and stops at the first failure, and work-management points to it
status: done
parent: S-0187
owner: arobson
created: 2026-10-01T09:26:25Z
updated: 2026-10-01T09:29:52Z
transitions:
  - to: ready
    at: 2026-10-01T09:26:41Z
    by: agent-S-0187
  - to: in-progress
    at: 2026-10-01T09:27:07Z
    by: agent-S-0187
  - to: done
    at: 2026-10-01T09:29:52Z
    by: agent-S-0187
stream: S-0187
tags: []
usage:
  source: log
  seconds: 165
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 46
      output: 13273
      cache_read: 1835994
      cache_write: 45096
      cost: 0.9467
---

# T-0660 A close-out script runs the checks before review and stops at the first failure, and work-management points to it

## Work

Write `scripts/close-out.sh`, here and in `template/root/scripts`, POSIX `sh` under `set -eu`: given a story ID, run lint, the tests, and `flai check --strict`, check the narrative's `## Current state` and `## Next steps` are not empty, then commit what is outstanding with the message given, and stop at the first step that fails. Add a `close-out` Make target and a line in `scripts/README.md`. The baseline `work-management.md` definition of done points to it.

## Done when

- `scripts/close-out.sh` exits non-zero at the first failing step and commits nothing then.
- Both copies of `work-management.md` point to it, and the design says what it runs.

## Notes
