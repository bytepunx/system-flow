---
id: T-2568
type: task
nature: improvement
title: Publish what is accepted and not yet released
status: done
parent: S-0588
owner: alex
created: 2026-11-07T09:24:05Z
updated: 2026-11-07T12:07:11Z
transitions:
  - to: ready
    at: 2026-11-07T12:06:41Z
    by: agent-S-0588
  - to: in-progress
    at: 2026-11-07T12:06:42Z
    by: agent-S-0588
  - to: done
    at: 2026-11-07T12:07:11Z
    by: agent-S-0588
stream: S-0588
tags: []
after: [T-2567]
usage:
  source: log
  seconds: 29
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 7
      output: 810
      cache_read: 447956
      cache_write: 18985
      cost: 0.2577
---
# T-2568 Publish what is accepted and not yet released

## Work

Run `flai release --pending --json` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --pending --json` exits 0 in the main checkout.

## Notes

- 2026-11-07: `flai release --pending --json` ran in the main checkout and exited 0, with `pushed: true`. It tagged `flai/v1.85.0` (minor, from 1.84.1; S-0406 minor, S-0560, S-0561, S-0568, S-0586, S-0589, S-0590 patch) and `flaiover/v0.57.2` (patch, from 0.57.1; S-0586), bumped the template to 1.1.25 (S-0406) and published it to `git@github.com:bytepunx/system-flow-template.git` as `v1.1.25` (6fbbb46), and pushed the branch and tags.
