---
id: T-2591
type: task
nature: improvement
title: Publish what is accepted and not yet released
status: done
parent: S-0594
owner: alex
created: 2026-11-07T21:47:12Z
updated: 2026-11-07T21:59:27Z
transitions:
  - to: ready
    at: 2026-11-07T21:58:10Z
    by: system-flow
  - to: in-progress
    at: 2026-11-07T21:58:17Z
    by: system-flow
  - to: done
    at: 2026-11-07T21:59:27Z
    by: agent-S-0594
stream: S-0594
tags: []
after: [T-2590]
usage:
  source: log
  seconds: 70
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 9
      output: 1238
      cache_read: 630283
      cache_write: 21986
      cost: 0.3267
---
# T-2591 Publish what is accepted and not yet released

## Work

Run `flai release --pending --json` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --pending --json` exits 0 in the main checkout.

## Notes

- 2026-11-07: `flai release --pending --json` ran in the main checkout (`scripts/flai.sh`) and exited 0 with `pushed: true`. It released flai 1.85.1 → 1.86.0, tag `flai/v1.86.0` (minor, S-0407), and the template 1.1.26 → 1.1.27 (patch, S-0407), published to `git@github.com:bytepunx/system-flow-template.git` as v1.1.27 (78fc638).
