---
id: T-3283
type: task
nature: improvement
title: Publish what is accepted and not yet released
status: done
parent: S-0765
owner: alex
created: 2026-11-13T09:17:34Z
updated: 2026-11-13T09:19:14Z
transitions:
  - to: ready
    at: 2026-11-13T09:19:01Z
    by: system-flow
  - to: in-progress
    at: 2026-11-13T09:19:01Z
    by: system-flow
  - to: done
    at: 2026-11-13T09:19:14Z
    by: agent-S-0765
stream: S-0765
tags: []
after: [T-3282]
usage:
  source: log
  seconds: 13
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 5
      output: 681
      cache_read: 327567
      cache_write: 17719
      cost: 0.2209
---
# T-3283 Publish what is accepted and not yet released

## Work

Run `flai release --pending --require-pending --json` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --pending --require-pending --json` exits 0 in the main checkout, or exits 2, which stops the run: record it with `flai formula stop S-0765 --step publish --exit 2`, giving it the step's output.

## Notes

- 2026-11-13: ran in the main checkout and exited 0. It released flai 1.87.0 → 1.88.0 (minor, for S-0735), tagged `flai/v1.88.0`, and pushed (`"pushed": true`).
