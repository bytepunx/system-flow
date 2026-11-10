---
id: T-2892
type: task
nature: improvement
title: Publish what is accepted and not yet released
status: done
parent: S-0664
owner: alex
created: 2026-11-10T08:33:20Z
updated: 2026-11-10T08:35:20Z
transitions:
  - to: ready
    at: 2026-11-10T08:35:11Z
    by: agent-S-0664
  - to: in-progress
    at: 2026-11-10T08:35:11Z
    by: agent-S-0664
  - to: done
    at: 2026-11-10T08:35:20Z
    by: agent-S-0664
stream: S-0664
tags: []
after: [T-2891]
usage:
  source: log
  seconds: 9
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 2
      output: 311
      cache_read: 161435
      cache_write: 5406
      cost: 0.0818
---
# T-2892 Publish what is accepted and not yet released

## Work

Run `flai release --pending --require-pending --json` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --pending --require-pending --json` exits 0 in the main checkout, or exits 2, which stops the run: record it with `flai formula stop S-0664 --step publish --exit 2`, giving it the step's output.

## Notes

Exited 0 in the main checkout, run as `scripts/flai.sh release --pending --require-pending --json`, the tree's flai, as for the evaluate step. It released flai 1.86.11 to 1.86.12 (patch), tag `flai/v1.86.12`, bundling S-0652 and S-0637, and pushed the branch and the tag (`pushed: true`). No other component had anything pending.
