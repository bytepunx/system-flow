---
id: T-3494
type: task
nature: improvement
title: Publish what is accepted and not yet released
status: done
parent: S-0818
owner: alex
created: 2026-11-17T02:27:29Z
updated: 2026-11-17T02:31:20Z
transitions:
  - to: ready
    at: 2026-11-17T02:30:31Z
    by: agent-S-0818
  - to: in-progress
    at: 2026-11-17T02:30:31Z
    by: agent-S-0818
  - to: done
    at: 2026-11-17T02:31:20Z
    by: agent-S-0818
stream: S-0818
tags: []
after: [T-3493]
usage:
  source: log
  seconds: 49
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 10
      output: 1148
      cache_read: 675218
      cache_write: 22960
      cost: 0.3417
---
# T-3494 Publish what is accepted and not yet released

## Work

Run `scripts/flai.sh release --pending --require-pending --json` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`scripts/flai.sh release --pending --require-pending --json` exits 0 in the main checkout, or exits 2, which stops the run: record it with `scripts/flai.sh formula stop S-0818 --step publish --exit 2`, giving it the step's output.

## Notes

- 2026-11-17: ran in the main checkout and exited 0. Released flai 1.101.3 → 1.101.4 (patch, S-0797), tag `flai/v1.101.4`; `pushed: true`.
