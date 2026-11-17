---
id: T-3515
type: task
nature: improvement
title: Publish what is accepted and not yet released
status: done
parent: S-0825
owner: alex
created: 2026-11-17T07:54:16Z
updated: 2026-11-17T08:16:21Z
transitions:
  - to: ready
    at: 2026-11-17T08:15:24Z
    by: agent-S-0825
  - to: in-progress
    at: 2026-11-17T08:15:25Z
    by: agent-S-0825
  - to: done
    at: 2026-11-17T08:16:21Z
    by: agent-S-0825
stream: S-0825
tags: []
after: [T-3514]
usage:
  source: log
  seconds: 56
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 5
      output: 690
      cache_read: 344970
      cache_write: 16082
      cost: 0.2115
---
# T-3515 Publish what is accepted and not yet released

## Work

Run `scripts/flai.sh release --pending --require-pending --json` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`scripts/flai.sh release --pending --require-pending --json` exits 0 in the main checkout, or exits 2, which stops the run: record it with `scripts/flai.sh formula stop S-0825 --step publish --exit 2`, giving it the step's output.

## Notes

- 2026-11-17: exited 0 in the main checkout, `pushed: true`, tags `flai/v1.103.0` and `flaiover/v0.62.0`.

| Component | From | To | Tag | Items |
|-----------|------|----|-----|-------|
| flai | 1.102.0 | 1.103.0 (minor) | `flai/v1.103.0` | S-0742 (minor), S-0753, S-0817, S-0819 (patch) |
| flaiover | 0.61.3 | 0.62.0 (minor) | `flaiover/v0.62.0` | S-0753 (minor) |
