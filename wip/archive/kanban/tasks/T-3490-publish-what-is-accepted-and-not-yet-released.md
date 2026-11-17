---
id: T-3490
type: task
nature: improvement
title: Publish what is accepted and not yet released
status: done
parent: S-0816
owner: alex
created: 2026-11-17T01:33:37Z
updated: 2026-11-17T01:37:10Z
transitions:
  - to: ready
    at: 2026-11-17T01:36:44Z
    by: agent-S-0816
  - to: in-progress
    at: 2026-11-17T01:36:45Z
    by: agent-S-0816
  - to: done
    at: 2026-11-17T01:37:10Z
    by: agent-S-0816
stream: S-0816
tags: []
after: [T-3489]
usage:
  source: log
  seconds: 25
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 8
      output: 961
      cache_read: 502568
      cache_write: 18877
      cost: 0.2708
---
# T-3490 Publish what is accepted and not yet released

## Work

Run `scripts/flai.sh release --pending --require-pending --json` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`scripts/flai.sh release --pending --require-pending --json` exits 0 in the main checkout, or exits 2, which stops the run: record it with `scripts/flai.sh formula stop S-0816 --step publish --exit 2`, giving it the step's output.

## Notes

- 2026-11-17: exited 0 in the main checkout, with `pushed: true`. It released S-0804 as a patch for two components:

| Component | From | To | Tag |
|-----------|------|----|-----|
| flai | 1.101.2 | 1.101.3 | `flai/v1.101.3` |
| flaiover | 0.61.2 | 0.61.3 | `flaiover/v0.61.3` |
