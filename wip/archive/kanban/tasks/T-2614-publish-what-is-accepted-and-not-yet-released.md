---
id: T-2614
type: task
nature: improvement
title: Publish what is accepted and not yet released
status: done
parent: S-0597
owner: alex
created: 2026-11-07T22:31:23Z
updated: 2026-11-07T23:07:40Z
transitions:
  - to: ready
    at: 2026-11-07T23:07:03Z
    by: system-flow
  - to: in-progress
    at: 2026-11-07T23:07:04Z
    by: system-flow
  - to: done
    at: 2026-11-07T23:07:40Z
    by: agent-S-0597
stream: S-0597
tags: []
after: [T-2613]
usage:
  source: log
  seconds: 36
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 11
      output: 1400
      cache_read: 784688
      cache_write: 26716
      cost: 0.3987
---
# T-2614 Publish what is accepted and not yet released

## Work

Run `flai release --pending --json` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --pending --json` exits 0 in the main checkout.

## Notes

- 2026-11-07: `scripts/flai.sh release --pending --json` ran in the main checkout and exited 0, with `pushed: true`. It committed `3ec051f28 chore: publish flai 1.86.1, flaiover 0.58.0` and pushed main and both tags; main is now level with `origin/main`.

| Tag | From | Level | Items |
|-----|------|-------|-------|
| `flai/v1.86.1` | 1.86.0 | patch | S-0592 (patch), S-0595 (patch) |
| `flaiover/v0.58.0` | 0.57.3 | minor | S-0595 (minor) |
