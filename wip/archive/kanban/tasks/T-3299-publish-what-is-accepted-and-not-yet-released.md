---
id: T-3299
type: task
nature: improvement
title: Publish what is accepted and not yet released
status: done
parent: S-0770
owner: alex
created: 2026-11-13T13:16:45Z
updated: 2026-11-13T13:19:01Z
transitions:
  - to: ready
    at: 2026-11-13T13:18:40Z
    by: agent-S-0770
  - to: in-progress
    at: 2026-11-13T13:18:40Z
    by: agent-S-0770
  - to: done
    at: 2026-11-13T13:19:01Z
    by: agent-S-0770
stream: S-0770
tags: []
after: [T-3298]
usage:
  source: log
  seconds: 21
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 5
      output: 783
      cache_read: 337681
      cache_write: 15765
      cost: 0.2093
---
# T-3299 Publish what is accepted and not yet released

## Work

Run `flai release --pending --require-pending --json` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --pending --require-pending --json` exits 0 in the main checkout, or exits 2, which stops the run: record it with `flai formula stop S-0770 --step publish --exit 2`, giving it the step's output.

## Notes

- 2026-11-13: ran as `scripts/flai.sh release --pending --require-pending --json` (the manifest's `flai.command`) in the main checkout; exited 0. It released flai 1.90.0 → 1.91.0 (minor, for S-0740), tagged `flai/v1.91.0`, and pushed the branch and the tag (`pushed: true`).
