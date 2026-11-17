---
id: T-3502
type: task
nature: improvement
title: Publish what is accepted and not yet released
status: done
parent: S-0820
owner: alex
created: 2026-11-17T02:48:43Z
updated: 2026-11-17T03:26:37Z
transitions:
  - to: ready
    at: 2026-11-17T03:26:15Z
    by: agent-S-0820
  - to: in-progress
    at: 2026-11-17T03:26:15Z
    by: agent-S-0820
  - to: done
    at: 2026-11-17T03:26:37Z
    by: agent-S-0820
stream: S-0820
tags: []
after: [T-3501]
usage:
  source: log
  seconds: 22
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 7
      output: 919
      cache_read: 499254
      cache_write: 19513
      cost: 0.2744
---
# T-3502 Publish what is accepted and not yet released

## Work

Run `scripts/flai.sh release --pending --require-pending --json` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`scripts/flai.sh release --pending --require-pending --json` exits 0 in the main checkout, or exits 2, which stops the run: record it with `scripts/flai.sh formula stop S-0820 --step publish --exit 2`, giving it the step's output.

## Notes

- 2026-11-17: exited 0 in the main checkout. It tagged `flai/v1.102.0` (minor, from 1.101.4) for S-0739 and S-0747, and pushed the branch and the tag (`pushed: true`, `published: null`).
