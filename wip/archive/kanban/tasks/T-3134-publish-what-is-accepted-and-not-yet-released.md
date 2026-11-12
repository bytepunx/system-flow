---
id: T-3134
type: task
nature: improvement
title: Publish what is accepted and not yet released
status: done
parent: S-0729
owner: alex
created: 2026-11-12T03:39:40Z
updated: 2026-11-12T03:42:05Z
transitions:
  - to: ready
    at: 2026-11-12T03:41:47Z
    by: agent-S-0729
  - to: in-progress
    at: 2026-11-12T03:41:47Z
    by: agent-S-0729
  - to: done
    at: 2026-11-12T03:42:05Z
    by: agent-S-0729
stream: S-0729
tags: []
after: [T-3133]
usage:
  source: log
  seconds: 18
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 5
      output: 583
      cache_read: 322270
      cache_write: 11968
      cost: 0.1719
---
# T-3134 Publish what is accepted and not yet released

## Work

Run `flai release --pending --require-pending --json` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --pending --require-pending --json` exits 0 in the main checkout, or exits 2, which stops the run: record it with `flai formula stop S-0729 --step publish --exit 2`, giving it the step's output.

## Notes

- 2026-11-12: exited 0 in the main checkout. It tagged `flai/v1.86.21` (patch from 1.86.20) for S-0711 and S-0712, and pushed the branch and the tag (`pushed: true`).
