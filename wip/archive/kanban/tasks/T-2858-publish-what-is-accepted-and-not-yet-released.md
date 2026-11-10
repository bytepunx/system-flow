---
id: T-2858
type: task
nature: improvement
title: Publish what is accepted and not yet released
status: done
parent: S-0658
owner: alex
created: 2026-11-10T07:56:53Z
updated: 2026-11-10T08:29:25Z
transitions:
  - to: ready
    at: 2026-11-10T08:29:10Z
    by: system-flow
  - to: in-progress
    at: 2026-11-10T08:29:11Z
    by: system-flow
  - to: done
    at: 2026-11-10T08:29:25Z
    by: agent-S-0658
stream: S-0658
tags: []
after: [T-2857]
usage:
  source: log
  seconds: 14
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 5
      output: 564
      cache_read: 331999
      cache_write: 9911
      cost: 0.157
---
# T-2858 Publish what is accepted and not yet released

## Work

Run `flai release --pending --require-pending --json` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --pending --require-pending --json` exits 0 in the main checkout, or exits 2, which stops the run: record it with `flai formula stop S-0658 --step publish --exit 2`, giving it the step's output.

## Notes

Exited 0 in the main checkout, run as `scripts/flai.sh release --pending --require-pending --json`, the tree's flai, as for the evaluate step. It released flai 1.86.10 to 1.86.11 (patch), tag `flai/v1.86.11`, bundling S-0651, S-0653, and S-0654, and pushed the branch and the tag (`pushed: true`). No other component had anything pending.
