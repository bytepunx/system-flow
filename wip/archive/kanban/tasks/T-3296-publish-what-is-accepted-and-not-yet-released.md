---
id: T-3296
type: task
nature: improvement
title: Publish what is accepted and not yet released
status: done
parent: S-0769
owner: alex
created: 2026-11-13T12:40:07Z
updated: 2026-11-13T12:54:44Z
transitions:
  - to: ready
    at: 2026-11-13T12:54:30Z
    by: agent-S-0769
  - to: in-progress
    at: 2026-11-13T12:54:30Z
    by: agent-S-0769
  - to: done
    at: 2026-11-13T12:54:44Z
    by: agent-S-0769
stream: S-0769
tags: []
after: [T-3295]
usage:
  source: log
  seconds: 14
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 5
      output: 679
      cache_read: 351782
      cache_write: 16213
      cost: 0.2137
---
# T-3296 Publish what is accepted and not yet released

## Work

Run `flai release --pending --require-pending --json` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --pending --require-pending --json` exits 0 in the main checkout, or exits 2, which stops the run: record it with `flai formula stop S-0769 --step publish --exit 2`, giving it the step's output.

## Notes

- 2026-11-13: `flai release --pending --require-pending --json` ran in the main checkout and exited 0, with `pushed: true`. It tagged `flai/v1.90.0` (minor, from 1.89.1; S-0763 minor, S-0767 patch) and pushed the branch and the tag. No other component had anything to release.
