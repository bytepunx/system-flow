---
id: T-2696
type: task
nature: improvement
title: Publish what is accepted and not yet released
status: done
parent: S-0635
owner: alex
created: 2026-11-08T01:14:12Z
updated: 2026-11-09T11:24:13Z
transitions:
  - to: ready
    at: 2026-11-09T11:24:01Z
    by: system-flow
  - to: in-progress
    at: 2026-11-09T11:24:02Z
    by: system-flow
  - to: done
    at: 2026-11-09T11:24:13Z
    by: agent-S-0635
stream: S-0635
tags: []
after: [T-2695]
usage:
  source: log
  seconds: 11
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 2
      output: 308
      cache_read: 151050
      cache_write: 6495
      cost: 0.0883
---
# T-2696 Publish what is accepted and not yet released

## Work

Run `flai release --pending --json` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --pending --json` exits 0 in the main checkout.

## Notes

- 2026-11-09: `flai release --pending --json` (through `scripts/flai.sh`) ran in the main checkout and exited 0: it tagged `flai/v1.86.3` (a patch from 1.86.2, over S-0603, S-0606, S-0608, S-0612, S-0613, S-0614, and S-0617), and pushed the branch and the tag (`pushed: true`). No other component had anything pending.
