---
id: T-2588
type: task
nature: improvement
title: Publish what is accepted and not yet released
status: done
parent: S-0593
owner: alex
created: 2026-11-07T15:17:08Z
updated: 2026-11-07T15:19:29Z
transitions:
  - to: ready
    at: 2026-11-07T15:19:03Z
    by: system-flow
  - to: in-progress
    at: 2026-11-07T15:19:04Z
    by: system-flow
  - to: done
    at: 2026-11-07T15:19:29Z
    by: agent-S-0593
stream: S-0593
tags: []
after: [T-2587]
usage:
  source: log
  seconds: 25
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 5
      output: 580
      cache_read: 300367
      cache_write: 12458
      cost: 0.1714
---
# T-2588 Publish what is accepted and not yet released

## Work

Run `flai release --pending --json` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --pending --json` exits 0 in the main checkout.

## Notes

- 2026-11-07: `flai release --pending --json` (run as `scripts/flai.sh`) in the main checkout exited 0, with `pushed: true`. It tagged `flai/v1.85.1` (patch, from 1.85.0; S-0587) and `flaiover/v0.57.3` (patch, from 0.57.2; S-0587), bumped the template to 1.1.26 (patch, from 1.1.25; S-0587) and published it to `git@github.com:bytepunx/system-flow-template.git` as `v1.1.26` (77ece98), and pushed the branch and tags.
