---
id: T-3336
type: task
nature: improvement
title: Publish what is accepted and not yet released
status: done
parent: S-0780
owner: alex
created: 2026-11-13T20:51:49Z
updated: 2026-11-13T20:55:26Z
transitions:
  - to: ready
    at: 2026-11-13T20:54:02Z
    by: agent-S-0780
  - to: in-progress
    at: 2026-11-13T20:54:02Z
    by: agent-S-0780
  - to: done
    at: 2026-11-13T20:55:26Z
    by: agent-S-0780
stream: S-0780
tags: []
after: [T-3335]
usage:
  source: log
  seconds: 84
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 9
      output: 1344
      cache_read: 679886
      cache_write: 25142
      cost: 0.364
---
# T-3336 Publish what is accepted and not yet released

## Work

Run `scripts/flai.sh release --pending --require-pending --json` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`scripts/flai.sh release --pending --require-pending --json` exits 0 in the main checkout, or exits 2, which stops the run: record it with `scripts/flai.sh formula stop S-0780 --step publish --exit 2`, giving it the step's output.

## Notes

- 2026-11-13: ran in the main checkout and exited 0. It released flai 1.95.1 → 1.95.2 (patch, for S-0777), tagged `flai/v1.95.2`, and pushed (`pushed: true`).
