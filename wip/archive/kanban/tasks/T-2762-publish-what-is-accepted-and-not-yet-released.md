---
id: T-2762
type: task
nature: improvement
title: Publish what is accepted and not yet released
status: done
parent: S-0642
owner: alex
created: 2026-11-09T17:28:50Z
updated: 2026-11-09T21:54:31Z
transitions:
  - to: ready
    at: 2026-11-09T21:54:06Z
    by: system-flow
  - to: in-progress
    at: 2026-11-09T21:54:06Z
    by: system-flow
  - to: done
    at: 2026-11-09T21:54:31Z
    by: agent-S-0642
stream: S-0642
tags: []
after: [T-2761]
usage:
  source: log
  seconds: 25
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 8
      output: 133
      cache_read: 640319
      cache_write: 3904
      cost: 0.2355
---
# T-2762 Publish what is accepted and not yet released

## Work

Run `flai release --pending --json` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --pending --json` exits 0 in the main checkout.

## Notes

- 2026-11-09: `flai release --pending --json` ran in the main checkout and exited 0, with `pushed: true`.
  - flai 1.86.5 → 1.86.6, patch, tag `flai/v1.86.6`: S-0619, S-0620, S-0611, S-0629, S-0621, S-0624
  - flaiover 0.58.0 → 0.58.1, patch, tag `flaiover/v0.58.1`: S-0626, S-0611, S-0627
  - template 1.1.28 → 1.1.29, patch, published to `bytepunx/system-flow-template` as v1.1.29 (5f117cf): S-0621
